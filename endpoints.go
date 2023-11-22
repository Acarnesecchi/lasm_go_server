package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "Server started at": serverStartTime})
}

func rentHistory(c *gin.Context) {
	// show rent history
}

func scooterList(c *gin.Context) {
	var scooters []Scooter
	result := DB.Find(&scooters)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
	}
	c.JSON(http.StatusOK, gin.H{"scooters": scooters})
}

func scooter(c *gin.Context) {
	sc := c.Param("uuid")
	var scooter Scooter
	result := DB.First(&scooter, "uuid = ?", sc)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
	}
	c.JSON(http.StatusOK, scooter)
}

func startRent(c *gin.Context) {
	// rent a scooter action
}

func stopRent(c *gin.Context) {
	// stop renting a scooter action
}
