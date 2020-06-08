// Package comm provides common methods.
package comm

import (
	"os"
	"os/exec"
	"path/filepath"
)

// IsWindows returns whether current os is windows.
func IsWindows() bool {
	return isWindows
}

// LineSeparator returns line separator of current os.
func LineSeparator() string {
	return lineSeparator
}

// ExecutablePath returns executable file path.
func ExecutablePath() (executablePath string, err error) {
	executablePath, err = exec.LookPath(os.Args[0])
	if err != nil {
		return
	}
	executablePath, err = filepath.Abs(executablePath)
	return
}

// StartupPath returns startup directory path.
func StartupPath() (string, error) {
	return filepath.Abs(".")
}
