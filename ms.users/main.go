package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fredele20/microservice-practice/ms.users/cache"
	"github.com/fredele20/microservice-practice/ms.users/config"
	"github.com/fredele20/microservice-practice/ms.users/core"
	"github.com/fredele20/microservice-practice/ms.users/db/mongod"
	"github.com/fredele20/microservice-practice/ms.users/handlers"
	"github.com/fredele20/microservice-practice/ms.users/libs/session"
	"github.com/fredele20/microservice-practice/ms.users/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	secrets := config.GetSecrets()
	db, err := mongod.MongoConnection(secrets.DatabaseURL, secrets.DatabaseName)
	logger := logrus.New()
	redis, err := cache.NewRedisConnection(secrets.RedisAddress)
	if err != nil {
		log.Fatal(err)
	}
	address := fmt.Sprintf("127.0.0.1:%s", secrets.Port)

	fileLogger := "logs.log"

	logFile, err := os.OpenFile(fileLogger, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return
	}

	defer logFile.Close()

	logrus.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	log.Println("log file created")

	router := gin.New()
	router.Use(gin.Logger())

	session := session.NewSessionManager(redis, db)
	core := core.NewUserService(redis, db, *session, logger)
	routes := routes.NewUserRoute(core)
	handler := handlers.NewUserHandler(routes)

	handlers.UserRoutes(router, *handler)
	handlers.AuthRoutes(router, *handler)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted for api-1"})
	})
	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(address)

}
