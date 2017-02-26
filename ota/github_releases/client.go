package github_releases

import (
	"context"
	"errors"
	"github.com/google/go-github/github"
	"github.com/jbonachera/homie-controller/config"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

type ghClient interface {
	GetLatestRelease(owner string, repo string) (*github.RepositoryRelease, error)
	GetReleaseByTag(owner string, repo string, tag string) (*github.RepositoryRelease, error)
	DownloadReleaseAsset(owner, repo string, id int) (rc io.ReadCloser, err error)
}

type defaultGhClient struct {
	client *github.Client
	ctx    context.Context
}

func GetDefaultGHClient() *defaultGhClient {
	c := &defaultGhClient{}
	c.ctx = context.Background()

	if token := config.Get("GITHUB_OAUTH_TOKEN"); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(c.ctx, ts)
		c.client = github.NewClient(tc)

	} else {
		c.client = github.NewClient(nil)
	}
	return c
}

func (c *defaultGhClient) GetReleaseByTag(owner string, repo string, tag string) (*github.RepositoryRelease, error) {
	release, _, err := c.client.Repositories.GetReleaseByTag(c.ctx, owner, repo, tag)
	return release, err
}

func (c *defaultGhClient) GetLatestRelease(owner string, repo string) (*github.RepositoryRelease, error) {
	release, _, err := c.client.Repositories.GetLatestRelease(c.ctx, owner, repo)
	return release, err
}
func (c *defaultGhClient) DownloadReleaseAsset(owner, repo string, id int) (rc io.ReadCloser, err error) {
	release, redirect, err := c.client.Repositories.DownloadReleaseAsset(c.ctx, owner, repo, id)
	if err != nil {
		return nil, err
	}
	if redirect != "" {
		response, err := http.Get(redirect)
		if err != nil {
			return nil, errors.New("could not GET " + redirect)
		} else {
			return response.Body, nil
		}
	} else {
		return release, nil
	}
}
