package ocrimpls_test

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/jmoiron/sqlx"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/multi_ocr3_helper"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmconfig "github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"    // Register EVM plugin config factories
	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana" // Register Solana plugin config factories
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	kschaintype "github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func Test_ContractTransmitter_TransmitWithoutSignatures(t *testing.T) {
	type testCase struct {
		name                string
		pluginType          uint8
		withSigs            bool
		expectedSigsEnabled bool
		report              []byte
	}

	testCases := []testCase{
		{
			"empty report with sigs",
			uint8(cctypes.PluginTypeCCIPCommit),
			true,
			true,
			[]byte{},
		},
		{
			"empty report without sigs",
			uint8(cctypes.PluginTypeCCIPExec),
			false,
			false,
			[]byte{},
		},
		{
			"report with data with sigs",
			uint8(cctypes.PluginTypeCCIPCommit),
			true,
			true,
			randomReport(t, 96),
		},
		{
			"report with data without sigs",
			uint8(cctypes.PluginTypeCCIPExec),
			false,
			false,
			randomReport(t, 96),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			testTransmitter(t, tc.pluginType, tc.withSigs, tc.expectedSigsEnabled, tc.report)
		})
	}
}

func testTransmitter(
	t *testing.T,
	pluginType uint8,
	withSigs bool,
	expectedSigsEnabled bool,
	report []byte,
) {
	uni := newTestUniverse(t, nil)

	c, err := uni.wrapper.LatestConfigDetails(nil, pluginType)
	require.NoError(t, err, "failed to get latest config details")
	configDigest := c.ConfigInfo.ConfigDigest
	require.Equal(t, expectedSigsEnabled, c.ConfigInfo.IsSignatureVerificationEnabled, "signature verification enabled setting not correct")

	// set the plugin type on the helper so it fetches the right config info.
	// the important aspect is whether signatures should be enabled or not.
	_, err = uni.wrapper.SetTransmitOcrPluginType(uni.deployer, pluginType)
	require.NoError(t, err, "failed to set plugin type")
	uni.backend.Commit()

	// create attributed sigs
	// only need f+1 which is 2 in this case
	rwi := ocr3types.ReportWithInfo[[]byte]{
		Report: report,
		Info:   []byte{},
	}
	seqNr := uint64(1)
	attributedSigs := uni.SignReport(t, configDigest, rwi, seqNr)

	account, err := uni.transmitterWithSigs.FromAccount(t.Context())
	require.NoError(t, err, "failed to get from account")
	require.Equal(t, ocrtypes.Account(uni.transmitters[0].Hex()), account, "from account mismatch")
	if withSigs {
		err = uni.transmitterWithSigs.Transmit(testutils.Context(t), configDigest, seqNr, rwi, attributedSigs)
	} else {
		err = uni.transmitterWithoutSigs.Transmit(testutils.Context(t), configDigest, seqNr, rwi, attributedSigs)
	}
	require.NoError(t, err, "failed to transmit")
	uni.backend.Commit()

	var txStatus uint64
	require.Eventually(t, func() bool {
		uni.backend.Commit()
		rows, err := uni.db.QueryContext(testutils.Context(t), `SELECT hash FROM evm.tx_attempts LIMIT 1`)
		require.NoError(t, err, "failed to query txes")
		defer rows.Close()
		var txHash []byte
		for rows.Next() {
			require.NoError(t, rows.Scan(&txHash), "failed to scan")
		}
		t.Log("txHash:", txHash)
		receipt, err := uni.simClient.TransactionReceipt(testutils.Context(t), common.BytesToHash(txHash))
		if err != nil {
			t.Log("tx not found yet:", hexutil.Encode(txHash))
			return false
		}
		t.Log("tx found:", hexutil.Encode(txHash), "status:", receipt.Status)
		txStatus = receipt.Status
		return true
	}, testutils.WaitTimeout(t), 1*time.Second)

	// wait for receipt to be written to the db
	require.Eventually(t, func() bool {
		uni.backend.Commit()
		var count uint32
		err := uni.db.GetContext(testutils.Context(t), &count, `SELECT count(*) as cnt FROM evm.receipts LIMIT 1`)
		require.NoError(t, err)
		if count == 1 {
			t.Log("tx receipt found in db")
		}
		return count == 1
	}, testutils.WaitTimeout(t), 2*time.Second)

	require.Equal(t, uint64(1), txStatus, "tx status should be success")

	// check that the event was emitted
	events := uni.TransmittedEvents(t)
	require.Len(t, events, 1, "expected 1 event")
	require.Equal(t, configDigest, events[0].ConfigDigest, "config digest mismatch")
	require.Equal(t, seqNr, events[0].SequenceNumber, "seq num mismatch")
}

func abiEncodeUint32(data uint32) ([]byte, error) {
	return utils.ABIEncode(`[{ "type": "uint32" }]`, data)
}

// Test EVM -> SVM extra data decoding in contract transmitter
func TestSVMExecCallDataFuncExtraDataDecoding(t *testing.T) {
	extraDataCodec := ccipcommon.ExtraDataCodec(map[string]ccipcommon.SourceChainExtraDataCodec{
		chainsel.FamilyEVM:    ccipevm.ExtraDataDecoder{},
		chainsel.FamilySolana: ccipsolana.ExtraDataDecoder{},
	})
	t.Run("fails when multiple reports are included", func(t *testing.T) {
		reports := []ccipocr3.ExecutePluginReportSingleChain{{}, {}}
		reportWithInfo := ccipocr3.ExecuteReportInfo{
			AbstractReports: reports,
		}

		encodedExecReport, err := reportWithInfo.Encode()
		require.NoError(t, err)

		rwi := ocr3types.ReportWithInfo[[]byte]{
			Report: randomReport(t, 96),
			Info:   encodedExecReport,
		}
		_, _, _, err = ocrimpls.SVMExecCalldataFunc([2][32]byte{}, rwi, nil, nil, [32]byte{}, extraDataCodec)
		require.Contains(t, err.Error(), "unexpected report length, expected 1, got 2")
	})
	t.Run("fails when multiple report contains multiple messages", func(t *testing.T) {
		reports := []ccipocr3.ExecutePluginReportSingleChain{{
			Messages: []ccipocr3.Message{{}, {}},
		}}
		reportWithInfo := ccipocr3.ExecuteReportInfo{
			AbstractReports: reports,
		}

		encodedExecReport, err := reportWithInfo.Encode()
		require.NoError(t, err)

		rwi := ocr3types.ReportWithInfo[[]byte]{
			Report: randomReport(t, 96),
			Info:   encodedExecReport,
		}
		_, _, _, err = ocrimpls.SVMExecCalldataFunc([2][32]byte{}, rwi, nil, nil, [32]byte{}, extraDataCodec)
		require.Contains(t, err.Error(), "unexpected message length, expected 1, got 2")
	})
	t.Run("fails with invalid extra args", func(t *testing.T) {
		// invalid encoded extra args
		encoded := []byte{1, 2, 3, 4}

		report := ccipocr3.ExecutePluginReportSingleChain{
			SourceChainSelector: 5009297550715157269,
			Messages: []ccipocr3.Message{{
				Header: ccipocr3.RampMessageHeader{
					// EVM
					SourceChainSelector: 5009297550715157269,
					// to SOL
					DestChainSelector: 124615329519749607,
				},
				ExtraArgs: encoded,
			}},
		}

		reportWithInfo := ccipocr3.ExecuteReportInfo{
			AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{report},
		}

		encodedExecReport, err := reportWithInfo.Encode()
		require.NoError(t, err)

		rwi := ocr3types.ReportWithInfo[[]byte]{
			Report: randomReport(t, 96),
			Info:   encodedExecReport,
		}

		_, _, _, err = ocrimpls.SVMExecCalldataFunc([2][32]byte{}, rwi, nil, nil, [32]byte{}, extraDataCodec)
		require.Contains(t, err.Error(), "unknown extra args tag")
	})
	t.Run("fails with invalid extra exec data", func(t *testing.T) {
		// invalid encoded extra args
		encoded := []byte{31, 59, 58, 186, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170}
		encodedExecData := []byte{1, 2, 3, 4}

		report := ccipocr3.ExecutePluginReportSingleChain{
			SourceChainSelector: 5009297550715157269,
			Messages: []ccipocr3.Message{{
				Header: ccipocr3.RampMessageHeader{
					// EVM
					SourceChainSelector: 5009297550715157269,
					// to SOL
					DestChainSelector: 124615329519749607,
				},
				ExtraArgs: encoded,
				TokenAmounts: []ccipocr3.RampTokenAmount{{
					DestExecData: encodedExecData,
				}},
			}},
		}

		reportWithInfo := ccipocr3.ExecuteReportInfo{
			AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{report},
		}

		encodedExecReport, err := reportWithInfo.Encode()
		require.NoError(t, err)

		rwi := ocr3types.ReportWithInfo[[]byte]{
			Report: randomReport(t, 96),
			Info:   encodedExecReport,
		}

		_, _, _, err = ocrimpls.SVMExecCalldataFunc([2][32]byte{}, rwi, nil, nil, [32]byte{}, extraDataCodec)
		require.Contains(t, err.Error(), "failed to decode token amount dest exec data: decode dest gas amount: abi decode uint32: abi: cannot marshal in to go type: length insufficient 4 require 32")
	})
	t.Run("Successfully decodes valid EVM -> SOL report", func(t *testing.T) {
		// hardcode abi encoded extra args for simplicity
		encoded := []byte{31, 59, 58, 186, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170}
		destGasAmount := uint32(10000)
		encodedExecData, err := abiEncodeUint32(destGasAmount)
		require.NoError(t, err)

		report := ccipocr3.ExecutePluginReportSingleChain{
			SourceChainSelector: 5009297550715157269,
			Messages: []ccipocr3.Message{{
				Header: ccipocr3.RampMessageHeader{
					// EVM
					SourceChainSelector: 5009297550715157269,
					// to SOL
					DestChainSelector: 124615329519749607,
				},
				ExtraArgs: encoded,
				TokenAmounts: []ccipocr3.RampTokenAmount{{
					DestExecData: encodedExecData,
				}},
			}},
		}

		reportWithInfo := ccipocr3.ExecuteReportInfo{
			AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{report},
		}

		encodedExecReport, err := reportWithInfo.Encode()
		require.NoError(t, err)

		rwi := ocr3types.ReportWithInfo[[]byte]{
			Report: randomReport(t, 96),
			Info:   encodedExecReport,
		}

		_, _, args, err := ocrimpls.SVMExecCalldataFunc([2][32]byte{}, rwi, nil, nil, [32]byte{}, extraDataCodec)
		require.NoError(t, err)

		expectedArgs, ok := args.(ocrimpls.SVMExecCallArgs)
		require.True(t, ok)

		require.Equal(t, uint64(0x4), expectedArgs.ExtraData.ExtraArgsDecoded["accountIsWritableBitmap"])
		require.Equal(t, [32]uint8{44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170}, expectedArgs.ExtraData.ExtraArgsDecoded["accounts"].([][32]byte)[0])
		require.False(t, expectedArgs.ExtraData.ExtraArgsDecoded["allowOutOfOrderExecution"].(bool))
		require.Equal(t, destGasAmount, expectedArgs.ExtraData.ExtraArgsDecoded["computeUnits"])
		require.Equal(t, [32]uint8{44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170}, expectedArgs.ExtraData.ExtraArgsDecoded["tokenReceiver"])
		require.Equal(t, destGasAmount, expectedArgs.ExtraData.DestExecDataDecoded[0]["destGasAmount"])
	})
}

type testUniverse[RI any] struct {
	simClient              *client.SimulatedBackendClient
	backend                *simulated.Backend
	deployer               *bind.TransactOpts
	transmitters           []common.Address
	signers                []common.Address
	wrapper                *multi_ocr3_helper.MultiOCR3Helper
	transmitterWithSigs    ocr3types.ContractTransmitter[RI]
	transmitterWithoutSigs ocr3types.ContractTransmitter[RI]
	keyrings               []ocr3types.OnchainKeyring[RI]
	f                      uint8
	db                     *sqlx.DB
	txm                    txmgr.TxManager
	gasEstimator           gas.EvmFeeEstimator
}

type keyringsAndSigners[RI any] struct {
	keyrings []ocr3types.OnchainKeyring[RI]
	signers  []common.Address
}

func newTestUniverse(t *testing.T, ks *keyringsAndSigners[[]byte]) *testUniverse[[]byte] {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	owner := evmtestutils.MustNewSimTransactor(t)

	keyStore := keystest.NewMemoryChainStore()
	// create many transmitters but only need to fund one, rest are to get
	// setOCR3Config to pass.
	chainStore := keys.NewChainStore(keyStore, big.NewInt(1337))
	var transmitters []common.Address
	for i := 0; i < 4; i++ {
		addr, err := keyStore.Create()
		require.NoError(t, err, "failed to create key")
		transmitters = append(transmitters, addr)
	}

	backend := simulated.NewBackend(types.GenesisAlloc{
		owner.From: types.Account{
			Balance: assets.Ether(1000).ToInt(),
		},
		transmitters[0]: types.Account{
			Balance: assets.Ether(1000).ToInt(),
		},
	}, simulated.WithBlockGasLimit(30e6))

	ocr3HelperAddr, _, _, err := multi_ocr3_helper.DeployMultiOCR3Helper(owner, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	wrapper, err := multi_ocr3_helper.NewMultiOCR3Helper(ocr3HelperAddr, backend.Client())
	require.NoError(t, err)

	// create the oracle identities for setConfig
	// need to create at least 4 identities otherwise setConfig will fail
	var (
		keyrings []ocr3types.OnchainKeyring[[]byte]
		signers  []common.Address
	)
	if ks != nil {
		keyrings = ks.keyrings
		signers = ks.signers
	} else {
		for i := 0; i < 4; i++ {
			kb, err2 := ocr2key.New(kschaintype.EVM)
			require.NoError(t, err2, "failed to create key")
			kr := ocrimpls.NewOnchainKeyring[[]byte](kb, logger.TestLogger(t))
			signers = append(signers, common.BytesToAddress(kr.PublicKey()))
			keyrings = append(keyrings, kr)
		}
	}
	f := uint8(1)
	commitConfigDigest := testutils.Random32Byte()
	execConfigDigest := testutils.Random32Byte()
	_, err = wrapper.SetOCR3Configs(
		owner,
		[]multi_ocr3_helper.MultiOCR3BaseOCRConfigArgs{
			{
				ConfigDigest:                   commitConfigDigest,
				OcrPluginType:                  uint8(cctypes.PluginTypeCCIPCommit),
				F:                              f,
				IsSignatureVerificationEnabled: true,
				Signers:                        signers,
				Transmitters: []common.Address{
					transmitters[0],
					transmitters[1],
					transmitters[2],
					transmitters[3],
				},
			},
			{
				ConfigDigest:                   execConfigDigest,
				OcrPluginType:                  uint8(cctypes.PluginTypeCCIPExec),
				F:                              f,
				IsSignatureVerificationEnabled: false,
				Signers:                        signers,
				Transmitters: []common.Address{
					transmitters[0],
					transmitters[1],
					transmitters[2],
					transmitters[3],
				},
			},
		},
	)
	require.NoError(t, err)
	backend.Commit()

	commitConfig, err := wrapper.LatestConfigDetails(nil, uint8(cctypes.PluginTypeCCIPCommit))
	require.NoError(t, err, "failed to get latest commit config")
	require.Equal(t, commitConfigDigest, commitConfig.ConfigInfo.ConfigDigest, "commit config digest mismatch")
	execConfig, err := wrapper.LatestConfigDetails(nil, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err, "failed to get latest exec config")
	require.Equal(t, execConfigDigest, execConfig.ConfigInfo.ConfigDigest, "exec config digest mismatch")

	simClient := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)

	// create the chain writer service
	txm, gasEstimator := makeTestEvmTxm(t, db, simClient, chainStore)
	require.NoError(t, txm.Start(testutils.Context(t)), "failed to start tx manager")
	t.Cleanup(func() { require.NoError(t, txm.Close()) })

	chainWriter, err := evm.NewChainWriterService(
		logger.TestLogger(t),
		simClient,
		txm,
		gasEstimator,
		chainWriterConfigRaw(transmitters[0], assets.GWei(1)))
	require.NoError(t, err, "failed to create chain writer")
	require.NoError(t, chainWriter.Start(testutils.Context(t)), "failed to start chain writer")
	t.Cleanup(func() { require.NoError(t, chainWriter.Close()) })

	lggr := logger.TestLogger(t)
	transmitterWithSigs := ocrimpls.XXXNewContractTransmitterTestsOnly(
		lggr,
		chainWriter,
		ocrtypes.Account(transmitters[0].Hex()),
		contractName,
		methodTransmitWithSignatures,
		ocr3HelperAddr.Hex(),
		ocrimpls.NewEVMCommitCalldataFunc(consts.MethodCommit),
	)
	transmitterWithoutSigs := ocrimpls.XXXNewContractTransmitterTestsOnly(
		lggr,
		chainWriter,
		ocrtypes.Account(transmitters[0].Hex()),
		contractName,
		methodTransmitWithoutSignatures,
		ocr3HelperAddr.Hex(),
		ocrimpls.EVMExecCallDataFunc,
	)

	return &testUniverse[[]byte]{
		simClient:              simClient,
		backend:                backend,
		deployer:               owner,
		transmitters:           transmitters,
		signers:                signers,
		wrapper:                wrapper,
		transmitterWithSigs:    transmitterWithSigs,
		transmitterWithoutSigs: transmitterWithoutSigs,
		keyrings:               keyrings,
		f:                      f,
		db:                     db,
		txm:                    txm,
		gasEstimator:           gasEstimator,
	}
}

func (uni testUniverse[RI]) SignReport(t *testing.T, configDigest ocrtypes.ConfigDigest, rwi ocr3types.ReportWithInfo[RI], seqNum uint64) []ocrtypes.AttributedOnchainSignature {
	var attributedSigs []ocrtypes.AttributedOnchainSignature
	for i := uint8(0); i < uni.f+1; i++ {
		t.Log("signing report with", hexutil.Encode(uni.keyrings[i].PublicKey()))
		sig, err := uni.keyrings[i].Sign(configDigest, seqNum, rwi)
		require.NoError(t, err, "failed to sign report")
		attributedSigs = append(attributedSigs, ocrtypes.AttributedOnchainSignature{
			Signature: sig,
			Signer:    commontypes.OracleID(i),
		})
	}
	return attributedSigs
}

func (uni testUniverse[RI]) TransmittedEvents(t *testing.T) []*multi_ocr3_helper.MultiOCR3HelperTransmitted {
	iter, err := uni.wrapper.FilterTransmitted(&bind.FilterOpts{
		Start: 0,
	}, nil)
	require.NoError(t, err, "failed to create filter iterator")
	var events []*multi_ocr3_helper.MultiOCR3HelperTransmitted
	for iter.Next() {
		event := iter.Event
		events = append(events, event)
	}
	return events
}

func randomReport(t *testing.T, len int) []byte {
	report := make([]byte, len)
	_, err := rand.Reader.Read(report)
	require.NoError(t, err, "failed to read random bytes")
	return report
}

const (
	contractName                    = "MultiOCR3Helper"
	methodTransmitWithSignatures    = "TransmitWithSignatures"
	methodTransmitWithoutSignatures = "TransmitWithoutSignatures"
)

func chainWriterConfigRaw(fromAddress common.Address, maxGasPrice *assets.Wei) evmrelaytypes.ChainWriterConfig {
	return evmrelaytypes.ChainWriterConfig{
		Contracts: map[string]*evmrelaytypes.ContractConfig{
			contractName: {
				ContractABI: multi_ocr3_helper.MultiOCR3HelperABI,
				Configs: map[string]*evmrelaytypes.ChainWriterDefinition{
					methodTransmitWithSignatures: {
						ChainSpecificName: "transmitWithSignatures",
						GasLimit:          1e6,
						FromAddress:       fromAddress,
					},
					methodTransmitWithoutSignatures: {
						ChainSpecificName: "transmitWithoutSignatures",
						GasLimit:          1e6,
						FromAddress:       fromAddress,
					},
				},
			},
		},
		MaxGasPrice: maxGasPrice,
	}
}

func makeTestEvmTxm(t *testing.T, db *sqlx.DB, ethClient client.Client, keyStore keys.ChainStore) (txmgr.TxManager, gas.EvmFeeEstimator) {
	config, dbConfig, evmConfig := MakeTestConfigs(t)

	estimator, err := gas.NewEstimator(logger.TestLogger(t), ethClient, config.ChainType(), ethClient.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(t, err, "failed to create gas estimator")

	lggr := logger.TestLogger(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RPCBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}

	chainID := big.NewInt(1337)
	headSaver := heads.NewSaver(
		logger.NullLogger,
		heads.NewORM(*chainID, db),
		evmConfig,
		evmConfig.HeadTrackerConfig,
	)

	broadcaster := heads.NewBroadcaster(logger.NullLogger)
	require.NoError(t, broadcaster.Start(testutils.Context(t)), "failed to start head broadcaster")
	t.Cleanup(func() { require.NoError(t, broadcaster.Close()) })

	ht := heads.NewTracker(
		logger.NullLogger,
		ethClient,
		evmConfig,
		evmConfig.HeadTrackerConfig,
		broadcaster,
		headSaver,
		mailbox.NewMonitor("contract_transmitter_test", logger.NullLogger),
	)
	require.NoError(t, ht.Start(testutils.Context(t)), "failed to start head tracker")
	t.Cleanup(func() { require.NoError(t, ht.Close()) })

	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, logger.NullLogger),
		ethClient, logger.NullLogger, ht, lpOpts)
	require.NoError(t, lp.Start(testutils.Context(t)), "failed to start log poller")
	t.Cleanup(func() { require.NoError(t, lp.Close()) })

	// logic for building components (from evm/evm_txm.go) -------
	lggr.Infow("Initializing EVM transaction manager",
		"bumpTxDepth", evmConfig.GasEstimator().BumpTxDepth(),
		"maxInFlightTransactions", config.EvmConfig.Transactions().MaxInFlight(),
		"maxQueuedTransactions", config.EvmConfig.Transactions().MaxQueued(),
		"nonceAutoSync", evmConfig.NonceAutoSync(),
		"limitDefault", evmConfig.GasEstimator().LimitDefault(),
	)

	txm, err := txmgr.NewTxm(
		db,
		config,
		config.EvmConfig.GasEstimator(),
		config.EvmConfig.Transactions(),
		nil,
		dbConfig,
		dbConfig.Listener(),
		ethClient,
		lggr,
		lp,
		keyStore,
		estimator,
		ht,
		nil)
	require.NoError(t, err, "can't create tx manager")

	_, unsub := broadcaster.Subscribe(txm)
	t.Cleanup(unsub)

	return txm, estimator
}

// Code below copied/pasted and slightly modified in order to work from core/chains/evm/txmgr/test_helpers.go.

func ptr[T any](t T) *T { return &t }

type TestDatabaseConfig struct {
	config.Database
	defaultQueryTimeout time.Duration
}

func (d *TestDatabaseConfig) DefaultQueryTimeout() time.Duration {
	return d.defaultQueryTimeout
}

func (d *TestDatabaseConfig) LogSQL() bool {
	return false
}

type TestListenerConfig struct {
	config.Listener
}

func (l *TestListenerConfig) FallbackPollInterval() time.Duration {
	return 1 * time.Minute
}

func (d *TestDatabaseConfig) Listener() config.Listener {
	return &TestListenerConfig{}
}

type TestHeadTrackerConfig struct{}

// FinalityTagBypass implements config.HeadTracker.
func (t *TestHeadTrackerConfig) FinalityTagBypass() bool {
	return false
}

// HistoryDepth implements config.HeadTracker.
func (t *TestHeadTrackerConfig) HistoryDepth() uint32 {
	return 50
}

// MaxAllowedFinalityDepth implements config.HeadTracker.
func (t *TestHeadTrackerConfig) MaxAllowedFinalityDepth() uint32 {
	return 100
}

// MaxBufferSize implements config.HeadTracker.
func (t *TestHeadTrackerConfig) MaxBufferSize() uint32 {
	return 100
}

// SamplingInterval implements config.HeadTracker.
func (t *TestHeadTrackerConfig) SamplingInterval() time.Duration {
	return 1 * time.Second
}

func (t *TestHeadTrackerConfig) PersistenceEnabled() bool {
	return true
}

var _ evmconfig.HeadTracker = (*TestHeadTrackerConfig)(nil)

type TestEvmConfig struct {
	evmconfig.EVM
	HeadTrackerConfig    evmconfig.HeadTracker
	MaxInFlight          uint32
	ReaperInterval       time.Duration
	ReaperThreshold      time.Duration
	ResendAfterThreshold time.Duration
	BumpThreshold        uint64
	MaxQueued            uint64
	Enabled              bool
	Threshold            uint32
	MinAttempts          uint32
	DetectionApiUrl      *url.URL
}

func (e *TestEvmConfig) FinalityTagEnabled() bool {
	return false
}

func (e *TestEvmConfig) FinalityDepth() uint32 {
	return 42
}

func (e *TestEvmConfig) FinalizedBlockOffset() uint32 {
	return 42
}

func (e *TestEvmConfig) BlockEmissionIdleWarningThreshold() time.Duration {
	return 10 * time.Second
}

func (e *TestEvmConfig) Transactions() evmconfig.Transactions {
	return &transactionsConfig{e: e, autoPurge: &autoPurgeConfig{}}
}

func (e *TestEvmConfig) NonceAutoSync() bool { return true }

func (e *TestEvmConfig) ChainType() chaintype.ChainType { return "" }

type TestGasEstimatorConfig struct {
	bumpThreshold uint64
}

func (g *TestGasEstimatorConfig) DAOracle() evmconfig.DAOracle {
	return &TestDAOracleConfig{}
}

type TestDAOracleConfig struct {
	evmconfig.DAOracle
}

func (d *TestDAOracleConfig) OracleType() *toml.DAOracleType {
	oracleType := toml.DAOracleOPStack
	return &oracleType
}

func (d *TestDAOracleConfig) OracleAddress() *evmtypes.EIP55Address {
	a, err := evmtypes.NewEIP55Address("0x420000000000000000000000000000000000000F")
	if err != nil {
		panic(err)
	}
	return &a
}

func (d *TestDAOracleConfig) CustomGasPriceCalldata() *string {
	return nil
}

func (g *TestGasEstimatorConfig) BlockHistory() evmconfig.BlockHistory {
	return &TestBlockHistoryConfig{}
}

func (g *TestGasEstimatorConfig) FeeHistory() evmconfig.FeeHistory {
	return &TestFeeHistoryConfig{}
}

func (g *TestGasEstimatorConfig) EIP1559DynamicFees() bool   { return false }
func (g *TestGasEstimatorConfig) LimitDefault() uint64       { return 1e6 }
func (g *TestGasEstimatorConfig) BumpPercent() uint16        { return 2 }
func (g *TestGasEstimatorConfig) BumpThreshold() uint64      { return g.bumpThreshold }
func (g *TestGasEstimatorConfig) BumpMin() *assets.Wei       { return assets.GWei(1) }
func (g *TestGasEstimatorConfig) FeeCapDefault() *assets.Wei { return assets.GWei(1) }
func (g *TestGasEstimatorConfig) PriceDefault() *assets.Wei  { return assets.GWei(1) }
func (g *TestGasEstimatorConfig) TipCapDefault() *assets.Wei { return assets.GWei(1) }
func (g *TestGasEstimatorConfig) TipCapMin() *assets.Wei     { return assets.GWei(1) }
func (g *TestGasEstimatorConfig) LimitMax() uint64           { return 0 }
func (g *TestGasEstimatorConfig) LimitMultiplier() float32   { return 1 }
func (g *TestGasEstimatorConfig) BumpTxDepth() uint32        { return 42 }
func (g *TestGasEstimatorConfig) LimitTransfer() uint64      { return 42 }
func (g *TestGasEstimatorConfig) PriceMax() *assets.Wei      { return assets.GWei(1) }
func (g *TestGasEstimatorConfig) PriceMin() *assets.Wei      { return assets.GWei(1) }
func (g *TestGasEstimatorConfig) Mode() string               { return "FixedPrice" }
func (g *TestGasEstimatorConfig) LimitJobType() evmconfig.LimitJobType {
	return &TestLimitJobTypeConfig{}
}
func (g *TestGasEstimatorConfig) PriceMaxKey(addr common.Address) *assets.Wei {
	return assets.GWei(1)
}
func (g *TestGasEstimatorConfig) EstimateLimit() bool                   { return false }
func (g *TestGasEstimatorConfig) SenderAddress() *evmtypes.EIP55Address { return nil }

func (e *TestEvmConfig) GasEstimator() evmconfig.GasEstimator {
	return &TestGasEstimatorConfig{bumpThreshold: e.BumpThreshold}
}

type TestLimitJobTypeConfig struct {
}

func (l *TestLimitJobTypeConfig) OCR() *uint32    { return ptr(uint32(0)) }
func (l *TestLimitJobTypeConfig) OCR2() *uint32   { return ptr(uint32(0)) }
func (l *TestLimitJobTypeConfig) DR() *uint32     { return ptr(uint32(0)) }
func (l *TestLimitJobTypeConfig) FM() *uint32     { return ptr(uint32(0)) }
func (l *TestLimitJobTypeConfig) Keeper() *uint32 { return ptr(uint32(0)) }
func (l *TestLimitJobTypeConfig) VRF() *uint32    { return ptr(uint32(0)) }

type TestBlockHistoryConfig struct {
	evmconfig.BlockHistory
}

func (b *TestBlockHistoryConfig) BatchSize() uint32                 { return 42 }
func (b *TestBlockHistoryConfig) BlockDelay() uint16                { return 42 }
func (b *TestBlockHistoryConfig) BlockHistorySize() uint16          { return 42 }
func (b *TestBlockHistoryConfig) EIP1559FeeCapBufferBlocks() uint16 { return 42 }
func (b *TestBlockHistoryConfig) TransactionPercentile() uint16     { return 42 }

type TestFeeHistoryConfig struct {
	evmconfig.FeeHistory
}

type transactionsConfig struct {
	evmconfig.Transactions
	e         *TestEvmConfig
	autoPurge evmconfig.AutoPurgeConfig
}

func (*transactionsConfig) ForwardersEnabled() bool                { return false }
func (t *transactionsConfig) MaxInFlight() uint32                  { return t.e.MaxInFlight }
func (t *transactionsConfig) MaxQueued() uint64                    { return t.e.MaxQueued }
func (t *transactionsConfig) ReaperInterval() time.Duration        { return t.e.ReaperInterval }
func (t *transactionsConfig) ReaperThreshold() time.Duration       { return t.e.ReaperThreshold }
func (t *transactionsConfig) ResendAfterThreshold() time.Duration  { return t.e.ResendAfterThreshold }
func (t *transactionsConfig) AutoPurge() evmconfig.AutoPurgeConfig { return t.autoPurge }

type autoPurgeConfig struct {
	evmconfig.AutoPurgeConfig
}

func (a *autoPurgeConfig) Enabled() bool { return false }

type MockConfig struct {
	EvmConfig           *TestEvmConfig
	RpcDefaultBatchSize uint32
	finalityDepth       uint32
	finalityTagEnabled  bool
}

func (c *MockConfig) EVM() evmconfig.EVM {
	return c.EvmConfig
}

func (c *MockConfig) NonceAutoSync() bool            { return true }
func (c *MockConfig) ChainType() chaintype.ChainType { return "" }
func (c *MockConfig) FinalityDepth() uint32          { return c.finalityDepth }
func (c *MockConfig) SetFinalityDepth(fd uint32)     { c.finalityDepth = fd }
func (c *MockConfig) FinalityTagEnabled() bool       { return c.finalityTagEnabled }
func (c *MockConfig) RPCDefaultBatchSize() uint32    { return c.RpcDefaultBatchSize }

func MakeTestConfigs(t *testing.T) (*MockConfig, *TestDatabaseConfig, *TestEvmConfig) {
	db := &TestDatabaseConfig{defaultQueryTimeout: utils.DefaultQueryTimeout}
	ec := &TestEvmConfig{
		HeadTrackerConfig: &TestHeadTrackerConfig{},
		BumpThreshold:     42,
		MaxInFlight:       uint32(42),
		MaxQueued:         uint64(0),
		ReaperInterval:    time.Duration(0),
		ReaperThreshold:   time.Duration(0),
	}
	config := &MockConfig{EvmConfig: ec}
	return config, db, ec
}
