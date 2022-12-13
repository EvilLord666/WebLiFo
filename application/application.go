package application

import (
	"WebLiFo/config"
	"gorm.io/gorm"
)

type WebLiFoAppRunner struct {
	cfg *config.ServerConfig
	db  *gorm.DB
}

func CreateApp(config string) AppRunner {
	return nil
}

func (w *WebLiFoAppRunner) Start() (bool, error) {
	return true, nil
}

func (w *WebLiFoAppRunner) Stop() (bool, error) {
	return true, nil
}

func (w *WebLiFoAppRunner) Init() (bool, error) {
	return true, nil
}
