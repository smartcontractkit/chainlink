package registrysyncer

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestRegistrySyncerORM_InsertAndRetrieval(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	orm := newORM(db, lggr)

	var states []LocalRegistry
	for i := 0; i < 11; i++ {
		state := generateState(t)
		err := orm.addLocalRegistry(ctx, state)
		require.NoError(t, err)
		states = append(states, state)
	}

	var count int
	err := db.Get(&count, `SELECT count(*) FROM registry_syncer_states`)
	require.NoError(t, err)
	assert.Equal(t, 10, count)

	state, err := orm.latestLocalRegistry(ctx)
	require.NoError(t, err)
	assert.Equal(t, states[10], *state)
}

func generateState(t *testing.T) LocalRegistry {
	dID := uint32(1)
	var pid ragetypes.PeerID
	err := pid.UnmarshalText([]byte("12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"))
	require.NoError(t, err)
	nodes := [][32]byte{
		pid,
		randomWord(),
		randomWord(),
		randomWord(),
	}
	capabilityID := randomWord()
	capabilityID2 := randomWord()
	capabilityIDStr := hex.EncodeToString(capabilityID[:])
	capabilityID2Str := hex.EncodeToString(capabilityID2[:])
	rtc := &capabilities.RemoteTriggerConfig{
		RegistrationRefresh:     20 * time.Second,
		MinResponsesToAggregate: 2,
		RegistrationExpiry:      60 * time.Second,
		MessageExpiry:           120 * time.Second,
	}

	fmt.Println(capabilityID2Str, capabilityIDStr)

	return LocalRegistry{
		IDsToDONs: map[DonID]DON{
			DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          toPeerIDs(nodes),
				},
				CapabilityConfigurations: map[string]capabilities.CapabilityConfiguration{
					capabilityIDStr: {
						DefaultConfig:       values.EmptyMap(),
						RemoteTriggerConfig: rtc,
					},
					capabilityID2Str: {
						DefaultConfig:       values.EmptyMap(),
						RemoteTriggerConfig: rtc,
					},
				},
			},
		},
		IDsToCapabilities: map[string]Capability{
			capabilityIDStr: {
				ID:             capabilityIDStr,
				CapabilityType: 0,
			},
			capabilityID2Str: {
				ID:             capabilityID2Str,
				CapabilityType: 3,
			},
		},
		IDsToNodes: map[types.PeerID]kcr.CapabilitiesRegistryNodeInfo{
			nodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[0],
				HashedCapabilityIds: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[1],
				HashedCapabilityIds: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[2],
				HashedCapabilityIds: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[3],
				HashedCapabilityIds: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
		},
	}
}
