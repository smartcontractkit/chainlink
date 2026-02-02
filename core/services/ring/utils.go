package ring

import (
	"errors"
	"slices"
	"strconv"

	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash/v2"
)

var errInvalidRing = errors.New("RingOCR invalid ring for consistent hashing")
var errInvalidMember = errors.New("RingOCR invalid member for consistent hashing")

func uniqueSorted(s []string) []string {
	result := slices.Clone(s)
	slices.Sort(result)
	return slices.Compact(result)
}

type xxhashHasher struct{}

func (h xxhashHasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

type ShardMember string

func (m ShardMember) String() string {
	return string(m)
}

func consistentHashConfig() consistent.Config {
	return consistent.Config{
		PartitionCount:    997, // Prime number for better distribution
		ReplicationFactor: 50,  // Number of replicas per node
		Load:              1.1, // Load factor for bounded loads
		Hasher:            xxhashHasher{},
	}
}

func newShardRing(healthyShards []uint32) *consistent.Consistent {
	if len(healthyShards) == 0 {
		return nil
	}
	members := make([]consistent.Member, len(healthyShards))
	for i, shardID := range healthyShards {
		members[i] = ShardMember(strconv.FormatUint(uint64(shardID), 10))
	}
	return consistent.New(members, consistentHashConfig())
}

func locateShard(ring *consistent.Consistent, workflowID string) (uint32, error) {
	if ring == nil {
		return 0, errInvalidRing
	}
	member := ring.LocateKey([]byte(workflowID))
	if member == nil {
		return 0, errInvalidMember
	}
	shardID, err := strconv.ParseUint(member.String(), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(shardID), nil
}
