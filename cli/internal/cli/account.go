package cli

import (
	"context"
	"log"
	"os"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/pkg/sshutil"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var accountRoot = &cobra.Command{
	Use: "account",
}

func mustConnectAccountService() account.AccountServiceClient {
	conn := mustConnectGrpc(accountApi)
	return account.NewAccountServiceClient(conn)
}

var (
	name        string
	description string
	wellKnownId string
)
var accountCreate = &cobra.Command{
	Use:   "create",
	Short: "iamcli account create --name NAME [--description DESCRIPTION] [--well-known-id WELL-KNOWN-ID]",

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccountService()

		resp, err := client.CreateAccount(context.Background(), &account.CreateAccountRequest{
			WellKnownId: wellKnownId,
			Name:        name,
			Description: description,
		})
		handleError(err)

		print(resp.GetAccount())
	},
}

var accountGet = &cobra.Command{
	Use:   "get",
	Short: "iamcli account get ID",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccountService()

		resp, err := client.GetAccount(context.Background(), &account.GetAccountRequest{
			AccountId: args[0],
		})
		handleError(err)

		print(resp.GetAccount())
	},
}

var accountUpdate = &cobra.Command{
	Use:   "update",
	Short: "iamcli account update ID [--name NAME] [--description DESCRIPTION]",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccountService()

		paths := []string{}
		if cmd.PersistentFlags().Changed("name") {
			paths = append(paths, "name")
		}
		if cmd.PersistentFlags().Changed("description") {
			paths = append(paths, "description")
		}

		updateMask, err := fieldmaskpb.New(&account.Account{}, paths...)
		handleError(err)

		resp, err := client.UpdateAccount(context.Background(), &account.UpdateAccountRequest{
			AccountId:   args[0],
			Name:        name,
			Description: description,
			UpdateMask:  updateMask,
		})
		handleError(err)

		print(resp.GetAccount())
	},
}

var accountDelete = &cobra.Command{
	Use:   "delete",
	Short: "iamcli account get ID",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccountService()

		resp, err := client.DeleteAccount(context.Background(), &account.DeleteAccountRequest{
			AccountId: args[0],
		})
		handleError(err)

		print(resp)
	},
}

var accountKeysRoot = &cobra.Command{
	Use:   "key",
	Short: "iamcli key",
}

var accountId string
var accountSshKeysList = &cobra.Command{
	Use:   "list",
	Short: "iamcli account key list --account-id ID",

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccountService()

		resp, err := client.ListSshKeys(context.Background(), &account.ListSshKeysRequest{
			AccountId: accountId,
		})
		handleError(err)

		print(resp.GetKeys())
	},
}

var filePath string
var accountSshKeyCreate = &cobra.Command{
	Use:   "create",
	Short: "iamcli account key create --account-id ID --file SSH_PUB_KEY_FILE_PATH",

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccountService()

		key, err := os.ReadFile(filePath)
		handleError(err)

		resp, err := client.CreateSshKey(context.Background(), &account.CreateSshKeyRequest{
			AccountId: accountId,
			PublicKey: string(key),
		})
		handleError(err)

		print(resp.GetKey())
	},
}

var fingerprint string
var accountSshKeyDelete = &cobra.Command{
	Use:   "delete",
	Short: "iamcli account key delete --account-id ID <--file SSH_PUB_KEY_FILE_PATH or --fingerprint FINGERPRINT>",

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccountService()

		if cmd.Flags().Changed("fingerprint") == cmd.Flags().Changed("file") {
			log.Fatal("expected either --fingerprint or --file")
		}

		if cmd.Flags().Changed("file") {
			key, err := os.ReadFile(filePath)
			handleError(err)

			parsedKey, err := sshutil.ParsePublicKey(key)
			handleError(err)

			fingerprint = sshutil.GetFingerprint(parsedKey)
		}

		resp, err := client.DeleteSshKey(context.Background(), &account.DeleteSshKeyRequest{
			AccountId:      accountId,
			KeyFingerprint: fingerprint,
		})
		handleError(err)

		print(resp)
	},
}

func init() {
	accountCreate.PersistentFlags().StringVar(&name, "name", "", "Account Name")
	accountCreate.MarkFlagRequired("name")
	accountCreate.PersistentFlags().StringVar(&description, "description", "", "Account Description")
	accountCreate.PersistentFlags().StringVar(&wellKnownId, "well-known-id", "", "Account Well-Known ID")

	accountUpdate.PersistentFlags().StringVar(&name, "name", "", "New name for the Account")
	accountUpdate.PersistentFlags().StringVar(&description, "description", "", "New description for the Account")

	accountKeysRoot.PersistentFlags().StringVar(&accountId, "account-id", "", "Account ID")
	accountKeysRoot.MarkFlagRequired("account-id")
	accountKeysRoot.AddCommand(accountSshKeysList)

	accountSshKeyCreate.PersistentFlags().StringVar(&filePath, "file", "", "Path to ssh public key")
	accountSshKeyCreate.MarkFlagRequired("file")
	accountKeysRoot.AddCommand(accountSshKeyCreate)

	accountSshKeyDelete.PersistentFlags().StringVar(&filePath, "file", "", "Path to ssh public key")
	accountSshKeyDelete.PersistentFlags().StringVar(&fingerprint, "fingerprint", "", "Ssh public key fingerprint")
	accountKeysRoot.AddCommand(accountSshKeyDelete)

	accountRoot.AddCommand(accountCreate)
	accountRoot.AddCommand(accountUpdate)
	accountRoot.AddCommand(accountGet)
	accountRoot.AddCommand(accountDelete)
}
