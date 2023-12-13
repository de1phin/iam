package cli

import (
	"context"
	"os"

	token "github.com/de1phin/iam/genproto/services/token/api"
	"github.com/spf13/cobra"
)

func mustConnectTokenService() token.TokenServiceClient {
	conn := mustConnectGrpc(tokenApi)
	return token.NewTokenServiceClient(conn)
}

var tokenRoot = &cobra.Command{
	Use: "token",
}

var keyFile string
var tokenCreate = &cobra.Command{
	Use:   "create",
	Short: "iamcli token create --key-file SSH_PUB_KEY_FILE_PATH",

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectTokenService()

		key, err := os.ReadFile(keyFile)
		handleError(err)

		resp, err := client.CreateToken(context.Background(), &token.CreateTokenRequest{
			SshPubKey: string(key),
		})
		handleError(err)

		print(resp)
	},
}

var tokenCheck = &cobra.Command{
	Use:   "exchange",
	Short: "iamcli token exchange TOKEN",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectTokenService()

		resp, err := client.ExchangeToken(context.Background(), &token.ExchangeTokenRequest{
			Token: args[0],
		})
		handleError(err)

		print(resp)
	},
}

var tokenRefresh = &cobra.Command{
	Use:   "refresh",
	Short: "iamcli token refresh TOKEN",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectTokenService()

		resp, err := client.RefreshToken(context.Background(), &token.RefreshTokenRequest{
			Token: args[0],
		})
		handleError(err)

		print(resp)
	},
}

var tokenDelete = &cobra.Command{
	Use:   "delete",
	Short: "iamcli token delete TOKEN",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectTokenService()

		resp, err := client.DeleteToken(context.Background(), &token.DeleteTokenRequest{
			Token: args[0],
		})
		handleError(err)

		print(resp)
	},
}

func init() {
	tokenCreate.PersistentFlags().StringVar(&keyFile, "key-file", "", "Path to ssh public key")

	tokenRoot.AddCommand(tokenCreate)
	tokenRoot.AddCommand(tokenCheck)
	tokenRoot.AddCommand(tokenRefresh)
	tokenRoot.AddCommand(tokenDelete)
}
