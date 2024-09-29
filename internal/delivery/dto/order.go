package dto

type OrderDTO struct {
	UserID uint64         `json:"-"`
	Items  []OrderItemDTO `json:"items"`
}

type OrderItemDTO struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
