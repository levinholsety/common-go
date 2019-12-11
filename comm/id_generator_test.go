package comm_test

import (
	"testing"

	"github.com/levinholsety/common-go/comm"
)

func Test_GenerateID(t *testing.T) {
	iterationCount := 10000
	threadCount := 100
	idCount := iterationCount * threadCount
	ch := make(chan int64, idCount)
	for i := 0; i < threadCount; i++ {
		go func(index int) {
			for i := 0; i < iterationCount; i++ {
				ch <- comm.GenerateID()
			}
		}(i)
	}
	idMap := make(map[int64]bool)
	for i := 0; i < idCount; i++ {
		idMap[<-ch] = false
	}
	close(ch)
	actualIDCount := len(idMap)
	if actualIDCount != idCount {
		t.Errorf("expected: %d, actrual: %d", idCount, actualIDCount)
	}
}
