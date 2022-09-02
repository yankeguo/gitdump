package gitdump

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestMirrorGit(t *testing.T) {
	dir := filepath.Join("testdata", "repo")
	defer os.RemoveAll(dir)
	err := MirrorGit(context.Background(), MirrorGitOptions{
		URL: "https://github.com/guoyk93/gitdump",
		Dir: dir,
	})
	require.NoError(t, err)
	repo, err := git.PlainOpen(dir)
	require.NoError(t, err)
	cmt, err := repo.Log(&git.LogOptions{
		From: plumbing.NewHash("9d3aa548cc6b649f633b25283582342620be6e17"),
	})
	require.NoError(t, err)
	cm, err := cmt.Next()
	require.NoError(t, err)
	require.Contains(t, cm.Message, "Initial")
}
