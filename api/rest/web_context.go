package rest

import (
	"WebLiFo/logging"
	"gorm.io/gorm"
)

type WebApiContext struct {
	DbContext *gorm.DB
	Logger    *logging.AppLogger
}
