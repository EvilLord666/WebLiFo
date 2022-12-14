package dto

type LifoInfo struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int    `json:"size"`
}

type Lifo struct {
	Id          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Size        int        `json:"size"`
	Items       []LifoItem `json:"items"`
}
