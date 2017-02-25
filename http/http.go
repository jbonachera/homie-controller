package http

import (
	"encoding/json"
	"github.com/husobee/vestigo"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/device"
	"net/http"
)

var router = vestigo.NewRouter()

type statusMessage struct {
	Error bool `json:"error"`
	Msg string `json:"msg"`
}

func Start() {
	router.Get("/devices/", ListDevicesHandler)
	router.Get("/device/:name", GetDeviceHandler)
	router.Post("/device/:name/implementation/:action", PostImplementationActionHandler)
	log.Info("starting HTTP API")
	log.Error(http.ListenAndServe(":8989", router).Error())
}

func GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	name := vestigo.Param(r, "name")
	myDevice, err := device.Get(name)
	if err != nil {
		w.WriteHeader(404)
	} else {
		json.NewEncoder(w).Encode(myDevice)
	}
}
func ListDevicesHandler(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(device.List())
	w.WriteHeader(200)
}
func PostImplementationActionHandler(w http.ResponseWriter, r *http.Request) {
	name := vestigo.Param(r, "name")
	myDevice, err := device.Get(name)
	if err != nil {
		w.WriteHeader(404)
	} else {
		action := vestigo.Param(r, "action")
		log.Info("will send "+action + " command to device "+name)
		err := myDevice.Implementation.Do(action)
		if err != nil {
			log.Info("command failed: "+err.Error())
			json.NewEncoder(w).Encode(statusMessage{Error: true, Msg: "command said: "+err.Error()})
		} else {
			log.Info("command sent")
			json.NewEncoder(w).Encode(statusMessage{Error: false, Msg: "command sent"})
		}
	}
}
