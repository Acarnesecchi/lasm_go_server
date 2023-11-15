package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

var (
	jwtKey          = []byte("llave-super-secreta")
	serverStartTime time.Time
)

type Credentials struct {
	gorm.Model
	Username string `gorm:"primaryKey" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	serverStartTime = time.Now()
	DBConnection()

	DB.AutoMigrate(&Credentials{})
	createDemoUsersIfEmpty()

	router := gin.Default()
	router.POST("/login", Login)
	router.GET("/status", status)
	authorized := router.Group("/")
	authorized.Use(authenticate())
	{
		authorized.GET("/protected", validate)
	}
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error while running the server:", err)
	}
}

func createDemoUsersIfEmpty() {
	var count int64
	DB.Model(&Credentials{}).Count(&count)
	if count == 0 {
		demoUsers := []Credentials{
			{Username: "alex", Password: "alex1234"},
			{Username: "tomas", Password: "firebaseSSO"},
		}
		for _, user := range demoUsers {
			DB.Create(&user)
		}
	}
}
