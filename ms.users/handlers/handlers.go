package handlers

import (
	"github.com/fredele20/microservice-practice/ms.users/routes"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", routes.ListUsers())
	// incomingRoutes.GET("/users")
	incomingRoutes.GET("/users/:user_id", routes.GetUserById())
}

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", routes.Signup())
	incomingRoutes.POST("users/login", routes.Login())
	incomingRoutes.DELETE("users/logout", routes.Logout())
	incomingRoutes.POST("users/forgot-password", routes.ForgotPassword())
	incomingRoutes.POST("users/reset-password", routes.ResetPassword())
}
