package checkin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/okb97/virtual-waiting-room/api"
)

type JoinResponse struct {
	TicketID string `json:"ticketId"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodDelete {
		HandleCheckIn(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	ticketId, err := PopFromQueue("queue")
	if err != nil {
		log.Println("PopFromQueue error:", err)
		http.Error(w, "failed to get ticket from queue", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(JoinResponse{TicketID: ticketId})
}

func PopFromQueue(queueName string) (string, error) {
	res, err := api.RedisCommand([]interface{}{"LPOP", queueName})
	if err != nil {
		return "", err
	}
	var popResult string
	if err := json.Unmarshal(res, &popResult); err != nil {
		if string(res) == "null" {
			return "", nil
		}
		return "", fmt.Errorf("Pop結果のJSONパースエラー: %v", err)
	}
	return popResult, nil
}
