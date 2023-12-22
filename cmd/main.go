package main

import (
	"github.com/feliux/blkchn/blockchain"
)

func init() {}

func main() {
	bc := blockchain.NewBlockchain()
	bc.Print()

	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(5, previousHash)
	bc.Print()

	previousHash = bc.LastBlock().Hash()
	bc.CreateBlock(2, previousHash)
	bc.Print()
}
