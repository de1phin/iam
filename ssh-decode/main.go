package main

import (
	"fmt"
	"log"
	"os"

	"github.com/de1phin/iam/pkg/sshutil"
)

func main() {

	private, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	decrypted, err := sshutil.DecryptWithPrivateKey([]byte(os.Args[2]), private)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(decrypted))

}
