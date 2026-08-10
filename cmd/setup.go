package cmd

import (
	"fmt"
	"os"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/driver"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Extracts and installs the embedded linuwu_sense kernel driver",
	Run: func(cmd *cobra.Command, args []string) {
		if err := driver.Install(); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing driver: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
