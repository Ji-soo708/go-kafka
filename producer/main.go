package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "test-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	messages := []kafka.Message{
		{
			Key:   []byte("key-1"),
			Value: []byte("안녕하세요, 첫 번째 Kafka 메시지입니다!"),
		},
		{
			Key:   []byte("key-2"),
			Value: []byte("두 번째 메시지입니다. 현재 시간: " + time.Now().Format(time.RFC3339)),
		},
		{
			Key:   []byte("key-3"),
			Value: []byte("세 번째 메시지입니다. KRaft 모드로 실행 중!"),
		},
	}

	for i, msg := range messages {
		_, err := conn.WriteMessages(msg)
		if err != nil {
			log.Printf("failed to write message %d: %v", i+1, err)
		} else {
			fmt.Printf("✅ 메시지 %d 전송 성공: %s\n", i+1, string(msg.Value))
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n모든 메시지 전송 완료!")
}
