package github_releases

import (
	"github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/model/homieMessage"
)

type firmware struct {
	id       string
	version  string
	checksum string
	repo     repoInfo
	payload  []byte
}

func (f *firmware) Name() string {
	return f.id
}

func (f *firmware) Version() string {
	return f.version
}

func (f *firmware) Checksum() string {
	return f.checksum
}

func (f *firmware) Payload() []byte {
	return f.payload
}
func (f *firmware) Announce() {
	nodePath := "devices/controller/" + f.Name()
	implemPath := "devices/controller/$implementation/firmware/" + f.Name()
	messaging.PublishState(homieMessage.HomieMessage{Topic: nodePath + "/" + f.Version(), Payload: f.Checksum()})
	messaging.AddHandler(nodePath+"/"+f.Version()+"/get", func(message homieMessage.HomieMessage) {
		if message.Payload == "true" {
			messaging.PublishFile(nodePath+"/"+f.Version()+"/checksum", f.Payload())
		}
	})
	messaging.AddHandler(implemPath+"/"+f.Version()+"/upgrade", func(message homieMessage.HomieMessage) {
		device := "devices/" + message.Payload + "/$implementation/ota/firmware"
		messaging.PublishFile(device, f.Payload())
	})
}
