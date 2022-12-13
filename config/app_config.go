package config

type AppConfig struct {
	ServerCfg  ServerConfig  `json:"server"`
	LoggingCfg LoggingConfig `json:"logging"`
	DbCfg      DbConfig      `json:"database"`
}
