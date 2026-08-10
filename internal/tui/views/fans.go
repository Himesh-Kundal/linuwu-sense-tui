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
		return lipgloss.NewStyle().Foreground(style.ColorWarning).Render("Fan speed control is not supported on this device.")
	}

	cpuBar := renderBar(m.CPUSpeed, 100, 20)
	gpuBar := renderBar(m.GPUSpeed, 100, 20)

	cpuMode := fmt.Sprintf("%d%%", m.CPUSpeed)
	if m.CPUSpeed == 0 {
		cpuMode = "Auto"
	}
	gpuMode := fmt.Sprintf("%d%%", m.GPUSpeed)
	if m.GPUSpeed == 0 {
		gpuMode = "Auto"
	}

	content := fmt.Sprintf(
		"CPU Fan Setting: %s %s\nGPU Fan Setting: %s %s\n\nLive Tachometer:\nCPU: %d RPM\nGPU: %d RPM",
		cpuBar, style.StyleValue.Render(cpuMode),
		gpuBar, style.StyleValue.Render(gpuMode),
		sensors.CPUFan, sensors.GPUFan,
	)

	box := style.MakeBox("Fan Speed Control", content)

	hints := lipgloss.NewStyle().Foreground(style.ColorMuted).Render(
		"Controls:\n" +
			"  [↑ / ↓]   Adjust CPU & GPU Speed (+5% / -5%)\n" +
			"  [A]       Auto Mode (0,0)\n" +
			"  [1]       Quiet Preset (30,30)\n" +
			"  [2]       Balanced Preset (60,60)\n" +
			"  [3]       Performance Preset (80,80)\n" +
			"  [M]       Maximum Speed (100,100)",
	)

	return lipgloss.JoinVertical(lipgloss.Left, box, strings.Repeat("\n", 1), hints)
}
