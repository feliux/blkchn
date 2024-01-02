package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/feliux/blkchn/transaction"
)

func NewBlock(nonce int, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	b := new(Block)
	b.Timestamp = 0 // time.Now().UnixNano()
	b.Nonce = nonce
	b.PreviousHash = previousHash
	b.Transactions = transactions
	return b
}

func (b *Block) Print(n int) {
	fmt.Printf("Chain %d ---> timestamp: %d | nonce: %d | previousHash: %x \n", n, b.Timestamp, b.Nonce, b.PreviousHash)
	for _, t := range b.Transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	jsonBlock, _ := json.Marshal(b)
	return sha256.Sum256([]byte(jsonBlock))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64                      `json:"timestamp"`
		Nonce        int                        `json:"nonce"`
		PreviousHash [32]byte                   `json:"previousHash"`
		Transactions []*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    b.Timestamp,
		Nonce:        b.Nonce,
		PreviousHash: b.PreviousHash,
		Transactions: b.Transactions,
	})
}
