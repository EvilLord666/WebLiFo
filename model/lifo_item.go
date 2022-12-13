package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type LifoItem struct {
	gorm.Model
	Value          string
	PreviousItemId sql.NullInt32
	LifoId         uint
	Lifo           Lifo
}
