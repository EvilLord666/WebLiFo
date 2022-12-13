package application

import (
	"WebLiFo/config"
	"WebLiFo/logging"
	"gorm.io/gorm"
)

type WebLiFoAppRunner struct {
	configFile *string
	cfg        *config.ServerConfig
	db         *gorm.DB
}

func CreateApp(config string) AppRunner {
	app := &WebLiFoAppRunner{configFile: &config}
	appRunner := AppRunner(app)
	return appRunner
}

func (w *WebLiFoAppRunner) Start() (bool, error) {
	return true, nil
}

func (w *WebLiFoAppRunner) Stop() (bool, error) {
	return true, nil
}

func (w *WebLiFoAppRunner) Init() (bool, error) {
	// 1. Read config, get settings
	// 2. Init logger
	// 3. Init database
	// 4. Init web api
	return true, nil
}

func (w *WebLiFoAppRunner) getConfigs() {

}

func (w *WebLiFoAppRunner) GetLogger() *logging.AppLogger {
	return nil
}
