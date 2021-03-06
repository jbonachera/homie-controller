package http

import (
	"encoding/json"
	"github.com/husobee/vestigo"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/device"
	"github.com/jbonachera/homie-controller/model/search"
	"github.com/jbonachera/homie-controller/ota"
	"github.com/rs/cors"
	"net/http"
	"sort"
)

var router = vestigo.NewRouter()

type statusMessage struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
}

func Start() {
	router.Get("/devices/", ListDevicesHandler)
	router.Post("/devices/search", SearchDevicesHandler)
	router.Get("/device/:name", GetDeviceHandler)
	router.Post("/device/:name/implementation/:action", PostImplementationActionHandler)
	router.Get("/firmware/:firmware", GetFirmwareHandler)
	router.Get("/firmware/:firmware/download", DownloadFirmwareHandler)
	log.Info("starting HTTP API")
	handler := cors.Default().Handler(router)
	log.Error(http.ListenAndServe(":8989", handler).Error())
}

func GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := vestigo.Param(r, "name")
	myDevice, err := device.Get(name)
	if err != nil {
		w.WriteHeader(404)
	} else {
		json.NewEncoder(w).Encode(myDevice)
	}
}

func GetFirmwareHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := vestigo.Param(r, "firmware")
	firmware, err := ota.LastFirmware(name)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(statusMessage{Error: true, Msg: "firmware not found"})
	} else {
		// TODO: implement Mashalable on Firmware interface ..
		json.NewEncoder(w).Encode(struct {
			Id       string `json:"id"`
			Version  string `json:"version"`
			Checksum string `json:"checksum"`
		}{name, firmware.Version(), firmware.Checksum()})
	}
}

func DownloadFirmwareHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := vestigo.Param(r, "firmware")
	firmware, err := ota.LastFirmware(name)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(statusMessage{Error: true, Msg: "firmware not found"})
	} else {
		w.Write(firmware.Payload())
	}
}

func ListDevicesHandler(w http.ResponseWriter, r *http.Request) {
	expand := r.URL.Query().Get("expand") == "true"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if expand {
		devicesMap := device.GetAll()
		var deviceList []*device.Device
		for _, dev := range devicesMap {
			deviceList = append(deviceList, dev)
		}
		sort.Slice(deviceList, func(i, j int) bool { return deviceList[i].Id < deviceList[j].Id })
		json.NewEncoder(w).Encode(deviceList)
	} else {
		json.NewEncoder(w).Encode(device.List())
	}
}
func PostImplementationActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := vestigo.Param(r, "name")
	myDevice, err := device.Get(name)
	if err != nil {
		w.WriteHeader(404)
	} else {
		action := vestigo.Param(r, "action")
		log.Info("will send " + action + " command to device " + name)
		err := myDevice.Implementation.Do(action)
		if err != nil {
			log.Info("command failed: " + err.Error())
			json.NewEncoder(w).Encode(statusMessage{Error: true, Msg: "command said: " + err.Error()})
		} else {
			log.Info("command sent")
			json.NewEncoder(w).Encode(statusMessage{Error: false, Msg: "command sent"})
		}
	}
}

func SearchDevicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	expand := r.URL.Query().Get("expand") == "true"
	decoder := json.NewDecoder(r.Body)
	var opts search.Options
	err := decoder.Decode(&opts)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(statusMessage{Error: true, Msg: err.Error()})
		return
	}

	devices := device.Search(opts)
	if expand {
		json.NewEncoder(w).Encode(devices)
	} else {
		devices_id := []string{}
		for id := range devices {
			devices_id = append(devices_id, id)
		}
		json.NewEncoder(w).Encode(devices_id)
	}
}
