package factory

import (
	"WebLiFo/dto"
	"WebLiFo/model"
)

func CreateLifoInfo(lifo *model.Lifo) dto.LifoInfo {
	return dto.LifoInfo{Id: lifo.ID, Name: lifo.Name, Description: lifo.Description}
}

func CreateLifoWithItems(lifo *model.Lifo) dto.Lifo {
	items := make([]dto.LifoItem, len(lifo.Items))
	for index, item := range lifo.Items {
		items[index] = CreateLifoItem(&item)
	}
	return dto.Lifo{Id: lifo.ID, Name: lifo.Name, Description: lifo.Description, Items: items}
}
