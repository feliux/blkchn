package main

import (
	"fmt"

	"github.com/feliux/blkchn/wallet"
)

func init() {}

func main() {
	w := wallet.NewWallet()
	fmt.Println("PrivateKey: ", w.PrivateKeyStr())
	fmt.Println("PublicKey: ", w.PublicKeyStr())
	fmt.Println("Bitcoin address: ", w.BlockchainAddress())

	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "B", 1.0)
	fmt.Printf("signature %s\n", t.GenerateSignature())
}
