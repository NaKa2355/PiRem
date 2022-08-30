package respjson

import (
	"encoding/json"
	"pirem/irdevice"
)

func DeviceToJson(dev irdevice.Device) ([]byte, error) {
	return json.Marshal(dev)
}

func DevicesToJson(devs map[string]irdevice.Device) ([]byte, error) {
	return json.Marshal(devs)
}
