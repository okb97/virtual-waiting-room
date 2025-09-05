package main

import (
	"fmt"

	"github.com/okb97/virtual-waiting-room/frontend/api"
)

func main() {
	queueName := "my-go-queue"

	// 1. キューに要素を追加
	fmt.Println("--- キューに要素を追加 ---")
	pushRes, err := api.PushToQueue(queueName, "task:1")
	if err != nil {
		fmt.Printf("Pushエラー: %v\n", err)
		return
	}
	fmt.Println(pushRes)

	// 2. キューの長さを取得
	fmt.Println("\n--- キューの長さを取得 ---")
	length, err := api.GetQueueLength(queueName)
	if err != nil {
		fmt.Printf("長さ取得エラー: %v\n", err)
		return
	}
	fmt.Printf("現在のキューの長さ: %d\n", length)

	// 3. キューから要素をポップ
	fmt.Println("\n--- キューから要素をポップ ---")
	popRes, err := api.PopFromQueue(queueName)
	if err != nil {
		fmt.Printf("Popエラー: %v\n", err)
		return
	}
	if popRes == "" {
		fmt.Println("Pop成功: キューは空でした。")
	} else {
		fmt.Printf("Pop成功: %s\n", popRes)
	}

	// 4. 再度、キューの長さを取得
	fmt.Println("\n--- ポップ後のキューの長さを取得 ---")
	length, err = api.GetQueueLength(queueName)
	if err != nil {
		fmt.Printf("長さ取得エラー: %v\n", err)
		return
	}
	fmt.Printf("ポップ後のキューの長さ: %d\n", length)

}
