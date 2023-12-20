package main

import "gorm.io/gorm"

type ScooterUser struct {
	gorm.Model
	Username string `gorm:"primaryKey" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

type Scooter struct {
	Uuid                string  `gorm:"primaryKey" json:"uuid"`
	Name                string  `gorm:"not null" json:"name"`
	Longitude           float64 `gorm:"not null" json:"longitude"`
	Latitude            float64 `gorm:"not null" json:"latitude"`
	BatteryLevel        int     `gorm:"not null" json:"battery_level"`
	MetersUsed          int     `gorm:"not null" json:"meters_used"`
	DateCreated         string  `gorm:"not null" json:"date_created"`
	DateLastMaintenance string  `gorm:"not null" json:"date_last_maintenance"`
	State               string  `gorm:"not null" json:"state"`
	Vacant              bool    `gorm:"not null" json:"vacant"`
}

type Rent struct {
	Uuid      string  `gorm:"primaryKey" json:"uuid"`
	ScooterID string  `gorm:"not null" json:"-"`
	Scooter   Scooter `gorm:"foreignKey:ScooterID"`
	DateStart string  `gorm:"not null" json:"date_start"`
	DateStop  string  `gorm:"not null" json:"date_stop"`
}
