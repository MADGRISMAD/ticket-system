package main

import (
	"log"
	"ticket-system/controllers"
	"ticket-system/database"
	"ticket-system/routes"
)

func main() {
	// Conectar a la base de datos
	database.ConnectDatabase()

	// Inicializar controladores después de conectar la BD
	controllers.InitAuthController()
	controllers.InitTicketController()

	// Iniciar servidor
	r := routes.SetupRouter()
	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
