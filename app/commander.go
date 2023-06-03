package app

import "io"

// Commander is a generic interface that executes commands
type Commander interface {
	Execute(cmd string, args ...string) ([]byte, error)
	ExecuteWithPipe(cmd string, args ...string) (io.ReadCloser, error)
}
