IMAGE_NAME = varnit/placementlog-server
IMAGE_TAG = latest
CONTAINER_NAME = placementlog-container

runServer:
	go run ./cmd/main.go

buildDocker:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

runDocker:
	docker run -p 8080:8080 --name $(CONTAINER_NAME) -v "$(shell pwd)/.env":/app/.env $(IMAGE_NAME):$(IMAGE_TAG)