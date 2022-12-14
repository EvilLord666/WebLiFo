package managers

import (
	"WebLiFo/dto"
	"WebLiFo/logging"
	"WebLiFo/model"
	"database/sql"
	"errors"
	"github.com/wissance/stringFormatter"
	"gorm.io/gorm"
)

var LifoIsFull error = errors.New("lifo is full, it is unable to add more items")
var LifoIsEmpty error = errors.New("there is no more items in lifo, nothing to pop")

func PushToLifo(lifoId uint, lifoItem *dto.LifoItem, db *gorm.DB, logger *logging.AppLogger) (model.LifoItem, error) {
	var newTopLifoItem model.LifoItem
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Check whether we could add one more item or not
		lifo, err := GetLifoByIdWithItems(lifoId, tx, logger)
		if err != nil {
			return err
		}
		if lifo.Size <= len(lifo.Items) {
			logger.Error(stringFormatter.Format("It is unable to add to lifo with id: {0} more items, lifo is full", lifoId))
			return LifoIsFull
		}
		// 2. Save new top lifo item
		newTopLifoItem = model.LifoItem{LifoId: lifoId, PreviousItemId: sql.NullInt32{Valid: false}, Previous: nil, Value: lifoItem.Value}
		//newTopLifoItem.ID = g.GetNextTableId(tx, newTopLifoItem.GetTableName())
		err = tx.Save(&newTopLifoItem).Error
		if err != nil {
			return err
		}
		// 3. Set to previous LifoItem previous_lifo_item_id to null
		var previousTopItem model.LifoItem
		if len(lifo.Items) > 0 {
			previousTopItem = selectTopItem(&lifo.Items)
			previousTopItem.PreviousItemId = sql.NullInt32{Valid: true, Int32: int32(newTopLifoItem.ID)}
			previousTopItem.Previous = &newTopLifoItem
			err = tx.Save(&previousTopItem).Error
		}
		return err
	})
	return newTopLifoItem, err
}

func PopFromLifo(lifoId uint, db *gorm.DB, logger *logging.AppLogger) (model.LifoItem, error) {
	var topItem model.LifoItem
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Grt Lifo, check Items
		lifo, err := GetLifoByIdWithItems(lifoId, tx, logger)
		if err != nil {
			return err
		}
		if len(lifo.Items) == 0 {
			return LifoIsEmpty
		}
		// 2. Find 2 items : 1 - top item , and item that have
		topItem = selectTopItem(&lifo.Items)

		if len(lifo.Items) > 1 {
			var nextTopItem model.LifoItem
			err = tx.Where("previous_item_id = ? AND lifo_id = ?", topItem.ID, lifoId).First(&nextTopItem).Error
			if err != nil {
				logger.Error("An unexpected error occurred during getting item that should point on the top of lifo")
				return err
			}
			nextTopItem.PreviousItemId = sql.NullInt32{Valid: false}
			nextTopItem.Previous = nil
			err = tx.Save(&nextTopItem).Error
			if err != nil {
				logger.Error(stringFormatter.Format("An error occurred during setting lifo item with id: {0}, to point on the top of lifo: {1}, error: {2}",
					nextTopItem.ID, lifoId, err.Error()))
				return err
			}
			err = tx.Unscoped().Delete(&topItem).Error
			if err != nil {
				logger.Error(stringFormatter.Format("An error occurred during removing lifo: {0}, top item, error: {1}",
					lifoId, err.Error()))
				return err
			}
		}
		return nil
	})
	return topItem, err
}

func selectTopItem(items *[]model.LifoItem) model.LifoItem {
	for _, item := range *items {
		if !item.PreviousItemId.Valid {
			return item
		}
	}
	return model.LifoItem{}
}
