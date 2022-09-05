package github

import (
	"context"
	"github.com/google/go-github/v47/github"
	"github.com/guoyk93/gitdump"
	"github.com/guoyk93/grace"
	"golang.org/x/oauth2"
	"path/filepath"
	"strings"
)

func init() {
	gitdump.SetHosting("github", Hosting{})
}

type Hosting struct{}

func (h Hosting) List(ctx context.Context, opts gitdump.HostingOptions) (out []gitdump.HostingRepo, err error) {
	defer grace.Guard(&err)

	opts.MustUsername()
	opts.MustPassword()

	const (
		hostname = "github.com"
	)

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: opts.Password},
	)))

	var orgNames []string
	{
		var i = 0
		for {
			i++

			data, _ := grace.Must2(client.Organizations.List(ctx, "", &github.ListOptions{
				Page:    i,
				PerPage: 50,
			}))

			if len(data) == 0 {
				break
			}

			for _, item := range data {
				orgNames = append(orgNames, item.GetLogin())
			}
		}
	}

	repos := map[string]*github.Repository{}

	{
		i := 0
		for {
			i++

			data, _ := grace.Must2(client.Repositories.List(ctx, "", &github.RepositoryListOptions{
				ListOptions: github.ListOptions{
					Page:    i,
					PerPage: 50,
				},
			}))

			if len(data) == 0 {
				break
			}

			for _, item := range data {
				repos[item.GetFullName()] = item
			}
		}
	}

	for _, orgName := range orgNames {

		{
			i := 0
			for {
				i++

				data, _ := grace.Must2(client.Repositories.ListByOrg(ctx, orgName, &github.RepositoryListByOrgOptions{
					ListOptions: github.ListOptions{
						Page:    i,
						PerPage: 50,
					},
				}))

				if len(data) == 0 {
					break
				}

				for _, item := range data {
					repos[item.GetFullName()] = item
				}
			}
		}

	}

	for _, repo := range repos {

		items := []string{hostname}
		items = append(items, strings.Split(repo.GetFullName(), "/")...)

		out = append(out, gitdump.HostingRepo{
			SubDir:   filepath.Join(items...),
			URL:      repo.GetCloneURL(),
			Username: opts.Username,
			Password: opts.Password,
		})
	}
	return
}
