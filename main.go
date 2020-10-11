package main

import (
	"api/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func requestHandle() {
	const port = 8092
	fmt.Printf("Server is up and running at port %d ", port)
	myRoutes := mux.NewRouter()
	myRoutes.HandleFunc("/api/v1/", routes.Home).Methods("GET")
	myRoutes.HandleFunc("/api/v1/generate/token", routes.TokenHandler).Methods("POST")
	myRoutes.HandleFunc("/api/v1/generate/UUID", routes.UUIDHandler).Methods("POST")

	myRoutes.HandleFunc("/", routes.Home).Methods("GET")

	err := http.ListenAndServe(":8092", myRoutes)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func main() {

	requestHandle()

}
