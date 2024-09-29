package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sigit14ap/order-service/config"
	"github.com/sigit14ap/order-service/helpers"
	delivery "github.com/sigit14ap/order-service/internal/delivery/http"
	"github.com/sigit14ap/order-service/internal/domain"
	repository "github.com/sigit14ap/order-service/internal/repository/mysql"
	"github.com/sigit14ap/order-service/internal/router"
	"github.com/sigit14ap/order-service/internal/services"
	"github.com/sigit14ap/order-service/internal/usecase"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	log := helpers.InitializeLogs()

	if log == nil {
		log.Fatal("Logger failed to started")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&domain.Order{}, &domain.OrderItem{})
	if err != nil {
		log.Fatalf("failed to auto-migrate Order model: %v", err)
	}

	userService := services.NewUserService(cfg.UserServiceUrl, cfg.AppSecret)

	orderRepo := repository.NewOrderRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderHandler := delivery.NewOrderHandler(orderUsecase)

	router := router.NewRouter(orderHandler, userService)

	go func() {
		for {
			log.Info("Release Stock Started")
			err := orderUsecase.ReleaseStock()

			if err != nil {
				log.Fatalf("Error when release stock product: %v", err)
			}

			// 5 Minutes
			time.Sleep(5 * time.Minute)
		}
	}()

	log.Info(router.Run(":" + os.Getenv("APP_PORT")))
}
