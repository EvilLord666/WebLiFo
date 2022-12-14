package model

import (
	"gorm.io/gorm"
)

type LifoItem struct {
	gorm.Model
	Value          string    `gorm:"type:varchar(4096);not null;default:'';"`
	Previous       *LifoItem `gorm:"foreignkey:PreviousItemId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PreviousItemId uint
	LifoId         uint
	Lifo           Lifo `gorm:"foreignkey:LifoId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
