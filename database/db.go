package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Cambia esto si usas una DB en la nube
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal("❌ Error al crear el cliente de MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("❌ Error al conectar a MongoDB:", err)
	}

	DB = client.Database("ticket_system") // Nombre de la base de datos

	log.Println("✅ Conectado a MongoDB correctamente.")
}
