package domain

import "time"

type Product struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	ShopID    uint64    `gorm:"not null;" json:"shop_id"`
	Name      string    `gorm:"type:varchar(255);not null;" json:"name"`
	Price     float64   `gorm:"type:int;not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
