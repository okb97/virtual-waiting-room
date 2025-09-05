package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Upstashの認証情報を設定
	os.Setenv("UPSTASH_REDIS_REST_URL", "あなたのRESTエンドポイントURL")
	os.Setenv("UPSTASH_REDIS_REST_TOKEN", "あなたのRESTトークン")
	os.Setenv("GO_ENV", "test")

	code := m.Run()
	os.Exit(code)
}

func TestJoinAndStatus(t *testing.T) {

	// 1. Joinのテスト
	reqJoin, err := http.NewRequest("POST", "/join", nil)
	if err != nil {
		t.Fatal(err)
	}
	rrJoin := httptest.NewRecorder()

	// handleJoinハンドラーを直接呼び出す
	handleJoin(rrJoin, reqJoin)

	if status := rrJoin.Code; status != http.StatusOK {
		t.Errorf("handleJoin returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var joinResponse JoinResponse
	if err := json.NewDecoder(rrJoin.Body).Decode(&joinResponse); err != nil {
		t.Fatal(err)
	}

	if joinResponse.TicketID == "" {
		t.Errorf("handleJoin returned empty TicketID")
	}
	log.Printf("ticketId : %s", joinResponse.TicketID)

	// 2. Statusのテスト
	reqStatus, err := http.NewRequest("GET", "/status?ticketId="+joinResponse.TicketID, nil)
	if err != nil {
		t.Fatal(err)
	}
	rrStatus := httptest.NewRecorder()

	// handleStatusハンドラーを直接呼び出す
	handleStatus(rrStatus, reqStatus)

	if status := rrStatus.Code; status != http.StatusOK {
		t.Errorf("handleStatus returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var statusResponse StatusResponse
	if err := json.NewDecoder(rrStatus.Body).Decode(&statusResponse); err != nil {
		t.Fatal(err)
	}

	// // テスト環境ではPositionが常に1になるはず
	// if statusResponse.Position != 1 {
	// 	t.Errorf("handleStatus returned wrong position: got %d want %d",
	// 		statusResponse.Position, 1)
	// }

	// // テスト後のクリーンアップ
	// // 厳密にはLPOPを呼び出してテストで追加した要素を削除すべきですが、
	// // 簡易的なテストとして省略します。
}
