package comm_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/comm"
)

func Test_GenerateID(t *testing.T) {
	threadCount := 100
	iterationCount := 10000
	idCount := threadCount * iterationCount
	ch := make(chan int64, idCount)
	for i := 0; i < threadCount; i++ {
		go func() {
			for i := 0; i < iterationCount; i++ {
				ch <- comm.GenerateID()
			}
		}()
	}
	idMap := make(map[int64]bool)
	for i := 0; i < idCount; i++ {
		idMap[<-ch] = false
	}
	assert.IntEqual(t, idCount, len(idMap))
}

func Benchmark_GenerateID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		comm.GenerateID()
	}
}
