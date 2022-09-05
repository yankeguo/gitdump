package github

import (
	"context"
	"github.com/guoyk93/gitdump"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHosting_List(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	h := Hosting{}
	opts := gitdump.HostingOptions{
		Username: os.Getenv("TEST_GITHUB_USERNAME"),
		Password: os.Getenv("TEST_GITHUB_PASSWORD"),
	}
	repos, err := h.List(context.Background(), opts)
	require.NoError(t, err)
	require.NotEmpty(t, repos)
	repo := repos[rand.Intn(len(repos))]
	require.True(t, strings.HasPrefix(repo.URL, "https://github.com"))
	require.True(t, strings.HasPrefix(repo.SubDir, "github.com"))
	require.True(t, repo.Username == opts.Username)
	require.True(t, repo.Password == opts.Password)
}
