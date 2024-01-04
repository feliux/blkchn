package main

import (
	"flag"

	"github.com/feliux/blkchn/server"
)

var port int

func init() {
	flag.IntVar(&port, "port", 5000, "TCP Port Number for Blockchain Server.")
	flag.Parse()
}

func main() {
	app := server.NewBlockchainServer(port)
	app.Run()
}
