package humidity

import (
	"github.com/jbonachera/homie-controller/model/node"
	"time"
)

type WeatherSensorFactory struct{}

func (WeatherSensorFactory) New(name string, parent string, baseTopic string) node.Type {
	return &WeatherSensorNode{name: name, baseTopic: baseTopic,
		Nodetype:        "weather_controller",
		Temperature:     0,
		Pressure:        0,
		BatteryVoltage:  0,
		ParentId:        parent,
		LastUpdate:      time.Now(),
		Room:            "",
		HumidityPercent: 0,
	}
}

func init() {
	node.RegisterNodeTypeFactory("weather_sensor", WeatherSensorFactory{})
}
