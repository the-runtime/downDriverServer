package main

import (
	"fmt"
	"net/http"
	handler "serverFordownDrive/handlers"
)

func main() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":8000"),
		Handler: handler.New(),
	}
	fmt.Printf("Starting HTTP server. Listening at " + server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("%s", err)
	} else {
		fmt.Println("Server cloased!")
	}

}
