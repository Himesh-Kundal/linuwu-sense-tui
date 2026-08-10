package views

import (
	"fmt"
	"strings"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

type FansModel struct {
	CPUSpeed int
	GPUSpeed int
	Loaded   bool
}

func (m *FansModel) Init(caps hardware.Capabilities) {
	if !m.Loaded && caps.HasFanSpeed {
		cpu, gpu, err := hardware.GetFanSpeed(caps)
		if err == nil {
			m.CPUSpeed = cpu
			m.GPUSpeed = gpu
		}
		m.Loaded = true
	}
}

func RenderFans(m *FansModel, caps hardware.Capabilities, sensors hardware.SensorData) string {
	if !caps.HasFanSpeed {
		return style.StyleWarning.Render("  Fan speed control is not supported on this device.")
	}

	cpuLabel := modeLabel(m.CPUSpeed)
	gpuLabel := modeLabel(m.GPUSpeed)

	lines := []string{
		fmtFan("CPU Fan", m.CPUSpeed, cpuLabel, sensors.CPUFan),
		"",
		fmtFan("GPU Fan", m.GPUSpeed, gpuLabel, sensors.GPUFan),
	}

	panel := style.SectionFocused("  Fan Speed Control", strings.Join(lines, "\n"))

	help := lipgloss.JoinHorizontal(lipgloss.Left,
		style.KeyHint("←/→", "Speed ±5%")+"   ",
		style.KeyHint("0", "Auto (0%)")+"   ",
		style.KeyHint("9", "Max (100%)"),
	)

	return lipgloss.JoinVertical(lipgloss.Left, panel, "", help)
}

func modeLabel(spd int) string {
	if spd == 0 {
		return "Auto"
	}
	return fmt.Sprintf("%d%%", spd)
}

func fmtFan(name string, speed int, label string, rpm int) string {
	bar := MiniBar(speed, 100, 24)
	pct := lipgloss.NewStyle().Bold(true).Foreground(style.ColorCyan).Render(fmt.Sprintf("%-8s", label))
	live := style.StyleMuted.Render(fmt.Sprintf("  %d RPM", rpm))
	return style.StyleLabel.Render(name) + " " + bar + " " + pct + live
}

func renderBar(val, max, width int) string {
	return MiniBar(val, max, width)
}

// keep symbol for compat with old dashboard if needed
var _ = strings.Repeat
