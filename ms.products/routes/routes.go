package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/fredele20/microservice-practice/ms.products/core"
	"github.com/fredele20/microservice-practice/ms.products/models"
	"github.com/gin-gonic/gin"
)

type RouteService struct {
	core *core.ProductService
}

func NewRouteService(core *core.ProductService) *RouteService {
	return &RouteService{
		core: core,
	}
}

func (r RouteService) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product

		product.OwnerID = c.GetString("userId")
		product.OwnerName = c.GetString("firstName") + " " + c.GetString("lastName")

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newProduct, err := r.core.CreateProduct(ctx, product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, newProduct)
	}
}

func (r RouteService) GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), time.Second*100)
		defer cancle()

		var filter models.ProductFilter
		if err := c.BindJSON(&filter); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		productList, err := r.core.GetProducts(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, productList)
	}
}
