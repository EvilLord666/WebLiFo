package factory

import (
	"WebLiFo/dto"
	"WebLiFo/model"
)

func CreateLifoItem(lifoItem *model.LifoItem) dto.LifoItem {
	return dto.LifoItem{Id: lifoItem.ID, LifoId: lifoItem.LifoId, PreviousItemId: uint(lifoItem.PreviousItemId.Int32),
		Value: lifoItem.Value}
}
