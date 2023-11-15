package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "Server started at": serverStartTime})
}

func rentHistory(c *gin.Context) {
	// show rent history
}

func scooterList(c *gin.Context) {
	// show scooter list
}

func scooter(c *gin.Context) {
	// show scooter info
}

func startRent(c *gin.Context) {
	// rent a scooter action
}

func stopRent(c *gin.Context) {
	// stop renting a scooter action
}
