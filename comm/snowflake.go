package comm

import (
	"errors"
	"sync"
	"time"
)

var (
	errOverflowedNodeID = errors.New("overflowed node id")
)

// NewSnowflake creates an instance of snowflake and returns it.
func NewSnowflake(nodeIDBits, sequenceIDBits uint8, epoch time.Time, nodeID int64) IDGenerator {
	if nodeID>>nodeIDBits > 0 {
		panic(errOverflowedNodeID)
	}
	return &snowflake{
		nodeIDBits:     nodeIDBits,
		sequenceIDBits: sequenceIDBits,
		epoch:          TimeMilli(epoch),
	}
}

// snowflake is a service used to generate unique IDs for objects within Twitter. This is an implementation.
type snowflake struct {
	lock           sync.Mutex
	nodeIDBits     uint8
	sequenceIDBits uint8
	epoch          int64
	nodeID         int64
	timestamp      int64
	sequenceID     int64
}

// NewID creates an unique ID.
func (p *snowflake) NewID() int64 {
	timestamp, sequenceID := p.next()
	id := timestamp - p.epoch
	id <<= p.nodeIDBits
	id |= p.nodeID
	id <<= p.sequenceIDBits
	id |= sequenceID
	return id
}

func (p *snowflake) next() (timestamp int64, sequenceID int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	timestamp = UnixMilli()
	if timestamp != p.timestamp {
		sequenceID = 0
	} else {
		sequenceID = p.sequenceID + 1
		if sequenceID>>p.sequenceIDBits > 0 {
			time.Sleep(time.Millisecond)
			timestamp++
			sequenceID = 0
		}
	}
	p.timestamp = timestamp
	p.sequenceID = sequenceID
	return
}
