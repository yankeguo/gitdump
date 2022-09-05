package coding

import (
	"context"
	"github.com/guoyk93/gitdump"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestHosting_List(t *testing.T) {
	h := Hosting{}
	repos, err := h.List(context.Background(), gitdump.HostingOptions{
		URL:      os.Getenv("TEST_CODING_URL"),
		Username: os.Getenv("TEST_CODING_USERNAME"),
		Password: os.Getenv("TEST_CODING_PASSWORD"),
	})
	require.NoError(t, err)
	require.NotEmpty(t, repos)
}
