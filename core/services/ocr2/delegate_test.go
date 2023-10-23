package ocr2_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
	ocr2validate "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestGetEVMEffectiveTransmitterID(t *testing.T) {
	customChainID := utils.NewBig(testutils.NewRandomEVMChainID())

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		enabled := true
		c.EVM = append(c.EVM, &evmcfg.EVMConfig{
			ChainID: customChainID,
			Chain:   evmcfg.Defaults(customChainID),
			Enabled: &enabled,
			Nodes:   evmcfg.EVMNodes{{}},
		})
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	require.NoError(t, keyStore.OCR2().Add(cltest.DefaultOCR2Key))
	lggr := logger.TestLogger(t)

	txManager := txmmocks.NewMockEvmTxManager(t)
	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config, KeyStore: keyStore.Eth(), TxManager: txManager})
	require.True(t, relayExtenders.Len() > 0)
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)

	type testCase struct {
		name                  string
		pluginType            types.OCR2PluginType
		transmitterID         null.String
		sendingKeys           []any
		expectedError         bool
		expectedTransmitterID string
		forwardingEnabled     bool
		getForwarderForEOAArg common.Address
		getForwarderForEOAErr bool
	}

	setTestCase := func(jb *job.Job, tc testCase, txManager *txmmocks.MockEvmTxManager) {
		jb.OCR2OracleSpec.PluginType = tc.pluginType
		jb.OCR2OracleSpec.TransmitterID = tc.transmitterID
		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = tc.sendingKeys
		jb.ForwardingAllowed = tc.forwardingEnabled

		if tc.forwardingEnabled && tc.getForwarderForEOAErr {
			txManager.Mock.On("GetForwarderForEOA", tc.getForwarderForEOAArg).Return(common.HexToAddress("0x0"), errors.New("random error")).Once()
		} else if tc.forwardingEnabled {
			txManager.Mock.On("GetForwarderForEOA", tc.getForwarderForEOAArg).Return(common.HexToAddress(tc.expectedTransmitterID), nil).Once()
		}
	}

	testCases := []testCase{
		{
			name:                  "mercury plugin should just return transmitterID",
			pluginType:            types.Mercury,
			transmitterID:         null.StringFrom("Mercury transmitterID"),
			expectedTransmitterID: "Mercury transmitterID",
		},
		{
			name:          "when transmitterID is not defined, it should validate that sending keys are defined",
			sendingKeys:   []any{},
			expectedError: true,
		},
		{
			name:          "when transmitterID is not defined, it should validate that plugin type is ocr2 vrf if more than 1 sending key is defined",
			sendingKeys:   []any{"0x7e57000000000000000000000000000000000001", "0x7e57000000000000000000000000000000000002", "0x7e57000000000000000000000000000000000003"},
			expectedError: true,
		},
		{
			name:                  "when transmitterID is not defined and plugin is ocr2vrf, it should allow>1 sendingKeys and set transmitterID to the first one",
			pluginType:            types.OCR2VRF,
			sendingKeys:           []any{"0x7e57000000000000000000000000000000000000", "0x7e57000000000000000000000000000000000001", "0x7e57000000000000000000000000000000000002"},
			expectedTransmitterID: "0x7e57000000000000000000000000000000000000",
		},
		{
			name:                  "when transmitterID is not defined, it should set transmitterID to first sendingKey",
			sendingKeys:           []any{"0x7e57000000000000000000000000000000000004"},
			expectedTransmitterID: "0x7e57000000000000000000000000000000000004",
		},
		{
			name:                  "when forwarders are enabled and when transmitterID is defined, it should default to using spec transmitterID to retrieve forwarder address",
			forwardingEnabled:     true,
			transmitterID:         null.StringFrom("0x7e57000000000000000000000000000000000001"),
			getForwarderForEOAArg: common.HexToAddress("0x7e57000000000000000000000000000000000001"),
			expectedTransmitterID: "0x7e58000000000000000000000000000000000000",
		},
		{
			name:                  "when forwarders are enabled and when transmitterID is not defined, it should use first sendingKey to retrieve forwarder address",
			pluginType:            types.OCR2VRF,
			forwardingEnabled:     true,
			sendingKeys:           []any{"0x7e57000000000000000000000000000000000001", "0x7e57000000000000000000000000000000000002"},
			getForwarderForEOAArg: common.HexToAddress("0x7e57000000000000000000000000000000000001"),
			expectedTransmitterID: "0x7e58000000000000000000000000000000000000",
		},
		{
			name:                  "when forwarders are enabled but forwarder address fails to be retrieved and when transmitterID is not defined, it should default to using first sendingKey",
			pluginType:            types.OCR2VRF,
			forwardingEnabled:     true,
			sendingKeys:           []any{"0x7e57000000000000000000000000000000000001", "0x7e57000000000000000000000000000000000002"},
			getForwarderForEOAArg: common.HexToAddress("0x7e57000000000000000000000000000000000001"),
			getForwarderForEOAErr: true,
			expectedTransmitterID: "0x7e57000000000000000000000000000000000001",
		},
		{
			name:                  "when forwarders are enabled but forwarder address fails to be retrieved and when transmitterID is defined, it should default to using spec transmitterID",
			forwardingEnabled:     true,
			transmitterID:         null.StringFrom("0x7e57000000000000000000000000000000000003"),
			getForwarderForEOAErr: true,
			getForwarderForEOAArg: common.HexToAddress("0x7e57000000000000000000000000000000000003"),
			expectedTransmitterID: "0x7e57000000000000000000000000000000000003",
		},
	}

	t.Run("when sending keys are not defined, the first one should be set to transmitterID", func(t *testing.T) {
		jb, err := ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
		require.NoError(t, err)
		jb.OCR2OracleSpec.TransmitterID = null.StringFrom("some transmitterID string")
		jb.OCR2OracleSpec.RelayConfig["sendingKeys"] = nil
		chain, err := legacyChains.Get(customChainID.String())
		require.NoError(t, err)
		effectiveTransmitterID, err := ocr2.GetEVMEffectiveTransmitterID(&jb, chain, lggr)
		require.NoError(t, err)
		require.Equal(t, "some transmitterID string", effectiveTransmitterID)
		require.Equal(t, []string{"some transmitterID string"}, jb.OCR2OracleSpec.RelayConfig["sendingKeys"].([]string))
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jb, err := ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
			require.NoError(t, err)
			setTestCase(&jb, tc, txManager)
			chain, err := legacyChains.Get(customChainID.String())
			require.NoError(t, err)

			effectiveTransmitterID, err := ocr2.GetEVMEffectiveTransmitterID(&jb, chain, lggr)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedTransmitterID, effectiveTransmitterID)
			// when forwarding is enabled effectiveTransmitter differs unless it failed to fetch forwarder address
			if !jb.ForwardingAllowed {
				require.Equal(t, jb.OCR2OracleSpec.TransmitterID.String, effectiveTransmitterID)
			}

		})
	}

	t.Run("when forwarders are enabled and chain retrieval fails, error should be handled", func(t *testing.T) {
		jb, err := ocr2validate.ValidatedOracleSpecToml(config.OCR2(), config.Insecure(), testspecs.GetOCR2EVMSpecMinimal())
		require.NoError(t, err)
		jb.ForwardingAllowed = true
		jb.OCR2OracleSpec.TransmitterID = null.StringFrom("0x7e57000000000000000000000000000000000001")
		chain, err := legacyChains.Get("not an id")
		require.Error(t, err)
		_, err = ocr2.GetEVMEffectiveTransmitterID(&jb, chain, lggr)
		require.Error(t, err)
	})
}
