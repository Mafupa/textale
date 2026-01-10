.PHONY: run build docker test clean

run:
	go run .

build:
	go build -o textale .

docker:
	docker compose up -d --build

docker-down:
	docker compose down

test:
	go test -v ./...

clean:
	rm -f textale textale.db
	rm -rf .ssh/

install:
	go mod download
	mkdir -p .ssh

dev:
	# Start Redis in background, then run the app
	docker run -d --name textale-redis -p 6379:6379 redis:7-alpine || true
	go run .
