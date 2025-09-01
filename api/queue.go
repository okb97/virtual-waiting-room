package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type JoinResponse struct {
	TicketID string `json:"ticketId"`
}

// VercelのサーバーレスFunction用エントリポイント
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		ticketId := generateTicketID()

		resp := JoinResponse{TicketID: ticketId}
		json.NewEncoder(w).Encode(resp)
		return
	}
	// 未対応のエンドポイントは404
	http.NotFound(w, r)
}

func generateTicketID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
