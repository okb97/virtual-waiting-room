package main

import (
	"log"
	"net/http"

	"github.com/okb97/virtual-waiting-room/frontend/api"
)

func main() {
	//_ = godotenv.Load("../../.env.local")
	api.InitRedis()

	http.HandleFunc("/api/queue", api.Handler) // パスをマッピング
	log.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
