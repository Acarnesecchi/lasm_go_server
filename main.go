package main

import (
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
	go startWebServer()
	startAPIServer()
}
