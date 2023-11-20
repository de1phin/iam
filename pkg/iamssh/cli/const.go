package cli

import (
	"os"
	"path"
)

const (
	BastionHost     = "bastion.iam.de1phin.ru"
	BastionUser     = "ubuntu"
	FolderId        = "b1gqe3skkuiko3bv671e"
	LockBoxKeysName = "iam_bastion"
)

func getAbsolutePath(file string) string {
	homedir, _ := os.UserHomeDir()
	return path.Join(homedir, file)
}

var (
	IamBastionSshKeyFile    = getAbsolutePath(".ssh/iam_bastion")
	IamBastionSshPubKeyFile = getAbsolutePath(".ssh/iam_bastion.pub")
)
