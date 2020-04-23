package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// NewCommand creates a Command.
func NewCommand(path string) *Command {
	return (&Command{Path: path}).ResetEnvVars()
}

// Command provides functions to execute command.
type Command struct {
	Path                 string
	EnvironmentVariables []string
}

// AddEnvVar adds a environment variable to current command.
func (p *Command) AddEnvVar(name, value string) *Command {
	p.EnvironmentVariables = append(p.EnvironmentVariables, fmt.Sprintf("%s=%s", name, value))
	return p
}

// ResetEnvVars resets environment variables of current command to system defaults.
func (p *Command) ResetEnvVars() *Command {
	p.EnvironmentVariables = os.Environ()
	return p
}

// ClearEnvVars clears environment variables of current command.
func (p *Command) ClearEnvVars() *Command {
	p.EnvironmentVariables = []string{}
	return p
}

// Execute executes command and redirect command output and error to standard output and error.
func (p *Command) Execute(args ...string) error {
	cmd := exec.Command(p.Path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecuteOutput executes command and returns output combined with command output and error.
func (p *Command) ExecuteOutput(args ...string) ([]byte, error) {
	return exec.Command(p.Path, args...).Output()
}
