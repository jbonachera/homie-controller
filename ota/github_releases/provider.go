package github_releases

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/ota"
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

type repoInfo struct {
	owner string
	repo  string
}

type GhOTAProvider struct {
	id              string
	releaseProvider ghClient
}

func (c *GhOTAProvider) Id() string {
	return c.id
}

func getGithubInfo(firmware string) repoInfo {
	infos := map[string]repoInfo{
		"vx-intercom-sensor": {
			owner: "jbonachera",
			repo:  "homie-intercom-sensor",
		},
	}
	return infos[firmware]
}

func (c *GhOTAProvider) GetLatest() ota.Firmware {
	repoInfo := getGithubInfo(c.Id())
	releases, err := c.releaseProvider.GetLatestRelease(repoInfo.owner, repoInfo.repo)
	if err != nil {
		log.Error("could not fetch releases from github: " + err.Error())
	} else {
		log.Info("Found release")

	}
	return &firmware{id: c.Id(), version: releases.GetName(), repo: repoInfo}
}
