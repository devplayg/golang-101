package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strings"
)

func main() {
	privA, pubA, _ := issueKeys()
	privB, pubB, _ := issueKeys()

	fmt.Printf("Private key (A) %x\n", privA.D)
	fmt.Printf("Private key (B) %x\n", privB.D)
	fmt.Printf("Public key  (A) %x, %x\n", pubA.X, pubA.Y)
	fmt.Printf("Public key  (B) %x, %x\n", pubB.X, pubB.Y)
	fmt.Printf("%s\n", strings.Repeat("=", 100))

	secretA, _ := pubB.Curve.ScalarMult(pubB.X, pubB.Y, privA.D.Bytes())
	secretB, _ := pubA.Curve.ScalarMult(pubA.X, pubA.Y, privB.D.Bytes())

	fmt.Printf("Secret key  (A) %x\n", secretA)
	fmt.Printf("Secret key  (B) %x\n", secretB)
}

func issueKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}
