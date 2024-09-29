package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/order-service/helpers"
)

func ServiceMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		serviceToken := context.GetHeader("X-Service-Token")
		if serviceToken == "" {
			helpers.ErrorResponse(context, http.StatusForbidden, "X-Service-Token required")
			context.Abort()
			return
		}

		if strings.TrimSpace(serviceToken) != os.Getenv("APP_SECRET") {
			helpers.ErrorResponse(context, http.StatusForbidden, "Invalid service token")
			context.Abort()
			return
		}

		context.Next()
	}
}
