package github_releases

import (
	"github.com/google/go-github/github"
	"io"
)

type MockGHClient struct{}

func (m *MockGHClient) GetReleaseByTag(owner string, repo string, tag string) (*github.RepositoryRelease, error) {
	version := tag
	id := 1
	release := &github.RepositoryRelease{TagName: &version, ID: &id}
	return release, nil
}

func (m *MockGHClient) GetLatestRelease(owner string, repo string) (*github.RepositoryRelease, error) {
	version := "1.0.1"
	id := 1
	release := &github.RepositoryRelease{TagName: &version, ID: &id}
	return release, nil
}
func (m *MockGHClient) DownloadReleaseAsset(owner, repo string, id int) (rc io.ReadCloser, err error) {
	return nil, err
}

/*
func TestGhOTAProvider_GetLatest(t *testing.T) {
	ota.AddFirmware("mock", "github_release")
	provider := &GhOTAProvider{id: "mock", releaseProvider: &MockGHClient{}}
	firmware := provider.GetLatest()
	if firmware.Name() != "mock" {
		t.Error("could not get the firmware name: got " + firmware.Name())
	}
	if firmware.Version() != "1.0.1" {
		t.Error("could not get the latest firmware version: got " + firmware.Version())
	}
}
*/
