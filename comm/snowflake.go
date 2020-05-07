package comm

import (
	"log"
	"sync"
	"time"
)

// Snowflake is a service used to generate unique IDs for objects within Twitter. This is an implementation.
type Snowflake struct {
	Epoch          int64
	NodeIDBits     int
	SequenceIDBits int
	lock           sync.Mutex
	timestamp      int64
	sequenceID     int64
}

// NewIDOfNode creates an unique ID of specified node.
func (p *Snowflake) NewIDOfNode(nodeID int64) int64 {
	if nodeID>>p.NodeIDBits > 0 {
		log.Fatalf("Node ID(0 to %d) is too large: %d", 1<<p.NodeIDBits-1, nodeID)
	}
	timestamp, sequenceID := p.next()
	id := timestamp - p.Epoch
	id <<= p.NodeIDBits
	id |= nodeID
	id <<= p.SequenceIDBits
	id |= sequenceID
	return id
}

// NewID creates an unique ID.
func (p *Snowflake) NewID() int64 {
	return p.NewIDOfNode(0)
}

func (p *Snowflake) next() (timestamp int64, sequenceID int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	timestamp = CurrentTimeMillis()
	if timestamp != p.timestamp {
		sequenceID = 0
	} else {
		sequenceID = p.sequenceID + 1
		if sequenceID>>p.SequenceIDBits > 0 {
			time.Sleep(time.Millisecond)
			timestamp++
			sequenceID = 0
		}
	}
	p.timestamp = timestamp
	p.sequenceID = sequenceID
	return
}
