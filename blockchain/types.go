package blockchain

import (
	"sync"

	"github.com/feliux/blkchn/block"
	"github.com/feliux/blkchn/transaction"
)

type Blockchain struct {
	transactionPool   []*transaction.Transaction
	chain             []*block.Block
	blockchainAddress string
	port              int
	mux               sync.Mutex
}
