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

var LifoIsFull error = errors.New("lifo is full, it is unable to add more items")

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
		newTopLifoItem = model.LifoItem{LifoId: lifoId, PreviousItemId: 0, Value: lifoItem.Value}
		newTopLifoItem.ID = g.GetNextTableId(tx, newTopLifoItem.GetTableName())
		err = tx.Create(&newTopLifoItem).Error
		if err != nil {
			return err
		}
		// 3. Set to previous LifoItem previous_lifo_item_id to null
		if len(lifo.Items) > 0 {
			previousTopItem := lifo.Items[0]
			previousTopItem.PreviousItemId = newTopLifoItem.ID
			previousTopItem.Previous = &newTopLifoItem
		}
		return nil
	})
	return newTopLifoItem, err
}
