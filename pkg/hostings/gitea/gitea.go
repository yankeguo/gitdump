package gitea

import (
	"context"
	"errors"
	"github.com/guoyk93/gitdump"
)

func init() {
	gitdump.SetHosting("gitea", Hosting{})
}

type Hosting struct{}

func (h Hosting) List(ctx context.Context, opts gitdump.HostingOptions) (repos []gitdump.HostingRepo, err error) {
	err = errors.New("not implemented")
	return
}
