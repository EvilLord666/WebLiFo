package model

import (
	"gorm.io/gorm"
)

type Lifo struct {
	gorm.Model
	Name        string
	Description string
	Size        int
	Items       []LifoItem
}
