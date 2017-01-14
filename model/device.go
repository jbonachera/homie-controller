package model

type DeviceStats struct {
	uptime   int
	signal   int
	interval int
}

type DeviceFirmware struct {
	name     string
	version  string
	checksum string
}

type Device struct {
	id             string
	online         bool
	name           string
	localip        string
	mac            string
	stats          DeviceStats
	fw             DeviceFirmware
	implementation string
}
