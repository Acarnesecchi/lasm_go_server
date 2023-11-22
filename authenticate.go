package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// this endpoint serves as a middleware to verify tokens. Should not return anything
func authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		// endpoint will crash if the request does not use Authorization field :)
		if len(authHeader) <= len(BearerSchema) || !strings.HasPrefix(authHeader, BearerSchema) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":     "Invalid or missing authorization header",
				"code":      "401",
				"timestamp": time.Now(),
			})
			c.Abort()
			return
		}

		// If it does not crash, verify the token
		tokenString := authHeader[len(BearerSchema):]

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":     "unauthorized",
				"code":      "401",
				"timestamp": time.Now(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
