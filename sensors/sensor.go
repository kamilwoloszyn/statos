package sensors

type SensorKind string

const (
	SensorCPU  SensorKind = "CPU"
	SensorDisk SensorKind = "DISK"
)

type Sensor struct {
	ID                 string
	Kind               string
	CurrentTemperature float32
	MaxTemperature     *float32
}
