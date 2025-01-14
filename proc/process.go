package proc

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	ProcessesLocation   = "/proc"
	ProcessPathTemplate = ProcessesLocation + "/%s/%s"

	AttrProcessStat = "stat"
)

// Process describe system process
type Process struct {
	PID   string
	Name  string
	Utime int64
	Stime int64
}

// Processes retrieves all processes from system
// with only filled pids. Some options are available
func RetrieveProcesses() ([]Process, error) {
	messedEntries, err := os.ReadDir(ProcessesLocation)
	if err != nil {
		return nil, err
	}
	processes := convAndFilterEntries(messedEntries)
	if len(processes) == 0 {
		return nil, errors.New("processes not found")
	}
	whiteSpaceSep := []byte(" ")
	for i := range len(processes) {
		statPath := fmt.Sprintf(ProcessPathTemplate, processes[i].PID, AttrProcessStat)
		raw, err := os.ReadFile(
			statPath,
		)
		if err != nil {
			return nil, err
		}

		rawProperties := bytes.Split(raw, whiteSpaceSep)
		if len(rawProperties) != 52 {
			return nil, fmt.Errorf("incopatibile stat file: %s", statPath)
		}
		utime, err := strconv.ParseInt(string(rawProperties[13]), 10, 64)
		if err != nil {
			return nil, err
		}
		stime, err := strconv.ParseInt(string(rawProperties[14]), 10, 64)
		if err != nil {
			return nil, err
		}
		processes[i].Utime = utime
		processes[i].Stime = stime
		// cutting braces from name
		processes[i].Name = string(rawProperties[1][1 : len(rawProperties[1])-1])
	}

	return processes, nil

}

func convAndFilterEntries(entries []os.DirEntry) []Process {
	result := make([]Process, 0, len(entries))

	for _, v := range entries {
		// processess are represented by dirs
		if !v.IsDir() {
			continue
		}
		fName := v.Name()

		// the name of dirs should be a number which is actually id
		_, err := strconv.ParseInt(fName, 10, 32)
		if err != nil {
			continue
		}
		result = append(
			result,
			Process{
				PID: fName,
			})
	}
	return result
}
