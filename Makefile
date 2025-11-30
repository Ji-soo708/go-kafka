.PHONY: help start stop restart logs ui producer consumer clean topic-list topic-create topic-delete

help: ## 도움말 표시
	@echo "사용 가능한 명령어:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

start: ## Kafka 시작 (KRaft 모드)
	docker-compose up -d
	@echo "✅ Kafka가 시작되었습니다!"
	@echo "📊 Kafka UI: http://localhost:8080"

stop: ## Kafka 중지
	docker-compose down
	@echo "🛑 Kafka가 중지되었습니다!"

restart: stop start ## Kafka 재시작

logs: ## Kafka 로그 보기
	docker-compose logs -f kafka

ui: ## 브라우저에서 Kafka UI 열기
	open http://localhost:8080 || xdg-open http://localhost:8080 || echo "브라우저에서 http://localhost:8080 열기"

producer: ## Producer 실행
	go run producer/main.go

consumer: ## Consumer 실행
	go run consumer/main.go

clean: ## 모든 컨테이너 및 볼륨 삭제
	docker-compose down -v
	@echo "🧹 모든 Kafka 데이터가 삭제되었습니다!"

topic-list: ## 토픽 목록 확인
	docker exec -it kafka kafka-topics.sh --list --bootstrap-server localhost:9092

topic-create: ## 테스트 토픽 생성
	docker exec -it kafka kafka-topics.sh --create \
		--topic test-topic \
		--bootstrap-server localhost:9092 \
		--partitions 3 \
		--replication-factor 1
	@echo "✅ test-topic이 생성되었습니다!"

topic-delete: ## 테스트 토픽 삭제
	docker exec -it kafka kafka-topics.sh --delete \
		--topic test-topic \
		--bootstrap-server localhost:9092
	@echo "🗑️ test-topic이 삭제되었습니다!"

topic-describe: ## 테스트 토픽 상세 정보
	docker exec -it kafka kafka-topics.sh --describe \
		--topic test-topic \
		--bootstrap-server localhost:9092

console-producer: ## Kafka 콘솔 producer 실행
	docker exec -it kafka kafka-console-producer.sh \
		--topic test-topic \
		--bootstrap-server localhost:9092

console-consumer: ## Kafka 콘솔 consumer 실행
	docker exec -it kafka kafka-console-consumer.sh \
		--topic test-topic \
		--bootstrap-server localhost:9092 \
		--from-beginning
