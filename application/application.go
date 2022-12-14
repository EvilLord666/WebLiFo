package application

import (
	"WebLiFo/api/rest"
	"WebLiFo/config"
	"WebLiFo/logging"
	"WebLiFo/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	r "github.com/wissance/gwuu/api/rest"
	"github.com/wissance/stringFormatter"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type WebLiFoAppRunner struct {
	configFile    *string
	cfg           *config.AppConfig
	logger        *logging.AppLogger
	modelContext  *model.AppModelContext
	httpHandler   *http.Handler
	webApiHandler *r.WebApiHandler
	webApiContext *rest.WebApiContext
}

func CreateApp(config string) AppRunner {
	app := &WebLiFoAppRunner{configFile: &config}
	appRunner := AppRunner(app)
	return appRunner
}

func (w *WebLiFoAppRunner) Start() (bool, error) {
	var err error
	go func() {
		err = w.startWebService()
		if err != nil {
			w.logger.Error(stringFormatter.Format("An error occurred during API Service Start"))
		}
	}()
	return err == nil, err
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
	w.modelContext = model.CreateAppModelContext(w.logger, &w.cfg.DbCfg)

	err = w.modelContext.Init()
	if err != nil {
		w.logger.Error(stringFormatter.Format("An error occurred during init ORM Db Context: {0}", err.Error()))
		return false, err
	}
	err = w.modelContext.Prepare()
	if err != nil {
		w.logger.Error(stringFormatter.Format("An error occurred during prepare Database (structure migrate & data init): {0}", err.Error()))
		return false, err
	}
	// 4. Init web api
	err = w.initRestApi()
	if err != nil {
		w.logger.Error(stringFormatter.Format("An error occurred during init Web Api: {0}", err.Error()))
		return false, err
	}
	return true, nil
}

func (w *WebLiFoAppRunner) GetLogger() *logging.AppLogger {
	return w.logger
}

func (w *WebLiFoAppRunner) initRestApi() error {
	w.webApiHandler = r.NewWebApiHandler(true, r.AnyOrigin)
	w.webApiContext = &rest.WebApiContext{DbContext: w.modelContext.GetContext(), Logger: w.logger}
	router := w.webApiHandler.Router
	// Setting up listener for logging
	w.webApiHandler.HandleFunc(router, "/api/lifo", w.webApiContext.GetAllLifo, http.MethodGet)
	w.webApiHandler.HandleFunc(router, "/api/lifo/{id}", w.webApiContext.GetLifoById, http.MethodGet)
	w.webApiHandler.HandleFunc(router, "/api/lifo/{id}/items", w.webApiContext.GetLifoByIdWithItems, http.MethodGet)
	w.webApiHandler.HandleFunc(router, "/api/lifo", w.webApiContext.CreateLifo, http.MethodPost)
	w.webApiHandler.HandleFunc(router, "/api/lifo/{id}", w.webApiContext.UpdateLifo, http.MethodPut)
	w.webApiHandler.HandleFunc(router, "/api/lifo/{id}", w.webApiContext.DeleteLifo, http.MethodDelete)
	w.webApiHandler.HandleFunc(router, "/api/lifo/{id}/push", w.webApiContext.PushLifo, http.MethodPost)
	w.webApiHandler.HandleFunc(router, "/api/lifo/{id}/pop", w.webApiContext.PopLifo, http.MethodPost)

	appenderIndex := w.logger.GetAppenderIndex(config.RollingFile, w.cfg.LoggingCfg.Appenders)
	if appenderIndex == -1 {
		w.logger.Info("The RollingFile appender is not found.")
		var resultRouter http.Handler = router
		w.httpHandler = &resultRouter
		return nil
	}
	w.httpHandler = w.createHttpLoggingHandler(appenderIndex, router)
	return nil
}

func (w *WebLiFoAppRunner) startWebService() error {
	var err error
	addressTemplate := "{0}:{1}"
	address := stringFormatter.Format(addressTemplate, w.cfg.ServerCfg.Address, w.cfg.ServerCfg.Port)
	switch w.cfg.ServerCfg.Schema {
	case config.HTTP:
		w.logger.Info(stringFormatter.Format("Starting \"HTTP\" WEB API Service on address: \"{0}\"", address))
		err = http.ListenAndServe(address, *w.httpHandler)
		if err != nil {
			w.logger.Error(stringFormatter.Format("An error occurred during attempt to start \"HTTP\" WEB API Service: {0}", err.Error()))
		}
	case config.HTTPS:
		return errors.New("not supported")
	}
	return err
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

func (w *WebLiFoAppRunner) createHttpLoggingHandler(index int, router *mux.Router) *http.Handler {
	var resultRouter http.Handler = router

	destination := w.cfg.LoggingCfg.Appenders[index].Destination
	lumberjackWriter := lumberjack.Logger{
		Filename:   string(destination.File),
		MaxSize:    destination.MaxSize,
		MaxAge:     destination.MaxAge,
		MaxBackups: destination.MaxBackups,
		LocalTime:  destination.LocalTime,
		Compress:   false,
	}

	if w.cfg.LoggingCfg.LogHTTP {
		if w.cfg.LoggingCfg.ConsoleOutHTTP {
			writer := io.MultiWriter(&lumberjackWriter, os.Stdout)
			resultRouter = handlers.LoggingHandler(writer, router)
		} else {
			resultRouter = handlers.LoggingHandler(&lumberjackWriter, router)
		}
	}
	return &resultRouter
}
