package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		tokenString := c.GetHeader("Authorization")[len(BearerSchema):]
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "code": "401", "timestamp": time.Now()})
			c.Abort()
			return
		}
		c.Next()
	}
}
