package main

import (
	"context"
	"log"

	"tempfunctiontools/controllers"

	"tempfunctiontools/internal/database"

	"github.com/gin-gonic/gin"
)

const (
	systemMsg = "You are a helpful assistant that can use tools to answer user queries."
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.LstdFlags)
	router := gin.Default()

	ctx := context.Background()
	dbConfig := database.DbConfig{}
	dbConfig.InitDb()

	agent := controllers.NewAgent(systemMsg, 3, &dbConfig)

	ctrl := controllers.NewChatController(ctx, agent, &dbConfig)

	router.POST("/api/chat", ctrl.GetChat)
	router.GET("/api/revenue/:quarter/:year", ctrl.GetQuarterlyRevenue)
	// router.GET("/api/revenue/:month/:year", ctrl.GetRevenue)

	// Define a new route for ProcessQuery

	router.Run(":8080")

	dbConfig.Close()
}
