package main

import (
	"log"
	"productanalyzer/api/api"
	"productanalyzer/api/config"
	"productanalyzer/api/db"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.BotDetectionMiddleware())
	apiRouter := router.Group("/api")
	api.SetupRoutes(apiRouter)
	return router
}

func main() {
	err := config.Config.Load()
	if err != nil {
		log.Panic(err)
	}
	err = db.Connection.Connect()
	if err != nil {
		log.Panic(err)
	}
	err = db.Connection.FetchCollections()
	if err != nil {
		log.Panic(err)
	}
	defer db.Connection.Close()
	db.Connection.Initialize()
	router := SetupRouter()
	if err := router.Run(config.Config.HOST); err != nil {
		log.Panic(err)
	}
}
