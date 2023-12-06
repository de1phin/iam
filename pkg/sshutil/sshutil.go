package sshutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

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
	var buf []byte
	base64.StdEncoding.Encode(buf, encryptedBytes)
	return buf, nil
}

func DecryptWithPrivateKey(data []byte, key []byte) ([]byte, error) {
	var base64decoded []byte
	_, err := base64.StdEncoding.Decode(base64decoded, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	pemBlock, _ := pem.Decode(key)

	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

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
