.PHONY: build build-all clean test run dev docker-build docker-run docker-dev docs

BINARY_NAME=surfmate.io
DIST_DIR=dist
DOCKER_IMAGE=surfmate.io

build:
	go build -o $(BINARY_NAME) ./main.go

run:
	go run ./main.go

test:
	go test -v ./...

clean:
	rm -rf $(DIST_DIR) $(BINARY_NAME)

build-all: clean
	mkdir -p $(DIST_DIR)
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./main.go
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./main.go
	GOOS=linux GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./main.go

deps:
	go mod tidy

lint:
	golangci-lint run

# Docker commands
docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE) -http -port 8080

docker-dev:
	docker compose up dev

docker-prod:
	docker compose up -d prod

docs:
	docker compose up -d docs
	@echo "Docs available at http://localhost:3000"

docker-down:
	docker compose down
