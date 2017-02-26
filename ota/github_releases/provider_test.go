package github_releases

import (
	"github.com/google/go-github/github"
	"testing"
)

type MockGHClient struct{}

func (m *MockGHClient) GetLatestRelease(owner string, repo string) (*github.RepositoryRelease, error) {
	version := "1.0.1"
	id := 1
	release := &github.RepositoryRelease{TagName: &version, ID: &id}
	return release, nil
}

func TestGhOTAProvider_GetLatest(t *testing.T) {
	provider := &GhOTAProvider{id: "mock", releaseProvider: &MockGHClient{}}
	firmware := provider.GetLatest()
	if firmware.Name() != "mock" {
		t.Error("could not get the firmware name: got " + firmware.Name())
	}
	if firmware.Version() != "1.0.1" {
		t.Error("could not get the latest firmware version: got " + firmware.Version())
	}
}
