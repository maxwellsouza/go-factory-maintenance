run:
	go run ./cmd/api

test:
	go test ./... -cover

docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down
