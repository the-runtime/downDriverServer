package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	handler "serverFordownDrive/handlers"
)

func init() {
	if err := godotenv.Load(); err != nil {
		println("No .env file found")
	}
}

func main() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":8000"),
		Handler: handler.New(),
	}
	log.Printf("Starting HTTP server. Listening at #{server.Addr")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("#{err")
	} else {
		log.Println("Server cloased!")
	}
}
