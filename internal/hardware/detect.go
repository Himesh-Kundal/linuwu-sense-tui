package hardware

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/sysfs"
)

const BaseSysfs = "/sys/module/linuwu_sense/drivers/platform:acer-wmi/acer-wmi"

func Detect() Capabilities {
	var caps Capabilities

	if !sysfs.Exists("/sys/module/linuwu_sense") {
		caps.ModuleLoaded = false
		return caps
	}
	caps.ModuleLoaded = true

	predPath := filepath.Join(BaseSysfs, "predator_sense")
	nitroPath := filepath.Join(BaseSysfs, "nitro_sense")

	if sysfs.Exists(predPath) {
		caps.Model = ModelPredator
		caps.SensePath = predPath
	} else if sysfs.Exists(nitroPath) {
		caps.Model = ModelNitro
		caps.SensePath = nitroPath
	} else {
		caps.Model = ModelUnknown
	}

	if caps.Model != ModelUnknown {
		caps.HasBacklightTimeout = sysfs.Exists(filepath.Join(caps.SensePath, "backlight_timeout"))
		caps.HasBatteryCalibration = sysfs.Exists(filepath.Join(caps.SensePath, "battery_calibration"))
		caps.HasBatteryLimiter = sysfs.Exists(filepath.Join(caps.SensePath, "battery_limiter"))
		caps.HasBootAnimationSound = sysfs.Exists(filepath.Join(caps.SensePath, "boot_animation_sound"))
		caps.HasFanSpeed = sysfs.Exists(filepath.Join(caps.SensePath, "fan_speed"))
		caps.HasLCDOverride = sysfs.Exists(filepath.Join(caps.SensePath, "lcd_override"))
		caps.HasUSBCharging = sysfs.Exists(filepath.Join(caps.SensePath, "usb_charging"))
	}

	kbPath := filepath.Join(BaseSysfs, "four_zoned_kb")
	if sysfs.Exists(kbPath) {
		caps.HasFourZonedKB = true
		caps.KBPath = kbPath
	}

	caps.HasPlatformProfile = sysfs.Exists("/sys/firmware/acpi/platform_profile")
	caps.HwmonPath = findHwmonPath()

	return caps
}

func findHwmonPath() string {
	entries, err := os.ReadDir("/sys/class/hwmon")
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		namePath := filepath.Join("/sys/class/hwmon", entry.Name(), "name")
		name, err := sysfs.ReadString(namePath)
		if err == nil && name == "acer" {
			return filepath.Join("/sys/class/hwmon", entry.Name())
		}
	}
	return ""
}

func ReadSensors(caps Capabilities) SensorData {
	var data SensorData
	if caps.HwmonPath == "" {
		return data
	}

	if t1, err := sysfs.ReadInt(filepath.Join(caps.HwmonPath, "temp1_input")); err == nil {
		data.CPUTemp = t1 / 1000
	}
	if t2, err := sysfs.ReadInt(filepath.Join(caps.HwmonPath, "temp2_input")); err == nil {
		data.GPUTemp = t2 / 1000
	}
	if t3, err := sysfs.ReadInt(filepath.Join(caps.HwmonPath, "temp3_input")); err == nil {
		data.SYSTemp = t3 / 1000
	}
	if f1, err := sysfs.ReadInt(filepath.Join(caps.HwmonPath, "fan1_input")); err == nil {
		data.CPUFan = f1
	}
	if f2, err := sysfs.ReadInt(filepath.Join(caps.HwmonPath, "fan2_input")); err == nil {
		data.GPUFan = f2
	}
	return data
}

func GetFanSpeed(caps Capabilities) (int, int, error) {
	if !caps.HasFanSpeed {
		return 0, 0, fmt.Errorf("fan speed not supported")
	}
	val, err := sysfs.ReadString(filepath.Join(caps.SensePath, "fan_speed"))
	if err != nil {
		return 0, 0, err
	}
	var cpu, gpu int
	fmt.Sscanf(val, "%d,%d", &cpu, &gpu)
	return cpu, gpu, nil
}

func SetFanSpeed(caps Capabilities, cpu, gpu int) error {
	if !caps.HasFanSpeed {
		return fmt.Errorf("fan speed not supported")
	}
	val := fmt.Sprintf("%d,%d", cpu, gpu)
	return sysfs.WriteString(filepath.Join(caps.SensePath, "fan_speed"), val)
}

func GetPlatformProfile() (string, error) {
	return sysfs.ReadString("/sys/firmware/acpi/platform_profile")
}

func GetPlatformProfileChoices() ([]string, error) {
	s, err := sysfs.ReadString("/sys/firmware/acpi/platform_profile_choices")
	if err != nil {
		return nil, err
	}
	return strings.Fields(s), nil
}

func SetPlatformProfile(profile string) error {
	path := "/sys/firmware/acpi/platform_profile"
	err := sysfs.WriteString(path, profile)
	if err != nil {
		// Fallback for unprivileged process: attempt to write via sudo tee
		cmd := exec.Command("sudo", "tee", path)
		cmd.Stdin = strings.NewReader(profile + "\n")
		return cmd.Run()
	}
	return nil
}
