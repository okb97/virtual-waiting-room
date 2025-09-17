package cron

import (
	"encoding/json"
	"net/http"

	"github.com/okb97/virtual-waiting-room/api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		AllowNextBatch()
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
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
