package daemon

import (
	"fmt"
	"pirem/defs"
	"pirem/irdevice"
)

func (d Daemon) getDeviceHandler(devName string) (irdevice.Device, error) {
	dev, exist := d.devices[devName]
	if !exist {
		return dev, fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}
	return dev, nil
}
