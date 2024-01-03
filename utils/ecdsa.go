package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/feliux/blkchn/signature"
)

func String2BigIntTuple(s string) (big.Int, big.Int) {
	// returns the bytes represented by the hexadecimal string
	bx, err := hex.DecodeString(s[:64])
	if err != nil {
		log.Printf("ERROR converting from hex to bytes: %s" + err.Error())
	}
	by, err := hex.DecodeString(s[64:])
	if err != nil {
		log.Printf("ERROR converting from hex to bytes: %s" + err.Error())
	}
	var bix big.Int
	var biy big.Int
	// interprets buf as the bytes of a big-endian unsigned integer
	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)

	return bix, biy
}

func SignatureFromString(s string) *signature.Signature {
	x, y := String2BigIntTuple(s)
	return &signature.Signature{&x, &y}
}

func PublicKeyFromString(s string) *ecdsa.PublicKey {
	x, y := String2BigIntTuple(s)
	return &ecdsa.PublicKey{elliptic.P256(), &x, &y}
}

func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, err := hex.DecodeString(s[:])
	if err != nil {
		log.Printf("ERROR converting from hex to bytes: %s" + err.Error())
	}
	var bi big.Int
	_ = bi.SetBytes(b)
	return &ecdsa.PrivateKey{*publicKey, &bi}
}
