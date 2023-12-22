package blockchain

import (
	"log"
	"strings"

	"github.com/feliux/blkchn/block"
	"github.com/feliux/blkchn/transaction"
)

func NewBlockchain() *Blockchain {
	b := block.Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *block.Block {
	b := block.NewBlock(nonce, previousHash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.TransactionPool = []*transaction.Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *block.Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.Chain {
		log.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	log.Printf("%s\n\n", strings.Repeat("*", 25))
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := transaction.NewTransaction(sender, recipient, value)
	bc.TransactionPool = append(bc.TransactionPool, t)
}
