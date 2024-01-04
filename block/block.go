package block

import (
	"crypto/sha256"
	"encoding/hex"
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

func (b *Block) PreviousHash() [32]byte {
	return b.PreviousHash
}

func (b *Block) Nonce() int {
	return b.Nonce
}

func (b *Block) Transactions() []*Transaction {
	return b.Transactions
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
		PreviousHash string                     `json:"previousHash"`
		Transactions []*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    b.Timestamp,
		Nonce:        b.Nonce,
		PreviousHash: fmt.Sprintf("%x", b.PreviousHash),
		Transactions: b.Transactions,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	str := &struct {
		Timestamp    *int64                      `json:"timestamp"`
		Nonce        *int                        `json:"nonce"`
		PreviousHash *string                     `json:"previous_hash"`
		Transactions *[]*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    &b.Timestamp,
		Nonce:        &b.Nonce,
		PreviousHash: &previousHash,
		Transactions: &b.Transactions,
	}
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	ph, _ := hex.DecodeString(*str.PreviousHash)
	copy(b.previousHash[:], ph[:32])
	return nil
}
