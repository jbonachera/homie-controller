package humidity

import (
	"github.com/jbonachera/homie-controller/influxdb"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"strconv"
	"strings"
	"time"
)

var nodeType string = "weather_sensor"

type WeatherSensorNode struct {
	name            string
	baseTopic       string
	Nodetype        string    `json:"type"`
	HumidityPercent float64   `json:"humidity"`
	Temperature     float64   `json:"temperature"`
	Pressure        float64   `json:"pressure"`
	BatteryVoltage  float64   `json:"battery"`
	Room            string    `json:"room"`
	LastUpdate      time.Time `json:"last_update"`
	ParentId        string    `json:"-"`
}

func (t WeatherSensorNode) GetName() string {
	return t.name
}
func (t WeatherSensorNode) GetType() string {
	return nodeType
}
func (t WeatherSensorNode) GetProperties() []string {
	return []string{"humidity", "temperature", "room", "pressure", "battery"}
}
func (t WeatherSensorNode) GetPoint() *metric.Metric {
	return metric.New("weather", map[string]string{"room": t.Room, "sensor": t.name},
		map[string]interface{}{
			"humidity":    t.HumidityPercent,
			"temperature": t.Temperature,
			"pressure":    t.Pressure,
			"battery":     t.BatteryVoltage,
		})
}

func (t *WeatherSensorNode) MessageHandler(message homieMessage.HomieMessage) {
	topicComponents := strings.Split(message.Path, "/")
	node, property := topicComponents[0], topicComponents[1]
	if node != t.name {
		log.Debug("received message for " + node + " but we are " + t.name)
		return
	}
	switch property {
	case "room":
		t.Room = message.Payload
	case "humidity":
		percent, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.HumidityPercent = percent
			t.LastUpdate = time.Now()
			if influxdb.Ready() {
				influxdb.PublishMetric(metric.New("humidity", map[string]string{"room": t.Room, "sensor": t.name},
					map[string]interface{}{
						"humidity": t.HumidityPercent,
					}))
			}
		}
	case "temperature":
		degrees, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.Temperature = degrees
			t.LastUpdate = time.Now()
			if influxdb.Ready() {
				influxdb.PublishMetric(metric.New("temperature", map[string]string{"room": t.Room, "sensor": t.name},
					map[string]interface{}{
						"temperature": t.Temperature,
					}))
			}
		}

	case "pressure":
		pressure, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.Pressure = pressure
			t.LastUpdate = time.Now()
			if influxdb.Ready() {
				influxdb.PublishMetric(metric.New("pressure", map[string]string{"room": t.Room, "sensor": t.name},
					map[string]interface{}{
						"pressure": t.Pressure,
					}))
			}
		}
	case "battery":
		voltage, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.BatteryVoltage = voltage
			t.LastUpdate = time.Now()
			if influxdb.Ready() {
				influxdb.PublishMetric(metric.New("battery", map[string]string{"room": t.Room, "sensor": t.name},
					map[string]interface{}{
						"battery": t.BatteryVoltage,
					}))
			}
		}

	}

}
