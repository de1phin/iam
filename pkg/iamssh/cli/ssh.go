package cli

import (
	"log"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "run shell on host via bastion",

	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, _ []string) {
		log.Fatal("Unimplemented")
	},
}
