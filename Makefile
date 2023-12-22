build:
	@go build -ldflags "-s -w" -o bin/main cmd/main.go

run: build
	@./bin/main
