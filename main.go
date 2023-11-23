package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	jwtKey          = []byte("llave-super-secreta")
	serverStartTime time.Time
)

const version = "1.0"

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {

	go startWebServer()

	serverStartTime = time.Now()
	DBConnection()

	err := DB.AutoMigrate(&Credentials{}, &Scooter{}, &Rent{})
	if err != nil {
		e := DB.Migrator().DropTable(&Credentials{}, &Scooter{}, &Rent{})
		if e != nil {
			log.Fatal("Error while dropping the database:", e)
		}
		log.Fatal("Error while migrating the database:", err)
	}

	router := gin.Default()
	// everything under router is accesible in the root `localhost/`
	endpointsGroup := router.Group("/endpoints")
	// everything under endpoindsGroup is accesible in `localhost/endpoints`
	endpointsGroup.POST("/login", Login)
	endpointsGroup.GET("/status", status)
	authorized := endpointsGroup.Group("/")
	authorized.Use(authenticate())
	{
		authorized.GET("/validate", validate)

		// Define the scooterGroup under authorized to ensure it requires authentication
		scooterGroup := authorized.Group("/scooter")
		{
			scooterGroup.GET("/", scooterList)
			scooterGroup.GET("/:uuid", scooter)
		}
		rentGroup := authorized.Group("/rent")
		{
			rentGroup.GET("/", rentHistory)
			rentGroup.POST("/start/:uuid", startRent)
		}
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error while running the server:", err)
	}

}

func startWebServer() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/html", h1)
	fmt.Println("Starting WebServer on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
