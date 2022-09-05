package gitdump

import (
	"context"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"os"
	"strings"
)

const (
	Upstream = "upstream"
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
	}

	var repo *git.Repository

	if repo, err = git.PlainOpen(opts.Dir); err != nil {
		if err == git.ErrRepositoryNotExists {
			err = nil
		} else {
			return
		}
	}

	if repo == nil {
		if repo, err = git.PlainInit(opts.Dir, true); err != nil {
			return
		}
	}

	if err = repo.DeleteRemote(Upstream); err != nil {
		if err == git.ErrRemoteNotFound {
			err = nil
		} else {
			return
		}
	}

	if _, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name:  Upstream,
		URLs:  []string{opts.URL},
		Fetch: []gitconfig.RefSpec{"+refs/*:refs/*"},
	}); err != nil {
		return
	}

	defer func() {
		_ = repo.DeleteRemote(Upstream)
	}()

	var auth transport.AuthMethod

	if opts.Username != "" || opts.Password != "" {
		auth = &http.BasicAuth{
			Username: opts.Username,
			Password: opts.Password,
		}
	}

	if err = repo.FetchContext(ctx, &git.FetchOptions{
		RemoteName: Upstream,
		Auth:       auth,
		Force:      true,
		Tags:       git.AllTags,
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			err = nil
		} else if strings.Contains(err.Error(), "remote repository is empty") {
			err = nil
		} else {
			return
		}
	}

	log.Println("done:", opts.URL)

	return
}
