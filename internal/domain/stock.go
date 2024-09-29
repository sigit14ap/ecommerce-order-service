package domain

import "time"

type Stock struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	WarehouseID uint64    `gorm:"not null;" json:"warehouse_id"`
	ProductID   uint64    `gorm:"not null;" json:"product_id"`
	Quantity    int       `gorm:"not null;" json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
