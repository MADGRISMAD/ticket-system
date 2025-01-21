package main

import (
	"log"
	"ticket-system/controllers"
	"ticket-system/database"
	"ticket-system/routes"

	_ "ticket-system/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Ticket System API
// @version 1.0
// @description API para gestionar tickets usando Gin y MongoDB
// @host localhost:8080
// @BasePath /
func main() {
	database.ConnectDatabase()
	controllers.InitTicketController()

	r := routes.SetupRouter()

	// Documentación Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
