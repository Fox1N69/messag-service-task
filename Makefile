dep:
	go mod tidy

run-message-service:
	go run cmd/message-service/main.go

test:
	go test -short -cover ./...

build-message:
	go build -o bin/server cmd/message-service/main.go

run-client:
	go run cmd/client/main.go

build-client:
	go build -o bin/client cmd/client/main.go

docker-image:
	docker build -t server:v1 .

docker-build:
	docker-compose up --build

docker-run:
	docker-compose up
