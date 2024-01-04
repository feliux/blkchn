package main

import (
	"flag"
	"fmt"

	"github.com/feliux/blkchn/utils"
)

var (
	host      string
	port      int
	startIp   int
	endIp     int
	startPort int
	endPort   int
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host running nodes.")
	flag.IntVar(&port, "port", 5000, "TCP Port where node is lestening.")
	flag.IntVar(&startIp, "startIp", 0, "Start IP (not included).")
	flag.IntVar(&endIp, "endIp", 3, "Final IP (included).")
	flag.IntVar(&startPort, "startPort", 5000, "Start port (included).")
	flag.IntVar(&endPort, "endPort", 5003, "Final port (included).")
	flag.Parse()
}

func main() {
	fmt.Println(utils.FindNeighbors(host, port, startIp, endIp, startPort, endPort))
}
