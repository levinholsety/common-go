package util

import (
	"sync"
	"time"

	"github.com/levinholsety/common-go/timeutil"
)

const (
	nodeIDBitLen           = 10
	sequenceIDBitLen       = 12
	timestampBitLen        = 63 - nodeIDBitLen - sequenceIDBitLen
	epoch            int64 = 1477958400000
)

// IDGenerator can generate unique ID.
type IDGenerator struct {
	nodeID     int64
	sequenceID int64
	timestamp  int64
	lock       sync.Mutex
}

// NewIDGenerator creates an IDGenerator instance with specified node ID.
func NewIDGenerator(nodeID int64) *IDGenerator {
	return &IDGenerator{
		nodeID: nodeID & (1<<nodeIDBitLen - 1),
	}
}

// GenerateID generates a new ID.
func (p *IDGenerator) GenerateID() int64 {
	timestamp, sequenceID := p.next()
	id := timestamp
	id <<= nodeIDBitLen
	id |= p.nodeID
	id <<= sequenceIDBitLen
	id |= sequenceID
	return id
}

func (p *IDGenerator) next() (timestamp int64, sequenceID int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	timestamp = timeutil.CurrentTimeMillis()
	sequenceID = p.sequenceID
	if timestamp != p.timestamp {
		sequenceID = 0
	} else {
		sequenceID++
		if sequenceID>>sequenceIDBitLen > 0 {
			time.Sleep(time.Millisecond)
			timestamp++
			sequenceID = 0
		}
	}
	p.timestamp = timestamp
	p.sequenceID = sequenceID
	return
}
