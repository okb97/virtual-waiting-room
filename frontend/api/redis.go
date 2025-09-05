package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	UPSTASH_REST_API_ENDPOINT = os.Getenv("UPSTASH_REDIS_REST_URL")
	UPSTASH_REST_API_TOKEN    = os.Getenv("UPSTASH_REDIS_REST_TOKEN")
)

type UpstashResponse struct {
	Result interface{} `json:"result"`
}

func redisCommand(command []interface{}) ([]byte, error) {
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

	res, err := redisCommand(command)
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
	res, err := redisCommand([]interface{}{"LPOP", queueName})
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

func GetQueueLength(queueName string) (int, error) {
	res, err := redisCommand([]interface{}{"LLEN", queueName})
	if err != nil {
		return 0, err
	}
	var length int
	if err := json.Unmarshal(res, &length); err != nil {
		return 0, fmt.Errorf("長さ結果のJSONパースエラー: %v", err)
	}
	return length, nil
}
