package cmd

import (
	"fmt"
	"os"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "linuwu-sense-tui",
	Short: "Terminal UI and control utility for Acer Predator/Nitro Linuwu-Sense driver",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
