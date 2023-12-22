package blockchain

import (
	"github.com/feliux/blkchn/block"
)

type Blockchain struct {
	TransactionPool []string //[]*Transaction
	Chain           []*block.Block
}
