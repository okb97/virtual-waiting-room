package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type JoinResponse struct {
	TicketID string `json:"ticketId"`
}

type StatusResponse struct {
	WaitTime int64 `json:"waitTime"`
	Position int64 `json:"position"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		handleJoin(w, r)
		return
	}
	if r.Method == http.MethodGet {
		handleStatus(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleJoin(w http.ResponseWriter, r *http.Request) {
	ticketId := uuid.New().String()
	// RedisRPush を PushToQueue に変更
	if _, err := PushToQueue("queue", ticketId); err != nil {
		log.Println("PushToQueue error:", err)
		http.Error(w, "failed to join queue", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(JoinResponse{TicketID: ticketId})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	ticketId := r.URL.Query().Get("ticketId")
	if ticketId == "" {
		http.Error(w, "ticketId required", http.StatusBadRequest)
		return
	}
	log.Printf("Received ticketId: %s", ticketId)

	queueLength, err := GetQueueLength("queue")
	if err != nil {
		log.Println("GetQueueLength error:", err)
		http.Error(w, "redis error", http.StatusInternalServerError)
		return
	}
	log.Printf("Current queue length: %d", queueLength)

	pos := int64(queueLength)
	log.Printf("Calculated position: %d", pos)
	waitTime := pos * 30 / 60

	json.NewEncoder(w).Encode(StatusResponse{
		Position: pos,
		WaitTime: waitTime,
	})
}
