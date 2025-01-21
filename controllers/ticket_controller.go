package controllers

import (
	"context"
	"net/http"
	"ticket-system/database"
	"ticket-system/models"
	"ticket-system/utils" // Importamos utils para respuestas estándar
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ticketCollection *mongo.Collection

func InitTicketController() {
	ticketCollection = database.DB.Collection("tickets")
}

// @Summary Obtener todos los tickets
// @Description Devuelve una lista de todos los tickets almacenados en la base de datos
// @Tags Tickets
// @Accept json
// @Produce json
// @Success 200 {array} models.Ticket
// @Failure 500 {object} utils.StandardResponse
// @Router /tickets/ [get]
func GetTickets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tickets []models.Ticket
	cursor, err := ticketCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.StandardResponse{Message: "Error interno", Error: err.Error()})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ticket models.Ticket
		cursor.Decode(&ticket)
		tickets = append(tickets, ticket)
	}

	c.JSON(http.StatusOK, tickets)
}

// @Summary Crear un ticket
// @Description Crea un nuevo ticket en la base de datos
// @Tags Tickets
// @Accept json
// @Produce json
// @Param ticket body models.Ticket true "Datos del Ticket"
// @Success 201 {object} models.Ticket
// @Failure 400 {object} utils.StandardResponse
// @Failure 500 {object} utils.StandardResponse
// @Router /tickets/ [post]
func CreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, utils.StandardResponse{Message: "Datos inválidos", Error: err.Error()})
		return
	}

	ticket.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := ticketCollection.InsertOne(ctx, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.StandardResponse{Message: "Error al crear el ticket", Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// @Summary Obtener un ticket por ID
// @Description Devuelve un ticket específico según su ID
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path string true "ID del Ticket"
// @Success 200 {object} models.Ticket
// @Failure 400 {object} utils.StandardResponse
// @Failure 404 {object} utils.StandardResponse
// @Router /tickets/{id} [get]
func GetTicketByID(c *gin.Context) {
	ticketID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.StandardResponse{Message: "ID inválido", Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var ticket models.Ticket
	err = ticketCollection.FindOne(ctx, bson.M{"_id": ticketID}).Decode(&ticket)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.StandardResponse{Message: "Ticket no encontrado"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// @Summary Actualizar un ticket
// @Description Actualiza los datos de un ticket existente
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path string true "ID del Ticket"
// @Param ticket body models.Ticket true "Datos actualizados del Ticket"
// @Success 200 {object} utils.StandardResponse
// @Failure 400 {object} utils.StandardResponse
// @Failure 500 {object} utils.StandardResponse
// @Router /tickets/{id} [put]
func UpdateTicket(c *gin.Context) {
	ticketID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.StandardResponse{Message: "ID inválido", Error: err.Error()})
		return
	}

	var updatedTicket models.Ticket
	if err := c.ShouldBindJSON(&updatedTicket); err != nil {
		c.JSON(http.StatusBadRequest, utils.StandardResponse{Message: "Datos inválidos", Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": updatedTicket}
	_, err = ticketCollection.UpdateOne(ctx, bson.M{"_id": ticketID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.StandardResponse{Message: "Error al actualizar el ticket", Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.StandardResponse{Message: "Ticket actualizado correctamente"})
}

// @Summary Eliminar un ticket
// @Description Elimina un ticket de la base de datos
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path string true "ID del Ticket"
// @Success 200 {object} utils.StandardResponse
// @Failure 400 {object} utils.StandardResponse
// @Failure 500 {object} utils.StandardResponse
// @Router /tickets/{id} [delete]
func DeleteTicket(c *gin.Context) {
	ticketID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.StandardResponse{Message: "ID inválido", Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = ticketCollection.DeleteOne(ctx, bson.M{"_id": ticketID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.StandardResponse{Message: "Error al eliminar el ticket", Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.StandardResponse{Message: "Ticket eliminado correctamente"})
}
