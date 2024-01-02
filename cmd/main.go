package main

import (
	"log"

	"github.com/feliux/blkchn/blockchain"
	"github.com/feliux/blkchn/wallet"
)

func init() {}

func main() {
	// generate wallets
	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// genesis on walletM as miner
	bc := blockchain.NewBlockchain(walletM.BlockchainAddress(), 5000)

	// send 1.0 from A to B
	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)
	// verify the transaction
	isAdded := bc.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, walletA.PublicKey(), t.GenerateSignature())
	log.Println("Added transaction? ", isAdded)
	// mining
	bc.Mining()
	bc.Print()

	log.Printf("Miner %.1f\n", bc.CalculateTotalAmount(walletM.BlockchainAddress()))
	log.Printf("A %.1f\n", bc.CalculateTotalAmount(walletA.BlockchainAddress()))
	log.Printf("B %.1f\n", bc.CalculateTotalAmount(walletB.BlockchainAddress()))
}
