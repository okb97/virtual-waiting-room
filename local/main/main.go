package main

import (
	"log"
	"net/http"

	"github.com/okb97/virtual-waiting-room/api"
	"github.com/okb97/virtual-waiting-room/api/checkin"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/queue", api.Handler)
	mux.HandleFunc("/api/checkin", checkin.Handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
