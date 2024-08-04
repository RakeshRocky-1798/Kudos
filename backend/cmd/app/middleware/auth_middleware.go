package middleware

import (
	"github.com/gin-gonic/gin"
	"kleos/cmd/app/clients"
	"net/http"
)

func AuthMiddleware(client *clients.MjolnirClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionToken := ctx.GetHeader("X-Session-Token")

		if sessionToken == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		response, err := client.GetSessionResponse(sessionToken)

		if err != nil || response.StatusCode != http.StatusOK {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", response)
		ctx.Set("email", response.EmailId)
		ctx.Set("roles", response.Roles)

		ctx.Next()
	}
}
