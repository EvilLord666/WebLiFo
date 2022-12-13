package application

import "WebLiFo/logging"

type AppRunner interface {
	Start() (bool, error)
	Stop() (bool, error)
	Init() (bool, error)
	GetLogger() *logging.AppLogger
}
