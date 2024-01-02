package transaction

import (
	"log"
	"strings"
)

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	log.Printf("%s\n", strings.Repeat("-", 40))
	log.Printf("sender_blockchain_address      %s\n", t.SenderBlockchainAddress)
	log.Printf("recipient_blockchain_address   %s\n", t.RecipientBlockchainAddress)
	log.Printf("value                          %.1f\n", t.Value)
}

/*
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.SenderBlockchainAddress,
		Recipient: t.RecipientBlockchainAddress,
		Value:     t.Value,
	})
}
*/
