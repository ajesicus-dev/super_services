package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	upFlag   bool
	downFlag bool
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

func init() {
	rootCmd.AddCommand(containerCmd)
	containerCmd.Flags().BoolVar(&upFlag, "up", false, "Start the container")
	containerCmd.Flags().BoolVar(&downFlag, "down", false, "Stop the container")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}
