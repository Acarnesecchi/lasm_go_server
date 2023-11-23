package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "Server started at": serverStartTime})
}

func rentHistory(c *gin.Context) {
	var rents []Rent
	result := DB.Preload("Scooter").Find(&rents)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
	}
	c.JSON(http.StatusOK, gin.H{"rents": rents})
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
	id := c.Param("uuid")
	var sc Scooter
	var r Rent
	result := DB.First(&sc, "uuid = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
	}
	if sc.Vacant {
		random_uuid, err := uuid.NewRandom()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong creating the UUID"})
			return
		}
		r = Rent{
			Uuid:      random_uuid.String(),
			ScooterID: id,
			Scooter:   sc,
			DateStart: time.Now().String(),
		}
		DB.Create(&r)
		sc.Vacant = false
		updateResult := DB.Save(&sc)
		if updateResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating scooter status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "OK",
			"rent": gin.H{
				"uuid":       r.Uuid,
				"date_start": r.DateStart,
			},
			"timestamp": time.Now(),
			"version":   version})
	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":      405,
			"msg":       "Scooter is not vacant",
			"rent":      gin.H{},
			"timestamp": time.Now(),
			"version":   version})
	}
}

func stopRent(c *gin.Context) {
	// stop renting a scooter action
}
