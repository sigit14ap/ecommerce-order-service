package domain

import "time"

type Order struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint64    `gorm:"not null;" json:"user_id"`
	Status    string    `gorm:"type:varchar(255);not null;" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	OrderID     uint64    `gorm:"not null;" json:"order_id"`
	WarehouseID uint64    `gorm:"not null;" json:"warehouse_id"`
	ProductID   uint64    `gorm:"not null;" json:"product_id"`
	Quantity    int       `gorm:"not null;" json:"quantity"`
	Price       float64   `gorm:"not null;" json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
