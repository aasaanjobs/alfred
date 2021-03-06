package main

import (
	"fmt"
	"log"
	"net/http"

	c "github.com/aasaanjobs/alfred/controllers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// Add Index URL
	router.Handle("/build", c.ViewHandler(c.WebHook)).Methods("POST")
	router.Handle("/", c.ViewHandler(c.Ping)).Methods("GET", "POST")
	serverURI := fmt.Sprintf("0.0.0.0:8000")
	log.Println("Starting Alfred server at 0.0.0.0:8000")
	err := http.ListenAndServe(serverURI, router)
	if err != nil {
		panic(err)
	}
}
