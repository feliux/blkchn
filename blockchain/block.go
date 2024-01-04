package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/feliux/blkchn/transaction"
)

func NewBlock(nonce int, previousHash [32]byte, transactions []*transaction.Transaction) *Block {
	b := new(Block)
	b.timestamp = 0 // time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func (b *Block) Nonce() int {
	return b.nonce
}

func (b *Block) Transactions() []*transaction.Transaction {
	return b.transactions
}

func (b *Block) Print(n int) {
	fmt.Printf("Chain %d ---> timestamp: %d | nonce: %d | previousHash: %x \n", n, b.timestamp, b.nonce, b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	jsonBlock, err := json.Marshal(b)
	if err != nil {
		log.Printf("ERROR marshaling data: %s", err.Error())
	}
	return sha256.Sum256([]byte(jsonBlock))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64                      `json:"timestamp"`
		Nonce        int                        `json:"nonce"`
		PreviousHash string                     `json:"previousHash"`
		Transactions []*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	str := &struct {
		Timestamp    *int64                      `json:"timestamp"`
		Nonce        *int                        `json:"nonce"`
		PreviousHash *string                     `json:"previousHash"`
		Transactions *[]*transaction.Transaction `json:"transactions"`
	}{
		Timestamp:    &b.timestamp,
		Nonce:        &b.nonce,
		PreviousHash: &previousHash,
		Transactions: &b.transactions,
	}
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	ph, err := hex.DecodeString(*str.PreviousHash)
	if err != nil {
		log.Printf("ERROR decoding hex to bytes: %s", err)
	}
	copy(b.previousHash[:], ph[:32])
	return nil
}
