package github_releases

import (
	"bytes"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/ota"
	"io"
	"strings"
)

type repoInfo struct {
	owner string
	repo  string
}

type GhOTAProvider struct {
	id              string
	releaseProvider ghClient
	version         map[string]*firmware
}

func (c *GhOTAProvider) Id() string {
	return c.id
}

func getGithubInfo(firmware string) repoInfo {
	infos := map[string]repoInfo{
		"intercom": {
			owner: "jbonachera",
			repo:  "homie-intercom",
		},
		"vx-temperature-sensor": {
			owner: "jbonachera",
			repo:  "homie-sensor",
		},
		"temperature-sensor": {
			owner: "jbonachera",
			repo:  "homie-sensor",
		},
		"mock": {
			owner: "mock",
			repo:  "mock",
		},
	}
	return infos[firmware]
}

func (c *GhOTAProvider) GetVersion(version string) ota.Firmware {
	repoInfo := getGithubInfo(c.Id())
	if wantedFirmware, ok := c.version[version]; ok {
		return wantedFirmware
	} else {
		releases, err := c.releaseProvider.GetReleaseByTag(repoInfo.owner, repoInfo.repo, version)
		if err != nil {
			log.Error("could not fetch releases from github for firmware " + c.Id() + ": " + err.Error())
			return &firmware{id: c.Id(), version: "unknown", repo: repoInfo, checksum: "", payload: []byte{}}
		}
		var checksum string
		var payload []byte
		for _, asset := range releases.Assets {
			if asset.GetName() == repoInfo.repo+".md5" {
				log.Debug("will fetch asset " + asset.GetName())
				reader, err := c.releaseProvider.DownloadReleaseAsset(repoInfo.owner, repoInfo.repo, asset.GetID())
				if err != nil {
					log.Error(err.Error())
				} else if reader != nil {
					checksum = fetchFromReader(reader, asset.GetSize()).String()
					checksum = strings.Split(checksum, " ")[0]
				}
			} else if asset.GetName() == repoInfo.repo {
				log.Debug("will fetch asset " + asset.GetName())
				reader, err := c.releaseProvider.DownloadReleaseAsset(repoInfo.owner, repoInfo.repo, asset.GetID())
				if err != nil {
					log.Error(err.Error())
				} else if reader != nil {
					payload = fetchFromReader(reader, asset.GetSize()).Bytes()
				}
			}
		}
		c.version[releases.GetTagName()] = &firmware{id: c.Id(), version: releases.GetTagName(), repo: repoInfo, checksum: checksum, payload: payload}
		c.version[releases.GetTagName()].Announce()
		versionString := "unit,provider"
		for version := range c.version {
			versionString = versionString + "," + version
		}
		messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/" + c.Id() + "/$properties", Payload: versionString})
		return c.version[releases.GetTagName()]
	}
}

func (c *GhOTAProvider) GetLatest() ota.Firmware {
	repoInfo := getGithubInfo(c.Id())
	releases, err := c.releaseProvider.GetLatestRelease(repoInfo.owner, repoInfo.repo)
	if err != nil {
		log.Error("could not fetch releases from github for firmware " + c.Id() + ": " + err.Error())
		return &firmware{id: c.Id(), version: "unknown", repo: repoInfo, checksum: "", payload: []byte{}}
	} else {
		log.Debug("Found release: latest is " + releases.GetTagName())

	}
	fw := c.GetVersion(releases.GetTagName())
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/" + fw.Name() + "/$latest", Payload: fw.Version()})
	return fw

}

func fetchFromReader(reader io.ReadCloser, size int) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.Grow(size)
	buf.ReadFrom(reader)
	reader.Close()
	return buf
}
