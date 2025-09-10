package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

var (
	UPSTASH_REST_API_ENDPOINT = os.Getenv("UPSTASH_REDIS_REST_URL")
	UPSTASH_REST_API_TOKEN    = os.Getenv("UPSTASH_REDIS_REST_TOKEN")
)

type JoinResponse struct {
	TicketID string `json:"ticketId"`
}
type StatusResponse struct {
	WaitTime int64 `json:"waitTime"`
	Position int64 `json:"position"`
}
type UpstashResponse struct {
	Result interface{} `json:"result"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		HandleJoin(w, r)
		return
	}
	if r.Method == http.MethodGet {
		HandleStatus(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	ticketId := uuid.New().String()
	// RedisRPush を PushToQueue に変更
	if _, err := PushToQueue("queue", ticketId); err != nil {
		log.Println("PushToQueue error:", err)
		http.Error(w, "failed to join queue", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(JoinResponse{TicketID: ticketId})
}

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	ticketId := r.URL.Query().Get("ticketId")
	if ticketId == "" {
		http.Error(w, "ticketId required", http.StatusBadRequest)
		return
	}
	log.Printf("Received ticketId: %s", ticketId)

	queueLength, err := GetQueuePosition("queue", ticketId)
	if err != nil {
		log.Println("GetQueuePosition error:", err)
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

func RedisCommand(command []interface{}) ([]byte, error) {
	if UPSTASH_REST_API_ENDPOINT == "" || UPSTASH_REST_API_TOKEN == "" {
		return nil, fmt.Errorf("Upstash APIの認証情報が設定されていません。")
	}

	payload, err := json.Marshal(command)
	if err != nil {
		return nil, fmt.Errorf("JSONマーシャリングエラー: %v", err)
	}

	req, err := http.NewRequest("POST", UPSTASH_REST_API_ENDPOINT, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+UPSTASH_REST_API_TOKEN)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストエラー: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み取りエラー: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	if bytes.HasPrefix(body, []byte("[")) {
		var result []interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("配列としてのJSONパースエラー: %v (受信: %s)", err, string(body))
		}
		if len(result) > 0 {
			return json.Marshal(result[0])
		}
	} else if bytes.HasPrefix(body, []byte("{")) {
		var upstashRes UpstashResponse
		if err := json.Unmarshal(body, &upstashRes); err != nil {
			return nil, fmt.Errorf("オブジェクトとしてのJSONパースエラー: %v (受信: %s)", err, string(body))
		}

		if upstashRes.Result != nil {
			if errString, ok := upstashRes.Result.(string); ok && errString == "ERR" {
				return nil, fmt.Errorf("Upstash APIエラー: %s", errString)
			}
			return json.Marshal(upstashRes.Result)
		}
	} else {
		return nil, fmt.Errorf("未知のレスポンス形式です: %s", string(body))
	}

	return nil, nil
}

func PushToQueue(queueName string, items ...string) (int, error) {
	command := []interface{}{"RPUSH", queueName}
	for _, item := range items {
		command = append(command, item)
	}

	res, err := RedisCommand(command)
	if err != nil {
		return 0, err
	}

	var result int
	if err := json.Unmarshal(res, &result); err != nil {
		return 0, fmt.Errorf("Push結果のJSONパースエラー: %v (受信: %s)", err, string(res))
	}
	return result, nil
}

func PopFromQueue(queueName string) (string, error) {
	res, err := RedisCommand([]interface{}{"LPOP", queueName})
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

func GetQueuePosition(queueName, ticketId string) (int, error) {
	res, err := RedisCommand([]interface{}{"LPOS", queueName, ticketId})
	if err != nil {
		return 0, err
	}
	var pos int
	if err := json.Unmarshal(res, &pos); err != nil {
		return 0, fmt.Errorf("長さ結果のJSONパースエラー: %v", err)
	}
	return pos, nil
}

func RemoveTicket(queueName, ticketId string) (int, error) {
	res, err := RedisCommand([]interface{}{"LREM", queueName, 0, ticketId})
	if err != nil {
		return 0, err
	}
	var result int
	if err := json.Unmarshal(res, &result); err != nil {
		return 0, fmt.Errorf("LREMのJSONパースエラー: %v", err)
	}
	return result, nil
}
