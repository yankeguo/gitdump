package gitdump

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/exec"

	"github.com/yankeguo/rg"
)

type MirrorGitOptions struct {
	Dir         string
	URL         string
	Username    string
	AskPassFile string
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
	if upstream.User == nil {
		upstream.User = url.User(opts.Username)
	}

	// fetch all refs
	{
		cmd := exec.CommandContext(
			ctx,
			"git",
			"fetch",
			"-t",
			"-f",
			upstream.String(),
			"+refs/*:refs/*",
		)
		cmd.Env = append(os.Environ(), "GIT_ASKPASS="+opts.AskPassFile)
		cmd.Dir = opts.Dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		rg.Must0(cmd.Run())
	}

	log.Println("done:", opts.URL)

	return
}
