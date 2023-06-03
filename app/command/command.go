package command

import (
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type CommandExecutor struct {
}

func (c *CommandExecutor) Execute(cmd string, args ...string) ([]byte, error) {
	readyToExec := exec.Command(cmd, args...)
	return readyToExec.Output()
}

func (c *CommandExecutor) ExecuteWithPipe(cmd string, args ...string) (io.ReadCloser, error) {
	readyToexec := exec.Command(cmd, args...)
	outPipe, err := readyToexec.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "ExecuteWithPipe")
	}
	return outPipe, readyToexec.Start()
}
