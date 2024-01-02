package blockchain

import (
	"github.com/feliux/blkchn/block"
	"github.com/feliux/blkchn/transaction"
)

type Blockchain struct {
	TransactionPool   []*transaction.Transaction
	Chain             []*block.Block
	BlockchainAddress string
	Port              int
}
