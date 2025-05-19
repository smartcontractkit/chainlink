package registrysyncer_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func TestRegistrySyncerORM_InsertAndRetrieval(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	orm := registrysyncer.NewORM(db, lggr)

	var states []registrysyncer.LocalRegistry
	for i := 0; i < 11; i++ {
		state := generateState(t)
		err := orm.AddLocalRegistry(ctx, state)
		require.NoError(t, err)
		states = append(states, state)
	}

	var count int
	err := db.Get(&count, `SELECT count(*) FROM registry_syncer_states`)
	require.NoError(t, err)
	assert.Equal(t, 10, count)

	state, err := orm.LatestLocalRegistry(ctx)
	require.NoError(t, err)
	assert.Equal(t, states[10], *state)
}

func generateState(t *testing.T) registrysyncer.LocalRegistry {
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

	config := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: durationpb.New(20 * time.Second),
				RegistrationExpiry:  durationpb.New(60 * time.Second),
				// F + 1
				MinResponsesToAggregate: uint32(1) + 1,
				MessageExpiry:           durationpb.New(120 * time.Second),
			},
		},
	}
	configb, err := proto.Marshal(config)
	require.NoError(t, err)

	return registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: capabilities.DON{
					ID:               dID,
					ConfigVersion:    uint32(0),
					F:                uint8(1),
					IsPublic:         true,
					AcceptsWorkflows: true,
					Members:          toPeerIDs(nodes),
				},
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					capabilityIDStr: {
						Config: configb,
					},
					capabilityID2Str: {
						Config: configb,
					},
				},
			},
		},
		IDsToCapabilities: map[string]registrysyncer.Capability{
			capabilityIDStr: {
				ID:             capabilityIDStr,
				CapabilityType: capabilities.CapabilityTypeAction,
			},
			capabilityID2Str: {
				ID:             capabilityID2Str,
				CapabilityType: capabilities.CapabilityTypeConsensus,
			},
		},
		IDsToNodes: map[types.PeerID]kcr.INodeInfoProviderNodeInfo{
			nodes[0]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[1]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[2]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[3]: {
				NodeOperatorId:      1,
				Signer:              randomWord(),
				P2pId:               nodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIds: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
		},
	}
}
