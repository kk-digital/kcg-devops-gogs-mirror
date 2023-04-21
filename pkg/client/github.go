package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/v51/github"
	"github.com/gregjones/httpcache"
	"golang.org/x/oauth2"
)

type githubClient struct {
	accessToken string
	client      *github.Client
}

func NewGithubClient(ctx context.Context, accessToken string) *githubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := &http.Client{
		Transport: &oauth2.Transport{
			Base:   httpcache.NewMemoryCacheTransport(),
			Source: ts,
		},
	}

	return &githubClient{
		accessToken: accessToken,
		client:      github.NewClient(tc),
	}
}

func (c *githubClient) ListOrgRepos(ctx context.Context, orgName string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := c.client.Repositories.ListByOrg(ctx, orgName, opt)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				log.Println("hit rate limit")
			}
			data, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("data:", string(data))
			fmt.Println("resp:", resp)
			return nil, fmt.Errorf("failed to list Github repositories: %w", err)
		}
		defer resp.Body.Close()

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

func (c *githubClient) Clone() {}
