package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type LifoItem struct {
	gorm.Model
	Value          string    `gorm:"type:varchar(4096);not null;default:'';"`
	Previous       *LifoItem `gorm:"foreignkey:PreviousItemId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PreviousItemId sql.NullInt32
	LifoId         uint
	Lifo           Lifo `gorm:"foreignkey:LifoId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (l *LifoItem) GetTableName() string {
	return "lifo_items"
}
