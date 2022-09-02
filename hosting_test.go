package gitdump

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type hostingDummy struct{}

func (h *hostingDummy) List(ctx context.Context, opts HostingOptions) (repos []HostingRepo, err error) {
	return
}

func TestGetHosting(t *testing.T) {
	d := &hostingDummy{}
	SetHosting("dummy", d)
	require.Equal(t, d, GetHosting("dummy"))
}
