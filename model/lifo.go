package model

import (
	"gorm.io/gorm"
)

type Lifo struct {
	gorm.Model
	Name        string     `gorm:"type:varchar(64);not null;unique;"`
	Description string     `gorm:"type:varchar(128);not null;default:'';"`
	Size        int        `gorm:"type:int;not null;default:16"`
	Items       []LifoItem `gorm:"ForeignKey:LifoId"`
}

func (l *Lifo) GetTableName() string {
	return "lifos"
}
