package main

import (
	"log"
	"net/http"

	"github.com/okb97/virtual-waiting-room/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/queue", api.Handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
