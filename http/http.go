package http

import (
	"github.com/husobee/vestigo"
	"github.com/jbonachera/homie-controller/log"
	"net/http"
	"github.com/jbonachera/homie-controller/model/device"
	"encoding/json"
)

var router = vestigo.NewRouter()

func Start() {
	router.Get("/devices/", ListDevicesHandler)
	router.Get("/device/:name", GetDeviceHandler)
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
