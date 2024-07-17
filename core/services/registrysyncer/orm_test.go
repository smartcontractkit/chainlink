package registrysyncer

import (
	"math/big"
	"testing"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	var states []State
	for i := 0; i < 11; i++ {
		state := generateState(t)
		err := orm.addState(ctx, state)
		require.NoError(t, err)
		states = append(states, state)
	}

	var count int
	err := db.Get(&count, `SELECT count(*) FROM registry_syncer_states`)
	require.NoError(t, err)
	assert.Equal(t, 10, count)

	state, err := orm.latestState(ctx)
	require.NoError(t, err)
	assert.Equal(t, states[10], *state)
}

func generateState(t *testing.T) State {
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
	addr := testutils.NewAddress()

	return State{
		IDsToDONs: map[DonID]kcr.CapabilitiesRegistryDONInfo{
			DonID(dID): {
				Id:               dID,
				ConfigCount:      uint32(0),
				F:                uint8(1),
				IsPublic:         true,
				AcceptsWorkflows: true,
				NodeP2PIds:       nodes,
				CapabilityConfigurations: []kcr.CapabilitiesRegistryCapabilityConfiguration{
					{
						CapabilityId: capabilityID,
						Config:       []byte(""),
					},
					{
						CapabilityId: capabilityID2,
						Config:       []byte(""),
					},
				},
			},
		},
		IDsToCapabilities: map[HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo{
			capabilityID: {
				HashedId:              capabilityID,
				LabelledName:          "label-1",
				Version:               "1.0.0",
				CapabilityType:        0,
				ResponseType:          0,
				ConfigurationContract: addr,
				IsDeprecated:          false,
			},
			capabilityID2: {
				HashedId:              capabilityID2,
				LabelledName:          "label-2",
				Version:               "1.0.0",
				CapabilityType:        3,
				ResponseType:          0,
				ConfigurationContract: addr,
				IsDeprecated:          false,
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
