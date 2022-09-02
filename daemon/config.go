package daemon

import (
	"pirem/irdevice"
)

type Config struct {
	ServerPort uint16           `json:"server_port"`
	Devices    irdevice.Devices `json:"devices"`
}
