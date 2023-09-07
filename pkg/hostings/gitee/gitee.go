package gitee

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/gitdump"
	"github.com/guoyk93/rg"
	"golang.org/x/exp/maps"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func init() {
	gitdump.SetHosting("gitee", Hosting{})
}

type Org struct {
	Login string `json:"login"`
}

type Repo struct {
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}

type Client struct {
	username string
	client   *resty.Client
}

func (c *Client) Invoke(ctx context.Context, fn func(r *resty.Request) (*resty.Response, error)) (err error) {
	var res *resty.Response
	if res, err = fn(c.client.R().SetContext(ctx)); err != nil {
		return
	}
	if res.IsError() {
		return errors.New(res.String())
	}
	return
}

func (c *Client) GetOrgs(ctx context.Context) (orgs []Org, err error) {
	s := ListService[Org]{client: c}
	return s.Invoke(ctx, func(r *resty.Request) (*resty.Response, error) {
		return r.SetPathParam("username", c.username).Get("users/{username}/orgs")
	})
}

func (c *Client) GetUserRepos(ctx context.Context) (repos []Repo, err error) {
	s := ListService[Repo]{client: c}
	return s.Invoke(ctx, func(r *resty.Request) (*resty.Response, error) {
		return r.Get("user/repos")
	})
}

func (c *Client) GetOrgRepos(ctx context.Context, orgName string) (repos []Repo, err error) {
	s := ListService[Repo]{client: c}
	return s.Invoke(ctx, func(r *resty.Request) (*resty.Response, error) {
		return r.SetPathParam("orgName", orgName).Get("orgs/{orgName}/repos")
	})
}

type ListService[T any] struct {
	client *Client
}

func (c *ListService[T]) Invoke(ctx context.Context, fn func(r *resty.Request) (*resty.Response, error)) (out []T, err error) {
	var i int

	for {
		i++

		var data []T

		if err = c.client.Invoke(ctx, func(r *resty.Request) (*resty.Response, error) {
			return fn(
				r.SetResult(&data).
					SetQueryParam("page", strconv.Itoa(i)).
					SetQueryParam("per_page", "50"),
			)
		}); err != nil {
			return
		}

		if len(data) == 0 {
			break
		}

		for _, item := range data {
			out = append(out, item)
		}
	}

	return
}

type Hosting struct{}

func (h Hosting) List(ctx context.Context, opts gitdump.HostingOptions) (out []gitdump.HostingRepo, err error) {
	defer rg.Guard(&err)

	opts.MustUsername()
	opts.MustPassword()

	const (
		hostname = "gitee.com"
	)

	client := &Client{
		username: opts.Username,
		client: resty.New().
			SetBaseURL("https://gitee.com/api/v5").
			SetAuthToken(opts.Password),
	}

	repos := map[string]Repo{}

	{
		userRepos := rg.Must(client.GetUserRepos(ctx))

		for _, repo := range userRepos {
			repos[repo.FullName] = repo
		}
	}

	{
		orgs := rg.Must(client.GetOrgs(ctx))

		for _, org := range orgs {

			orgRepos := rg.Must(client.GetOrgRepos(ctx, org.Login))

			for _, repo := range orgRepos {
				repos[repo.FullName] = repo
			}
		}
	}

	names := maps.Keys(repos)
	sort.Strings(names)

	for _, name := range names {
		repo := repos[name]

		items := []string{hostname}
		items = append(items, strings.Split(repo.FullName, "/")...)

		out = append(out, gitdump.HostingRepo{
			SubDir:   filepath.Join(items...),
			URL:      repo.HTMLURL,
			Username: opts.Username,
			Password: opts.Password,
		})
	}

	return
}
