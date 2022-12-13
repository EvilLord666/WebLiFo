package config

type DbConfig struct {
	Hostname           string `json:"hostname"`
	Port               int    `json:"port"`
	Database           string `json:"database"`
	User               string `json:"user"`
	Password           string `json:"password"`
	MaxConnections     int    `json:"max_open_connections"`
	MaxIdleConnections int    `json:"max_idle_connections"`
}
