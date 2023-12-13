package cli

import (
	"context"

	access "github.com/de1phin/iam/genproto/services/access/api"
	"github.com/spf13/cobra"
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

var roleGrantAccountId string
var roleGrantResource string
var roleGrant = &cobra.Command{
	Use:   "grant",
	Short: "iamcli role grant ROLE --account-id ACCOUNT_ID --resource RESOURCE",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		client := mustConnectAccessService()

		_, err := client.AddAccessBinding(context.Background(), &access.AddAccessBindingRequest{
			AccessBinding: &access.AccessBinding{
				Resource:  roleGrantResource,
				RoleName:  args[0],
				AccountId: roleGrantAccountId,
			},
		})
		handleError(err)
	},
}

var authorizeCmd = &cobra.Command{
	Use: "authorize",
}

func init() {

	roleCreate.PersistentFlags().StringVar(&roleName, "name", "", "role name")
	roleCreate.PersistentFlags().StringArrayVar(&rolePerms, "permission", nil, "role permission")
	roleRoot.AddCommand(roleCreate)
	roleRoot.AddCommand(roleGet)
	roleRoot.AddCommand(roleDelete)

	roleGrant.PersistentFlags().StringVar(&roleGrantAccountId, "account-id", "", "account id to grant role to")
	roleGrant.PersistentFlags().StringVar(&roleGrantResource, "resource", "", "resource to grant role for")
	roleRoot.AddCommand(roleGrant)
}
