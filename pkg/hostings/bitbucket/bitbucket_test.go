package bitbucket

import (
	"context"
	"github.com/guoyk93/gitdump"
	"github.com/guoyk93/rg"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHosting_List(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	h := Hosting{}
	opts := gitdump.HostingOptions{
		Username: os.Getenv("TEST_BITBUCKET_USERNAME"),
		Password: os.Getenv("TEST_BITBUCKET_PASSWORD"),
	}
	u := rg.Must(url.Parse(opts.URL))
	repos, err := h.List(context.Background(), opts)
	require.NoError(t, err)
	require.NotEmpty(t, repos)
	repo := repos[rand.Intn(len(repos))]
	require.True(t, strings.HasPrefix(repo.URL, "https://bitbucket.org"))
	require.True(t, strings.HasPrefix(repo.SubDir, u.Hostname()))
	require.True(t, repo.Username == opts.Username)
	require.True(t, repo.Password == opts.Password)
}
