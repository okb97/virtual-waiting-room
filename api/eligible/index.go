package eligible

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/okb97/virtual-waiting-room/api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		HandleEligible(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func HandleEligible(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("ticketId")
	log.Printf("id: %s", id)
	exists, _ := api.RedisCommand([]interface{}{"SISMEMBER", "eligible", id})
	var ok int
	json.Unmarshal(exists, &ok)
	json.NewEncoder(w).Encode(map[string]bool{"canPurchase": ok == 1})
}
