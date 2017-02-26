package esp8266

type wifiConfig struct {
	Ssid     string `json:"ssid"`
	Password string `json:"password,omitempty"`
}

type mqttConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	BaseTopic string `json:"base_topic"`
	Auth      bool   `json:"auth"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}

type OTAConfig struct {
	Enabled bool `json:"enabled"`
}

type config struct {
	Wifi     wifiConfig  `json:"wifi"`
	Mqtt     mqttConfig  `json:"mqtt"`
	Ota      OTAConfig   `json:"ota"`
	Settings interface{} `json:"settings,omitempty"`
}
