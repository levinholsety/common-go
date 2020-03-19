// +build !windows,!darwin

package comm

const (
	// IsWindows represents whether current os is windows.
	IsWindows = false
	// LineSeparator represents line separator of current os.
	LineSeparator = "\n"
)
