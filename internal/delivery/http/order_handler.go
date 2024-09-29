package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/sigit14ap/order-service/helpers"
	"github.com/sigit14ap/order-service/internal/delivery/dto"
	"github.com/sigit14ap/order-service/internal/usecase"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

var validate *validator.Validate

func NewOrderHandler(orderUsecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

func (handler *OrderHandler) Checkout(context *gin.Context) {
	userIDValue, exists := context.Get("userID")

	if !exists {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "User ID not found")
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		helpers.ErrorResponse(context, http.StatusUnauthorized, "invalid User ID")
		return
	}

	validate = validator.New()

	var orderData dto.OrderDTO

	err := context.BindJSON(&orderData)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(orderData)
	if err != nil {
		helpers.ErrorValidationResponse(context, err)
		return
	}

	orderData.UserID = userID
	order, err := handler.orderUsecase.Checkout(orderData)
	if err != nil {
		helpers.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SuccessResponse(context, order)
}
