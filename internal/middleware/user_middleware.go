package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/order-service/helpers"
	"github.com/sigit14ap/order-service/internal/services"
)

func UserMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if token == "" {
			helpers.ErrorResponse(context, http.StatusUnauthorized, "Authorization required")
			context.Abort()
			return
		}

		user, err := userService.UserDetail(context)
		if err != nil {
			helpers.ErrorResponse(context, http.StatusUnauthorized, "User does not allowed")
			return
		}

		context.Set("userID", user.ID)
		context.Next()
	}
}
