package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/de1phin/iam/pkg/yccli"
	"github.com/spf13/cobra"
)

func getAbsolutePath(file string) string {
	homedir, _ := os.UserHomeDir()
	return path.Join(homedir, file)
}

const (
	LockBoxKeysName = "iam_bastion"
)

var (
	IamBastionSshKeyFile    = getAbsolutePath(".ssh/iam_bastion")
	IamBastionSshPubKeyFile = getAbsolutePath(".ssh/iam_bastion.pub")
)

type sshKey struct {
	Private []byte `json:"private"`
	Public  []byte `json:"public"`
}

func (key sshKey) writePair(privatePath, publicPath string) error {
	return errors.Join(
		os.WriteFile(privatePath, key.Private, 0600),
		os.WriteFile(publicPath, key.Public, 0600),
	)
}

var loadKeyCmd = &cobra.Command{
	Use:   "load-key",
	Short: "Load bastion key from lockbox. Requires configured YC CLI",

	Run: func(_ *cobra.Command, _ []string) {

		key, err := loadKeyFromLockbox()
		if err != nil {
			log.Fatal("failed to load key from lockbox: ", err)
		}

		err = key.writePair(IamBastionSshKeyFile, IamBastionSshPubKeyFile)
		if err != nil {
			log.Fatal("failed to save ssh keys: ", err)
		}
	},
}

func loadKeyFromLockbox() (*sshKey, error) {

	secrets, err := yccli.ListLockboxSecrets(FolderId)
	if err != nil {
		return nil, fmt.Errorf("failed to list lockbox secrets: %w", err)
	}

	secretId := ""
	for _, secret := range secrets {
		if secret.Name == LockBoxKeysName {
			secretId = secret.Id
			break
		}
	}

	if secretId == "" {
		return nil, fmt.Errorf("no secret with name %s in folder %s", LockBoxKeysName, FolderId)
	}

	secret, err := yccli.LockboxSecretGet(secretId)
	if err != nil {
		return nil, fmt.Errorf("failed to get lockbox secret: %w", err)
	}

	sshKey := &sshKey{}
	for _, entry := range secret.Entries {
		if entry.Key == "private" {
			sshKey.Private = []byte(entry.TextValue + "\n")
		}
		if entry.Key == "public" {
			sshKey.Public = []byte(entry.TextValue + "\n")
		}
	}

	return sshKey, nil
}
