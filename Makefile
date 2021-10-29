hello:
	echo "Hello"

build:
	go build -o bin/cli cmd/cli/main.go

run:
	go run cmd/cli/main.go