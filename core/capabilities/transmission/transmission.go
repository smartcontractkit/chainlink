package transmission

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/permutation"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	"golang.org/x/crypto/sha3"
)

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

	// Default if no schedule and deltaStage is provided
	if len(tc.Schedule) == 0 && len(tc.DeltaStage) == 0 {
		return TransmissionConfig{
			Schedule:   Schedule_AllAtOnce,
			DeltaStage: 0,
		}, nil
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
// before transmitting the capability request. If a node is not in the map, it should not transmit.
func GetPeerIDToTransmissionDelay(donPeerIDs []types.PeerID, req capabilities.CapabilityRequest) (map[types.PeerID]time.Duration, error) {
	tc, err := ExtractTransmissionConfig(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transmission config from request: %w", err)
	}

	workflowExecutionID := req.Metadata.WorkflowExecutionID
	if err := validation.ValidateWorkflowOrExecutionID(workflowExecutionID); err != nil {
		return nil, fmt.Errorf("workflow or execution ID is invalid: %w", err)
	}

	return GetPeerIDToTransmissionDelaysForConfig(donPeerIDs, workflowExecutionID, tc)
}

func GetPeerIDToTransmissionDelaysForConfig(donPeerIDs []types.PeerID, transmissionID string, tc TransmissionConfig) (map[types.PeerID]time.Duration, error) {
	donMemberCount := len(donPeerIDs)
	key := transmissionScheduleSeed(transmissionID)
	schedule, err := createTransmissionSchedule(tc.Schedule, donMemberCount)
	if err != nil {
		return nil, err
	}

	picked := permutation.Permutation(donMemberCount, key)

	peerIDToTransmissionDelay := map[types.PeerID]time.Duration{}
	for i, peerID := range donPeerIDs {
		delay := delayFor(i, schedule, picked, tc.DeltaStage)
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

func createTransmissionSchedule(scheduleType string, N int) ([]int, error) {
	switch scheduleType {
	case Schedule_AllAtOnce:
		return []int{N}, nil
	case Schedule_OneAtATime:
		sch := []int{}
		for i := 0; i < N; i++ {
			sch = append(sch, 1)
		}
		return sch, nil
	}
	return nil, fmt.Errorf("unknown schedule type %s", scheduleType)
}

func transmissionScheduleSeed(transmissionID string) [16]byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(transmissionID))
	var key [16]byte
	copy(key[:], hash.Sum(nil))
	return key
}
