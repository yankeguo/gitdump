package bitbucket

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/gitdump"
	"github.com/guoyk93/rg"
	"log"
	"path/filepath"
	"strings"
)

func init() {
	gitdump.SetHosting("bitbucket", Hosting{})
}

type WorkspaceList struct {
	Next   string `json:"next"`
	Values []struct {
		Workspace struct {
			Slug string `json:"slug"`
		} `json:"workspace"`
	} `json:"values"`
}

type RepositoryList struct {
	Next   string `json:"next"`
	Values []struct {
		Links struct {
			Clone []struct {
				Href string `json:"href"`
				Name string `json:"name"`
			} `json:"clone"`
		} `json:"links"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"values"`
}

type Hosting struct{}

func (h Hosting) List(ctx context.Context, opts gitdump.HostingOptions) (out []gitdump.HostingRepo, err error) {
	defer rg.Guard(&err)

	opts.MustUsername()
	opts.MustPassword()

	client := resty.New().SetBaseURL("https://api.bitbucket.org/2.0").SetBasicAuth(opts.Username, opts.Password)

	workspaceSlugs := map[string]struct{}{}
	{

		link := "https://api.bitbucket.org/2.0/user/permissions/workspaces"

		for {
			var data WorkspaceList

			res := rg.Must(client.R().SetContext(ctx).SetResult(&data).Get(link))

			if res.IsError() {
				log.Println(res.String())
				break
			}

			for _, item := range data.Values {
				workspaceSlugs[item.Workspace.Slug] = struct{}{}
			}

			if data.Next == "" {
				break
			}

			link = data.Next
		}
	}

	for workspace := range workspaceSlugs {

		link := "https://api.bitbucket.org/2.0/repositories/" + workspace

		for {
			var data RepositoryList
			res := rg.Must(client.R().SetContext(ctx).SetResult(&data).Get(link))

			if res.IsError() {
				log.Println(res.String())
				break
			}

			for _, item := range data.Values {
				var cloneHref = ""
				for _, cloneItem := range item.Links.Clone {
					if cloneItem.Name == "https" {
						cloneHref = cloneItem.Href
					}
				}
				if cloneHref == "" {
					continue
				}

				dirItems := []string{"bitbucket.org"}
				dirItems = append(dirItems, strings.Split(item.FullName, "/")...)

				out = append(out, gitdump.HostingRepo{
					SubDir:   filepath.Join(dirItems...),
					URL:      cloneHref,
					Username: opts.Username,
					Password: opts.Password,
				})
			}

			if data.Next == "" {
				break
			}

			link = data.Next
		}

	}

	return
}
