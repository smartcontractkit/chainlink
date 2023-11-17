package test

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

type StaticKeystore struct{}

func (s StaticKeystore) Accounts(ctx context.Context) (accounts []string, err error) {
	return []string{string(account)}, nil
}

func (s StaticKeystore) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	if string(account) != id {
		return nil, fmt.Errorf("expected id %q but got %q", account, id)
	}
	if !bytes.Equal(encoded, data) {
		return nil, fmt.Errorf("expected encoded data %x but got %x", encoded, data)
	}
	return signed, nil
}

type StaticPluginRelayer struct{}

func (s StaticPluginRelayer) NewRelayer(ctx context.Context, config string, keystore types.Keystore) (internal.Relayer, error) {
	if config != ConfigTOML {
		return nil, fmt.Errorf("expected config %q but got %q", ConfigTOML, config)
	}
	keys, err := keystore.Accounts(ctx)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual([]string{string(account)}, keys) {
		return nil, fmt.Errorf("expected keys %v but got %v", []string{string(account)}, keys)
	}
	gotSigned, err := keystore.Sign(ctx, string(account), encoded)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(signed, gotSigned) {
		return nil, fmt.Errorf("expected signed bytes %x but got %x", signed, gotSigned)
	}
	return staticRelayer{}, nil
}

type staticRelayer struct{}

func (s staticRelayer) Start(ctx context.Context) error { return nil }

func (s staticRelayer) Close() error { return nil }

func (s staticRelayer) Ready() error { panic("unimplemented") }

func (s staticRelayer) Name() string { panic("unimplemented") }

func (s staticRelayer) HealthReport() map[string]error { panic("unimplemented") }

func (s staticRelayer) NewConfigProvider(ctx context.Context, r types.RelayArgs) (types.ConfigProvider, error) {
	if !equalRelayArgs(r, RelayArgs) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
	}
	return staticConfigProvider{}, nil
}

func (s staticRelayer) NewMedianProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.MedianProvider, error) {
	ra := newRelayArgsWithProviderType(types.Median)
	if !equalRelayArgs(r, ra) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
	}
	if !reflect.DeepEqual(PluginArgs, p) {
		return nil, fmt.Errorf("expected plugin args %v but got %v", PluginArgs, p)
	}
	return StaticMedianProvider{}, nil
}

func (s staticRelayer) NewPluginProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.PluginProvider, error) {
	ra := newRelayArgsWithProviderType(types.Median)
	if !equalRelayArgs(r, ra) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
	}
	if !reflect.DeepEqual(PluginArgs, p) {
		return nil, fmt.Errorf("expected plugin args %v but got %v", PluginArgs, p)
	}
	return StaticPluginProvider{}, nil
}

func (s staticRelayer) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	return chain, nil
}

func (s staticRelayer) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) ([]types.NodeStatus, string, int, error) {

	if limit != pageSize {
		return nil, "", -1, fmt.Errorf("expected page_size %d but got %d", limit, pageSize)
	}
	if pageToken != "" {
		return nil, "", -1, fmt.Errorf("expected empty page_token but got %q", pageToken)

	}
	return nodes, "", total, nil
}

func (s staticRelayer) Transact(ctx context.Context, f, t string, a *big.Int, b bool) error {
	if f != from {
		return fmt.Errorf("expected from %s but got %s", from, f)
	}
	if t != to {
		return fmt.Errorf("expected to %s but got %s", to, t)
	}
	if amount.Cmp(a) != 0 {
		return fmt.Errorf("expected amount %s but got %s", amount, a)
	}
	if b != balanceCheck { //nolint:gosimple
		return fmt.Errorf("expected balance check %t but got %t", balanceCheck, b)
	}
	return nil
}

func equalRelayArgs(a, b types.RelayArgs) bool {
	return a.ExternalJobID == b.ExternalJobID &&
		a.JobID == b.JobID &&
		a.ContractID == b.ContractID &&
		a.New == b.New &&
		bytes.Equal(a.RelayConfig, b.RelayConfig)
}

func newRelayArgsWithProviderType(_type types.OCR2PluginType) types.RelayArgs {
	return types.RelayArgs{
		ExternalJobID: RelayArgs.ExternalJobID,
		JobID:         RelayArgs.JobID,
		ContractID:    RelayArgs.ContractID,
		New:           RelayArgs.New,
		RelayConfig:   RelayArgs.RelayConfig,
		ProviderType:  string(_type),
	}
}

func TestPluginRelayer(t *testing.T, p internal.PluginRelayer) {
	ctx := tests.Context(t)

	t.Run("Relayer", func(t *testing.T) {
		relayer, err := p.NewRelayer(ctx, ConfigTOML, StaticKeystore{})
		require.NoError(t, err)
		require.NoError(t, relayer.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, relayer.Close()) })
		TestRelayer(t, relayer)
	})
}

func TestRelayer(t *testing.T, relayer internal.Relayer) {
	ctx := tests.Context(t)

	t.Run("ConfigProvider", func(t *testing.T) {
		t.Parallel()
		configProvider, err := relayer.NewConfigProvider(ctx, RelayArgs)
		require.NoError(t, err)
		require.NoError(t, configProvider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, configProvider.Close()) })

		t.Run("OffchainConfigDigester", func(t *testing.T) {
			t.Parallel()
			ocd := configProvider.OffchainConfigDigester()
			gotConfigDigestPrefix, err := ocd.ConfigDigestPrefix()
			require.NoError(t, err)
			assert.Equal(t, configDigestPrefix, gotConfigDigestPrefix)
			gotConfigDigest, err := ocd.ConfigDigest(contractConfig)
			require.NoError(t, err)
			assert.Equal(t, configDigest, gotConfigDigest)
		})
		t.Run("ContractConfigTracker", func(t *testing.T) {
			t.Parallel()
			cct := configProvider.ContractConfigTracker()
			gotBlockHeight, err := cct.LatestBlockHeight(ctx)
			require.NoError(t, err)
			assert.Equal(t, blockHeight, gotBlockHeight)
			gotChangedInBlock, gotConfigDigest, err := cct.LatestConfigDetails(ctx)
			require.NoError(t, err)
			assert.Equal(t, changedInBlock, gotChangedInBlock)
			assert.Equal(t, configDigest, gotConfigDigest)
			gotContractConfig, err := cct.LatestConfig(ctx, changedInBlock)
			require.NoError(t, err)
			assert.Equal(t, contractConfig, gotContractConfig)
		})
	})

	t.Run("MedianProvider", func(t *testing.T) {
		t.Parallel()
		ra := newRelayArgsWithProviderType(types.Median)
		p, err := relayer.NewPluginProvider(ctx, ra, PluginArgs)
		provider := p.(types.MedianProvider)
		require.NoError(t, err)
		require.NoError(t, provider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, provider.Close()) })

		t.Run("ReportingPluginProvider", func(t *testing.T) {
			t.Parallel()

			t.Run("OffchainConfigDigester", func(t *testing.T) {
				t.Parallel()
				ocd := provider.OffchainConfigDigester()
				gotConfigDigestPrefix, err := ocd.ConfigDigestPrefix()
				require.NoError(t, err)
				assert.Equal(t, configDigestPrefix, gotConfigDigestPrefix)
				gotConfigDigest, err := ocd.ConfigDigest(contractConfig)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
			})
			t.Run("ContractConfigTracker", func(t *testing.T) {
				t.Parallel()
				cct := provider.ContractConfigTracker()
				gotBlockHeight, err := cct.LatestBlockHeight(ctx)
				require.NoError(t, err)
				assert.Equal(t, blockHeight, gotBlockHeight)
				gotChangedInBlock, gotConfigDigest, err := cct.LatestConfigDetails(ctx)
				require.NoError(t, err)
				assert.Equal(t, changedInBlock, gotChangedInBlock)
				assert.Equal(t, configDigest, gotConfigDigest)
				gotContractConfig, err := cct.LatestConfig(ctx, changedInBlock)
				require.NoError(t, err)
				assert.Equal(t, contractConfig, gotContractConfig)
			})
			t.Run("ContractTransmitter", func(t *testing.T) {
				t.Parallel()
				ct := provider.ContractTransmitter()
				gotAccount, err := ct.FromAccount()
				require.NoError(t, err)
				assert.Equal(t, account, gotAccount)
				gotConfigDigest, gotEpoch, err := ct.LatestConfigDigestAndEpoch(ctx)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
				assert.Equal(t, epoch, gotEpoch)
				err = ct.Transmit(ctx, reportContext, report, sigs)
				require.NoError(t, err)
			})
			t.Run("ReportCodec", func(t *testing.T) {
				t.Parallel()
				rc := provider.ReportCodec()
				gotReport, err := rc.BuildReport(pobs)
				require.NoError(t, err)
				assert.Equal(t, report, gotReport)
				gotMedianValue, err := rc.MedianFromReport(report)
				require.NoError(t, err)
				assert.Equal(t, medianValue, gotMedianValue)
				gotMax, err := rc.MaxReportLength(n)
				require.NoError(t, err)
				assert.Equal(t, max, gotMax)
			})
			t.Run("MedianContract", func(t *testing.T) {
				t.Parallel()
				mc := provider.MedianContract()
				gotConfigDigest, gotEpoch, gotRound, err := mc.LatestRoundRequested(ctx, lookbackDuration)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
				assert.Equal(t, epoch, gotEpoch)
				assert.Equal(t, round, gotRound)
				gotConfigDigest, gotEpoch, gotRound, gotLatestAnswer, gotLatestTimestamp, err := mc.LatestTransmissionDetails(ctx)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
				assert.Equal(t, epoch, gotEpoch)
				assert.Equal(t, round, gotRound)
				assert.Equal(t, latestAnswer, gotLatestAnswer)
				assert.WithinDuration(t, latestTimestamp, gotLatestTimestamp, time.Second)
			})
			t.Run("OnchainConfigCodec", func(t *testing.T) {
				t.Parallel()
				occ := provider.OnchainConfigCodec()
				gotEncoded, err := occ.Encode(onchainConfig)
				require.NoError(t, err)
				assert.Equal(t, encoded, gotEncoded)
				gotDecoded, err := occ.Decode(encoded)
				require.NoError(t, err)
				assert.Equal(t, onchainConfig, gotDecoded)
			})
		})
	})

	t.Run("PluginProvider", func(t *testing.T) {
		t.Parallel()
		ra := newRelayArgsWithProviderType(types.GenericPlugin)
		provider, err := relayer.NewPluginProvider(ctx, ra, PluginArgs)
		require.NoError(t, err)
		require.NoError(t, provider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, provider.Close()) })

		t.Run("ReportingPluginProvider", func(t *testing.T) {
			t.Parallel()

			t.Run("OffchainConfigDigester", func(t *testing.T) {
				t.Parallel()
				ocd := provider.OffchainConfigDigester()
				gotConfigDigestPrefix, err := ocd.ConfigDigestPrefix()
				require.NoError(t, err)
				assert.Equal(t, configDigestPrefix, gotConfigDigestPrefix)
				gotConfigDigest, err := ocd.ConfigDigest(contractConfig)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
			})
			t.Run("ContractConfigTracker", func(t *testing.T) {
				t.Parallel()
				cct := provider.ContractConfigTracker()
				gotBlockHeight, err := cct.LatestBlockHeight(ctx)
				require.NoError(t, err)
				assert.Equal(t, blockHeight, gotBlockHeight)
				gotChangedInBlock, gotConfigDigest, err := cct.LatestConfigDetails(ctx)
				require.NoError(t, err)
				assert.Equal(t, changedInBlock, gotChangedInBlock)
				assert.Equal(t, configDigest, gotConfigDigest)
				gotContractConfig, err := cct.LatestConfig(ctx, changedInBlock)
				require.NoError(t, err)
				assert.Equal(t, contractConfig, gotContractConfig)
			})
			t.Run("ContractTransmitter", func(t *testing.T) {
				t.Parallel()
				ct := provider.ContractTransmitter()
				gotAccount, err := ct.FromAccount()
				require.NoError(t, err)
				assert.Equal(t, account, gotAccount)
				gotConfigDigest, gotEpoch, err := ct.LatestConfigDigestAndEpoch(ctx)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
				assert.Equal(t, epoch, gotEpoch)
				err = ct.Transmit(ctx, reportContext, report, sigs)
				require.NoError(t, err)
			})
		})
	})

	t.Run("GetChainStatus", func(t *testing.T) {
		t.Parallel()
		gotChain, err := relayer.GetChainStatus(ctx)
		require.NoError(t, err)
		assert.Equal(t, chain, gotChain)
	})

	t.Run("ListNodeStatuses", func(t *testing.T) {
		t.Parallel()
		gotNodes, gotNextToken, gotCount, err := relayer.ListNodeStatuses(ctx, limit, "")
		require.NoError(t, err)
		assert.Equal(t, nodes, gotNodes)
		assert.Equal(t, total, gotCount)
		assert.Empty(t, gotNextToken)
	})

	t.Run("Transact", func(t *testing.T) {
		t.Parallel()
		err := relayer.Transact(ctx, from, to, amount, balanceCheck)
		require.NoError(t, err)
	})
}
