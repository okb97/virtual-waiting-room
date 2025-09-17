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

func AllowNextBatch() error {
	const BATCHSIZE = 10
	res, err := api.RedisCommand([]interface{}{"LRANGE", "queue", 0, BATCHSIZE - 1})
	if err != nil {
		return err
	}
	var ids []string
	json.Unmarshal(res, &ids)

	if len(ids) == 0 {
		return nil
	}

	cmd := []interface{}{"SADD", "eligible"}
	for _, id := range ids {
		cmd = append(cmd, id)
	}
	if _, err := api.RedisCommand(cmd); err != nil {
		return err
	}

	rem := []interface{}{"LTRIM", "queue", len(ids), -1}
	_, err = api.RedisCommand(rem)
	return err
}
