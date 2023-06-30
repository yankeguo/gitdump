package gitdump

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/guoyk93/rg"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

type MirrorGitOptions struct {
	Dir      string
	URL      string
	Username string
	Password string
}

func MirrorGit(ctx context.Context, opts MirrorGitOptions) (err error) {
	defer rg.Guard(&err)

	var stat os.FileInfo

	if stat, err = os.Stat(opts.Dir); err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			log.Println("failed checking directory:", opts.Dir)
			return
		}
	} else if !stat.IsDir() {
		log.Println("not a directory, deleting:", opts.Dir)
		rg.Must0(os.RemoveAll(opts.Dir))
		stat = nil
	}

	if stat == nil {
		log.Println("creating directory:", opts.Dir)
		rg.Must0(os.MkdirAll(opts.Dir, 0755))

		log.Println("initializing git repository:", opts.Dir)
		{
			cmd := exec.CommandContext(ctx, "git", "init", "--bare")
			cmd.Dir = opts.Dir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			rg.Must0(cmd.Run())
		}
	}

	// embed username in URL
	upstream := rg.Must(url.Parse(opts.URL))
	upstream.User = url.User(opts.Username)

	// create askPass script
	urlDigest := md5.Sum([]byte(opts.URL))
	credentialScript := filepath.Join(os.TempDir(), "gitdumpcredentials-"+hex.EncodeToString(urlDigest[:])+".sh")
	defer os.RemoveAll(credentialScript)

	rg.Must0(
		os.WriteFile(
			credentialScript,
			[]byte(
				fmt.Sprintf(
					"#!/bin/sh\necho '%s' | base64 -d",
					base64.StdEncoding.EncodeToString([]byte(opts.Password)),
				),
			),
			0755,
		),
	)

	// fetch all refs
	{
		cmd := exec.CommandContext(
			ctx,
			"git",
			"fetch",
			"-v",
			"-t",
			"-f",
			upstream.String(),
			"+refs/*:refs/*",
		)
		cmd.Env = append(os.Environ(), "GIT_ASKPASS="+credentialScript)
		cmd.Dir = opts.Dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		rg.Must0(cmd.Run())
	}

	log.Println("done:", opts.URL)

	return
}
