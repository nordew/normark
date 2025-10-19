package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/user/normark/internal/app"
)

// @title           Normark Trading Journal API
// @version         1.0
// @description     A comprehensive trading journal API for tracking and analyzing trading performance
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@normark.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
}

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
