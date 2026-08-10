package cmd

import (
	"fmt"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints current hardware and module status",
	Run: func(cmd *cobra.Command, args []string) {
		caps := hardware.Detect()
		fmt.Println("=== Linuwu-Sense Status ===")
		fmt.Printf("Module Loaded: %v\n", caps.ModuleLoaded)
		if !caps.ModuleLoaded {
			return
		}
		fmt.Printf("Model Detected: %s\n", caps.Model)
		if caps.HasPlatformProfile {
			prof, _ := hardware.GetPlatformProfile()
			fmt.Printf("Thermal Profile: %s\n", prof)
		}
		sensors := hardware.ReadSensors(caps)
		fmt.Printf("Temps: CPU %d°C | GPU %d°C\n", sensors.CPUTemp, sensors.GPUTemp)
		fmt.Printf("Fans:  CPU %d RPM | GPU %d RPM\n", sensors.CPUFan, sensors.GPUFan)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
