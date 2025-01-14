package cpu

// #include <syscall.h>
import "C"

import (
	"bytes"
	"cmp"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kamilwoloszyn/statos/proc"
)

const (
	pathCpuInfo = "/cpuinfo"
	// not sure if that works on another machine eg. amd
	// TODO: find more reliable way to figure out max cpu freq
	absPathMaxFreq = "/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq"

	// contains already known labels from cpuinfo
	// most of them probably won't be used, however
	// its good to support all for later improvements

	labelProcessor       = "processor"
	labelVendorID        = "vendor_id"
	labelCPUFamily       = "cpu family"
	labelModel           = "model"
	labelModelName       = "model name"
	labelStepping        = "stepping"
	labelMicroCode       = "microcode"
	labelCPUMhz          = "cpu MHz"
	labelCacheSize       = "cache size"
	labelPhysicalID      = "physical id"
	labelSiblings        = "siblings"
	labelCoreID          = "core id"
	labelCPUCores        = "cpu cores"
	labelAPICid          = "apicid"
	labelInitialAPICid   = "initial apicid"
	labelFPU             = "fpu"
	labelFPUException    = "fpu_exception"
	labelCpuidLevel      = "cpuid level"
	labelWP              = "wp"
	labelFlags           = "flags"
	labelVmxFlags        = "vmx flags"
	labelBugs            = "bugs"
	labelBogomips        = "bogomips"
	labelClflushSize     = "clflush size"
	labelCacheAlignment  = "cache_alignment"
	labelAddressSizes    = "address sizes"
	labelPowerManagement = "power management"
)

type CPU struct {
	CurrentHiClock float32
	MaxClock       float32
	NumCores       int
	CLKTCK         int
}

// GetCPUInfo reads informations directly from cpuinfo
// Making more independent from terminal command availability.
func GetCPUInfo() (CPU, error) {
	f, err := os.ReadFile(proc.ProcessesLocation + pathCpuInfo)
	if err != nil {
		return CPU{}, err
	}

	fullSchema, err := parseCpuInfo(f)
	if err != nil {
		return CPU{}, fmt.Errorf("parse %s: %v", proc.ProcessesLocation+pathCpuInfo, err)
	}

	clktck := C.auxval()

	// we don't care if the file doesn't exist
	// that's optional
	maxClockRaw, err := os.ReadFile(absPathMaxFreq)
	var maxClock float32
	if err == nil {
		maxFreq, err := strconv.ParseFloat(
			strings.TrimSpace(string(maxClockRaw)), 32,
		)
		if err == nil {
			maxClock = float32(maxFreq)
		}
	}
	var coreClocks []float32 = make([]float32, 0, 8)
	for _, v := range fullSchema {
		coreClocks = append(coreClocks, v.CpuMHZ)
	}

	return CPU{
		CurrentHiClock: Max(coreClocks...),
		CLKTCK:         int(clktck),
		MaxClock:       maxClock,
		NumCores:       len(coreClocks) - 1,
	}, nil
}

type CPUInfoSchema struct {
	Processor       string
	VendorID        string
	CPUFamily       string
	Model           string
	ModelName       string
	Stepping        string
	Microcode       string
	CpuMHZ          float32
	CacheSize       string
	PhysicalID      string
	Siblings        string
	CoreID          string
	CPUCores        string
	APICid          string
	InitialAPICid   string
	Fpu             string
	FpuExtenstion   string
	CPUIDLevel      string
	WP              string
	Flags           string
	VmxFlags        string
	Bugs            string
	Bogomips        string
	ClflushSize     string
	CacheAlignment  string
	AddressSizes    string
	PowerManagement string
}

func (s CPUInfoSchema) Validate() error {
	return nil
}

func parseCpuInfo(rawData []byte) ([]CPUInfoSchema, error) {
	var result []CPUInfoSchema
	var currentItem CPUInfoSchema
	delimiter := []byte("\n")
	rows := bytes.Split(rawData, delimiter)
	for _, row := range rows {
		if len(row) == 0 {
			// assumed that new row occurs here
			result = append(result, currentItem)
			currentItem = CPUInfoSchema{}
			continue
		}
		delimiter = []byte(":")
		kv := bytes.Split(row, delimiter)

		keyStr := string(kv[0])
		var valueStr string
		if len(kv) > 1 {
			valueStr = string(kv[1])
		}

		switch {
		case strings.Contains(keyStr, labelProcessor):
			currentItem.Processor = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelVendorID):
			currentItem.VendorID = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelCPUFamily):
			currentItem.CPUFamily = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelModel):
			currentItem.Model = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelModelName):
			currentItem.ModelName = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelStepping):
			currentItem.Stepping = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelMicroCode):
			currentItem.Microcode = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelCPUMhz):
			v, err := strconv.ParseFloat(
				strings.TrimSpace(valueStr), 32)
			if err != nil {
				return nil, err
			}
			currentItem.CpuMHZ = float32(v)
		case strings.Contains(keyStr, labelCacheSize):
			currentItem.CacheSize = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelPhysicalID):
			currentItem.PhysicalID = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelSiblings):
			currentItem.Siblings = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelCoreID):
			currentItem.CoreID = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelCPUCores):
			currentItem.CPUCores = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelAPICid):
			currentItem.APICid = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelInitialAPICid):
			currentItem.InitialAPICid = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelFPU):
			currentItem.Fpu = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelFPUException):
			currentItem.FpuExtenstion = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelCpuidLevel):
			currentItem.CPUIDLevel = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelWP):
			currentItem.WP = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelFlags):
			currentItem.Flags = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelVmxFlags):
			currentItem.VmxFlags = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelBugs):
			currentItem.Bugs = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelBogomips):
			currentItem.Bogomips = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelClflushSize):
			currentItem.ClflushSize = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelCacheAlignment):
			currentItem.CacheAlignment = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelAddressSizes):
			currentItem.AddressSizes = strings.TrimSpace(valueStr)
		case strings.Contains(keyStr, labelPowerManagement):
			currentItem.PowerManagement = strings.TrimSpace(valueStr)
		}
	}
	return result, nil
}

func Max[T cmp.Ordered](args ...T) T {
	if len(args) == 0 {
		return *new(T)
	}
	max := args[0]
	for _, v := range args[1:] {
		if max < v {
			max = v
		}
	}

	return max
}
