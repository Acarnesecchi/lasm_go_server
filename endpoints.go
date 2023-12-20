package main

import (
	"errors"
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"rents": rents})
	}
}

func scooterList(c *gin.Context) {
	var scooters []Scooter
	result := DB.Find(&scooters)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"scooters": scooters})
	}
}

func scooter(c *gin.Context) {
	sc := c.Param("uuid")
	var scooter Scooter
	result := DB.First(&scooter, "uuid = ?", sc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
	} else {
		c.JSON(http.StatusOK, scooter)
	}
}

func startRent(c *gin.Context) {
	id := c.Param("uuid")
	var sc Scooter
	var r Rent
	result := DB.First(&sc, "uuid = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	if sc.Vacant {
		tx := DB.Begin() // starts a transaction

		randomUuid, err := uuid.NewRandom()
		if err != nil {
			tx.Rollback() // probably unnecessary, no DB changed have been made
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong creating the UUID"})
			return
		}
		r = Rent{
			Uuid:      randomUuid.String(),
			ScooterID: id,
			Scooter:   sc,
			DateStart: time.DateTime,
		}
		err = tx.Create(&r).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create rent record"})
			return
		}
		sc.Vacant = false
		err = tx.Save(&sc).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update scooter status"})
			return
		}

		tx.Commit() // If nothing fails, we commit the transaction

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "OK",
			"rent": gin.H{
				"uuid":       r.Uuid,
				"date_start": r.DateStart,
			},
			"timestamp": time.DateTime,
			"version":   version,
		})
	} else {
		// scooter not vacant
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":      405,
			"msg":       "Scooter is not vacant",
			"rent":      gin.H{},
			"timestamp": time.DateTime,
			"version":   version,
		})
	}
}

func stopRent(c *gin.Context) {
	id := c.Param("uuid")
	var sc Scooter
	var r Rent

	result := DB.First(&sc, "uuid = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	if sc.Vacant {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":      405,
			"msg":       "Scooter is not rented",
			"rent":      gin.H{},
			"timestamp": time.DateTime,
			"version":   version,
		})
		return
	}

	tx := DB.Begin()

	result = tx.First(&r, "scooter_id = ?", id)
	if result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rent record not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}

	r.DateStop = time.Now().Format(time.RFC3339)
	if result = tx.Save(&r); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update rent record"})
		return
	}

	sc.Vacant = true
	if result = tx.Save(&sc); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update scooter record"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
		"rent": gin.H{
			"uuid":       r.Uuid,
			"date_start": r.DateStart,
		},
		"timestamp": time.DateTime,
		"version":   version,
	})
}
