module github.com/feliux/blkchn

go 1.21.4

replace github.com/feliux/blkchn/blockchain => ./blockchain

replace github.com/feliux/blkchn/transaction => ./transaction

replace github.com/feliux/blkchn/wallet => ./wallet

replace github.com/feliux/blkchn/signature => ./signature

replace github.com/feliux/blkchn/server => ./server

replace github.com/feliux/blkchn/walletserver => ./walletserver

replace github.com/feliux/blkchn/utils => ./utils

require (
	github.com/btcsuite/btcutil v1.0.2
	golang.org/x/crypto v0.17.0
)
