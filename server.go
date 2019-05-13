package main

import (
	"fmt"
	routes "gemmails/app"
	"gemmails/app/models"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Global variable
var topics []string

func init() {
	if os.Getenv("MODE") != "production" {
		err := godotenv.Load()
		log.Println(err)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	models.InitDB()
}

func main() {
	// Routes
	routes := routes.NewRouter()

	// Cors domain
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000", "http://localhost:8080", "http://localhost:8081"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler(routes)

	// Run server
	port := os.Getenv("PORT")
	fmt.Println("Server run port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
