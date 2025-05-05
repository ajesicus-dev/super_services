package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	readFlag  bool
	writeFlag bool
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage Vault",
	Run: func(cmd *cobra.Command, args []string) {
		// Mutual exclusivity check
		if readFlag && writeFlag {
			fmt.Fprintln(os.Stderr, "Error: --up and --down cannot be used together")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)
	containerCmd.Flags().BoolVar(&readFlag, "read", false, "read secret key")
	containerCmd.Flags().BoolVar(&writeFlag, "write", false, "write secret key")
}
