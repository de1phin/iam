package cli

import (
	"context"
	"fmt"

	access "github.com/de1phin/iam/genproto/services/access/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mustConnectAccessService() access.AccessServiceClient {
	conn := mustConnectGrpc(accessApi)
	return access.NewAccessServiceClient(conn)
}

var roleRoot = &cobra.Command{
	Use: "role",
}

var roleName string
var rolePerms []string
var roleCreate = &cobra.Command{
	Use:   "create",
	Short: "iamcli role create --name NAME --permission PERM1 --permission PERM2 ...",

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccessService()

		_, err := client.AddRole(context.Background(), &access.AddRoleRequest{
			Role: &access.Role{
				Name:        roleName,
				Permissions: rolePerms,
			},
		})
		handleError(err)
	},
}

var roleGet = &cobra.Command{
	Use:   "get",
	Short: "iamcli role get NAME",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccessService()

		resp, err := client.GetRole(context.Background(), &access.GetRoleRequest{
			Name: args[0],
		})
		handleError(err)

		print(resp.GetRole())
	},
}

var roleDelete = &cobra.Command{
	Use:   "delete",
	Short: "iamcli role delete NAME",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccessService()

		_, err := client.DeleteRole(context.Background(), &access.DeleteRoleRequest{
			Name: args[0],
		})
		handleError(err)
	},
}

var roleAccountId string
var roleResource string
var roleGrant = &cobra.Command{
	Use:   "grant",
	Short: "iamcli role grant ROLE --account-id ACCOUNT_ID --resource RESOURCE",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccessService()

		_, err := client.AddAccessBinding(context.Background(), &access.AddAccessBindingRequest{
			AccessBinding: &access.AccessBinding{
				Resource:  roleResource,
				RoleName:  args[0],
				AccountId: roleAccountId,
			},
		})
		handleError(err)
	},
}

var roleRevoke = &cobra.Command{
	Use:   "revoke",
	Short: "iamcli role revoke ROLE --account-id ACCOUNT_ID --resorce RESOURCE",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccessService()

		_, err := client.DeleteAccessBinding(context.Background(), &access.DeleteAccessBindingRequest{
			AccessBinding: &access.AccessBinding{
				Resource:  roleResource,
				RoleName:  args[0],
				AccountId: roleAccountId,
			},
		})
		handleError(err)
	},
}

var authorizeToken string
var authorizePermission string
var authorizeResource string
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "iamcli authorize --permission PERM --role ROLE --token TOKEN",

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccessService()

		_, err := client.CheckPermission(context.Background(), &access.CheckPermissionRequest{
			Token:      authorizeToken,
			Resource:   authorizeResource,
			Permission: authorizePermission,
		})

		status := status.Convert(err)
		if status != nil &&
			(status.Code() == codes.Unauthenticated || status.Code() == codes.PermissionDenied) {
			fmt.Println(status.Err().Error())
			return
		}
		handleError(err)
	},
}

func init() {
	roleCreate.PersistentFlags().StringVar(&roleName, "name", "", "role name")
	roleCreate.PersistentFlags().StringArrayVar(&rolePerms, "permission", nil, "role permission")
	roleCreate.MarkFlagRequired("name")
	roleCreate.MarkFlagRequired("permission")
	roleRoot.AddCommand(roleCreate)
	roleRoot.AddCommand(roleGet)
	roleRoot.AddCommand(roleDelete)

	roleRoot.PersistentFlags().StringVar(&roleAccountId, "account-id", "", "account id to grant role to")
	roleRoot.PersistentFlags().StringVar(&roleResource, "resource", "", "resource to grant role for")
	roleRoot.AddCommand(roleGrant)
	roleRoot.AddCommand(roleRevoke)
	roleGrant.MarkFlagRequired("resource")
	roleGrant.MarkFlagRequired("account-id")
	roleRevoke.MarkFlagRequired("resource")
	roleRevoke.MarkFlagRequired("account-id")

	authorizeCmd.PersistentFlags().StringVar(&authorizeToken, "token", "", "account auth token")
	authorizeCmd.PersistentFlags().StringVar(&authorizePermission, "permission", "", "permission to authorize")
	authorizeCmd.PersistentFlags().StringVar(&authorizeResource, "resource", "", "resource to authorize permission for")
	authorizeCmd.MarkFlagRequired("token")
	authorizeCmd.MarkFlagRequired("permission")
	authorizeCmd.MarkFlagRequired("resource")

	roleRoot.AddCommand(authorizeCmd)
}
