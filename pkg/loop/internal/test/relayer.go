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

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal"
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

func (s StaticPluginRelayer) NewRelayer(ctx context.Context, config string, keystore internal.Keystore) (internal.Relayer, error) {
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
	if !equalRelayArgs(r, rargs) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", rargs, r)
	}
	return staticConfigProvider{}, nil
}

func (s staticRelayer) NewMedianProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.MedianProvider, error) {
	if !equalRelayArgs(r, rargs) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", rargs, r)
	}
	if !reflect.DeepEqual(pargs, p) {
		return nil, fmt.Errorf("expected plugin args %v but got %v", pargs, p)
	}
	return staticMedianProvider{}, nil
}

func (s staticRelayer) NewMercuryProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MercuryProvider, error) {
	panic("unimplemented")
}

func (s staticRelayer) ChainStatus(ctx context.Context, id string) (types.ChainStatus, error) {
	if id != chainID {
		return types.ChainStatus{}, fmt.Errorf("expected id %s but got %s", chainID, id)
	}
	return chain, nil
}

func (s staticRelayer) ChainStatuses(ctx context.Context, o, l int) ([]types.ChainStatus, int, error) {
	if offset != o {
		return nil, -1, fmt.Errorf("expected offset %d but got %d", offset, o)
	}
	if limit != l {
		return nil, -1, fmt.Errorf("expected limit %d but got %d", limit, l)
	}
	return chains, count, nil
}

func (s staticRelayer) NodeStatuses(ctx context.Context, o, l int, cs ...string) ([]types.NodeStatus, int, error) {
	if offset != o {
		return nil, -1, fmt.Errorf("expected offset %d but got %d", offset, o)
	}
	if limit != l {
		return nil, -1, fmt.Errorf("expected limit %d but got %d", limit, l)
	}
	if !reflect.DeepEqual(chainIDs, cs) {
		return nil, -1, fmt.Errorf("expected chain IDs %v but got %v", chainIDs, cs)
	}
	return nodes, count, nil
}

func (s staticRelayer) SendTx(ctx context.Context, id, f, t string, a *big.Int, b bool) error {
	if id != chainID {
		return fmt.Errorf("expected id %s but got %s", chainID, id)
	}
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

func TestPluginRelayer(t *testing.T, p internal.PluginRelayer) {
	ctx := utils.Context(t)

	t.Run("Relayer", func(t *testing.T) {
		relayer, err := p.NewRelayer(ctx, ConfigTOML, StaticKeystore{})
		require.NoError(t, err)
		require.NoError(t, relayer.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, relayer.Close()) })
		TestRelayer(t, relayer)
	})
}

func TestRelayer(t *testing.T, relayer internal.Relayer) {
	ctx := utils.Context(t)

	t.Run("ConfigProvider", func(t *testing.T) {
		t.Parallel()
		configProvider, err := relayer.NewConfigProvider(ctx, rargs)
		require.NoError(t, err)
		require.NoError(t, configProvider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, configProvider.Close()) })

		t.Run("OffchainConfigDigester", func(t *testing.T) {
			t.Parallel()
			ocd := configProvider.OffchainConfigDigester()
			gotConfigDigestPrefix := ocd.ConfigDigestPrefix()
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
		provider, err := relayer.NewMedianProvider(ctx, rargs, pargs)
		require.NoError(t, err)
		require.NoError(t, provider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, provider.Close()) })

		t.Run("ReportingPluginProvider", func(t *testing.T) {
			t.Parallel()

			t.Run("OffchainConfigDigester", func(t *testing.T) {
				t.Parallel()
				ocd := provider.OffchainConfigDigester()
				gotConfigDigestPrefix := ocd.ConfigDigestPrefix()
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
				gotAccount := ct.FromAccount()
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
				gotMax := rc.MaxReportLength(n)
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

	t.Run("ChainStatus", func(t *testing.T) {
		t.Parallel()
		gotChain, err := relayer.ChainStatus(ctx, chainID)
		require.NoError(t, err)
		assert.Equal(t, chain, gotChain)
	})

	t.Run("ChainStatuses", func(t *testing.T) {
		t.Parallel()
		gotChains, gotCount, err := relayer.ChainStatuses(ctx, offset, limit)
		require.NoError(t, err)
		assert.Equal(t, chains, gotChains)
		assert.Equal(t, count, gotCount)
	})

	t.Run("NodeStatuses", func(t *testing.T) {
		t.Parallel()
		gotNodes, gotCount, err := relayer.NodeStatuses(ctx, offset, limit, chainIDs...)
		require.NoError(t, err)
		assert.Equal(t, nodes, gotNodes)
		assert.Equal(t, count, gotCount)
	})

	t.Run("SendTx", func(t *testing.T) {
		t.Parallel()
		err := relayer.SendTx(ctx, chainID, from, to, amount, balanceCheck)
		require.NoError(t, err)
	})
}
