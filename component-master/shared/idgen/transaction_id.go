package idgen

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

var sf *Snowflake

func init() {
	workerID := 1 // Example: Assign a unique worker ID for each worker
	nodeID, err := getNodeID(workerID)
	if err != nil {
		panic(err)
	}
	sf, err = NewSnowflake(nodeID)
	if err != nil {
		panic(err)
	}
}

// getNodeID retrieves the node ID from the configuration (or environment variables)
func getNodeID(workerID int) (int64, error) {
	nodeIDStr := os.Getenv("NODE_ID") // Retrieve node ID from the environment variable
	if nodeIDStr == "" {
		nodeIDStr = fmt.Sprintf("%d", workerID) // Default to workerID if NODE_ID is not configured
	}
	nodeID, err := strconv.ParseInt(nodeIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return nodeID, nil
}

// GenID generates a unique ID
// exp: id := idgen.GenID()
func GenID() int64 {
	workerID := 1
	nodeID, err := getNodeID(workerID)
	sf, err := NewSnowflake(nodeID)
	if err != nil {
		slog.Error("Failed to create Snowflake instance", "error", err)
		return 0
	}
	return sf.Generate()
}
