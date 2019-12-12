package comm_test

import (
	"fmt"
	"testing"

	"github.com/levinholsety/common-go/comm"
)

func TestCommonPaths(t *testing.T) {
	fmt.Println(comm.StartupPath)
	fmt.Println(comm.ExecutablePath)
	fmt.Println(comm.ExecutableDir)
}
