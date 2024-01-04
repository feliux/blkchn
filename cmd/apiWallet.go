package main

import (
	"flag"

	ws "github.com/feliux/blkchn/walletserver"
)

var (
	port    int
	gateway string
)

func init() {
	flag.IntVar(&port, "port", 8080, "TCP Port Number for Wallet Server.")
	flag.StringVar(&gateway, "gateway", "http://127.0.0.1:5000", "Blockchain Gateway endpoint to connect.")
	flag.Parse()
}

func main() {
	app := ws.NewWalletServer(port, gateway)
	app.Run()
}
