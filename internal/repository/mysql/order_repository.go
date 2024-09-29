package repository

import (
	"errors"
	"time"

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

	GetStockProductByWarehouse(tx *gorm.DB, productID uint64, warehouseID uint64) (domain.Stock, error)
	GetUnpaidOrders(tx *gorm.DB, duration time.Duration) ([]domain.Order, error)
	GetOrderItems(tx *gorm.DB, orderID uint64) ([]domain.OrderItem, error)
	UpdateOrderStatus(tx *gorm.DB, orderID uint64, status string) error
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

func (repository *orderRepository) GetUnpaidOrders(tx *gorm.DB, duration time.Duration) ([]domain.Order, error) {
	var orders []domain.Order
	cutoffTime := time.Now().Add(-duration)

	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("status = ? AND created_at < ?", helpers.OrderStatusPending, cutoffTime).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (repository *orderRepository) GetStockProductByWarehouse(tx *gorm.DB, productID uint64, warehouseID uint64) (domain.Stock, error) {
	var stock domain.Stock
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&stock, "product_id = ? AND warehouse_id = ?", productID, warehouseID).Error

	return stock, err
}

func (repository *orderRepository) GetOrderItems(tx *gorm.DB, orderID uint64) ([]domain.OrderItem, error) {
	var items []domain.OrderItem

	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (repository *orderRepository) UpdateOrderStatus(tx *gorm.DB, orderID uint64, status string) error {
	return tx.Model(&domain.Order{}).Where("id = ?", orderID).Update("status", status).Error
}
