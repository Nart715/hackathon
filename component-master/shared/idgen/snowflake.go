package idgen

import (
	"fmt"
	"sync"
	"time"
)

const (
	nodeBits     = 10                      // Number of bits allocated for Node ID
	sequenceBits = 12                      // Number of bits allocated for Sequence
	maxNodeID    = (1 << nodeBits) - 1     // Maximum value for Node ID
	maxSequence  = (1 << sequenceBits) - 1 // Maximum value for Sequence
	timeShift    = nodeBits + sequenceBits // Time shift for generating the ID
	nodeShift    = sequenceBits            // Node shift for generating the ID
	epoch        = int64(1609459200000)    // Epoch: 01/01/2021 00:00:00 UTC
)

type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	nodeID    int64
	sequence  int64
}

// NewSnowflake initializes a new Snowflake generator
func NewSnowflake(nodeID int64) (*Snowflake, error) {
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, fmt.Errorf("node ID must be between 0 and %d", maxNodeID)
	}
	return &Snowflake{
		timestamp: 0,
		nodeID:    nodeID,
		sequence:  0,
	}, nil
}

// Generate generates a unique ID
func (sf *Snowflake) Generate() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < sf.timestamp {
		panic("Clock moved backwards!")
	}

	if now == sf.timestamp {
		sf.sequence = (sf.sequence + 1) & maxSequence
		if sf.sequence == 0 {
			// Wait for the next millisecond if sequence overflows
			for now <= sf.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		sf.sequence = 0
	}

	sf.timestamp = now
	id := ((now - epoch) << timeShift) | (sf.nodeID << nodeShift) | sf.sequence

	// Add logging to check the generated ID
	// slog.Info("Generated ID", "node_id", sf.nodeID, "timestamp", now, "sequence", sf.sequence, "generated_id", id)

	return id
}
