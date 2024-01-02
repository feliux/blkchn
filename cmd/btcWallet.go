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
}
