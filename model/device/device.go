package device

type DeviceStats struct {
	Uptime   int
	Signal   int
	Interval int
}

type DeviceFirmware struct {
	Name     string
	Version  string
	Checksum string
}

type Device struct {
	Id             string
	Online         bool
	Name           string
	Localip        string
	Mac            string
	Stats          DeviceStats
	Fw             DeviceFirmware
	Implementation string
}

func New(id string) Device {
	return Device{id, false, "", "", "", DeviceStats{0, 0, 0}, DeviceFirmware{"", "", ""}, ""}
}
