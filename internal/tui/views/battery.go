package views

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/sysfs"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
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

	if caps.HasBacklightTimeout {
		val, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "backlight_timeout"))
		lines = append(lines, fmt.Sprintf("[4] Backlight Timeout (30s):   %s\n    Turns off keyboard backlight after 30 seconds idle.", renderToggle(val == 1)))
	}

	if caps.HasLCDOverride {
		val, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "lcd_override"))
		lines = append(lines, fmt.Sprintf("[5] LCD Override:              %s\n    Reduces LCD latency and ghosting.", renderToggle(val == 1)))
	}

	if caps.HasBootAnimationSound {
		val, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "boot_animation_sound"))
		lines = append(lines, fmt.Sprintf("[6] Boot Animation Sound:      %s\n    Enables/disables Predator boot animation sound effect.", renderToggle(val == 1)))
	}

	box := style.MakeBox("Power & Display Features", strings.Join(lines, "\n\n"))

	hints := lipgloss.NewStyle().Foreground(style.ColorMuted).Render("Press [1]-[6] to toggle/cycle settings, or click any option with your mouse.")
	return lipgloss.JoinVertical(lipgloss.Left, box, "\n", hints)
}
