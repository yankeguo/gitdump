package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/yankeguo/conc"
	"github.com/yankeguo/gitdump"
	"github.com/yankeguo/rg"
	"gopkg.in/yaml.v3"

	_ "github.com/yankeguo/gitdump/pkg/hostings"
)

type Options struct {
	Dir         string `yaml:"dir" default:"." validate:"required"`
	Concurrency int    `yaml:"concurrency" default:"3" validate:"gt=0"`
	Accounts    []struct {
		Vendor   string `yaml:"vendor" validate:"required"`
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"accounts"`
}

func main() {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()
	defer rg.Guard(&err)

	var ctx context.Context
	{
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(context.Background())
		go func() {
			chSig := make(chan os.Signal, 1)
			signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)
			sig := <-chSig
			log.Println("received signal:", sig.String())
			cancel()
		}()
	}

	var (
		optConf string
	)
	flag.StringVar(&optConf, "conf", "config.yaml", "config file")
	flag.Parse()

	buf := rg.Must(os.ReadFile(optConf))

	var opts Options
	rg.Must0(yaml.Unmarshal(buf, &opts))
	rg.Must0(defaults.Set(&opts))
	rg.Must0(validator.New().Struct(&opts))

	var (
		tasksOptions []gitdump.MirrorGitOptions
		duplicated   = map[string]struct{}{}
		askPassFiles []string
	)

	defer func() {
		for _, askPassFile := range askPassFiles {
			_ = os.Remove(askPassFile)
		}
	}()

	for i, account := range opts.Accounts {
		var name string

		if account.URL == "" {
			name = fmt.Sprintf("#%d %s / %s", i+1, account.Vendor, account.Username)
		} else {
			name = fmt.Sprintf("#%d %s (%s) / %s", i+1, account.Vendor, account.URL, account.Username)
		}

		log.Println("scanning on:", name)

		hosting := gitdump.GetHosting(account.Vendor)
		if hosting == nil {
			err = errors.New("hosting not supported: " + account.Vendor)
			return
		}

		// create askPass script
		askPassDigest := md5.Sum([]byte(account.Vendor + account.URL + account.Username + account.Password))
		askPassFile := filepath.Join(os.TempDir(), "gitdumpcredentials-"+hex.EncodeToString(askPassDigest[:])+".sh")
		askPassFiles = append(askPassFiles, askPassFile)

		rg.Must0(
			os.WriteFile(
				askPassFile,
				[]byte(
					fmt.Sprintf(
						"#!/bin/sh\necho '%s' | base64 -d",
						base64.StdEncoding.EncodeToString([]byte(account.Password)),
					),
				),
				0755,
			),
		)

		var repos []gitdump.HostingRepo
		if repos, err = hosting.List(ctx, gitdump.HostingOptions{
			URL:      account.URL,
			Username: account.Username,
			Password: account.Password,
		}); err != nil {
			log.Println("failed to scan repos:", name, ":", err.Error())
			err = nil
			continue
		}

		for _, repo := range repos {
			if _, ok := duplicated[repo.URL]; ok {
				continue
			}
			duplicated[repo.URL] = struct{}{}

			log.Println("found:", repo.URL)

			tasksOptions = append(tasksOptions, gitdump.MirrorGitOptions{
				Dir:         filepath.Join(opts.Dir, repo.SubDir),
				URL:         repo.URL,
				Username:    repo.Username,
				AskPassFile: askPassFile,
			})
		}
	}

	var (
		tasks  []conc.Task
		taskID int64
	)

	for _, _taskOptions := range tasksOptions {
		var taskOptions = _taskOptions

		tasks = append(tasks, conc.TaskFunc(func(ctx context.Context) error {
			log.Printf(
				"[%d/%d] working on: %s",
				atomic.AddInt64(&taskID, 1),
				len(tasksOptions),
				taskOptions.URL,
			)
			return gitdump.MirrorGit(ctx, taskOptions)
		}))
	}

	err = conc.ParallelFailSafeWithLimit(opts.Concurrency, tasks...).Do(ctx)
}
