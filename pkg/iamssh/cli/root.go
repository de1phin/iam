package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "iamssh",
	Short: "ssh wrapper for iam bastion",

	Args: sshCmd.Args,
	Run:  sshCmd.Run,
}

func Run() {
	root.AddCommand(sshCmd)
	root.AddCommand(loadKeyCmd)

	root.Execute()
}

func interactiveCheck(msg string) bool {
	fmt.Println(msg, "[y/n]?")
	var s string
	fmt.Scanln(&s)
	return strings.Trim(strings.ToLower(s), " \n\t") == "y"
}
