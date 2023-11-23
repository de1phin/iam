package server

import "golang.org/x/crypto/ssh"

func GetFingerprint(pubKey ssh.PublicKey) string {
	return ssh.FingerprintSHA256(pubKey)
}
