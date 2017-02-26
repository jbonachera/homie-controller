package github_releases

import (
	"context"
	"github.com/google/go-github/github"
)

type ghClient interface {
	GetLatestRelease(owner string, repo string) (*github.RepositoryRelease, error)
}

type defaultGhClient struct {
	client *github.Client
	ctx    context.Context
}

func GetDefaultGHClient() *defaultGhClient {
	c := &defaultGhClient{}
	c.client = github.NewClient(nil)
	c.ctx = context.Background()
	return c
}

func (c *defaultGhClient) GetLatestRelease(owner string, repo string) (*github.RepositoryRelease, error) {
	release, _, err := c.client.Repositories.GetLatestRelease(c.ctx, owner, repo)
	return release, err
}
