package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type success struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type failure struct {
	Error failureInfo `json:"error"`
}

type failureValidation struct {
	Error failureInfoValidation `json:"error"`
}

type failureInfo struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"invalid request body"`
}

type dataValidation struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

type failureInfoValidation struct {
	Code int              `json:"code" example:"400"`
	Data []dataValidation `json:"error"  example:"invalid request body"`
}

func SuccessResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusOK, success{Data: data, Message: "Success"})
}

func CreatedResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusCreated, success{Data: data, Message: "Success"})
}

func ErrorResponse(context *gin.Context, statusCode int, message string) {
	context.AbortWithStatusJSON(statusCode, failure{Error: failureInfo{
		Code:    statusCode,
		Message: message,
	}})
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	default:
		return "Invalid value"
	}
}

func ErrorValidationResponse(context *gin.Context, err error) {

	data := []dataValidation{}

	for _, err := range err.(validator.ValidationErrors) {
		data = append(data, dataValidation{
			Key:     err.Field(),
			Message: getErrorMessage(err),
		})
	}

	context.AbortWithStatusJSON(http.StatusUnprocessableEntity, failureValidation{
		Error: failureInfoValidation{
			Code: http.StatusUnprocessableEntity,
			Data: data,
		},
	})
}
