package application

import "gorm.io/gorm"

type WebLiFoAppRunner struct {
	db *gorm.DB
}

func CreateApp(config string) AppRunner {
	return nil
}
