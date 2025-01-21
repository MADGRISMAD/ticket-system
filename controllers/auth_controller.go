package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"ticket-system/database"
	"ticket-system/models"
	"ticket-system/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func InitAuthController() {
	userCollection = database.DB.Collection("users")
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Estructura de credenciales para login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Registro de usuario
// @Description Registra un nuevo usuario con nombre, email y contraseña
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "Datos del usuario"
// @Success 201 {object} utils.StandardResponse
// @Failure 400 {object} utils.StandardResponse
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("🔴 Error en ShouldBindJSON:", err)
		c.JSON(http.StatusBadRequest, utils.StandardResponse{Message: "Datos inválidos", Error: err.Error()})
		return
	}

	// 📌 Verificar si los datos se están recibiendo correctamente
	log.Println("🔵 Datos recibidos:", user)

	// Normalizar email en minúsculas
	user.Email = strings.ToLower(user.Email)

	// 📌 Verificar si la contraseña está vacía
	if user.Password == "" {
		log.Println("🔴 Error: La contraseña no fue enviada correctamente en la solicitud")
		c.JSON(http.StatusBadRequest, utils.StandardResponse{Message: "La contraseña es obligatoria"})
		return
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.StandardResponse{Message: "Error al registrar usuario", Error: err.Error()})
		return
	}
	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.StandardResponse{Message: "Error al registrar usuario", Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, utils.StandardResponse{Message: "Usuario registrado exitosamente"})
}

// @Summary Inicio de sesión
// @Description Permite a un usuario iniciar sesión con email y contraseña
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credenciales de usuario"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.StandardResponse
// @Failure 401 {object} utils.StandardResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, utils.StandardResponse{Message: "Datos inválidos", Error: err.Error()})
		return
	}

	// Normalizar email en minúsculas
	loginReq.Email = strings.ToLower(loginReq.Email)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": loginReq.Email}).Decode(&user)
	if err != nil {
		log.Println("🔴 Usuario no encontrado en la base de datos:", loginReq.Email)
		c.JSON(http.StatusUnauthorized, utils.StandardResponse{Message: "Credenciales incorrectas"})
		return
	}

	// Imprimir hash y contraseña ingresada
	log.Println("🔵 Hash en BD:", user.Password)
	log.Println("🔵 Contraseña ingresada:", loginReq.Password)

	// Comparar contraseñas con bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		log.Println("🔴 Error al comparar contraseñas:", err)
		c.JSON(http.StatusUnauthorized, utils.StandardResponse{Message: "Credenciales incorrectas"})
		return
	}

	// Generar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.StandardResponse{Message: "Error al generar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
