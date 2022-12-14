package factory

import (
	"WebLiFo/dto"
	"WebLiFo/model"
)

func CreateLifoItem(lifoItem *model.LifoItem) dto.LifoItem {
	return dto.LifoItem{Id: lifoItem.ID, LifoId: lifoItem.LifoId, PreviousItemId: lifoItem.PreviousItemId,
		Value: lifoItem.Value}
}
