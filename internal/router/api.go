package router

import (
	delivery "github.com/sigit14ap/order-service/internal/delivery/http"
	"github.com/sigit14ap/order-service/internal/middleware"
	"github.com/sigit14ap/order-service/internal/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(orderHandler *delivery.OrderHandler, userService *services.UserService) *gin.Engine {
	router := gin.New()
	v1 := router.Group("/api/v1")
	v1.Use(middleware.ServiceMiddleware())

	order := v1.Group("order")
	order.Use(middleware.UserMiddleware(userService))
	{
		order.POST("/checkout", orderHandler.Checkout)
	}

	return router
}
