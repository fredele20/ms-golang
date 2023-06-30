package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fredele20/microservice-practice/ms.products/cache"
	"github.com/fredele20/microservice-practice/ms.products/config"
	"github.com/fredele20/microservice-practice/ms.products/core"
	"github.com/fredele20/microservice-practice/ms.products/database/mongod"
	"github.com/fredele20/microservice-practice/ms.products/handlers"
	"github.com/fredele20/microservice-practice/ms.products/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	secrets := config.GetSecrets()
	logger := logrus.New()
	redis := *cache.NewRedisConnection()

	address := fmt.Sprintf("127.0.0.1:%s", secrets.Port)

	fileLogger := "logs.log"

	logFile, err := os.OpenFile(fileLogger, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println("error opening file: ", err)
		return
	}

	defer logFile.Close()

	logrus.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Println("log file created")

	router := gin.New()
	router.Use(gin.Logger())

	// var client *mongo.Client
	db, _ := mongod.DBInstance(secrets.DatabaseURL, secrets.DatabaseName)

	core := core.NewProductService(db, logger, redis)

	routes := routes.NewRouteService(core)

	handler := handlers.NewHandlers(routes)

	handlers.RouteHandlers(router, *handler)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted for api-1"})
	})
	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(address)
}
