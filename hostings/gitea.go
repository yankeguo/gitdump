package hostings

import (
	"context"
	"errors"
	"github.com/guoyk93/gitdump"
)

func init() {
	gitdump.SetHosting("gitea", hostingGitea{})
}

type hostingGitea struct{}

func (h hostingGitea) List(ctx context.Context, opts gitdump.HostingOptions) (repos []gitdump.HostingRepo, err error) {
	err = errors.New("not implemented")
	return
}
