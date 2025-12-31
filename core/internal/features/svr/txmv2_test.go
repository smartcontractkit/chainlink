package svr

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	clcommonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
)

// Copied from core/services/ocr2/plugins/llo/integration_test.go + svr-contracts/test/e2e/svr_test.go + integration-tests/smoke/ocr2_test.go

/*
	Steps to run:
	* `docker run --name cl-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dbname -p 5432:5432 -d postgres`
	* `make setup-testdb` (password is 'postgres')
	* `CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -run ^TestIntegration_secondary_feed_transmission$ github.com/smartcontractkit/chainlink/v2/core/internal/features/svr -v`
*/

var (
	nNodes = 4
)

// TODO(gg): restart node or txmgr so that the second EOA can also be used to send txs from (see also https://github.com/smartcontractkit/chainlink/pull/16967)

// TODO(gg): redirect core node logs to separate file

// TODO(gg): update bootstrap config to have lower finalityDepth
//    logger.go:146: 2025-03-24T14:33:08.030Z	DEBUG	bootstrap_svr.EVM.1337.HeadSaver	heads/saver.go:72	chain shorter than FinalityDepth	{"version": "unset@unset", "chainLen": 8, "evmFinalityDepth": 10}

// TODO(gg) use framework.Context from "github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework" for Contexts

//TODO(gg): potential speedup: deploy + configure config contract first before starting the jobs (maybe we already do? We see a lot of "    logger.go:146: 2025-11-28T11:53:17.787Z     WARN    core_node_0.OCR2.offchainreporting2.2c534b8d-584c-4d54-85fd-bf7993542261        managed/track_config.go:110     TrackConfig: LatestConfigDetails() returned a zero configDigest. Looks like the contract has not been configured        {"version": "unset@unset", "jobID": 1, "jobName": "SVR job 0", "contractID":" though)

func TestIntegration_secondary_feed_transmission(t *testing.T) {
	const salt = 99

	// create node keys
	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := range nNodes {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	transactOpts, backend := setupBlockchain(t)
	fromBlock := 1

	mockFlashbotsServer := setupFlashbotsMock(t, backend)
	defer mockFlashbotsServer.Close()
	t.Logf("Mock Flashbots server started at %s", mockFlashbotsServer.URL())

	// Setup bootstrap node
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_node", backend, bootstrapCSAKey)
	bootstrapNode := node{app: appBootstrap, keyBundle: bootstrapKb}
	t.Logf("created bootstrap node with id %q and public key %#v", bootstrapPeerID, bootstrapNode.keyBundle.OnChainPublicKey())

	// Setup oracle nodes
	oracles, nodes := setupNodes(t, nNodes, transactOpts, backend, clientCSAKeys)
	t.Logf("created %d oracle nodes", len(nodes))

	// for i, node := range nodes {
	// 	// set up the secondary transmitter key
	// 	transmitterKey2, err := node.app.GetKeyStore().Eth().Create(testutils.Context(t), testutils.SimulatedChainID)
	// 	require.NoErrorf(t, err, "could not create transmitter key for node %d", i)
	// 	err = fundAddress(transmitterKey2.Address, transactOpts, backend)
	// 	require.NoError(t, err, "Funding transmitter shouldn't fail for node %d", i)
	// 	backend.Commit()
	// 	t.Logf("Funded secondary transmitter for node %d", i)
	// }

	var allPrimaryTransmitterAddresses []common.Address
	var allSecondaryTransmitterAddresses []common.Address
	for _, node := range nodes {
		allPrimaryTransmitterAddresses = append(allPrimaryTransmitterAddresses, node.transmitter)
		allSecondaryTransmitterAddresses = append(allSecondaryTransmitterAddresses, node.secondaryTransmitter)
	}
	t.Logf("allPrimaryTransmitterAddresses: %v", allPrimaryTransmitterAddresses)
	t.Logf("allSecondaryTransmitterAddresses: %v", allSecondaryTransmitterAddresses)

	var allForwarderAddresses []common.Address
	for _, node := range nodes {
		allForwarderAddresses = append(allForwarderAddresses, node.effectiveTransmitter)
	}
	t.Logf("allForwarderAddresses: %v", allForwarderAddresses)

	// 8. Deploy dual aggregator contract
	abi, err := DualAggregatorMetaData.GetAbi()
	require.NoError(t, err, "Failed to get dual aggregator ABI")

	dualAggAddress, _, _, err := bind.DeployContract(transactOpts, *abi, common.FromHex(DualAggregatorMetaData.Bin), backend.Client(),
		common.HexToAddress(transactOpts.From.Hex()), // TODO(gg): actually linkAddress
		big.NewInt(1),                 // MinimumAnswer
		big.NewInt(50000000000000000), // MaximumAnswer
		common.Address{},              // BillingAccessController
		common.Address{},              // RequesterAccessController
		uint8(8),                      // Decimals
		"SVR test",
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // secondary proxy
		uint32(30), // cutOffTime
		uint32(20), // maxSyncIterations
	)
	require.NoError(t, err, "Failed to deploy dual aggregator contract")
	backend.Commit()
	dualAggregatorInstance, err := NewDualAggregator(dualAggAddress, backend.Client())
	require.NoError(t, err, "Failed to create new dual aggregator instance")
	t.Logf("Deployed dual aggregator contract at %s", dualAggAddress.String())

	// 9. Configure the dual aggregator contracts
	s := []int{1, 2, 2, 2}

	signerKeys, transmitters, f, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		20*time.Second,       // deltaProgress time.Duration,
		20*time.Second,       // deltaResend time.Duration,
		1*time.Second,        // deltaRound time.Duration,
		250*time.Millisecond, // deltaGrace time.Duration,
		20*time.Second,       // deltaStage time.Duration,
		3,                    // rMax uint8,
		s,                    // s []int,
		oracles,
		median.OffchainConfig{
			AlphaReportInfinite: false,
			AlphaReportPPB:      1,
			AlphaAcceptInfinite: false,
			AlphaAcceptPPB:      1,
			DeltaC:              500 * time.Millisecond,
		}.Encode(), // reportingPluginConfig []byte,
		nil,
		5*time.Second, // maxDurationQuery time.Duration,
		5*time.Second, // maxDurationObservation time.Duration,
		5*time.Second, // maxDurationReport time.Duration,
		5*time.Second, // maxDurationShouldAcceptFinalizedReport time.Duration,
		5*time.Second, // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,             // f int,
		nil,           // The median reporting plugin has an empty onchain config
	)
	require.NoError(t, err, "Failed to create contract configuration")

	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(big.NewInt(1), big.NewInt(50000000000000000)) // MinimumAnswer MaximumAnswer
	require.NoError(t, err, "Failed to generate default ocr2 on-chain configuration")

	// Convert signers to addresses
	var signerAddresses []common.Address
	for _, signer := range signerKeys {
		signerAddresses = append(signerAddresses, common.BytesToAddress(signer))
	}

	// Convert transmitters to addresses
	var transmitterAddresses []common.Address
	for _, transmitter := range transmitters {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(transmitter)))
	}
	t.Logf("TransmitterAddresses: %v", transmitterAddresses)
	t.Logf("allPrimaryTransmitterAddresses: %v", allPrimaryTransmitterAddresses)

	_, err = dualAggregatorInstance.SetConfig(transactOpts, signerAddresses, allForwarderAddresses, f, onchainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		// decode the revert reason
		cerr, ok := err.(rpc.DataError)
		if !ok {
			t.Fatalf("Failed to configure dual aggregator contract: %v", err)
		}
		if cerr.ErrorData() != nil {
			t.Logf("Decoding custom ABI error from tx error")
			for k, abiError := range abi.Errors {
				data, err := hex.DecodeString(cerr.ErrorData().(string)[2:])
				if err != nil {
					t.Fatalf("Failed to decode error data: %v", err)
				}
				if len(data) < 4 {
					t.Fatalf("Error data too short: %v", data)
				}
				if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
					// Found a matching error
					v, err := abiError.Unpack(data)
					if err != nil {
						t.Fatalf("Failed to unpack error data: %v", err)
					}
					t.Fatalf("Failed to configure dual aggregator contract due to revert of type %s: %v", k, v)
				}
			}
		}
	}
	require.NoError(t, err, "Failed to configure dual aggregator contract")
	backend.Commit()
	backend.Commit()
	backend.Commit()
	t.Logf("Configured dual aggregator contract")

	t.Logf("Creating bootstrap job")
	bootstrapJob := job.Job{
		Type:          job.Bootstrap,
		SchemaVersion: 1,
		Name:          null.StringFrom("SVR bootstrap"),
		ExternalJobID: uuid.New(),
		BootstrapSpec: &job.BootstrapSpec{
			ContractID: dualAggAddress.Hex(),
			Relay:      "evm",
			RelayConfig: map[string]any{
				"chainID": testutils.SimulatedChainID.String(),
			},
		},
	}
	err = bootstrapNode.app.AddJobV2(context.Background(), &bootstrapJob)
	require.NoError(t, err, "Failed to create bootstrap job")
	t.Logf("Created bootstrap job")

	t.Logf("Creating job for feed %s", dualAggAddress.String())

	pl, err := pipeline.Parse(observationSource)
	require.NoErrorf(t, err, "Failed to parse observation source")

	for i, node := range nodes {
		keys, err := node.app.GetKeyStore().Eth().EnabledAddressesForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoErrorf(t, err, "could not get eth keys for node %d", i)
		t.Logf("keys from node %d: %v", i, keys)
		require.Len(t, keys, 2)

		// create the job
		// TODO(gg): maybe use oevJobSpec() instead if possible?
		jb := &job.Job{
			Type:              job.OffchainReporting2,
			SchemaVersion:     1,
			Name:              null.StringFrom(fmt.Sprintf("SVR job %d", i)),
			CronSpec:          &job.CronSpec{CronSchedule: "@every 1s"},
			PipelineSpec:      &pipeline.Spec{},
			Pipeline:          *pl,
			ExternalJobID:     uuid.New(),
			ForwardingAllowed: true,
			MaxTaskDuration:   *sqlutil.NewInterval(0 * time.Second),
			OCR2OracleSpec: &job.OCR2OracleSpec{
				ContractID:           dualAggAddress.Hex(),
				Relay:                "evm",
				OCRKeyBundleID:       null.StringFrom(node.keyBundle.ID()),
				PluginType:           clcommonTypes.Median,
				TransmitterID:        null.StringFrom(node.transmitter.Hex()),
				AllowNoBootstrappers: true,                                                                         // TODO(gg): maybe we can get away with this?
				P2PV2Bootstrappers:   []string{fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort)}, // TODO(gg) bootstrapPeerID.Data[0].Attributes.PeerID, needed?
				RelayConfig: map[string]any{
					"chainID":                testutils.SimulatedChainID.String(),
					"fromBlock":              fromBlock,
					"enableDualTransmission": true,
					"dualTransmission": map[string]any{
						"contractAddress":    dualAggAddress.Hex(),
						"transmitterAddress": node.secondaryTransmitter.Hex(),
						"meta": map[string]any{
							"hint":   []any{"full"},
							"refund": []any{"0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90"},
						},
						"endpoint": mockFlashbotsServer.URL(), // Configure secondary endpoint to use mock server
					},
				},
				PluginConfig: map[string]any{
					"juelsPerFeeCoinSource": "juels_per_fee_coin [type=\"sum\" values=<[0]>]",
				},
			},
		}
		err = node.app.AddJobV2(context.Background(), jb)
		require.NoError(t, err, "Failed to create feed job")
	}
	t.Logf("Created jobs for feed %s", dualAggAddress.String())

	// Get the deployment block number for FilterLogs queries
	deploymentBlock, err := backend.Client().BlockNumber(testutils.Context(t))
	require.NoError(t, err, "Failed to get deployment block number")
	t.Logf("Dual aggregator deployment block: %d", deploymentBlock)

	// Start block ticker to continue committing blocks
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	go func() {
		for range tick.C {
			backend.Commit()
		}
	}()

	// Use Eventually to wait for both primary and secondary transmissions in logs and on-chain
	primaryFound := false
	secondaryFound := false // tracks secondary transactions landing on-chain
	flashbotsFound := false // tracks secondary transactions sent to Flashbots server
	var primaryTxHash, secondaryTxHash string

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		// Commit a block to ensure any pending transactions are included
		backend.Commit()

		// Check logs from all oracle nodes
		// Note: The "Created primary/secondary transaction" logs come from the EVM relayer logger,
		// which may not be captured in node.observedLogs. We primarily rely on on-chain detection.
		for i, node := range nodes {
			if node.observedLogs == nil {
				continue
			}
			logs := node.observedLogs.All()
			for _, log := range logs {
				msg := strings.ToLower(log.Message)
				// Check for primary transmission (to chain)
				// Look for "Created primary transaction" - this is the actual log message from dual contract transmitter
				if !primaryFound && strings.Contains(msg, "created primary transaction") {
					primaryFound = true
					t.Logf("Node %d: Found primary transmission log: %s", i, log.Message)
				}
				// Note: We don't set secondaryFound from logs because secondary transactions go to Flashbots first,
				// then the mock server submits them to chain. We only set secondaryFound when we see it on-chain.
			}
		}

		// Check on-chain for both transactions using FilterLogs to find NewTransmission events
		// This is more efficient than checking every block
		ctx := testutils.Context(t)
		latestBlockNum, err := backend.Client().BlockNumber(ctx)
		if err == nil && latestBlockNum > 0 {
			transmittedEventSig := common.HexToHash("0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a") // NewTransmission event

			// Check blocks from deployment onwards for NewTransmission events
			startBlock := big.NewInt(0)
			endBlock := big.NewInt(int64(latestBlockNum))

			query := ethereum.FilterQuery{
				Addresses: []common.Address{dualAggAddress},
				Topics:    [][]common.Hash{{transmittedEventSig}},
				FromBlock: startBlock,
				ToBlock:   endBlock,
			}

			logs, err := backend.Client().FilterLogs(ctx, query)
			if err != nil {
				t.Logf("Failed to filter logs for NewTransmission event: %v", err)
				// FilterLogs might not be fully supported by simulated backend, fall back to checking blocks
				// This is expected and we'll rely on the block checking below
			} else {
				if len(logs) > 0 && (!primaryFound || !secondaryFound) {
					t.Logf("Found %d NewTransmission event(s) from dual aggregator (blocks %d-%d), primaryFound: %v, secondaryFound: %v", len(logs), startBlock.Int64(), endBlock.Int64(), primaryFound, secondaryFound)
				}
				// Process all events to find both primary and secondary transmissions
				// We need to check all events, not just stop at the first match
				for _, log := range logs {
					// Get the transaction that emitted this event
					tx, _, err := backend.Client().TransactionByHash(ctx, log.TxHash)
					if err != nil {
						continue
					}
					txHash := tx.Hash().Hex()

					// Check if it was submitted by Flashbots mock server (secondary)
					if mockFlashbotsServer.WasSubmitted(txHash) {
						if !secondaryFound {
							secondaryFound = true
							secondaryTxHash = txHash
							t.Logf("Found secondary transmission via NewTransmission event (submitted by Flashbots mock server): %s", secondaryTxHash)
						}
						continue // Skip further checks for Flashbots-submitted transactions
					}

					// Check if transaction is from a primary or secondary transmitter address
					// Get the sender address from the transaction using types.Sender
					isSecondaryTransmitter := false
					isPrimaryTransmitter := false
					chainID := tx.ChainId()
					senderStr := "unknown"
					if chainID != nil {
						signer := gethtypes.NewEIP155Signer(chainID)
						sender, err := gethtypes.Sender(signer, tx)
						if err == nil {
							senderStr = sender.Hex()
							// Check if sender matches any secondary transmitter address
							for _, secondaryAddr := range allSecondaryTransmitterAddresses {
								if sender == secondaryAddr {
									isSecondaryTransmitter = true
									break
								}
							}
							// Check if sender matches any primary transmitter address
							if !isSecondaryTransmitter {
								for _, primaryAddr := range allPrimaryTransmitterAddresses {
									if sender == primaryAddr {
										isPrimaryTransmitter = true
										break
									}
								}
							}
						}
					}

					// Also check if transaction is directly to dual aggregator with transmitSecondary method
					data := tx.Data()
					isSecondaryMethod := false
					if tx.To() != nil && *tx.To() == dualAggAddress && len(data) >= 4 {
						methodSig := hex.EncodeToString(data[:4])
						if methodSig == "ba0cb29e" {
							isSecondaryMethod = true
						}
					}

					if isSecondaryTransmitter || isSecondaryMethod {
						// Secondary transmission (from secondary transmitter or transmitSecondary method)
						if !secondaryFound {
							secondaryFound = true
							secondaryTxHash = txHash
							t.Logf("Found secondary transmission via NewTransmission event (from secondary transmitter: %v, method: %v, sender: %s): %s", isSecondaryTransmitter, isSecondaryMethod, senderStr, secondaryTxHash)
						}
					} else {
						// Not a secondary transmission, so it must be a primary transmission
						// (either from primary transmitter, or through forwarder which we can't distinguish)
						// Mark as primary if we haven't found one yet
						if !primaryFound {
							primaryFound = true
							primaryTxHash = txHash
							txToStr := "nil"
							if tx.To() != nil {
								txToStr = tx.To().Hex()
							}
							t.Logf("Found primary transmission via NewTransmission event (from primary transmitter: %v, sender: %s, txTo: %s): %s", isPrimaryTransmitter, senderStr, txToStr, primaryTxHash)
						}
					}

					// Early exit if we've found both
					if primaryFound && secondaryFound {
						break
					}
				}
			}

			// Fallback: Also check recent blocks directly for transactions to forwarders or dual aggregator
			// This works even if FilterLogs doesn't work with simulated backend
			if !primaryFound || !secondaryFound {
				fallbackStartBlock := latestBlockNum
				if fallbackStartBlock > 20 {
					fallbackStartBlock = fallbackStartBlock - 20
				} else {
					fallbackStartBlock = 0
				}

				for blockNum := fallbackStartBlock; blockNum <= latestBlockNum; blockNum++ {
					block, err := backend.Client().BlockByNumber(ctx, big.NewInt(int64(blockNum)))
					if err != nil || block == nil {
						continue
					}

					for _, tx := range block.Transactions() {
						txTo := tx.To()
						if txTo == nil {
							continue
						}

						// Check if transaction is to dual aggregator or a forwarder
						isRelevant := *txTo == dualAggAddress
						if !isRelevant {
							for _, forwarderAddr := range allForwarderAddresses {
								if *txTo == forwarderAddr {
									isRelevant = true
									break
								}
							}
						}

						if isRelevant {
							receipt, err := backend.Client().TransactionReceipt(ctx, tx.Hash())
							if err == nil && receipt != nil {
								// NewTransmission event signature: 0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a
								newTransmissionEventSig := common.HexToHash("0xc797025feeeaf2cd924c99e9205acb8ec04d5cad21c41ce637a38fb6dee6016a")
								for _, log := range receipt.Logs {
									if log.Address == dualAggAddress && len(log.Topics) > 0 && log.Topics[0] == newTransmissionEventSig {
										txHash := tx.Hash().Hex()
										if mockFlashbotsServer.WasSubmitted(txHash) {
											if !secondaryFound {
												secondaryFound = true
												secondaryTxHash = txHash
												t.Logf("Found secondary transmission via fallback block check (submitted by Flashbots): %s", secondaryTxHash)
											}
										} else {
											// Check if transaction is from a primary or secondary transmitter address
											// Get the sender address from the transaction using types.Sender
											isSecondaryTransmitter := false
											isPrimaryTransmitter := false
											chainID := tx.ChainId()
											if chainID != nil {
												signer := gethtypes.NewEIP155Signer(chainID)
												sender, err := gethtypes.Sender(signer, tx)
												if err == nil {
													// Check if sender matches any secondary transmitter address
													for _, secondaryAddr := range allSecondaryTransmitterAddresses {
														if sender == secondaryAddr {
															isSecondaryTransmitter = true
															break
														}
													}
													// Check if sender matches any primary transmitter address
													if !isSecondaryTransmitter {
														for _, primaryAddr := range allPrimaryTransmitterAddresses {
															if sender == primaryAddr {
																isPrimaryTransmitter = true
																break
															}
														}
													}
												}
											}

											// Also check if transaction is directly to dual aggregator with transmitSecondary method
											data := tx.Data()
											isSecondaryMethod := false
											if *txTo == dualAggAddress && len(data) >= 4 {
												methodSig := hex.EncodeToString(data[:4])
												if methodSig == "ba0cb29e" {
													isSecondaryMethod = true
												}
											}

											if isSecondaryTransmitter || isSecondaryMethod {
												// Secondary transmission (from secondary transmitter or transmitSecondary method)
												if !secondaryFound {
													secondaryFound = true
													secondaryTxHash = txHash
													t.Logf("Found secondary transmission via fallback block check (from secondary transmitter: %v, method: %v): %s", isSecondaryTransmitter, isSecondaryMethod, secondaryTxHash)
												}
											} else {
												// Not a secondary transmission, so it must be a primary transmission
												// (either from primary transmitter, or through forwarder which we can't distinguish)
												if !primaryFound {
													primaryFound = true
													primaryTxHash = txHash
													t.Logf("Found primary transmission via fallback block check (from primary transmitter: %v): %s", isPrimaryTransmitter, primaryTxHash)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		// Also check if mock server received secondary transactions
		mockServerCount := mockFlashbotsServer.TransactionCount()

		if mockServerCount > 0 && !flashbotsFound {
			// Mock server received transactions, mark as found
			flashbotsFound = true
			t.Logf("Mock Flashbots server received %d secondary transaction(s)", mockServerCount)
		}

		secondaryAnyFound := secondaryFound || flashbotsFound
		if primaryFound && secondaryAnyFound {
			t.Logf("Found both primary and secondary transmissions - Primary: %v (tx: %s), Secondary on-chain: %v (tx: %s), Flashbots: %v (mock server count: %d)",
				primaryFound, primaryTxHash, secondaryFound, secondaryTxHash, flashbotsFound, mockServerCount)
		} else {
			t.Logf("Waiting for transmissions - primary: %v, secondary on-chain: %v, flashbots: %v (mock server count: %d)",
				primaryFound, secondaryFound, flashbotsFound, mockServerCount)
		}

		return primaryFound && secondaryAnyFound
	}, 2*time.Minute, 1*time.Second).Should(gomega.BeTrue(),
		"Expected both primary (to chain) and secondary transmissions (at least one secondary path: on-chain or Flashbots). Primary found: %v, Secondary on-chain: %v, Flashbots: %v, Mock server count: %d",
		primaryFound, secondaryFound, flashbotsFound, mockFlashbotsServer.TransactionCount())

	t.Logf("primaryFound: %v, secondaryFound: %v, flashbotsFound: %v, mockServerCount: %d", primaryFound, secondaryFound, flashbotsFound, mockFlashbotsServer.TransactionCount())

	// Final assertions - separate checks for each secondary path
	// These allow the test to separately succeed/fail based on txs landing on chain vs. being sent to Flashbots
	require.True(t, primaryFound, "Primary transmission should be found in logs or on-chain")
	require.True(t, secondaryFound || flashbotsFound, "Secondary transmission should be found either on-chain or sent to Flashbots")

	// Log which secondary path(s) succeeded for visibility
	if secondaryFound {
		t.Logf("✓ Secondary transmission found on-chain")
	} else {
		t.Logf("⚠ Secondary transmission was NOT found on-chain")
	}
	if flashbotsFound {
		t.Logf("✓ Secondary transmission sent to Flashbots server")
	} else {
		t.Logf("⚠ Secondary transmission was NOT sent to Flashbots server")
	}

	if primaryTxHash != "" {
		t.Logf("✓ Primary transaction confirmed on-chain: %s", primaryTxHash)
	}
	if secondaryTxHash != "" {
		t.Logf("✓ Secondary transaction confirmed on-chain: %s", secondaryTxHash)
	} else {
		count := mockFlashbotsServer.TransactionCount()
		if count > 0 {
			t.Logf("✓ Secondary transaction received by mock server and submitted to chain (count: %d)", count)
		}
	}
}

func setupBlockchain(t *testing.T) (*bind.TransactOpts, *simulated.Backend) {
	contractOwner := evmtestutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := gethtypes.GenesisAlloc{contractOwner.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := simulated.NewBackend(genesisData, simulated.WithBlockGasLimit(ethconfig.Defaults.Miner.GasCeil))
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	return contractOwner, backend
}

func fundAddress(address common.Address, contractOwner *bind.TransactOpts, backend *simulated.Backend) error {
	wei := new(big.Int)
	amount := big.NewFloat(0.2)
	amountWei := new(big.Float).Mul(amount, big.NewFloat(1e18))
	amountWei.Int(wei)

	backend.Client().PendingNonceAt(context.Background(), address)
	nonce, err := backend.Client().PendingNonceAt(context.Background(), contractOwner.From)
	if err != nil {
		return fmt.Errorf("failed to fetch nonce: %w", err)
	}

	gasPrice, err := backend.Client().SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch gas price: %w", err)
	}
	gasLimit := uint64(21000) // Standard gas limit for ETH transfer

	tx := gethtypes.NewTransaction(nonce, address, wei, gasLimit, gasPrice, nil)

	signedTx, err := contractOwner.Signer(contractOwner.From, tx)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = backend.Client().SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	backend.Commit()

	_, err = bind.WaitMined(context.Background(), backend.Client(), signedTx)
	return err
}

var observationSource = `
//randomness
   val1 [type="memo" value="10"]
   val2 [type="memo" value="20"]
   val3 [type="memo" value="30"]
   val4 [type="memo" value="40"]
   val5 [type="memo" value="50"]
   val6 [type="memo" value="60"]
   val7 [type="memo" value="70"]
   val8 [type="memo" value="80"]
   val9 [type="memo" value="90"]

   random1 [type="any"]
   random2 [type="any"]
   random3 [type="any"]

   val1 -> random1
   val2 -> random2
   val3 -> random3
   val4 -> random1
   val5 -> random2
   val6 -> random3
   val7 -> random1
   val8 -> random2
   val9 -> random3


   // data source 1
   ds1_multiply [type="multiply" times=100]

	// data source 2
   ds2_multiply [type="multiply" times=100]


   // data source 3
   ds3_multiply [type="multiply" times=100]


   random1 -> ds1_multiply -> answer
   random2 -> ds2_multiply -> answer
   random3 -> ds3_multiply -> answer

   answer [type=median]
`

var oevJobSpec = `
type = "offchainreporting2"
schemaVersion = 1
name = "OEV job %d"
externalJobID = "%s"
forwardingAllowed = true
maxTaskDuration = "0s"
contractID = "%s"
relay = "%s"
ocrKeyBundleID = "%s"
pluginType = "median"
transmitterID = "%s"
p2pv2Bootstrappers = ["%s@%s"]

observationSource = """
 //randomness
    val1 [type="memo" value="10"]
    val2 [type="memo" value="20"]
    val3 [type="memo" value="30"]
    val4 [type="memo" value="40"]
    val5 [type="memo" value="50"]
    val6 [type="memo" value="60"]
    val7 [type="memo" value="70"]
    val8 [type="memo" value="80"]
    val9 [type="memo" value="90"]

    random1 [type="any"]
    random2 [type="any"]
    random3 [type="any"]

    val1 -> random1
    val2 -> random2
    val3 -> random3
    val4 -> random1
    val5 -> random2
    val6 -> random3
    val7 -> random1
    val8 -> random2
    val9 -> random3


    // data source 1
    ds1_multiply [type="multiply" times=100]

     // data source 2
    ds2_multiply [type="multiply" times=100]


    // data source 3
    ds3_multiply [type="multiply" times=100]


    random1 -> ds1_multiply -> answer
    random2 -> ds2_multiply -> answer
    random3 -> ds3_multiply -> answer

    answer [type=median]
"""

[relayConfig]
chainID = %s
fromBlock = %d
enableDualTransmission = true

[relayConfig.dualTransmission]
contractAddress = "%s"
transmitterAddress = "%s"

[relayConfig.dualTransmission.meta]
hint = [ "calldata" ]
refund = [ "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90" ]

[pluginConfig]
juelsPerFeeCoinSource = """
juels_per_fee_coin [type="sum" values=<[0]>];
"""
`
