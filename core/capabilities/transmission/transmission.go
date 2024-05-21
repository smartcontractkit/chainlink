package transmission

import (
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/libocr/permutation"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// TODO determine location for this code

var (
	// S = [N]
	Schedule_AllAtOnce = "allAtOnce"
	// S = [1 * N]
	Schedule_OneAtATime = "oneAtATime"
)

type TransmissionConfig struct {
	Schedule   string
	DeltaStage time.Duration
}

func ExtractTransmissionConfig(config *values.Map) (TransmissionConfig, error) {
	var tc struct {
		DeltaStage string
		Schedule   string
	}
	err := config.UnwrapTo(&tc)
	if err != nil {
		return TransmissionConfig{}, fmt.Errorf("failed to unwrap tranmission config from value map: %w", err)
	}

	duration, err := time.ParseDuration(tc.DeltaStage)
	if err != nil {
		return TransmissionConfig{}, fmt.Errorf("failed to parse DeltaStage %s as duration: %w", tc.DeltaStage, err)
	}

	return TransmissionConfig{
		Schedule:   tc.Schedule,
		DeltaStage: duration,
	}, nil
}

// GetPeerIDToTransmissionDelay returns a map of PeerID to the time.Duration that the node with that PeerID should wait
// before transmitting. If a node is not in the map, it should not transmit.
func GetPeerIDToTransmissionDelay(donPeerIDs []ragep2ptypes.PeerID, sharedSecret [16]byte, workflowID string,
	workflowExecutionID string, tc TransmissionConfig) (map[p2ptypes.PeerID]time.Duration, error) {
	donMemberCount := len(donPeerIDs)
	key := scheduleSeed(sharedSecret, workflowID, workflowExecutionID)
	sched, err := schedule(tc.Schedule, donMemberCount)
	if err != nil {
		return nil, err
	}

	picked := permutation.Permutation(donMemberCount, key)

	peerIDToTransmissionDelay := map[p2ptypes.PeerID]time.Duration{}
	for i, peerID := range donPeerIDs {
		delay := delayFor(i, sched, picked, tc.DeltaStage)
		if delay != nil {
			peerIDToTransmissionDelay[peerID] = *delay
		}
	}
	return peerIDToTransmissionDelay, nil
}

func delayFor(position int, schedule []int, permutation []int, deltaStage time.Duration) *time.Duration {
	sum := 0
	for i, s := range schedule {
		sum += s
		if permutation[position] < sum {
			result := time.Duration(i) * deltaStage
			return &result
		}
	}

	return nil
}

func schedule(sched string, N int) ([]int, error) {
	switch sched {
	case Schedule_AllAtOnce:
		return []int{N}, nil
	case Schedule_OneAtATime:
		sch := []int{}
		for i := 0; i < N; i++ {
			sch = append(sch, 1)
		}
		return sch, nil
	}
	return nil, fmt.Errorf("unknown schedule %s", sched)
}

// scheduleSeed uses a shared secret, combined with a workflowID and a workflowExecutionID to generate
// a secret that can later be used to pseudo-randomly determine a schedule for a set of nodes in a DON.
// The addition of the workflowExecutionID -- which nodes don't know ahead of time -- additionally guarantees
// that a malicious coalition of nodes can't "game" the schedule.
// IMPORTANT: changing this function should happen carefully to maintain the guarantee that all nodes
// arrive at the same secret.
func scheduleSeed(sharedSecret [16]byte, workflowID, workflowExecutionID string) [16]byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(sharedSecret[:])
	hash.Write([]byte(workflowID))
	hash.Write([]byte(workflowExecutionID))

	var key [16]byte
	copy(key[:], hash.Sum(nil))
	return key
}
