package comm

import (
	"io"
	"os/exec"
)

// ExecutePiplineCommands executes commands in pipline mode.
func ExecutePiplineCommands(w io.Writer, commands ...*exec.Cmd) (err error) {
	if len(commands) == 0 {
		return
	}
	src := commands[0]
	for _, dst := range commands[1:] {
		dst.Stdin, src.Stdout = io.Pipe()
		src = dst
	}
	src.Stdout = w
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
	return nil
}
