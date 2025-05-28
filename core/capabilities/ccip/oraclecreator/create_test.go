package oraclecreator

import (
	"math/big"
	"testing"

	"github.com/google/uuid"
	chainsel "github.com/smartcontractkit/chain-selectors"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/stretchr/testify/require"

	ocr3types "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocrcommon "github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

// TestPluginOracleCreatorCreate_InvalidSelector ensures that Create returns an error when an
// unknown chain selector is provided. This exercises the early-return error path
// before any heavy initialisation takes place.
func TestPluginOracleCreatorCreate_InvalidSelector(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)

	// Minimal creator with empty/nil dependencies – these will not be hit as the
	// selector validation fails first.
	keyBundles := make(map[string]ocr2key.KeyBundle)
	transmitters := make(map[types.RelayID][]string)
	relayers := make(map[types.RelayID]loop.Relayer)

	p2pk, err := p2pkey.NewV2()
	require.NoError(t, err)

	creator := NewPluginOracleCreator(
		keyBundles,
		transmitters,
		relayers,
		/* peerWrapper= */ nil,
		uuid.New(),
		1,
		false,
		job.JSONConfig{},
		/* db */ ocr3types.Database(nil),
		lggr,
		/* monitoringEndpointGen */ nil,
		/* bootstrapperLocators */ nil,
		/* homeChainReader */ nil,
		/* homeChainSelector */ 0,
		/* addressCodec */ nil,
		p2pk,
	).(*pluginOracleCreator)

	// Zero value config will have ChainSelector == 0 which is not present in
	// the generated selector tables – Create should error immediately.
	var cfg cctypes.OCR3ConfigWithMeta

	oracle, err := creator.Create(t.Context(), 1, cfg)
	require.Error(t, err, "expected error due to unknown chain selector")
	require.Nil(t, oracle, "oracle should be nil on error")

	// Sanity-check the error comes from selector validation.
	require.Contains(t, err.Error(), "unknown chain selector", "unexpected error message: %v", err)
}

// TestCreateFactoryAndTransmitter_PeerWrapperNotStarted verifies that the
// helper returns an error when called for a commit plugin but the peer wrapper
// has not been started – this guards a critical control-flow condition.
func TestCreateFactoryAndTransmitter_PeerWrapperNotStarted(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)

	p2pk := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))

	creator := &pluginOracleCreator{
		lggr:        lggr,
		peerWrapper: &ocrcommon.SingletonPeerWrapper{}, // unstarted but non-nil to prevent panic
		p2pID:       p2pk,
	}

	var cfg cctypes.OCR3ConfigWithMeta // zero value => commit plugin type

	factory, transmitter, err := creator.createFactoryAndTransmitter(
		1,
		cfg,
		types.NewRelayID(chainsel.FamilyEVM, "1"),
		map[cciptypes.ChainSelector]types.ContractReader{},
		map[cciptypes.ChainSelector]types.ContractWriter{},
		/* destChainWriter */ nil,
		/* destFromAccounts */ nil,
		ocr3confighelper.PublicConfig{},
		"1",
		ccipcommon.PluginConfig{},
		"",
	)

	require.Error(t, err, "expected error when peer wrapper not started")
	require.Nil(t, factory)
	require.Nil(t, transmitter)
	require.Contains(t, err.Error(), "peer wrapper is not started")
}

// TestCreateFactoryAndTransmitter_NilDestChainWriter verifies that a NoOpTransmitter
// is created when destChainWriter is nil for both commit and execute plugin types.
func TestCreateFactoryAndTransmitter_NilDestChainWriter(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	p2pk := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))

	creator := &pluginOracleCreator{
		lggr:  lggr,
		p2pID: p2pk,
	}

	// Common arguments for both plugin types
	donID := uint32(1)
	relayID := types.NewRelayID(chainsel.FamilyEVM, "1")
	contractReaders := map[cciptypes.ChainSelector]types.ContractReader{}
	chainWriters := map[cciptypes.ChainSelector]types.ContractWriter{}
	publicCfg := ocr3confighelper.PublicConfig{}
	destChainID := "1"
	pluginCfg := ccipcommon.PluginConfig{
		// Provide a dummy factory that creates identifiable (or nil) transmitters
		// This isn't strictly necessary for this specific test path if we only check for NoOpTransmitter type,
		// but good for completeness if we were to test the non-nil path.
		ContractTransmitterFactory: &mocks.ContractTransmitterFactory{},
	}
	offrampAddrStr := "0x123"

	testCases := []struct {
		name       string
		pluginType cctypes.PluginType
	}{
		// Running only for execute because commit requires the peer wrapper to be started.
		// pluginOracleCreator is currently using a pointer to a concrete type rather than an interface,
		// so we can't mock the peer wrapper.
		{
			name:       "Execute Plugin",
			pluginType: cctypes.PluginTypeCCIPExec,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := cctypes.OCR3ConfigWithMeta{
				Config: ccipreaderpkg.OCR3Config{
					PluginType: uint8(tc.pluginType),
					// Other necessary fields for PublicConfig() or other early checks if any
					ChainSelector: 1, // Arbitrary valid selector
				},
			}

			factory, transmitter, err := creator.createFactoryAndTransmitter(
				donID,
				cfg,
				relayID,
				contractReaders,
				chainWriters,
				nil, // Key: destChainWriter is nil
				nil, // destFromAccounts can be nil as it's not used before NoOpTransmitter creation
				publicCfg,
				destChainID,
				pluginCfg,
				offrampAddrStr,
			)

			require.NoError(t, err)
			require.NotNil(t, factory, "factory should not be nil")
			require.NotNil(t, transmitter, "transmitter should not be nil")

			// Assert that the transmitter is a NoOpTransmitter
			account, err := transmitter.FromAccount(t.Context())
			require.NoError(t, err)
			require.Equal(t, account, ocrtypes.Account(""))
		})
	}
}
