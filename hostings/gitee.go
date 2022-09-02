package hostings

import (
	"context"
	"errors"
	"github.com/guoyk93/gitdump"
)

func init() {
	gitdump.SetHosting("gitee", hostingGitee{})
}

type hostingGitee struct{}

func (h hostingGitee) List(ctx context.Context, opts gitdump.HostingOptions) (repos []gitdump.HostingRepo, err error) {
	err = errors.New("not implemented")
	return
}
