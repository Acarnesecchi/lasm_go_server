package main

import (
	"gorm.io/gorm"
	"time"
)

type Actor struct {
	gorm.Model
	ID          uint       `gorm:"primaryKey" json:"id"`
	First_Name  string     `gorm:"not null" json:"first_name"`
	Second_Name string     `gorm:"not null" json:"second_name"`
	Age         int        `gorm:"not null" json:"age"`
	Gender      string     `gorm:"not null" json:"gender"`
	CreatedAt   time.Time  `gorm:"-"`
	UpdatedAt   time.Time  `gorm:"-"`
	DeletedAt   *time.Time `gorm:"-"`
}
