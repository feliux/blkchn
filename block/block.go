package block

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"time"

	"github.com/feliux/blkchn/transaction"
)

func NewBlock(nonce int, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	b := new(Block)
	b.Timestamp = time.Now().UnixNano()
	b.Nonce = nonce
	b.PreviousHash = previousHash
	b.Transactions = transactions
	return b
}

func (b *Block) Print() {
	log.Printf("timestamp       %d\n", b.Timestamp)
	log.Printf("nonce           %d\n", b.Nonce)
	log.Printf("previousHash   %x\n", b.PreviousHash)
	log.Printf("transactions    %s\n\n", b.Transactions)
	for _, t := range b.Transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	jsonBlock, _ := json.Marshal(b)
	return sha256.Sum256([]byte(jsonBlock))
}

/*
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64    `json:"timestamp"`
		Nonce        int      `json:"nonce"`
		PreviousHash [32]byte `json:"previousHash"`
		Transactions []*transactions.Transaction `json:"transactions"`
	}{
		Timestamp:    b.Timestamp,
		Nonce:        b.Nonce,
		PreviousHash: b.PreviousHash,
		Transactions: b.Transactions,
	})
}
*/
