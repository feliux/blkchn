package block

import (
	"github.com/feliux/blkchn/transaction"
)

type Block struct {
	Timestamp    int64
	Nonce        int
	PreviousHash [32]byte
	Transactions []*transaction.Transaction
}
