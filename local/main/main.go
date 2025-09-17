package main

import (
	"log"
	"net/http"

	"github.com/okb97/virtual-waiting-room/api"
	"github.com/okb97/virtual-waiting-room/api/cron"
	"github.com/okb97/virtual-waiting-room/api/eligible"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/queue", api.Handler)
	mux.HandleFunc("/api/eligible", eligible.Handler)

	if err := cron.AllowNextBatch(); err != nil {
		log.Println("AllowNextBatch error:", err)
	} else {
		log.Println("AllowNextBatch executed successfully")
	}

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
