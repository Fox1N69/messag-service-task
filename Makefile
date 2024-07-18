dep:
	go mod tidy

run-algosync:
	go run cmd/algosync-service/main.go

test:
	go test -short -cover ./...

build-algosync:
	go build -o bin/server cmd/algosync-service/main.go

docker-image:
	docker build -t server:v1 .

docker-build:
	docker-compose up --build

docker-run:
	docker-compose up
