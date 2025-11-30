package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "test-topic"
	groupID := "test-consumer-group"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,    // 1 byte (즉시 읽기)
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	fmt.Println("🎧 Consumer 시작... (Ctrl+C로 종료)")
	fmt.Println("----------------------------------------")

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\n👋 Consumer 종료 중...")
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := r.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("메시지 읽기 실패: %v", err)
				continue
			}

			fmt.Printf("\n📨 새 메시지 수신!\n")
			fmt.Printf("   Topic: %s\n", m.Topic)
			fmt.Printf("   Partition: %d\n", m.Partition)
			fmt.Printf("   Offset: %d\n", m.Offset)
			fmt.Printf("   Key: %s\n", string(m.Key))
			fmt.Printf("   Value: %s\n", string(m.Value))
			fmt.Printf("   Time: %s\n", m.Time.Format("2006-01-02 15:04:05"))
			fmt.Println("----------------------------------------")

			if err := r.CommitMessages(ctx, m); err != nil {
				log.Printf("메시지 커밋 실패: %v", err)
			}
		}
	}
}
