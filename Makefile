.PHONY: run build up down client-strong client-eventual

build:
	go build -o bin/node ./cmd/node
	go build -o bin/client ./cmd/client

run-nodeA:
	NODE_NAME=nodeA ADDR=:8081 PEER_URL=http://localhost:8082 go run ./cmd/node

run-nodeB:
	NODE_NAME=nodeB ADDR=:8082 PEER_URL=http://localhost:8081 go run ./cmd/node

up:
	docker compose up --build

down:
	docker compose down -v

client-eventual:
	go run -race ./cmd/client -a http://localhost:8081 -b http://localhost:8082 -mode eventual -rps 80 -dur 15s -readlag 300ms -seed 42

client-strong:
	go run ./cmd/client -a http://localhost:8081 -b http://localhost:8082 -mode strong -rps 80 -dur 15s -readlag 300ms -seed 42
