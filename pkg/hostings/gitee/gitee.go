package gitee

import (
	"context"
	"errors"
	"github.com/guoyk93/gitdump"
)

func init() {
	gitdump.SetHosting("gitee", Gitee{})
}

type Gitee struct{}

func (h Gitee) List(ctx context.Context, opts gitdump.HostingOptions) (repos []gitdump.HostingRepo, err error) {
	err = errors.New("not implemented")
	return
}
