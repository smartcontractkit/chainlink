package generic

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

func TestRelayerSet_List(t *testing.T) {
	testRelayersMap := map[types.RelayID]loop.Relayer{}
	testRelayersMap[types.RelayID{Network: "N1", ChainID: "C1"}] = &TestRelayer{}
	testRelayersMap[types.RelayID{Network: "N2", ChainID: "C2"}] = &TestRelayer{}
	testRelayersMap[types.RelayID{Network: "N3", ChainID: "C3"}] = &TestRelayer{}

	testGetter := TestRelayGetter{relayers: testRelayersMap}

	relayerSet, err := NewRelayerSet(testGetter, uuid.New(), 1, true)
	assert.NoError(t, err)
	relayers, err := relayerSet.List(context.Background())
	assert.NoError(t, err)

	assert.Len(t, relayers, 3)

	relayers, err = relayerSet.List(context.Background(), types.RelayID{Network: "N1", ChainID: "C1"}, types.RelayID{Network: "N3", ChainID: "C3"})
	assert.NoError(t, err)

	assert.Len(t, relayers, 2)

	_, ok := relayers[types.RelayID{Network: "N1", ChainID: "C1"}]
	assert.True(t, ok)

	_, ok = relayers[types.RelayID{Network: "N3", ChainID: "C3"}]
	assert.True(t, ok)
}

func TestRelayerSet_Get(t *testing.T) {
	testRelayersMap := map[types.RelayID]loop.Relayer{}
	testRelayersMap[types.RelayID{Network: "N1", ChainID: "C1"}] = &TestRelayer{}
	testRelayersMap[types.RelayID{Network: "N2", ChainID: "C2"}] = &TestRelayer{}
	testRelayersMap[types.RelayID{Network: "N3", ChainID: "C3"}] = &TestRelayer{}

	testGetter := TestRelayGetter{relayers: testRelayersMap}

	relayerSet, err := NewRelayerSet(testGetter, uuid.New(), 1, true)
	assert.NoError(t, err)

	_, err = relayerSet.Get(context.Background(), types.RelayID{Network: "N1", ChainID: "C1"})
	assert.NoError(t, err)

	_, err = relayerSet.Get(context.Background(), types.RelayID{Network: "N4", ChainID: "C4"})
	assert.Error(t, err)
}

func TestRelayerSet_NewPluginProvider(t *testing.T) {
	testRelayersMap := map[types.RelayID]loop.Relayer{}
	testRelayer := &TestRelayer{}
	testRelayersMap[types.RelayID{Network: "N1", ChainID: "C1"}] = testRelayer
	testRelayersMap[types.RelayID{Network: "N2", ChainID: "C2"}] = &TestRelayer{}
	testRelayersMap[types.RelayID{Network: "N3", ChainID: "C3"}] = &TestRelayer{}

	testGetter := TestRelayGetter{relayers: testRelayersMap}

	externalJobID := uuid.New()
	relayerSet, err := NewRelayerSet(testGetter, externalJobID, 1, true)
	assert.NoError(t, err)

	relayer, err := relayerSet.Get(context.Background(), types.RelayID{Network: "N1", ChainID: "C1"})
	assert.NoError(t, err)

	_, err = relayer.NewPluginProvider(context.Background(), core.RelayArgs{
		ContractID:   "c1",
		RelayConfig:  []byte("relayconfig"),
		ProviderType: "p1",
		MercuryCredentials: &types.MercuryCredentials{
			LegacyURL: "legacy",
			URL:       "url",
			Username:  "user",
			Password:  "pass",
		},
	}, core.PluginArgs{
		TransmitterID: "t1",
		PluginConfig:  []byte("pluginconfig"),
	})
	assert.NoError(t, err)

	assert.Equal(t, types.RelayArgs{
		ExternalJobID: externalJobID,
		JobID:         1,
		ContractID:    "c1",
		New:           true,
		RelayConfig:   []byte("relayconfig"),
		ProviderType:  "p1",
		MercuryCredentials: &types.MercuryCredentials{
			LegacyURL: "legacy",
			URL:       "url",
			Username:  "user",
			Password:  "pass",
		},
	}, testRelayer.relayArgs)

	assert.Equal(t, types.PluginArgs{
		TransmitterID: "t1",
		PluginConfig:  []byte("pluginconfig"),
	}, testRelayer.pluginArgs)
}

type TestRelayGetter struct {
	relayers map[types.RelayID]loop.Relayer
}

func (t TestRelayGetter) Get(id types.RelayID) (loop.Relayer, error) {
	if relayer, ok := t.relayers[id]; ok {
		return relayer, nil
	}

	return nil, fmt.Errorf("relayer with id %s not found", id)
}

func (t TestRelayGetter) GetIDToRelayerMap() map[types.RelayID]loop.Relayer {
	return t.relayers
}

type TestRelayer struct {
	relayArgs  types.RelayArgs
	pluginArgs types.PluginArgs
}

func (t *TestRelayer) NewPluginProvider(ctx context.Context, args types.RelayArgs, args2 types.PluginArgs) (types.PluginProvider, error) {
	t.relayArgs = args
	t.pluginArgs = args2

	return nil, nil
}

func (t *TestRelayer) Name() string { panic("implement me") }

func (t *TestRelayer) Start(ctx context.Context) error { panic("implement me") }

func (t *TestRelayer) Close() error { panic("implement me") }

func (t *TestRelayer) Ready() error { panic("implement me") }

func (t *TestRelayer) HealthReport() map[string]error { panic("implement me") }

func (t *TestRelayer) NewContractWriter(_ context.Context, _ []byte) (types.ContractWriter, error) {
	panic("implement me")
}

func (t *TestRelayer) NewContractReader(_ context.Context, _ []byte) (types.ContractReader, error) {
	panic("implement me")
}

func (t *TestRelayer) EVM() (types.EVMService, error) {
	panic("implement me")
}

func (t *TestRelayer) LatestHead(_ context.Context) (types.Head, error) {
	panic("implement me")
}

func (t *TestRelayer) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	panic("implement me")
}

func (t *TestRelayer) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []types.NodeStatus, nextPageToken string, total int, err error) {
	panic("implement me")
}

func (t *TestRelayer) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	panic("implement me")
}

func (t *TestRelayer) Replay(ctx context.Context, fromBlock string, args map[string]any) error {
	panic("implement me")
}

func (t *TestRelayer) NewConfigProvider(ctx context.Context, args types.RelayArgs) (types.ConfigProvider, error) {
	panic("implement me")
}

func (t *TestRelayer) NewLLOProvider(ctx context.Context, args types.RelayArgs, args2 types.PluginArgs) (types.LLOProvider, error) {
	panic("implement me")
}
