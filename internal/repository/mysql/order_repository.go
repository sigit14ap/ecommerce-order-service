package repository

import (
	"errors"

	"github.com/sigit14ap/order-service/helpers"
	"github.com/sigit14ap/order-service/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	BeginTransaction() *gorm.DB
	CreateOrder(tx *gorm.DB, order domain.Order, items []domain.OrderItem) (domain.Order, error)
	GetProductByID(productID uint64) (domain.Product, error)
	GetStockByProductID(tx *gorm.DB, productID uint64, quantity int) (domain.Stock, error)
	UpdateStock(tx *gorm.DB, stock domain.Stock) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (repository *orderRepository) BeginTransaction() *gorm.DB {
	return repository.db.Begin()
}

func (repository *orderRepository) CreateOrder(tx *gorm.DB, order domain.Order, items []domain.OrderItem) (domain.Order, error) {
	err := tx.Create(&order).Error

	if err != nil {
		return domain.Order{}, err
	}

	for i := range items {
		items[i].OrderID = order.ID
	}

	err = tx.Create(&items).Error

	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (repository *orderRepository) GetProductByID(productID uint64) (domain.Product, error) {
	var product domain.Product
	err := repository.db.First(&product, "id = ?", productID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Product{}, helpers.StockNotFound
	}

	return product, err
}

func (repository *orderRepository) GetStockByProductID(tx *gorm.DB, productID uint64, quantity int) (domain.Stock, error) {
	var stock domain.Stock
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&stock, "product_id = ? AND quantity >= ?", productID, quantity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Stock{}, helpers.StockNotFound
	}

	return stock, err
}

func (repository *orderRepository) UpdateStock(tx *gorm.DB, stock domain.Stock) error {
	return tx.Save(&stock).Error
}
