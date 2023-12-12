package cli

import (
	"context"
	"log"

	"github.com/de1phin/iam/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tokenApi   string
	accountApi string
	accessApi  string
)

var root = &cobra.Command{
	Use: "iamcli",
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustConnectGrpc(addr string) *grpc.ClientConn {
	conn, err := grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("grpc connect failed", zap.String("endpoint", addr), zap.Error(err))
	}
	return conn
}

func Run() {
	root.PersistentFlags().StringVar(&accessApi, "access", "access.iam.de1phin.ru:80", "access-service api endpoint")
	root.PersistentFlags().StringVar(&tokenApi, "token", "token.iam.de1phin.ru:80", "token-service api endpoint")
	root.PersistentFlags().StringVar(&accountApi, "account", "account.iam.de1phin.ru:80", "account-service api endpoint")

	root.AddCommand(accountRoot)

	root.Execute()
}
