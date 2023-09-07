package coding

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/gitdump"
	"github.com/guoyk93/rg"
	"log"
	"net/url"
	"path/filepath"
)

type M = map[string]any

func init() {
	gitdump.SetHosting("coding", Hosting{})
}

type ResultProject struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

type ResultRepo struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	HttpsUrl string `json:"HttpsUrl"`
}

type Result struct {
	Response struct {
		Error *struct {
			Message string `json:"Message"`
		} `json:"Error"`
		User struct {
			ID int `json:"Id"`
		} `json:"User"`
		ProjectList []ResultProject `json:"ProjectList"`
		DepotData   struct {
			Depots []ResultRepo `json:"Depots"`
		} `json:"DepotData"`
	} ` json:"Response"`
}

type Client struct {
	client *resty.Client
}

func (c *Client) Invoke(ctx context.Context, action string, body M) (result Result, err error) {
	if body == nil {
		body = M{}
	}
	body["Action"] = action
	var res *resty.Response
	if res, err = c.client.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(body).
		Post("open-api"); err != nil {
		return
	}
	if res.IsError() {
		err = errors.New(res.String())
		return
	}
	if result.Response.Error != nil {
		err = errors.New(result.Response.Error.Message)
		return
	}
	return
}

func (c *Client) GetUserID(ctx context.Context) (userID int, err error) {
	var result Result
	if result, err = c.Invoke(ctx, "DescribeCodingCurrentUser", nil); err != nil {
		err = fmt.Errorf("failed to get user id: %w", err)
		return
	}
	userID = result.Response.User.ID
	return
}

func (c *Client) GetUserProjectIDs(ctx context.Context, userID int) (projects []ResultProject, err error) {
	var result Result
	if result, err = c.Invoke(ctx, "DescribeUserProjects", M{"UserId": userID}); err != nil {
		err = fmt.Errorf("failed to get user %d projects: %w", userID, err)
		return
	}
	projects = result.Response.ProjectList
	return
}

func (c *Client) GetProjectRepos(ctx context.Context, projectID int) (repositories []ResultRepo, err error) {
	var result Result
	if result, err = c.Invoke(ctx, "DescribeProjectDepotInfoList", M{"ProjectId": projectID}); err != nil {
		err = fmt.Errorf("failed to get project %d repos: %w", projectID, err)
		return
	}
	repositories = result.Response.DepotData.Depots
	return
}

type Hosting struct{}

func (h Hosting) List(ctx context.Context, opts gitdump.HostingOptions) (out []gitdump.HostingRepo, err error) {
	defer rg.Guard(&err)

	opts.MustURL()
	opts.MustUsername()
	opts.MustPassword()

	u := rg.Must(url.Parse(opts.URL))

	var hostname string

	if hostname = u.Hostname(); hostname == "" {
		panic(errors.New("missing hostname (opts.URL)"))
	}

	client := &Client{
		client: resty.New().
			SetBaseURL(opts.URL).
			SetBasicAuth(opts.Username, opts.Password),
	}

	userId := rg.Must(client.GetUserID(ctx))

	projects := rg.Must(client.GetUserProjectIDs(ctx, userId))

	for _, project := range projects {
		repos, err1 := client.GetProjectRepos(ctx, project.ID)

		if err1 != nil {
			log.Printf("failed to get project %s repos: %s", project.Name, err1.Error())
			continue
		}

		for _, repo := range repos {
			out = append(out, gitdump.HostingRepo{
				SubDir:   filepath.Join(hostname, project.Name, repo.Name),
				URL:      repo.HttpsUrl,
				Username: opts.Username,
				Password: opts.Password,
			})
		}
	}

	return
}
