package arbiter

// Status constants for shard replicas.
const (
	StatusReady = "READY"
)

// ShardReplica represents the status and message of a single shard replica.
// Used by State to track current replicas.
type ShardReplica struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Metrics map[string]float64 `json:"metrics,omitempty"`
}

// RoutableShardsInfo contains routing info for Ring OCR.
// This is used internally; the actual protobuf response is built in arbiter_scaler.go.
type RoutableShardsInfo struct {
	ReadyCount int                    // Count of shards with Status == READY
	ShardInfo  map[uint32]ShardHealth // Per-shard health status (shard_id → health)
}

// ShardHealth represents the health status of a single shard.
type ShardHealth struct {
	IsHealthy bool
}
