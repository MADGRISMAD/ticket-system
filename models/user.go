package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User representa el modelo de usuario para la autenticación
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
}
