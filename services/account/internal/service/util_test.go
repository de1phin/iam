package service_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

func mustGenerateSshKey() ([]byte, ssh.PublicKey) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	sshPubKey, err = ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		panic(err)
	}
	block, err := ssh.MarshalPrivateKey(rsaKey, "")
	if err != nil {
		panic(err)
	}
	sshKey := pem.EncodeToMemory(block)
	return sshKey, sshPubKey
}
