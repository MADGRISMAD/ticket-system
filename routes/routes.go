package routes

import (
	"ticket-system/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Rutas de autenticación
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Rutas de tickets
	tickets := r.Group("/tickets")
	{
		tickets.GET("/", controllers.GetTickets)
		tickets.POST("/", controllers.CreateTicket)
		tickets.GET("/:id", controllers.GetTicketByID)
		tickets.PUT("/:id", controllers.UpdateTicket)
		tickets.DELETE("/:id", controllers.DeleteTicket)
	}

	return r
}
