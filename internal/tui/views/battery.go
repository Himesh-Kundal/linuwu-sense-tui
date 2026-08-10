package views

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Himesh-Kundal/linuwu-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-tui/internal/sysfs"
	"github.com/Himesh-Kundal/linuwu-tui/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

func RenderBattery(caps hardware.Capabilities) string {
	var lines []string

	if caps.HasBatteryLimiter {
		val, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "battery_limiter"))
		lines = append(lines, fmt.Sprintf("[1] Battery Limiter (80%% Cap): %s\n    Preserves battery health when plugged in.", renderToggle(val == 1)))
	}

	if caps.HasBatteryCalibration {
		val, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "battery_calibration"))
		lines = append(lines, fmt.Sprintf("[2] Battery Calibration Mode:  %s\n    Calibrates fuel gauge (100%% -> 0%% -> 100%%). Keep AC connected!", renderToggle(val == 1)))
	}

	if caps.HasUSBCharging {
		val, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "usb_charging"))
		lines = append(lines, fmt.Sprintf("[3] USB Charge Cutoff:         [%d%%]\n    Allows charging devices when powered off until threshold.", val))
	}

	box := style.MakeBox("Power & Battery Management", strings.Join(lines, "\n\n"))

	hints := lipgloss.NewStyle().Foreground(style.ColorMuted).Render("Press [1], [2], or [3] to toggle/cycle settings.")
	return lipgloss.JoinVertical(lipgloss.Left, box, "\n", hints)
}
