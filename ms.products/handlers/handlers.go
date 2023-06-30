package handlers

import (
	"github.com/fredele20/microservice-practice/ms.products/middlewares"
	"github.com/fredele20/microservice-practice/ms.products/routes"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	handler *routes.RouteService
}

func NewHandlers(handler *routes.RouteService) *Handlers {
	return &Handlers{
		handler: handler,
	}
}

func RouteHandlers(incomingRoutes *gin.Engine, h Handlers) {
	incomingRoutes.GET("/products", h.handler.GetProducts())
	incomingRoutes.Use(middlewares.Authentication())
	incomingRoutes.POST("/products", h.handler.CreateProduct())
}
