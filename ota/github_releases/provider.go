package github_releases

import (
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/ota"
)

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
			repo:  "homie-intercom",
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
		log.Debug("Found release: latest is " + releases.GetTagName())

	}
	return &firmware{id: c.Id(), version: releases.GetTagName(), repo: repoInfo}
}
