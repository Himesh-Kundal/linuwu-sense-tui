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

func RenderDashboard(caps hardware.Capabilities, sensors hardware.SensorData) string {
	if !caps.ModuleLoaded {
		return style.StyleDanger.Render("  ⚠  linuwu_sense module is not loaded\n\n  Run: sudo modprobe linuwu_sense")
	}

	profileStr, _ := hardware.GetPlatformProfile()
	if profileStr == "" {
		profileStr = "N/A"
	}

	// ── System Info panel ──
	infoLines := []string{
		row("Laptop Model", string(caps.Model), style.ColorCyan),
		row("Driver Type", string(caps.Model)+"Sense", style.ColorPurple),
		row("Thermal Profile", strings.ToUpper(profileStr), style.ColorGreen),
	}
	infoPanel := style.Section("  System", strings.Join(infoLines, "\n"))

	// ── Sensors panel ──
	sensorLines := []string{
		rowBar("CPU Temp", sensors.CPUTemp, 100, "°C"),
		rowBar("GPU Temp", sensors.GPUTemp, 100, "°C"),
		rowBar("SYS Temp", sensors.SYSTemp, 100, "°C"),
		"",
		row("CPU Fan", fmt.Sprintf("%d RPM", sensors.CPUFan), style.ColorCyan),
		row("GPU Fan", fmt.Sprintf("%d RPM", sensors.GPUFan), style.ColorCyan),
	}
	sensorsPanel := style.Section("  Live Sensors", strings.Join(sensorLines, "\n"))

	// ── Quick status panel ──
	var statusLines []string
	if caps.HasBatteryLimiter {
		v, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "battery_limiter"))
		statusLines = append(statusLines, row("Battery Limiter (80%)", "", style.ColorMuted))
		statusLines[len(statusLines)-1] = rowToggle("Battery Limiter (80%)", v == 1)
	}
	if caps.HasBacklightTimeout {
		v, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "backlight_timeout"))
		statusLines = append(statusLines, rowToggle("Backlight Timeout", v == 1))
	}
	if caps.HasLCDOverride {
		v, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "lcd_override"))
		statusLines = append(statusLines, rowToggle("LCD Override", v == 1))
	}
	if caps.HasUSBCharging {
		v, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "usb_charging"))
		statusLines = append(statusLines, row("USB Charge Cutoff", fmt.Sprintf("%d%%", v), style.ColorYellow))
	}
	if caps.HasBootAnimationSound {
		v, _ := sysfs.ReadInt(filepath.Join(caps.SensePath, "boot_animation_sound"))
		statusLines = append(statusLines, rowToggle("Boot Sound", v == 1))
	}
	var statusPanel string
	if len(statusLines) > 0 {
		statusPanel = style.Section("  Quick Status", strings.Join(statusLines, "\n"))
	}

	panels := []string{infoPanel, sensorsPanel}
	if statusPanel != "" {
		panels = append(panels, statusPanel)
	}
	return lipgloss.JoinVertical(lipgloss.Left, panels...)
}

// ── Shared helpers ─────────────────────────────────────────────────────────

func row(label, value string, valueColor lipgloss.Color) string {
	l := style.StyleLabel.Render(label)
	v := lipgloss.NewStyle().Bold(true).Foreground(valueColor).Render(value)
	return l + " " + v
}

func rowToggle(label string, active bool) string {
	l := style.StyleLabel.Render(label)
	return l + " " + style.OnOff(active)
}

func rowBar(label string, val, max int, unit string) string {
	l := style.StyleLabel.Render(label)
	bar := MiniBar(val, max, 16)
	v := lipgloss.NewStyle().Foreground(tempColor(val)).Render(fmt.Sprintf("%d%s", val, unit))
	return l + " " + bar + " " + v
}

func tempColor(t int) lipgloss.Color {
	switch {
	case t >= 85:
		return style.ColorRed
	case t >= 70:
		return style.ColorYellow
	default:
		return style.ColorGreen
	}
}

func MiniBar(val, max, width int) string {
	if val < 0 {
		val = 0
	}
	if val > max {
		val = max
	}
	fill := (val * width) / max
	if fill < 0 {
		fill = 0
	}
	empty := width - fill
	color := tempColor(val)
	filled := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", fill))
	rest := lipgloss.NewStyle().Foreground(style.ColorMuted).Render(strings.Repeat("░", empty))
	return filled + rest
}
