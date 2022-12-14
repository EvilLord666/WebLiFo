package dto

type OperationType string

const (
	Push OperationType = "PUSH"
	Pop  OperationType = "POP"
)

type LifoOperation struct {
	Operation OperationType `json:"operation"`
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	LifoItem  LifoItem      `json:"item"`
}
