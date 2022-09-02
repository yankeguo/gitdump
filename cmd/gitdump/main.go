package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/guoyk93/gitdump"
	"github.com/guoyk93/grace"
	"github.com/guoyk93/grace/graceconf"
	"github.com/guoyk93/grace/gracemain"
	"github.com/guoyk93/grace/gracenotify"
	"github.com/guoyk93/grace/gracesync"
	"github.com/guoyk93/grace/gracetrack"
	"log"
	"path/filepath"
	"strconv"

	_ "github.com/guoyk93/gitdump/hostings"
)

func main() {
	var (
		err error

		ctx, _ = gracemain.WithSignalCancel(
			gracetrack.Init(context.Background()),
		)
	)

	defer gracemain.Exit(&err)
	defer gracenotify.Notify("[GITDUMP]", &ctx, &err)
	defer grace.Guard(&err)

	opts := grace.Must(graceconf.LoadYAMLFlagConf[gitdump.Options]())

	_ = gracemain.WriteLastRun(opts.Dir)

	var tasks []gitdump.MirrorGitOptions

	for i, account := range opts.Accounts {
		var name string
		if account.URL == "" {
			name = fmt.Sprintf("%s / %s", account.Vendor, account.Username)
		} else {
			name = fmt.Sprintf("%s (%s) / %s", account.Vendor, account.URL, account.Username)
		}

		tg := gracetrack.Group(ctx, fmt.Sprintf("account-%d", i)).SetName(name)
		log.Println("working on:", name)

		hosting := gitdump.GetHosting(account.Vendor)
		if hosting == nil {
			err = errors.New("hosting not supported: " + account.Vendor)
			tg.Add("INVALID VENDOR")
			return
		}

		var repos []gitdump.HostingRepo
		if repos, err = hosting.List(ctx, gitdump.HostingOptions{
			URL:      account.URL,
			Username: account.Username,
			Password: account.Password,
		}); err != nil {
			log.Println("failed to listing repos:", name, ":", err.Error())
			tg.Add("ERROR: " + err.Error())
			err = nil
			continue
		}

		for _, repo := range repos {
			tasks = append(tasks, gitdump.MirrorGitOptions{
				Dir:      filepath.Join(opts.Dir, repo.SubDir),
				URL:      repo.URL,
				Username: repo.Username,
				Password: repo.Password,
			})
		}
	}

	tg := gracetrack.Group(ctx, "result").SetName("RESULT")

	if errs := gracesync.DoPara(ctx, tasks, opts.Concurrency, gitdump.MirrorGit); errs != nil {
		errs := errs.(grace.Errors)
		for i, err := range errs {
			if err != nil {
				tg.Add(tasks[i].URL + ": " + err.Error())
			}
		}
	}

	tg.Add("TOTAL: " + strconv.Itoa(len(tasks)))
}
