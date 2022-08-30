package configjson

import "encoding/json"

type Device struct {
	Name       string          `json:"name"`
	PluginPath string          `json:"plugin_path"`
	Config     json.RawMessage `json:"config"`
}

type Config struct {
	ServerPort uint16   `json:"server_port"`
	Devices    []Device `json:"devices"`
}
