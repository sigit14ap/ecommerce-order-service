package usecase

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sigit14ap/order-service/helpers"
	"github.com/sigit14ap/order-service/internal/delivery/dto"
	"github.com/sigit14ap/order-service/internal/domain"
	repository "github.com/sigit14ap/order-service/internal/repository/mysql"
)

type OrderUsecase struct {
	orderRepository repository.OrderRepository
	mutex           sync.Mutex
}

func NewOrderUsecase(orderRepository repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{
		orderRepository: orderRepository,
	}
}

func (uc *OrderUsecase) Checkout(orderData dto.OrderDTO) (*domain.Order, error) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	tx := uc.orderRepository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order := domain.Order{
		UserID: orderData.UserID,
		Status: helpers.OrderStatusPending,
	}

	items := []domain.OrderItem{}

	for _, item := range orderData.Items {
		product, err := uc.orderRepository.GetProductByID(item.ProductID)

		if err != nil {
			if errors.Is(err, helpers.StockNotFound) {
				tx.Rollback()
				return nil, fmt.Errorf("product %d is not found", item.ProductID)
			}

			tx.Rollback()
			return nil, err
		}

		stock, err := uc.orderRepository.GetStockByProductID(tx, item.ProductID, item.Quantity)

		if err != nil {
			if errors.Is(err, helpers.StockNotFound) {
				tx.Rollback()
				return nil, fmt.Errorf("product %s stock is not enough order request", product.Name)
			}

			tx.Rollback()
			return nil, err
		}

		items = append(items, domain.OrderItem{
			WarehouseID: stock.WarehouseID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			Price:       product.Price,
		})

		stock.Quantity -= item.Quantity

		err = uc.orderRepository.UpdateStock(tx, stock)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	orderCreated, err := uc.orderRepository.CreateOrder(tx, order, items)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &orderCreated, nil
}
