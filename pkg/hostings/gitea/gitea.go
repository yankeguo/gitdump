package gitea

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/gitdump"
	"github.com/guoyk93/grace"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func init() {
	gitdump.SetHosting("gitea", Hosting{})
}

type Org struct {
	Username string `json:"username"`
}

type Repo struct {
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
}

type Client struct {
	client *resty.Client
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
		return r.Get("orgs")
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
					SetQueryParam("limit", "50"),
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
	defer grace.Guard(&err)

	opts.MustURL()
	opts.MustUsername()
	opts.MustPassword()

	u := grace.Must(url.Parse(opts.URL))
	var hostname string
	if hostname = u.Hostname(); hostname == "" {
		err = errors.New("missing hostname (opts.URL)")
		return
	}

	client := &Client{
		client: resty.New().
			SetBaseURL(strings.TrimSuffix(opts.URL, "/") + "/api/v1").
			SetAuthToken(opts.Password),
	}

	repos := map[string]Repo{}

	{
		userRepos := grace.Must(client.GetUserRepos(ctx))

		for _, repo := range userRepos {
			repos[repo.FullName] = repo
		}
	}

	{
		orgs := grace.Must(client.GetOrgs(ctx))

		for _, org := range orgs {

			orgRepos := grace.Must(client.GetOrgRepos(ctx, org.Username))

			for _, repo := range orgRepos {
				repos[repo.FullName] = repo
			}
		}
	}

	names := grace.MapKeys(repos)
	sort.Strings(names)

	for _, name := range names {
		repo := repos[name]

		items := []string{hostname}
		items = append(items, strings.Split(repo.FullName, "/")...)

		out = append(out, gitdump.HostingRepo{
			SubDir:   filepath.Join(items...),
			URL:      repo.CloneURL,
			Username: opts.Username,
			Password: opts.Password,
		})
	}

	return
}
