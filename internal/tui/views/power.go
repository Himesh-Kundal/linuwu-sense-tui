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

type powerItem struct {
	label string
	desc  string
	value string
}

func RenderPower(caps hardware.Capabilities, cursor int) string {
	if !caps.ModuleLoaded {
		return style.StyleDanger.Render("  ⚠  Module not loaded.")
	}

	items := buildPowerItems(caps)
	if len(items) == 0 {
		return style.StyleMuted.Render("  No power features detected on this device.")
	}

	var lines []string
	for i, item := range items {
		label := fmt.Sprintf("  %-30s %s", item.label, item.value)
		hint := style.StyleMuted.Render("  " + item.desc)

		if i == cursor {
			lines = append(lines, style.StyleRowSelected.Render(fmt.Sprintf("▶ %-30s %s", item.label, item.value)))
			lines = append(lines, style.StyleRowHint.Render("  └─ "+item.desc))
		} else {
			lines = append(lines, style.StyleRowNormal.Render(label))
			_ = hint
		}
		lines = append(lines, "")
	}

	panel := style.SectionFocused("  Power & Display Settings", strings.Join(lines, "\n"))
	help := style.KeyHint("↑/↓", "Navigate") + "   " + style.KeyHint("Space/Enter", "Toggle / Cycle")

	return lipgloss.JoinVertical(lipgloss.Left, panel, "", help)
}

func buildPowerItems(caps hardware.Capabilities) []powerItem {
	var items []powerItem
	sp := caps.SensePath

	if caps.HasBatteryLimiter {
		v, _ := sysfs.ReadInt(filepath.Join(sp, "battery_limiter"))
		items = append(items, powerItem{
			label: "Battery Limiter (80% cap)",
			desc:  "Caps charge at 80% to prolong battery lifespan",
			value: style.OnOff(v == 1),
		})
	}
	if caps.HasBatteryCalibration {
		v, _ := sysfs.ReadInt(filepath.Join(sp, "battery_calibration"))
		items = append(items, powerItem{
			label: "Battery Calibration",
			desc:  "100% → 0% → 100% cycle. Keep AC connected!",
			value: style.OnOff(v == 1),
		})
	}
	if caps.HasUSBCharging {
		v, _ := sysfs.ReadInt(filepath.Join(sp, "usb_charging"))
		items = append(items, powerItem{
			label: "USB Charge Cutoff",
			desc:  "Cycles: 0% → 10% → 20% → 30% (Press Enter)",
			value: lipgloss.NewStyle().Bold(true).Foreground(style.ColorYellow).Render(fmt.Sprintf("%d%%", v)),
		})
	}
	if caps.HasBacklightTimeout {
		v, _ := sysfs.ReadInt(filepath.Join(sp, "backlight_timeout"))
		items = append(items, powerItem{
			label: "Backlight Timeout (30s)",
			desc:  "Turns off keyboard backlight after 30s of inactivity",
			value: style.OnOff(v == 1),
		})
	}
	if caps.HasLCDOverride {
		v, _ := sysfs.ReadInt(filepath.Join(sp, "lcd_override"))
		items = append(items, powerItem{
			label: "LCD Override",
			desc:  "Reduces LCD input latency and ghosting artifacts",
			value: style.OnOff(v == 1),
		})
	}
	if caps.HasBootAnimationSound {
		v, _ := sysfs.ReadInt(filepath.Join(sp, "boot_animation_sound"))
		items = append(items, powerItem{
			label: "Boot Animation Sound",
			desc:  "Plays Predator boot animation sound on startup",
			value: style.OnOff(v == 1),
		})
	}
	return items
}
