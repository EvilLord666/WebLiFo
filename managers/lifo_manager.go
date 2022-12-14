package managers

import (
	"WebLiFo/logging"
	"WebLiFo/model"
	"gorm.io/gorm"
)

func GetLifoList(db *gorm.DB, logger *logging.AppLogger) ([]model.Lifo, error) {
	var lifoList []model.Lifo
	err := db.Find(&lifoList).Error
	return lifoList, err
}

func GetLifoById(id uint, db *gorm.DB, logger *logging.AppLogger) (model.Lifo, error) {
	var lifo model.Lifo
	err := db.First(&lifo, "id = ?", id).Error
	return lifo, err
}

func GetLifoByIdWithItems(id uint, db *gorm.DB, logger *logging.AppLogger) (model.Lifo, error) {
	var lifo model.Lifo
	err := db.First(&lifo, "id = ?", id).Error
	return lifo, err
}
