package registrysyncer_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func TestRegistrySyncerORM_InsertAndRetrieval(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	orm := registrysyncer.NewORM(db, lggr)

	var states []registrysyncer.LocalRegistry
	for range 11 {
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
	var pid types.PeerID
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
		IDsToNodes: map[types.PeerID]registrysyncer.NodeInfo{
			nodes[0]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2pID:               nodes[0],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIDs: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[1]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2pID:               nodes[1],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIDs: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[2]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2pID:               nodes[2],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIDs: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
			nodes[3]: {
				NodeOperatorID:      1,
				Signer:              randomWord(),
				P2pID:               nodes[3],
				EncryptionPublicKey: randomWord(),
				HashedCapabilityIDs: [][32]byte{capabilityID, capabilityID2},
				CapabilitiesDONIds:  []*big.Int{},
			},
		},
	}
}

func TestRegistrySyncerORM_AddLocalRegistry_DuplicateHandling(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	orm := registrysyncer.NewORM(db, lggr)

	t.Run("duplicate_handling", func(t *testing.T) {
		// Generate original state
		originalState := generateState(t)

		// First insertion - should succeed
		err := orm.AddLocalRegistry(ctx, originalState)
		require.NoError(t, err)

		// Get the initial ID and hash
		var initialID int
		var initialHash string
		err = db.Get(&initialID, `SELECT id FROM registry_syncer_states ORDER BY id DESC LIMIT 1`)
		require.NoError(t, err)
		err = db.Get(&initialHash, `SELECT data_hash FROM registry_syncer_states ORDER BY id DESC LIMIT 1`)
		require.NoError(t, err)

		// Second insertion with same data - should not insert (no new row)
		err = orm.AddLocalRegistry(ctx, originalState)
		require.NoError(t, err)

		// Check that latest ID hasn't changed (no new insertion)
		id, hash := latestHashAndID(t, db)

		assert.Equal(t, initialID, id, "ID should not change when inserting duplicate data")
		assert.Equal(t, initialHash, hash, "Hash should not change when inserting duplicate data")

		// Generate new state with different data
		newState := generateState(t)

		// Third insertion with new data - should succeed and increment ID
		err = orm.AddLocalRegistry(ctx, newState)
		require.NoError(t, err)

		// Check that latest ID is incremented and hash is new
		newID, newHash := latestHashAndID(t, db)
		assert.Greater(t, newID, id, "ID should increment when inserting new data")
		assert.NotEqual(t, newHash, hash, "Hash should change when inserting new data")

		// Fourth insertion with original data again - should succeed and increment ID
		err = orm.AddLocalRegistry(ctx, originalState)
		require.NoError(t, err)

		// Check that latest ID is incremented again and hash is back to original
		finalID, finalHash := latestHashAndID(t, db)

		assert.Greater(t, finalID, newID, "ID should increment when re-inserting original data")
		assert.Equal(t, finalHash, initialHash, "Hash should match original data hash when re-inserting original data")

		// Verify total count - should have 3 records (original, new, original again)
		var totalCount int
		err = db.Get(&totalCount, `SELECT count(*) FROM registry_syncer_states`)
		require.NoError(t, err)
		assert.Equal(t, 3, totalCount, "Should have exactly 3 records")
	})
}

func latestHashAndID(t *testing.T, db *sqlx.DB) (int, string) {
	var id int
	var hash string
	err := db.Get(&id, `SELECT id FROM registry_syncer_states ORDER BY id DESC LIMIT 1`)
	require.NoError(t, err)
	err = db.Get(&hash, `SELECT data_hash FROM registry_syncer_states ORDER BY id DESC LIMIT 1`)
	require.NoError(t, err)
	return id, hash
}
