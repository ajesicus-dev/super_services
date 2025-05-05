package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	upFlag    bool
	downFlag  bool
	vaultFlag bool
)

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Manage containers",
	Run: func(cmd *cobra.Command, args []string) {
		// Mutual exclusivity check
		if upFlag && downFlag {
			fmt.Fprintln(os.Stderr, "Error: --up and --down cannot be used together")
			os.Exit(1)
		}

		if upFlag {
			fmt.Println("Running container up...")
			run("docker-compose", "up", "-d")
		} else if downFlag {
			fmt.Println("Running container down...")
			run("docker-compose", "down")
		} else {
			fmt.Fprintln(os.Stderr, "Please use --up or --down")
			cmd.Help()
		}
	},
}

var vaultComposeFiles = []string{
	"-f", "./deployments/vault/docker-compose.yml",
	"-f", "./deployments/vault/docker-compose.dev.yml",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start containers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting containers...")

		dockerArgs := []string{"compose"}
		if vaultFlag {
			dockerArgs = append(dockerArgs, vaultComposeFiles...)
		}

		dockerArgs = append(dockerArgs, "up", "-d")

		run("docker", dockerArgs...)
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop containers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping containers...")
		dockerArgs := []string{"compose"}
		if vaultFlag {
			dockerArgs = append(dockerArgs, vaultComposeFiles...)
		}

		dockerArgs = append(dockerArgs, "down")

		run("docker", dockerArgs...)
	},
}

func init() {
	containerCmd.AddCommand(upCmd)
	upCmd.Flags().BoolVar(&vaultFlag, "vault", false, "Start only containers of vault")

	containerCmd.AddCommand(downCmd)
	downCmd.Flags().BoolVar(&vaultFlag, "vault", false, "Stop only containers of vault")

	rootCmd.AddCommand(containerCmd)
	containerCmd.Flags().BoolVar(&upFlag, "up", false, "Start the container")
	containerCmd.Flags().BoolVar(&downFlag, "down", false, "Stop the container")
}
