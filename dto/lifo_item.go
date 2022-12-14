package dto

type LifoItem struct {
	Id             uint   `json:"id"`
	PreviousItemId uint   `json:"previous_item_id"`
	LifoId         uint   `json:"lifo_id"`
	Value          string `json:"value"`
	//Value          json.RawMessage `json:"value"`
}
