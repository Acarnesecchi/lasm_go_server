package main

import (
	"errors"
	"log"
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
	status := "success"
	result := DB.Find(&scooters)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		status = "failed"
	} else {
		c.JSON(http.StatusOK, gin.H{"scooters": scooters})
	}
	l := buildLog(c, status)
	err := SendLog(l)
	if err != nil {
		log.Fatal("Error while sending log:", err)
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
	status := "success"
	result := DB.First(&sc, "uuid = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		status = "failed"
	}
	if sc.Vacant {
		tx := DB.Begin() // starts a transaction

		randomUuid, err := uuid.NewRandom()
		if err != nil {
			tx.Rollback() // probably unnecessary, no DB changed have been made
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong creating the UUID"})
			status = "failed"
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
			status = "failed"
		}
		sc.Vacant = false
		err = tx.Save(&sc).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update scooter status"})
			status = "failed"
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
		status = "failure"
	}
	l := buildLog(c, status)
	err := SendLog(l)
	if err != nil {
		log.Fatal("Error while sending log:", err)
	}
}

func stopRent(c *gin.Context) {
	id := c.Param("uuid")
	var sc Scooter
	var r Rent
	status := "success"
	result := DB.First(&sc, "uuid = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scooter not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		status = "failed"
	}

	if sc.Vacant {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":      405,
			"msg":       "Scooter is not rented",
			"rent":      gin.H{},
			"timestamp": time.DateTime,
			"version":   version,
		})
		status = "failed"
	} else {
		tx := DB.Begin()

		result = tx.First(&r, "scooter_id = ?", id)
		if result.Error != nil {
			tx.Rollback()
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Rent record not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			}
			status = "failed"
		}

		r.DateStop = time.Now().Format(time.RFC3339)
		if result = tx.Save(&r); result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update rent record"})
			status = "failed"
		}

		sc.Vacant = true
		if result = tx.Save(&sc); result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update scooter record"})
			status = "failed"
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

	l := buildLog(c, status)
	err := SendLog(l)
	if err != nil {
		log.Fatal("Error while sending log:", err)
	}
}

func buildLog(c *gin.Context, s string) Log {
	l := Log{
		Endpoint:  c.Request.URL.Path,
		Ip:        c.ClientIP(),
		Timestamp: time.Now().Format(time.RFC3339),
		Client:    c.Request.UserAgent(),
		Status:    s,
	}
	return l
}
