package managers

import (
	"WebLiFo/dto"
	"WebLiFo/logging"
	"WebLiFo/model"
	"errors"
	g "github.com/wissance/gwuu/gorm"
	"github.com/wissance/stringFormatter"
	"gorm.io/gorm"
)

var BadSizeError error = errors.New("bad LIFO size")
var BadLifoName error = errors.New("bad LIFO name")

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
	err := db.Preload("Items").First(&lifo, "id = ?", id).Error
	return lifo, err
}

func CreateLifo(lifo *dto.LifoInfo, db *gorm.DB, logger *logging.AppLogger) (model.Lifo, error) {
	if len(lifo.Name) == 0 {
		msg := "lifo name can not be null or empty"
		logger.Warn(stringFormatter.Format("It is unable to create new lifo, reason: {0}", msg))
		return model.Lifo{}, errors.New(msg)
	}
	newItem := model.Lifo{Name: lifo.Name, Description: lifo.Description, Size: lifo.Size}
	newItem.ID = g.GetNextTableId(db, newItem.GetTableName())
	err := db.Create(&newItem).Error
	return newItem, err
}

func UpdateLifo(lifo dto.LifoInfo, id uint, db *gorm.DB, logger *logging.AppLogger) (model.Lifo, error) {
	existingLifo, err := GetLifoByIdWithItems(id, db, logger)
	if err != nil {
		return existingLifo, err
	}
	if len(lifo.Name) > 0 {
		existingLifo.Name = lifo.Name
	}
	existingLifo.Description = lifo.Description
	// check size
	if lifo.Size < len(existingLifo.Items) {
		msg := stringFormatter.Format("Size can not be set to {0} items, because lifo already contains: {1} items")
		logger.Error(msg)
		return existingLifo, BadSizeError
	}
	existingLifo.Size = lifo.Size
	err = db.Save(&existingLifo).Error
	return existingLifo, err
}

func DeleteLifo(id uint, db *gorm.DB, logger *logging.AppLogger) (bool, error) {
	_, err := GetLifoById(id, db, logger)
	if err != nil {
		return false, err
	}
	err = db.Unscoped().Delete(&model.Lifo{}, id).Error
	return err == nil, err
}
