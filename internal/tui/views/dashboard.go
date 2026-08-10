package views

import (
	"fmt"
	"strings"

	"github.com/Himesh-Kundal/linuwu-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-tui/internal/sysfs"
	"github.com/Himesh-Kundal/linuwu-tui/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

func RenderDashboard(caps hardware.Capabilities, sensors hardware.SensorData) string {
	if !caps.ModuleLoaded {
		return lipgloss.NewStyle().
			Foreground(style.ColorDanger).
			Bold(true).
			Render("⚠️  ERROR: linuwu_sense kernel module is not loaded!\nRun 'linuwu-tui setup' or 'sudo modprobe linuwu_sense' to start.")
	}

	// 1. Overview Box
	profileStr, _ := hardware.GetPlatformProfile()
	if profileStr == "" {
		profileStr = "N/A"
	}

	overviewContent := fmt.Sprintf(
		"%s %s\n%s %s",
		style.StyleLabel.Render("Model Detected:"), style.StyleValue.Render(string(caps.Model)),
		style.StyleLabel.Render("Thermal Profile:"), style.StyleValue.Render(strings.ToUpper(profileStr)),
	)
	overviewBox := style.MakeBox("System Overview", overviewContent)

	// 2. Sensors Box
	cpuBar := renderBar(sensors.CPUTemp, 100, 15)
	gpuBar := renderBar(sensors.GPUTemp, 100, 15)

	sensorContent := fmt.Sprintf(
		"CPU Temp: %s %d°C\nGPU Temp: %s %d°C\nSYS Temp: %s %d°C\n\nCPU Fan:  %s RPM\nGPU Fan:  %s RPM",
		cpuBar, sensors.CPUTemp,
		gpuBar, sensors.GPUTemp,
		renderBar(sensors.SYSTemp, 100, 15), sensors.SYSTemp,
		style.StyleValue.Render(fmt.Sprintf("%d", sensors.CPUFan)),
		style.StyleValue.Render(fmt.Sprintf("%d", sensors.GPUFan)),
	)
	sensorsBox := style.MakeBox("Live Sensors", sensorContent)

	// 3. Quick Status Box
	var statusLines []string
	if caps.HasBatteryLimiter {
		val, _ := sysfs.ReadInt(caps.SensePath + "/battery_limiter")
		statusLines = append(statusLines, fmt.Sprintf("%s %s", style.StyleLabel.Render("Battery Limiter (80%):"), renderToggle(val == 1)))
	}
	if caps.HasUSBCharging {
		val, _ := sysfs.ReadInt(caps.SensePath + "/usb_charging")
		statusLines = append(statusLines, fmt.Sprintf("%s %s", style.StyleLabel.Render("USB Charging Cutoff:"), style.StyleValue.Render(fmt.Sprintf("%d%%", val))))
	}
	if caps.HasBacklightTimeout {
		val, _ := sysfs.ReadInt(caps.SensePath + "/backlight_timeout")
		statusLines = append(statusLines, fmt.Sprintf("%s %s", style.StyleLabel.Render("Backlight Timeout:"), renderToggle(val == 1)))
	}
	if caps.HasLCDOverride {
		val, _ := sysfs.ReadInt(caps.SensePath + "/lcd_override")
		statusLines = append(statusLines, fmt.Sprintf("%s %s", style.StyleLabel.Render("LCD Override:"), renderToggle(val == 1)))
	}
	if caps.HasBootAnimationSound {
		val, _ := sysfs.ReadInt(caps.SensePath + "/boot_animation_sound")
		statusLines = append(statusLines, fmt.Sprintf("%s %s", style.StyleLabel.Render("Boot Sound:"), renderToggle(val == 1)))
	}

	statusBox := style.MakeBox("Quick Status", strings.Join(statusLines, "\n"))

	return lipgloss.JoinVertical(lipgloss.Left, overviewBox, sensorsBox, statusBox)
}

func renderToggle(active bool) string {
	if active {
		return lipgloss.NewStyle().Foreground(style.ColorSuccess).Bold(true).Render("[ON]")
	}
	return lipgloss.NewStyle().Foreground(style.ColorMuted).Render("[OFF]")
}

func renderBar(val, max, width int) string {
	if val < 0 {
		val = 0
	}
	if val > max {
		val = max
	}
	fill := (val * width) / max
	empty := width - fill

	fillStr := strings.Repeat("█", fill)
	emptyStr := strings.Repeat("░", empty)

	color := style.ColorSuccess
	if val > 80 {
		color = style.ColorDanger
	} else if val > 65 {
		color = style.ColorWarning
	}

	return lipgloss.NewStyle().Foreground(color).Render(fillStr) + lipgloss.NewStyle().Foreground(style.ColorMuted).Render(emptyStr)
}
