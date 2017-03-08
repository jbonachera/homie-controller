package github_releases

import (
	"github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/model/homieMessage"
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
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/" + name + "/provider", Payload: f.Id()})
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/" + name + "/$type", Payload: "firmwareProvider"})
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/" + name + "/unit", Payload: "md5"})

	return &GhOTAProvider{id: name, releaseProvider: GetDefaultGHClient(), version: map[string]*firmware{}}
}
