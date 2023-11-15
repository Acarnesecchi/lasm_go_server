package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedCredentials Credentials
	result := DB.Where("username = ?", credentials.Username).First(&storedCredentials)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credentials incorrect"})
		return
	}

	if storedCredentials.Password != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credentials incorrect"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Protected route handler
func Protected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"content": "This view is protected. Only authorized clients can access it."})
}

func status() {

}
