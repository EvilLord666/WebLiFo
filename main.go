package main

import (
	"WebLiFo/application"
	"fmt"
	"github.com/wissance/stringFormatter"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	osSignal := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	app := application.CreateApp("./config.json")
	res, initErr := app.Init()
	logger := app.GetLogger()
	if initErr != nil {
		logger.Error("An error occurred during app init, terminating the app")
		os.Exit(-1)
	} else {
		logger.Info("Application was successfully initialized")
	}

	res, err := app.Start()
	if !res {
		msg := stringFormatter.Format("An error occurred during starting application, error is: {0}", err.Error())
		fmt.Println(msg)
	} else {
		logger.Info("Application was successfully started")
	}

	go func() {
		sig := <-osSignal
		logger.Info(stringFormatter.Format("Got signal from OS: {0}", sig))
		logger.Info(stringFormatter.Format("Got signal from OS: \"{0}\", stopping", sig))
		done <- true
	}()
	<-done

	res, err = app.Stop()
	if !res {
		msg := stringFormatter.Format("An error occurred during stopping application, error is: {0}", err.Error())
		fmt.Println(msg)
	} else {
		logger.Info("Application was successfully stopped")
	}
}
