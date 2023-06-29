package gitdump

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/exec"
)

type MirrorGitOptions struct {
	Dir      string
	URL      string
	Username string
	Password string
}

func MirrorGit(ctx context.Context, opts MirrorGitOptions) (err error) {
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
		if err = os.RemoveAll(opts.Dir); err != nil {
			return
		}
		stat = nil
	}

	if stat == nil {
		log.Println("creating directory:", opts.Dir)
		if err = os.MkdirAll(opts.Dir, 0755); err != nil {
			return
		}

		log.Println("initializing git repository:", opts.Dir)
		{
			cmd := exec.CommandContext(ctx, "git", "init", "--bare")
			cmd.Dir = opts.Dir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				return
			}
		}
	}

	// embed credentials in URL
	var upstream *url.URL
	if upstream, err = url.Parse(opts.URL); err != nil {
		return
	}
	upstream.User = url.UserPassword(opts.Username, opts.Password)

	// fetch all refs
	{
		cmd := exec.CommandContext(ctx, "git", "fetch", "-t", "-f", upstream.String(), "+refs/*:refs/*")
		cmd.Dir = opts.Dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return
		}
	}

	log.Println("done:", opts.URL)

	return
}
