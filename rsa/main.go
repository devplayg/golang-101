package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func main() {
	priv, pub, _ := issueKeys()

	msg := []byte("hello")
	EncryptOAEP(msg, priv, pub)
	EncryptPkcs(msg, priv, pub)
}

func EncryptPkcs(message []byte, priv *rsa.PrivateKey, pub *rsa.PublicKey) {
	encrypted, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, message)
	decrypted, _ := rsa.DecryptPKCS1v15(nil, priv, encrypted)
	fmt.Printf("%x\n", message)
	fmt.Printf("%x\n", encrypted)
	fmt.Printf("%x\n", decrypted)


}

func EncryptOAEP(message []byte, priv *rsa.PrivateKey, pub *rsa.PublicKey) {
	encrypted, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, message, nil)
	decrypted, _ := priv.Decrypt(nil, encrypted, &rsa.OAEPOptions{Hash: crypto.SHA256})
	fmt.Printf("%x\n", message)
	fmt.Printf("%x\n", encrypted)
	fmt.Printf("%x\n", decrypted)
}

func issueKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}
