OUT_DIR = ./bin

.PHONY: run dev prod migrate

run: prod

prod: 
	go build -o ${OUT_DIR}/server ./cmd/server/main.go
	go build -o ${OUT_DIR}/client ./cmd/client/main.go
	${OUT_DIR}/server&
	${OUT_DIR}/client

dev:
	go run ./cmd/server/main.go&
	go run ./cmd/client/main.go

migrate: 
	go run ./cmd/migrations/main.go