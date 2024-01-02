build-main:
	@go build -ldflags "-s -w" -o bin/main cmd/main.go

build-wallet:
	@go build -ldflags "-s -w" -o bin/btcWallet cmd/btcWallet.go

run: build-main
	@./bin/main

wallet: build-wallet
	@./bin/btcWallet
