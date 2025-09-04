package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var redisURL string
var redisToken string

func InitRedis() {
	log.Println("環境変数を取得")
	redisURL = os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken = os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		log.Fatal("UPSTASH_REDIS_REST_URL or UPSTASH_REDIS_REST_TOKEN is not set")
	}
	log.Println("Upstash Redis REST config loaded")
}

func redisCmd(cmd string, args ...interface{}) (map[string]interface{}, error) {
	bodyMap := map[string]interface{}{
		"cmd":  cmd,
		"args": args,
	}
	bodyBytes, _ := json.Marshal(bodyMap)

	req, _ := http.NewRequest("POST", redisURL, bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+redisToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
