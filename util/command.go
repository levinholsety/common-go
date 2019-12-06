package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ExecutePiplineCommands executes commands in pipline mode.
func ExecutePiplineCommands(output io.Writer, commands ...*exec.Cmd) (err error) {
	if len(commands) > 0 {
		src := commands[0]
		for _, dst := range commands[1:] {
			dst.Stdin, src.Stdout = io.Pipe()
			src = dst
		}
		src.Stdout = output
		for _, cmd := range commands {
			if err = cmd.Start(); err != nil {
				return
			}
		}
		for _, cmd := range commands {
			if err = cmd.Wait(); err != nil {
				return
			}
			if obj, ok := cmd.Stdin.(io.Closer); ok {
				obj.Close()
			}
			if obj, ok := cmd.Stdout.(io.Closer); ok {
				obj.Close()
			}
		}
	}
	return nil
}

// Command provides functions to execute command.
type Command struct {
	Path                 string
	EnvironmentVariables []string
}

// NewCommand creates a Command.
func NewCommand(path string) *Command {
	return (&Command{Path: path}).ResetEnvVars()
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
	return cmd.Run()
	// console.RedirectCommandOutput(cmd)
	// console.RedirectCommandError(cmd)
	// return cmd.Run()
}

// ExecuteOutput executes command and returns output combined with command output and error.
func (p *Command) ExecuteOutput(args ...string) ([]byte, error) {
	return exec.Command(p.Path, args...).Output()
}
