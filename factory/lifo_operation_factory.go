package factory

import (
	"WebLiFo/dto"
	"WebLiFo/model"
)

func CreateLifoOperation(operation dto.OperationType, success bool, message string, lifoItem *model.LifoItem) dto.LifoOperation {
	return dto.LifoOperation{Operation: operation, Success: success, Message: message, LifoItem: CreateLifoItem(lifoItem)}
}
