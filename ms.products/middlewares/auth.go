package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			ctx.Abort()
			return
		}

		claims, err := ValidateToken(clientToken)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			fmt.Println(err)
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("firstName", claims.FirstName)
		ctx.Set("lastName", claims.LastName)
		ctx.Set("userId", claims.UserId)
		// ctx.Set("expiresAt", claims.ExpiresAt)
		ctx.Next()
	}
}
