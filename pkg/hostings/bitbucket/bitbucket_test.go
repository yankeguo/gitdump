package bitbucket

import (
	"context"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/gitdump"
	"github.com/yankeguo/rg"
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
