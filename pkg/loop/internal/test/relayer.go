package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

type StaticPluginRelayer struct {
	StaticChecks bool
}

func (s StaticPluginRelayer) NewRelayer(ctx context.Context, config string, keystore types.Keystore) (internal.Relayer, error) {
	if s.StaticChecks && config != ConfigTOML {
		return nil, fmt.Errorf("expected config %q but got %q", ConfigTOML, config)
	}
	keys, err := keystore.Accounts(ctx)
	if err != nil {
		return nil, err
	}
	if s.StaticChecks && !reflect.DeepEqual([]string{string(account)}, keys) {
		return nil, fmt.Errorf("expected keys %v but got %v", []string{string(account)}, keys)
	}
	gotSigned, err := keystore.Sign(ctx, string(account), encoded)
	if err != nil {
		return nil, err
	}
	if s.StaticChecks && !bytes.Equal(signed, gotSigned) {
		return nil, fmt.Errorf("expected signed bytes %x but got %x", signed, gotSigned)
	}
	return s, nil
}

func (s StaticPluginRelayer) Start(ctx context.Context) error { return nil }

func (s StaticPluginRelayer) Close() error { return nil }

func (s StaticPluginRelayer) Ready() error { panic("unimplemented") }

func (s StaticPluginRelayer) Name() string { panic("unimplemented") }

func (s StaticPluginRelayer) HealthReport() map[string]error { panic("unimplemented") }

func (s StaticPluginRelayer) NewConfigProvider(ctx context.Context, r types.RelayArgs) (types.ConfigProvider, error) {
	if s.StaticChecks && !equalRelayArgs(r, RelayArgs) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
	}
	return staticConfigProvider{}, nil
}

func (s StaticPluginRelayer) NewMedianProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.MedianProvider, error) {
	if s.StaticChecks {
		ra := newRelayArgsWithProviderType(types.Median)
		if !equalRelayArgs(r, ra) {
			return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
		}
		if !reflect.DeepEqual(PluginArgs, p) {
			return nil, fmt.Errorf("expected plugin args %v but got %v", PluginArgs, p)
		}
	}

	return StaticMedianProvider{}, nil
}

func (s StaticPluginRelayer) NewPluginProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.PluginProvider, error) {
	if s.StaticChecks {
		ra := newRelayArgsWithProviderType(types.Median)
		if !equalRelayArgs(r, ra) {
			return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
		}
		if !reflect.DeepEqual(PluginArgs, p) {
			return nil, fmt.Errorf("expected plugin args %v but got %v", PluginArgs, p)
		}
	}

	return StaticPluginProvider{}, nil
}

func (s StaticPluginRelayer) NewLLOProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.LLOProvider, error) {
	return nil, errors.New("not implemented")
}

func (s StaticPluginRelayer) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	return chain, nil
}

func (s StaticPluginRelayer) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) ([]types.NodeStatus, string, int, error) {
	if s.StaticChecks && limit != pageSize {
		return nil, "", -1, fmt.Errorf("expected page_size %d but got %d", limit, pageSize)
	}
	if pageToken != "" {
		return nil, "", -1, fmt.Errorf("expected empty page_token but got %q", pageToken)
	}
	return nodes, "", total, nil
}

func (s StaticPluginRelayer) Transact(ctx context.Context, f, t string, a *big.Int, b bool) error {
	if s.StaticChecks {
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

func RunPluginRelayer(t *testing.T, p internal.PluginRelayer) {
	ctx := tests.Context(t)

	t.Run("Relayer", func(t *testing.T) {
		relayer, err := p.NewRelayer(ctx, ConfigTOML, StaticKeystore{})
		require.NoError(t, err)
		require.NoError(t, relayer.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, relayer.Close()) })
		RunRelayer(t, relayer)
	})
}

func RunRelayer(t *testing.T, relayer internal.Relayer) {
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

func RunFuzzPluginRelayer(f *testing.F, relayerFunc func(*testing.T) internal.PluginRelayer) {
	f.Add("ABC\xa8\x8c\xb3G\xfc", "", true, []byte{}, true, true, "")
	f.Add(ConfigTOML, string(account), false, signed, false, false, "")

	f.Fuzz(func(
		t *testing.T, fConfig string, fAccts string, fAcctErr bool,
		fSigned []byte, fSignErr bool, fValsWErr bool, fErr string,
	) {
		keystore := fuzzerKeystore{
			accounts:      []string{fAccts}, // fuzzer does not support []string type
			acctErr:       fAcctErr,
			signed:        fSigned,
			signErr:       fSignErr,
			valuesWithErr: fValsWErr,
			errStr:        fErr,
		}

		ctx := tests.Context(t)
		_, err := relayerFunc(t).NewRelayer(ctx, fConfig, keystore)

		grpcUnavailableErr(t, err)
	})
}

func RunFuzzRelayer(f *testing.F, relayerFunc func(*testing.T) internal.Relayer) {
	validRaw := [16]byte(RelayArgs.ExternalJobID)
	validRawBytes := make([]byte, 16)

	copy(validRawBytes, validRaw[:])

	f.Add([]byte{}, int32(-1), "ABC\xa8\x8c\xb3G\xfc", false, []byte{}, "", "", []byte{})
	f.Add(validRawBytes, int32(123), "testcontract", true, []byte(ConfigTOML), string(types.Median), "testtransmitter", []byte{100: 88})

	f.Fuzz(func(
		t *testing.T, fExtJobID []byte, fJobID int32, fContractID string, fNew bool,
		fConfig []byte, fType string, fTransmID string, fPlugConf []byte,
	) {
		var rawBytes [16]byte

		copy(rawBytes[:], fExtJobID)

		relayer := relayerFunc(t)
		ctx := tests.Context(t)
		fRelayArgs := types.RelayArgs{
			ExternalJobID: uuid.UUID(rawBytes),
			JobID:         fJobID,
			ContractID:    fContractID,
			New:           fNew,
			RelayConfig:   fConfig,
			ProviderType:  fType,
		}

		_, err := relayer.NewConfigProvider(ctx, fRelayArgs)

		grpcUnavailableErr(t, err)

		pArgs := types.PluginArgs{
			TransmitterID: fTransmID,
			PluginConfig:  fPlugConf,
		}

		provider, err := relayer.NewPluginProvider(ctx, fRelayArgs, pArgs)
		// require.NoError(t, provider.Start(ctx))
		t.Log("provider created")
		t.Cleanup(func() {
			t.Log("cleanup called")
			if provider != nil {
				assert.NoError(t, provider.Close())
			}
		})

		grpcUnavailableErr(t, err)
		t.Logf("error tested: %s", err)
	})
}

type FuzzableProvider[K any] func(context.Context, types.RelayArgs, types.PluginArgs) (K, error)

func RunFuzzProvider[K any](f *testing.F, providerFunc func(*testing.T) FuzzableProvider[K]) {
	validRaw := [16]byte(RelayArgs.ExternalJobID)
	validRawBytes := make([]byte, 16)

	copy(validRawBytes, validRaw[:])

	f.Add([]byte{}, int32(-1), "ABC\xa8\x8c\xb3G\xfc", false, []byte{}, "", "", []byte{})                                                    // bad inputs
	f.Add(validRawBytes, int32(123), "testcontract", true, []byte(ConfigTOML), string(types.Median), "testtransmitter", []byte{100: 88})     // valid for MedianProvider
	f.Add(validRawBytes, int32(123), "testcontract", true, []byte(ConfigTOML), string(types.Mercury), "testtransmitter", []byte{100: 88})    // valid for MercuryProvider
	f.Add(validRawBytes, int32(123), "testcontract", true, []byte(ConfigTOML), string(types.Functions), "testtransmitter", []byte{100: 88})  // valid for FunctionsProvider
	f.Add(validRawBytes, int32(123), "testcontract", true, []byte(ConfigTOML), string(types.OCR2Keeper), "testtransmitter", []byte{100: 88}) // valid for AutomationProvider

	f.Fuzz(func(
		t *testing.T, fExtJobID []byte, fJobID int32, fContractID string, fNew bool,
		fConfig []byte, fType string, fTransmID string, fPlugConf []byte,
	) {
		var rawBytes [16]byte

		copy(rawBytes[:], fExtJobID)

		provider := providerFunc(t)
		ctx := tests.Context(t)
		fRelayArgs := types.RelayArgs{
			ExternalJobID: uuid.UUID(rawBytes),
			JobID:         fJobID,
			ContractID:    fContractID,
			New:           fNew,
			RelayConfig:   fConfig,
			ProviderType:  fType,
		}

		pArgs := types.PluginArgs{
			TransmitterID: fTransmID,
			PluginConfig:  fPlugConf,
		}

		_, err := provider(ctx, fRelayArgs, pArgs)

		grpcUnavailableErr(t, err)
	})
}

func grpcUnavailableErr(t *testing.T, err error) {
	t.Helper()

	if code := status.Code(err); code == codes.Unavailable {
		t.FailNow()
	}
}

type fuzzerKeystore struct {
	accounts      []string
	acctErr       bool
	signed        []byte
	signErr       bool
	valuesWithErr bool
	errStr        string
}

func (k fuzzerKeystore) Accounts(ctx context.Context) ([]string, error) {
	if k.acctErr {
		err := fmt.Errorf(k.errStr)

		if k.valuesWithErr {
			return k.accounts, err
		}

		return nil, err
	}

	return k.accounts, nil
}

// Sign returns data signed by account.
// nil data can be used as a no-op to check for account existence.
func (k fuzzerKeystore) Sign(ctx context.Context, account string, data []byte) ([]byte, error) {
	if k.signErr {
		err := fmt.Errorf(k.errStr)

		if k.valuesWithErr {
			return k.signed, err
		}

		return nil, err
	}

	return k.signed, nil
}
