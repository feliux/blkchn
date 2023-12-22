package main

import (
	"github.com/feliux/blkchn/blockchain"
)

func init() {}

func main() {
	// genesis
	bc := blockchain.NewBlockchain()
	bc.Print()

	bc.AddTransaction("A", "B", 1.0)
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(5, previousHash)
	bc.Print()

	bc.AddTransaction("C", "D", 2.0)
	bc.AddTransaction("X", "Y", 3.0)
	previousHash = bc.LastBlock().Hash()
	bc.CreateBlock(2, previousHash)
	bc.Print()
}
