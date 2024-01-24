package main

import (
	"log"
	"os"
	"time"

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
	app, ok := os.LookupEnv("APP")
	if !ok {
		log.Fatalf("error undefined app")
	}

	switch app {
	case "api-server-logger":
		startAPIServer()
	case "messaging-server":
		startMessagingServer()
	default:
		log.Fatalf("error unknown app")
	}
}
