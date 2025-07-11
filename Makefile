runServer:
	go run ./cmd/main.go

buildDocker:
	docker build -t placementlog-server .

runDocker:
	docker run -p 8080:8080 -v "$(shell pwd)/.env":/app/.env placementlog-server
