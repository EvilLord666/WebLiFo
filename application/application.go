package application

import (
	"WebLiFo/config"
	"WebLiFo/logging"
	"WebLiFo/model"
	"encoding/json"
	"fmt"
	"github.com/wissance/stringFormatter"
	"gorm.io/gorm"
	"io/ioutil"
	"path/filepath"
)

type WebLiFoAppRunner struct {
	configFile   *string
	cfg          *config.AppConfig
	db           *gorm.DB
	logger       *logging.AppLogger
	modelContext *model.AppModelContext
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
	err := w.readAppConfig()
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during reading app config file: {0}", err.Error()))
		return false, err
	}
	// 2. Init logger
	w.logger = logging.CreateLogger(&w.cfg.LoggingCfg)
	w.logger.Init()
	// 3. Init database
	// 4. Init web api
	return true, nil
}

func (w *WebLiFoAppRunner) GetLogger() *logging.AppLogger {
	return nil
}

func (w *WebLiFoAppRunner) readAppConfig() error {
	absPath, err := filepath.Abs(*w.configFile)
	if err != nil {
		w.logger.Error(stringFormatter.Format("An error occurred during getting config file abs path: {0}", err.Error()))
		return err
	}

	fileData, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during config file reading: {0}", err.Error()))
		return err
	}

	w.cfg = &config.AppConfig{}
	if err = json.Unmarshal(fileData, w.cfg); err != nil {
		fmt.Println(stringFormatter.Format("An error occurred during config file unmarshal: {0}", err.Error()))
		return err
	}

	return nil
}
