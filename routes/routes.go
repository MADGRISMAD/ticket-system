package routes

import (
	"ticket-system/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

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
