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

// VercelのサーバーレスFunction用エントリポイント
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		handleJoin(w, r)
	}
	if r.Method == http.MethodGet {
		handleStatus(w, r)
	}
}

func handleJoin(w http.ResponseWriter, r *http.Request) {
	log.Println("handleJoin called")
	ticketId := uuid.New().String()

	_, err := redisCmd("RPUSH", "queue", ticketId)
	if err != nil {
		log.Printf("Redis REST error: %v\n", err)
		http.Error(w, "failed to join queue", http.StatusInternalServerError)
		return
	}

	log.Println("ticket issued:", ticketId)
	json.NewEncoder(w).Encode(JoinResponse{TicketID: ticketId})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	ticketId := r.URL.Query().Get("ticketId")
	if ticketId == "" {
		http.Error(w, "ticketId required", http.StatusBadRequest)
		return
	}
	res, err := redisCmd("LPOS", "queue", ticketId)
	if err != nil {
		log.Printf("Redis REST error: %v\n", err)
		http.Error(w, "redis error", http.StatusInternalServerError)
		return
	}

	// LPOSの結果は float64 になる
	posF, ok := res["result"].(float64)
	if !ok {
		http.Error(w, "not in queue", http.StatusNotFound)
		return
	}
	pos := int64(posF)

	// 仮で「1人=30秒待ち」
	waitTime := pos * 30 / 60

	json.NewEncoder(w).Encode(StatusResponse{
		Position: pos,
		WaitTime: waitTime,
	})
}
