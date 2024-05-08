package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keystoretest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/keystore/test"
	cciptest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip/test"
	mediantest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/median/test"
	mercurytest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/test"
	ocr3capabilitytest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ocr3capability/test"
	ocr2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2/test"
	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	looptypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

var chainStatus = types.ChainStatus{
	ID:      "some_chain",
	Enabled: true,
}

type transactionRequest struct {
	from         string
	to           string
	amount       *big.Int
	balanceCheck bool
}

type nodeRequest struct {
	pageSize  int32
	pageToken string
}

type nodeResponse struct {
	nodes    []types.NodeStatus
	nextPage string
	total    int
}
type staticPluginRelayerConfig struct {
	StaticChecks           bool
	relayArgs              types.RelayArgs
	pluginArgs             types.PluginArgs
	medianProvider         testtypes.MedianProviderTester
	agnosticProvider       testtypes.PluginProviderTester
	mercuryProvider        mercurytest.MercuryProviderTester
	executionProvider      cciptest.ExecProviderTester
	commitProvider         cciptest.CommitProviderTester
	configProvider         ocr2test.ConfigProviderTester
	ocr3CapabilityProvider testtypes.OCR3CapabilityProviderTester
	// Note: add other Provider testers here when we implement them
	// eg Functions, Automation, etc
	nodeRequest        nodeRequest
	nodeResponse       nodeResponse
	transactionRequest transactionRequest
	chainStatus        types.ChainStatus
}

func NewRelayerTester(staticChecks bool) testtypes.RelayerTester {
	return staticPluginRelayer{
		staticPluginRelayerConfig: staticPluginRelayerConfig{
			StaticChecks:           staticChecks,
			relayArgs:              RelayArgs,
			pluginArgs:             PluginArgs,
			medianProvider:         mediantest.MedianProvider,
			mercuryProvider:        mercurytest.MercuryProvider,
			executionProvider:      cciptest.ExecutionProvider,
			agnosticProvider:       ocr2test.AgnosticProvider,
			configProvider:         ocr2test.ConfigProvider,
			ocr3CapabilityProvider: ocr3capabilitytest.OCR3CapabilityProvider,
			nodeRequest: nodeRequest{
				pageSize:  137,
				pageToken: "",
			},
			nodeResponse: nodeResponse{
				nodes:    nodes,
				nextPage: "",
				total:    len(nodes),
			},
			transactionRequest: transactionRequest{
				from:         "me",
				to:           "you",
				amount:       big.NewInt(97),
				balanceCheck: true,
			},
			chainStatus: chainStatus,
		},
	}
}

type staticPluginRelayer struct {
	staticPluginRelayerConfig
}

func (s staticPluginRelayer) NewRelayer(ctx context.Context, config string, keystore core.Keystore) (looptypes.Relayer, error) {
	if s.StaticChecks && config != ConfigTOML {
		return nil, fmt.Errorf("expected config %q but got %q", ConfigTOML, config)
	}
	keys, err := keystore.Accounts(ctx)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("expected at least one key but got none")
	}

	return s, nil
}

func (s staticPluginRelayer) Start(ctx context.Context) error { return nil }

func (s staticPluginRelayer) Close() error { return nil }

func (s staticPluginRelayer) Ready() error { panic("unimplemented") }

func (s staticPluginRelayer) Name() string { panic("unimplemented") }

func (s staticPluginRelayer) HealthReport() map[string]error { panic("unimplemented") }

func (s staticPluginRelayer) NewConfigProvider(ctx context.Context, r types.RelayArgs) (types.ConfigProvider, error) {
	if s.StaticChecks && !equalRelayArgs(r, s.relayArgs) {
		return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", s.relayArgs, r)
	}
	return s.configProvider, nil
}

func (s staticPluginRelayer) NewMedianProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.MedianProvider, error) {
	if s.StaticChecks {
		ra := newRelayArgsWithProviderType(types.Median)
		if !equalRelayArgs(r, ra) {
			return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
		}
		if !reflect.DeepEqual(PluginArgs, p) {
			return nil, fmt.Errorf("expected plugin args %v but got %v", PluginArgs, p)
		}
	}

	return s.medianProvider, nil
}

func (s staticPluginRelayer) NewPluginProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.PluginProvider, error) {
	if s.StaticChecks {
		ra := newRelayArgsWithProviderType(types.Median)
		if !equalRelayArgs(r, ra) {
			return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
		}
		if !reflect.DeepEqual(PluginArgs, p) {
			return nil, fmt.Errorf("expected plugin args %v but got %v", PluginArgs, p)
		}
	}
	return s.agnosticProvider, nil
}

func (s staticPluginRelayer) NewOCR3CapabilityProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.OCR3CapabilityProvider, error) {
	if s.StaticChecks {
		ra := newRelayArgsWithProviderType(types.OCR3Capability)
		if !equalRelayArgs(r, ra) {
			return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", RelayArgs, r)
		}
		if !reflect.DeepEqual(PluginArgs, p) {
			return nil, fmt.Errorf("expected plugin args %v but got %v", PluginArgs, p)
		}
	}
	return s.ocr3CapabilityProvider, nil
}

func (s staticPluginRelayer) NewMercuryProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.MercuryProvider, error) {
	if s.StaticChecks {
		if !equalRelayArgs(r, mercurytest.RelayArgs) {
			return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", mercurytest.RelayArgs, r)
		}
		if !reflect.DeepEqual(mercurytest.PluginArgs, p) {
			return nil, fmt.Errorf("expected plugin args %v but got %v", mercurytest.PluginArgs, p)
		}
	}
	return s.mercuryProvider, nil
}

func (s staticPluginRelayer) NewExecutionProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.CCIPExecProvider, error) {
	if s.StaticChecks {
		if !equalRelayArgs(r, cciptest.ExecutionRelayArgs) {
			return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", cciptest.ExecutionRelayArgs, r)
		}
		if !reflect.DeepEqual(cciptest.ExecutionPluginArgs, p) {
			return nil, fmt.Errorf("expected plugin args %v but got %v", cciptest.ExecutionPluginArgs, p)
		}
	}
	return s.executionProvider, nil
}

func (s staticPluginRelayer) NewCommitProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.CCIPCommitProvider, error) {
	if s.StaticChecks {
		if !equalRelayArgs(r, cciptest.CommitRelayArgs) {
			return nil, fmt.Errorf("expected relay args:\n\t%v\nbut got:\n\t%v", cciptest.CommitRelayArgs, r)
		}
		if !reflect.DeepEqual(cciptest.CommitPluginArgs, p) {
			return nil, fmt.Errorf("expected plugin args %v but got %v", cciptest.CommitPluginArgs, p)
		}
	}
	return s.commitProvider, nil
}

func (s staticPluginRelayer) NewLLOProvider(ctx context.Context, r types.RelayArgs, p types.PluginArgs) (types.LLOProvider, error) {
	return nil, errors.New("not implemented")
}

func (s staticPluginRelayer) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	return s.chainStatus, nil
}

func (s staticPluginRelayer) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) ([]types.NodeStatus, string, int, error) {
	if s.StaticChecks && s.nodeRequest.pageSize != pageSize {
		return nil, "", -1, fmt.Errorf("expected page_size %d but got %d", s.nodeRequest.pageSize, pageSize)
	}
	if pageToken != "" {
		return nil, "", -1, fmt.Errorf("expected empty page_token but got %q", pageToken)
	}
	return s.nodeResponse.nodes, s.nodeResponse.nextPage, s.nodeResponse.total, nil
}

func (s staticPluginRelayer) Transact(ctx context.Context, f, t string, a *big.Int, b bool) error {
	if s.StaticChecks {
		if f != s.transactionRequest.from {
			return fmt.Errorf("expected from %s but got %s", s.transactionRequest.from, f)
		}
		if t != s.transactionRequest.to {
			return fmt.Errorf("expected to %s but got %s", s.transactionRequest.to, t)
		}
		if s.transactionRequest.amount.Cmp(a) != 0 {
			return fmt.Errorf("expected amount %s but got %s", s.transactionRequest.amount, a)
		}
		if b != s.transactionRequest.balanceCheck { //nolint:gosimple
			return fmt.Errorf("expected balance check %t but got %t", s.transactionRequest.balanceCheck, b)
		}
	}

	return nil
}

func (s staticPluginRelayer) AssertEqual(ctx context.Context, t *testing.T, relayer looptypes.Relayer) {
	t.Run("ConfigProvider", func(t *testing.T) {
		t.Parallel()
		ctx := tests.Context(t)
		configProvider, err := relayer.NewConfigProvider(ctx, RelayArgs)
		require.NoError(t, err)
		require.NoError(t, configProvider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, configProvider.Close()) })

		s.configProvider.AssertEqual(ctx, t, configProvider)
	})

	t.Run("MedianProvider", func(t *testing.T) {
		t.Parallel()
		ctx := tests.Context(t)
		ra := newRelayArgsWithProviderType(types.Median)
		p, err := relayer.NewPluginProvider(ctx, ra, PluginArgs)
		require.NoError(t, err)
		require.NotNil(t, p)
		provider := p.(types.MedianProvider)
		require.NoError(t, provider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, provider.Close()) })

		t.Run("ReportingPluginProvider", func(t *testing.T) {
			t.Parallel()
			s.medianProvider.AssertEqual(ctx, t, provider)
		})
	})

	t.Run("PluginProvider", func(t *testing.T) {
		t.Parallel()
		ctx := tests.Context(t)
		ra := newRelayArgsWithProviderType(types.GenericPlugin)
		provider, err := relayer.NewPluginProvider(ctx, ra, PluginArgs)
		require.NoError(t, err)
		require.NoError(t, provider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, provider.Close()) })
		t.Run("ReportingPluginProvider", func(t *testing.T) {
			t.Parallel()
			ctx := tests.Context(t)
			s.agnosticProvider.AssertEqual(ctx, t, provider)
		})
	})

	t.Run("GetChainStatus", func(t *testing.T) {
		t.Parallel()
		ctx := tests.Context(t)
		gotChain, err := relayer.GetChainStatus(ctx)
		require.NoError(t, err)
		assert.Equal(t, s.chainStatus, gotChain)
	})

	t.Run("ListNodeStatuses", func(t *testing.T) {
		t.Parallel()
		ctx := tests.Context(t)
		gotNodes, gotNextToken, gotCount, err := relayer.ListNodeStatuses(ctx, s.nodeRequest.pageSize, s.nodeRequest.pageToken)
		require.NoError(t, err)
		assert.Equal(t, s.nodeResponse.nodes, gotNodes)
		assert.Equal(t, s.nodeResponse.total, gotCount)
		assert.Empty(t, s.nodeResponse.nextPage, gotNextToken)
	})

	t.Run("Transact", func(t *testing.T) {
		t.Parallel()
		ctx := tests.Context(t)
		err := relayer.Transact(ctx, s.transactionRequest.from, s.transactionRequest.to, s.transactionRequest.amount, s.transactionRequest.balanceCheck)
		require.NoError(t, err)
	})
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

func RunPlugin(t *testing.T, p looptypes.PluginRelayer) {
	t.Run("Relayer", func(t *testing.T) {
		ctx := tests.Context(t)
		relayer, err := p.NewRelayer(ctx, ConfigTOML, keystoretest.Keystore)
		require.NoError(t, err)
		require.NoError(t, relayer.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, relayer.Close()) })
		Run(t, relayer)
	})
}

func Run(t *testing.T, relayer looptypes.Relayer) {
	ctx := tests.Context(t)
	expectedRelayer := NewRelayerTester(false)
	expectedRelayer.AssertEqual(ctx, t, relayer)
}

func RunFuzzPluginRelayer(f *testing.F, relayerFunc func(*testing.T) looptypes.PluginRelayer) {
	var (
		account = "testaccount"
		signed  = []byte{5: 11}
	)
	f.Add("ABC\xa8\x8c\xb3G\xfc", "", true, []byte{}, true, true, "")
	f.Add(ConfigTOML, account, false, signed, false, false, "")

	// nolint: gocognit
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

func RunFuzzRelayer(f *testing.F, relayerFunc func(*testing.T) looptypes.Relayer) {
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
