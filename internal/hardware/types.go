package hardware

type ModelType string

const (
	ModelUnknown  ModelType = "Unknown"
	ModelPredator ModelType = "Predator"
	ModelNitro    ModelType = "Nitro"
)

type Capabilities struct {
	ModuleLoaded           bool
	Model                  ModelType
	SensePath              string
	KBPath                 string
	HwmonPath              string
	HasBacklightTimeout    bool
	HasBatteryCalibration  bool
	HasBatteryLimiter      bool
	HasBootAnimationSound  bool
	HasFanSpeed            bool
	HasLCDOverride         bool
	HasUSBCharging         bool
	HasFourZonedKB         bool
	HasPlatformProfile     bool
}

type SensorData struct {
	CPUTemp int // in Celsius
	GPUTemp int // in Celsius
	SYSTemp int // in Celsius
	CPUFan  int // RPM
	GPUFan  int // RPM
}
