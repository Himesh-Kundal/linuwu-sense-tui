package tui

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/sysfs"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg time.Time

type Model struct {
	ActiveTab int
	Caps      hardware.Capabilities
	Sensors   hardware.SensorData
	Fans      views.FansModel
}

func NewModel() Model {
	caps := hardware.Detect()
	sensors := hardware.ReadSensors(caps)
	var fans views.FansModel
	fans.Init(caps)

	return Model{
		ActiveTab: 0,
		Caps:      caps,
		Sensors:   sensors,
		Fans:      fans,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.Sensors = hardware.ReadSensors(m.Caps)
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.ActiveTab = (m.ActiveTab + 1) % 5
		case "shift+tab":
			m.ActiveTab = (m.ActiveTab + 4) % 5
		case "f1":
			m.ActiveTab = 0
		case "f2":
			m.ActiveTab = 1
		case "f3":
			m.ActiveTab = 2
		case "f4":
			m.ActiveTab = 3
		case "f5":
			m.ActiveTab = 4

		case "1":
			m.ActiveTab = 0
		case "2":
			m.ActiveTab = 1
		case "3":
			m.ActiveTab = 2
		case "4":
			m.ActiveTab = 3
		case "5":
			m.ActiveTab = 4

		case "a", "A":
			if m.ActiveTab == 1 && m.Caps.HasFanSpeed {
				m.Fans.CPUSpeed, m.Fans.GPUSpeed = 0, 0
				hardware.SetFanSpeed(m.Caps, 0, 0)
			} else if m.ActiveTab == 4 && m.Caps.HasPlatformProfile {
				hardware.SetPlatformProfile("quiet")
			}
		case "b", "B":
			if m.ActiveTab == 3 && m.Caps.HasFourZonedKB {
				m.updateKBParam(2, 101)
			} else if m.ActiveTab == 4 && m.Caps.HasPlatformProfile {
				hardware.SetPlatformProfile("balanced")
			}
		case "c", "C":
			if m.ActiveTab == 4 && m.Caps.HasPlatformProfile {
				hardware.SetPlatformProfile("performance")
			}
		case "d", "D":
			if m.ActiveTab == 3 && m.Caps.HasFourZonedKB {
				m.updateKBParam(3, 2)
			} else if m.ActiveTab == 4 && m.Caps.HasPlatformProfile {
				hardware.SetPlatformProfile("turbo")
			}

		// Fan & Keyboard controls
		case "up":
			if m.ActiveTab == 1 && m.Caps.HasFanSpeed {
				m.Fans.CPUSpeed = clamp(m.Fans.CPUSpeed+5, 0, 100)
				m.Fans.GPUSpeed = clamp(m.Fans.GPUSpeed+5, 0, 100)
				hardware.SetFanSpeed(m.Caps, m.Fans.CPUSpeed, m.Fans.GPUSpeed)
			}
		case "down":
			if m.ActiveTab == 1 && m.Caps.HasFanSpeed {
				m.Fans.CPUSpeed = clamp(m.Fans.CPUSpeed-5, 0, 100)
				m.Fans.GPUSpeed = clamp(m.Fans.GPUSpeed-5, 0, 100)
				hardware.SetFanSpeed(m.Caps, m.Fans.CPUSpeed, m.Fans.GPUSpeed)
			}
		case "m", "M":
			if m.ActiveTab == 1 && m.Caps.HasFanSpeed {
				m.Fans.CPUSpeed, m.Fans.GPUSpeed = 100, 100
				hardware.SetFanSpeed(m.Caps, 100, 100)
			} else if m.ActiveTab == 3 && m.Caps.HasFourZonedKB {
				m.updateKBParam(0, 8)
			}
		case "s", "S":
			if m.ActiveTab == 3 && m.Caps.HasFourZonedKB {
				m.updateKBParam(1, 10)
			}
		}
	}
	return m, nil
}

func (m *Model) toggleBattery(attr string) {
	if !m.Caps.ModuleLoaded {
		return
	}
	p := filepath.Join(m.Caps.SensePath, attr)
	val, _ := sysfs.ReadInt(p)
	next := 0
	if val == 0 {
		next = 1
	}
	sysfs.WriteInt(p, next)
}

func (m *Model) cycleUSBCharging() {
	if !m.Caps.HasUSBCharging {
		return
	}
	p := filepath.Join(m.Caps.SensePath, "usb_charging")
	val, _ := sysfs.ReadInt(p)
	next := 10
	if val == 10 {
		next = 20
	} else if val == 20 {
		next = 30
	} else if val == 30 {
		next = 0
	}
	sysfs.WriteInt(p, next)
}

func (m *Model) updateKBParam(idx int, max int) {
	if !m.Caps.HasFourZonedKB {
		return
	}
	p := filepath.Join(m.Caps.KBPath, "four_zone_mode")
	val, err := sysfs.ReadString(p)
	if err != nil {
		return
	}
	var params [7]int
	fmt.Sscanf(val, "%d,%d,%d,%d,%d,%d,%d", &params[0], &params[1], &params[2], &params[3], &params[4], &params[5], &params[6])

	if idx == 3 {
		if params[3] == 1 {
			params[3] = 2
		} else {
			params[3] = 1
		}
	} else if idx == 2 {
		params[2] = (params[2] + 25) % 125
		if params[2] > 100 {
			params[2] = 0
		}
	} else {
		params[idx] = (params[idx] + 1) % max
	}

	newVal := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d", params[0], params[1], params[2], params[3], params[4], params[5], params[6])
	sysfs.WriteString(p, newVal)
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (m Model) View() string {
	headerTitle := style.StyleHeaderTitle.Render("🐾 LINUWU-SENSE TUI")
	headerModel := style.StyleHeaderModel.Render(fmt.Sprintf("%s  ", m.Caps.Model))
	header := lipgloss.JoinHorizontal(lipgloss.Left, headerTitle, lipgloss.NewStyle().Width(40).Render(""), headerModel)

	tabsList := []string{"[1] Dashboard", "[2] Fans", "[3] Battery", "[4] Keyboard", "[5] Profiles"}
	var renderedTabs []string
	for i, t := range tabsList {
		if i == m.ActiveTab {
			renderedTabs = append(renderedTabs, style.StyleTabActive.Render(t))
		} else {
			renderedTabs = append(renderedTabs, style.StyleTabInactive.Render(t))
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...)

	var body string
	switch m.ActiveTab {
	case 0:
		body = views.RenderDashboard(m.Caps, m.Sensors)
	case 1:
		body = views.RenderFans(&m.Fans, m.Caps, m.Sensors)
	case 2:
		body = views.RenderBattery(m.Caps)
	case 3:
		body = views.RenderKeyboard(m.Caps)
	case 4:
		body = views.RenderProfile(m.Caps)
	}

	footer := style.StyleStatusKey.Render("\n[Tab / F1-F5]: Switch View | [q]: Quit")

	return lipgloss.JoinVertical(lipgloss.Left, header, tabBar, "\n", body, footer)
}
