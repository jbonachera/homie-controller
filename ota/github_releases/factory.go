package github_releases

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/jbonachera/homie-controller/ota"
)

type Factory struct {
	id string
}

func init() {
	ota.RegisterFactory("github_release", &Factory{id: "ghRelease"})
}

func (f *Factory) Id() string {
	return f.id
}

func (f *Factory) New(name string) ota.FirmwareProvider {
	return &GhOTAProvider{id: name, client: github.NewClient(nil), ctx: context.Background()}
}
