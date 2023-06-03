package system

import (
	"strconv"

	"github.com/kamilwoloszyn/statos/app"
	"github.com/kamilwoloszyn/statos/domain"
	"github.com/sirupsen/logrus"
)

type System struct {
	commandExecutor app.Commander
	l               logrus.FieldLogger
}

func NewSystemAdapter(cmdExecutor app.Commander) *System {
	return &System{
		commandExecutor: cmdExecutor,
	}
}

func (s *System) Stats() domain.SystemStats {
	return domain.SystemStats{
		MemoryUsage: s.memoryStats(),
	}
}

func (s *System) memoryStats() domain.Memory {
	memoryTotal, err := s.commandExecutor.Execute("free -m |grep Mem | awk '{print $2}'")
	if err != nil {
		s.l.Errorf("executor error occurred: %w", err)
	}
	if err == nil && len(memoryTotal) == 0 {
		s.l.Error("got empty memory total")
	}

	memoryUsed, err := s.commandExecutor.Execute("free -m |grep Mem | awk '{print $3}'")
	if err != nil {
		s.l.Error("executor error occurred")
	}
	if err == nil && len(memoryUsed) == 0 {
		s.l.Errorf("got empty memory usage")
	}

	memTotal, err := strconv.Atoi(string(memoryTotal))
	memUsed, err := strconv.Atoi(string(memoryUsed))

	// memory available can be calculated,there is no need to run another command
	memAvailable := memTotal - memUsed

	return domain.Memory{
		MemoryTotal:     memTotal,
		MemoryAvailable: memAvailable,
		MemoryUsed:      memUsed,
	}
}
