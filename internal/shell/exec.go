package shell

import (
	"bytes"
	"os/exec"
	"sync"
)

//ShellToUse bash shell
const ShellToUse = "bash"

var shellExec *CommandExec
var shellExecSync sync.Once

//Executor defines the interface for shell command executor
//exec.go
//go:generate mockgen -destination=../mocks/mock_Executor.go -package=mocks . Executor
type Executor interface {
	Exec(command string) (*CommandResult, error)
}

//CommandExec object
type CommandExec struct {
}

//NewShellExec return new instance of shell executor
func NewShellExec() Executor {
	shellExecSync.Do(func() {
		shellExec = &CommandExec{}
	})
	return shellExec
}

//CommandResult return data
type CommandResult struct {
	Stdout string
	Stderr string
}

//Exec execute shell command
// #nosec
func (ce CommandExec) Exec(command string) (*CommandResult, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return &CommandResult{Stdout: stdout.String(), Stderr: stderr.String()}, err
}
