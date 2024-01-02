package main

import (
	"log"

	"github.com/feliux/blkchn/blockchain"
)

func init() {}

func main() {
	// genesis
	myBlockchainAddress := "my_blockchain_address"
	bc := blockchain.NewBlockchain(myBlockchainAddress)
	bc.Print()

	bc.AddTransaction("A", "B", 1.0)
	bc.Mining()
	bc.Print()

	bc.AddTransaction("C", "D", 2.0)
	bc.AddTransaction("X", "Y", 3.0)
	bc.Mining()
	bc.Print()

	log.Printf("my %.1f\n", bc.CalculateTotalAmount("my_blockchain_address"))
	log.Printf("C %.1f\n", bc.CalculateTotalAmount("C"))
	log.Printf("D %.1f\n", bc.CalculateTotalAmount("D"))
}
