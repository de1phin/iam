package sshutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func GetFingerprint(pubKey ssh.PublicKey) string {
	return ssh.FingerprintSHA256(pubKey)
}

func EncryptWithPublicKey(data []byte, key []byte) ([]byte, error) {
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}

	parsedCryptoKey, ok := parsed.(ssh.CryptoPublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected: parsed key is not convertable to ssh.CryptoPublicKey")
	}

	publicCryptoKey := parsedCryptoKey.CryptoPublicKey()

	rsaKey, ok := publicCryptoKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected: parsed key is not convertable to rsa.PublicKey")
	}

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		rsaKey,
		data,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode data: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(encryptedBytes)
	return []byte(encoded), nil
}
func DecryptWithPrivateKey(data []byte, key []byte) ([]byte, error) {
	base64decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	sshPrivateKey, err := ssh.ParseRawPrivateKey(key)
	if sshErr, ok := err.(*ssh.PassphraseMissingError); ok && sshErr != nil {
		fmt.Println("Private Key is passphrase protected.")
		fmt.Print("Passphrase: ")
		passphrase, _ := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		sshPrivateKey, err = ssh.ParseRawPrivateKeyWithPassphrase(key, passphrase)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	privateKey := sshPrivateKey.(*rsa.PrivateKey)

	decrypted, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		base64decoded,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	return decrypted, nil
}

func ParsePublicKey(data []byte) (ssh.PublicKey, error) {
	key, err := ssh.ParsePublicKey(data)
	if err == nil {
		return key, nil
	}

	pk, _, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePublicKey(pk.Marshal())
}
