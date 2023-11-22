package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	jwtKey          = []byte("llave-super-secreta")
	serverStartTime time.Time
)

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
	createDemoUsersIfEmpty()

	var rents []Rent
	result := DB.Find(&rents)
	if result.Error != nil {
		log.Fatal("Error while fetching the database:", result.Error)
	}
	jsonData, err := json.Marshal(rents)
	if err != nil {
		log.Fatal("Failed to marshal rents data to JSON")
	}

	// Print the JSON or write it to a file
	fmt.Println(string(jsonData))

	//router := gin.Default()
	//router.POST("/login", Login)
	//router.GET("/status", status)
	//authorized := router.Group("/")
	//authorized.Use(authenticate())
	//{
	//	authorized.GET("/protected", validate)
	//}
	//if err := router.Run(":8080"); err != nil {
	//	log.Fatal("Error while running the server:", err)
	//}

}

func createDemoUsersIfEmpty() {
	var count int64
	DB.Model(&Credentials{}).Count(&count)
	if count == 0 {
		demoUsers := []Credentials{
			{Username: "alex", Password: "alex1234"},
			{Username: "tomas", Password: "firebaseSSO"},
		}
		for _, user := range demoUsers {
			DB.Create(&user)
		}
	}
	demoScooter := Scooter{
		uuid:                "3",
		Name:                "Scooter 1",
		Longitude:           0,
		Latitude:            0,
		BatteryLevel:        100,
		MetersUsed:          0,
		DateCreated:         "2021-02-01",
		DateLastMaintenance: "2021-11-01",
		State:               "vacant",
		Vacant:              true,
	}
	DB.Create(&demoScooter)
	demoRent := Rent{
		uuid:        "1",
		ScooterUUID: "3",
		DateStart:   "2022-01-01",
		DateStop:    "2023-01-01",
	}
	DB.Create(&demoRent)
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
