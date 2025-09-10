package integration_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/okb97/virtual-waiting-room/api"
)

func TestMain(m *testing.M) {
	os.Setenv("UPSTASH_REDIS_REST_URL", "あなたのRESTエンドポイントURL")
	os.Setenv("UPSTASH_REDIS_REST_TOKEN", "あなたのRESTトークン")
	os.Setenv("GO_ENV", "test")

	code := m.Run()
	os.Exit(code)
}

func clearQueue(t *testing.T, queueName string) {
	_, err := api.RedisCommand([]interface{}{"DEL", queueName})
	if err != nil {
		t.Fatalf("failed to clear queue:%v", err)
	}
}

func TestPushToQueueIntegration(t *testing.T) {
	queueName := "test_queue"

	clearQueue(t, queueName)
	newLen, err := api.PushToQueue(queueName, "item1")
	if err != nil {
		t.Fatalf("PushToQueue error: %v", err)
	}
	if newLen != 1 {
		t.Errorf("expected length=1, got %d", newLen)
	}

	res, err := api.RedisCommand([]interface{}{"LRANGE", queueName, 0, -1})
	if err != nil {
		t.Fatalf("LRANGE error: %v", err)
	}
	var items []string
	if err := json.Unmarshal(res, &items); err != nil {
		t.Fatalf("JSON parse error: %v", err)
	}

	expected := []string{"item1"}
	for i, v := range expected {
		if items[i] != v {
			t.Errorf("expected %s at position %d, got %s", v, i, items[i])
		}
	}
}

func TestPopFromQueueIntegration(t *testing.T) {
	queueName := "integration_test_queue"
	ticket1 := "ticket_1"

	clearQueue(t, queueName)

	_, err := api.RedisCommand([]interface{}{"RPUSH", queueName, ticket1})
	if err != nil {
		t.Fatalf("RPUSHエラー: %v", err)
	}

	val, err := api.PopFromQueue(queueName)
	if err != nil {
		t.Fatalf("PopFromQueueエラー: %v", err)
	}

	expected := ticket1
	if val != expected {
		t.Errorf("Pop結果が不正 got=%s want=%s", val, expected)
	}
}

func TestGetQueuePositionIntegration(t *testing.T) {
	queueName := "integration_test_queue"
	ticket1 := "ticket_1"
	ticket2 := "ticket_2"

	clearQueue(t, queueName)

	if _, err := api.PushToQueue(queueName, ticket1, ticket2); err != nil {
		t.Fatalf("PushToQueue failed: %v", err)
	}

	pos, err := api.GetQueuePosition(queueName, ticket1)
	if err != nil {
		t.Fatalf("GetQueuePosition failed: %v", err)
	}
	if pos != 0 {
		t.Errorf("expected position 0, got %d", pos)
	}
	pos2, err := api.GetQueuePosition(queueName, ticket2)
	if err != nil {
		t.Fatalf("GetQueuePosition failed: %v", err)
	}
	if pos2 != 1 {
		t.Errorf("expected position 1, got %d", pos2)
	}
}

func TestRemoveTicket(t *testing.T) {
	queueName := "integration_test_queue"
	ticket1 := "ticket_1"
	ticket2 := "ticket_2"
	ticket3 := "ticket_3"

	clearQueue(t, queueName)

	if _, err := api.PushToQueue(queueName, ticket1, ticket2, ticket3); err != nil {
		t.Fatalf("PushToQueue failed: %v", err)
	}

	res, err := api.RemoveTicket(queueName, ticket1)
	if err != nil {
		t.Fatalf("RemoveTicket failed: %v", err)
	}
	if res != 1 {
		t.Errorf("expected removed count 1, got %d", res)
	}
}
