package commands

import (
	"bytes"
	"os/exec"
	"sync"
)

const ShellToUse = "bash"

var shellExec *ShellCommand
var shellExecSync sync.Once

//ShellCommand object
type ShellCommand struct {
}

//NewShellExec return new instance of shell executor
func NewShellExec() *ShellCommand {
	shellExecSync.Do(func() {
		shellExec = &ShellCommand{}
	})
	return shellExec
}

//ShellCommandResult return data
type ShellCommandResult struct {
	Stdout string
	Stderr string
}

//Exec execute shell command
func (ce ShellCommand) Exec(command string) (*ShellCommandResult,error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return &ShellCommandResult{Stdout: stdout.String(), Stderr: stderr.String()},err
}
