package blockchain

import (
	"sync"

	"github.com/feliux/blkchn/transaction"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*transaction.Transaction
}

type Blockchain struct {
	transactionPool   []*transaction.Transaction
	chain             []*Block
	blockchainAddress string
	port              int
	mux               sync.Mutex
	neighbors         []string
	muxNeighbors      sync.Mutex
}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}
