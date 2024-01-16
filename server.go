package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func startWebServer() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("send.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("receive.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Handle("/firebase-messaging-sw.js", http.FileServer(http.Dir(".")))
	http.HandleFunc("/send-message", handleSendNotification)
	http.HandleFunc("/send", h1)
	http.HandleFunc("/receive", h2)
	fmt.Println("Starting WebServer on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func startAPIServer() {
	serverStartTime = time.Now()
	DBConnection()

	err := DB.AutoMigrate(&ScooterUser{}, &Scooter{}, &Rent{})
	if err != nil {
		e := DB.Migrator().DropTable(&ScooterUser{}, &Scooter{}, &Rent{})
		if e != nil {
			log.Fatal("Error while dropping the database:", e)
		}
		log.Fatal("Error while migrating the database:", err)
	}

	router := gin.Default()

	endpointsGroup := router.Group("/endpoints")
	endpointsGroup.POST("/login", Login)
	endpointsGroup.GET("/status", status)

	authorized := endpointsGroup.Group("/")
	authorized.Use(authenticate())
	{
		authorized.GET("/validate", validate)

		scooterGroup := authorized.Group("/scooter")
		{
			scooterGroup.GET("/", scooterList)
			scooterGroup.GET("/:uuid", scooter)
		}
		rentGroup := authorized.Group("/rent")
		{
			rentGroup.GET("/", rentHistory)
			rentGroup.POST("/start/:uuid", startRent)
			rentGroup.POST("/stop/:uuid", stopRent)
		}
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error while running the server:", err)
	}
}
