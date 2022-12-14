package factory

import (
	"WebLiFo/dto"
	"WebLiFo/model"
)

func CreateLifoInfo(lifo *model.Lifo) dto.LifoInfo {
	return dto.LifoInfo{Id: lifo.ID, Name: lifo.Name, Description: lifo.Description}
}
