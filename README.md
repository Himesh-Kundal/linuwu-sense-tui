# 🐾 linuwu-tui

A modern, fast, and standalone Terminal UI (TUI) and control utility for Acer Predator & Nitro laptops using the [Linuwu-Sense](https://github.com/0x7375646F/Linuwu-Sense) Linux kernel module.

Built in **Go** using [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and [Cobra](https://github.com/spf13/cobra).

---

## ⚡ Features

- **Embedded Kernel Driver**: Built-in driver source embedded via `//go:embed`. Install both driver and TUI seamlessly with `go install` and `linuwu-tui setup`.
- **Live Thermal & Fan Monitoring**: Real-time CPU, GPU, and System temperature readings and fan tachometer speeds.
- **Thermal Profiles**: One-key switching between `quiet`, `balanced`, `performance`, and `turbo`.
- **Fan Control**: Set custom fan speeds or quick presets (`auto`, `quiet`, `balanced`, `performance`, `max`).
- **Power & Battery Health**: Toggle battery charging limiter (80%), run battery calibration mode, and set USB power-off charging threshold.
- **4-Zone RGB Keyboard**: Full control over lighting effects, speed, brightness, and direction.

---

## 🚀 Quick Start

### 1. Installation via Go

```bash
go install github.com/Himesh-Kundal/linuwu-tui@latest
```

### 2. Driver Setup

If you haven't loaded/installed the `linuwu_sense` kernel module yet, run the embedded setup command:

```bash
linuwu-tui setup
```

### 3. Launch TUI

```bash
linuwu-tui
```

---

## 💻 Command Line Usage

`linuwu-tui` also supports headless/scripting commands:

```bash
# Print one-shot hardware status dump
linuwu-tui status

# Extract and build kernel module
linuwu-tui setup
```

---

## 📄 License

GPL-3.0 (matches Linuwu-Sense kernel module).
