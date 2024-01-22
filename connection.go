package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnection() {
	viper.SetConfigName("postgres")
	viper.SetConfigType("ini")
	viper.AddConfigPath("credentials")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	host := viper.GetString("postgresql.host")
	port := viper.GetInt("postgresql.port")
	user := viper.GetString("postgresql.user")
	password := viper.GetString("postgresql.password")
	dbName := viper.GetString("postgresql.database")

	connect := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	DB, err = gorm.Open(postgres.Open(connect), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
}
