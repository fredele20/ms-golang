package handlers

import (
	"github.com/fredele20/microservice-practice/ms.users/routes"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	routes *routes.UserRoutes
}

func NewUserHandler(routes *routes.UserRoutes) *UserHandler {
	return &UserHandler{
		routes: routes,
	}
}

func UserRoutes(incomingRoutes *gin.Engine, u UserHandler) {
	// incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", u.routes.ListUsers())
	// incomingRoutes.GET("/users")
	// incomingRoutes.GET("/users/:user_id", routes.GetUserById())
}

func AuthRoutes(incomingRoutes *gin.Engine, u UserHandler) {
	incomingRoutes.POST("users/signup", u.routes.Signup())
	incomingRoutes.POST("users/login", u.routes.Login())
	incomingRoutes.DELETE("users/logout", u.routes.Logout())
	incomingRoutes.POST("users/forgot-password", u.routes.ForgotPassword())
	incomingRoutes.POST("users/reset-password", u.routes.ResetPassword())
}
