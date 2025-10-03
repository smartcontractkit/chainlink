package testhelpers

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"

	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/offramp"
	tonlploader "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/backend/loader/account"
	tonlptypes "github.com/smartcontractkit/chainlink-ton/pkg/logpoller/types"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/hash"
)

const tonTickerInterval = 2500 * time.Millisecond
const safeLookbackBlocks = 50

// CCIPEventTopics maps TON event topics (CRC32 hashes) to their event names.
// Only includes events that are actually emitted as external out messages.
var CCIPEventTopics = map[uint32]string{
	// onramp: emitted when a CCIP message is sent from TON to another chain
	hash.CRC32("CCIPMessageSent"): "CCIPMessageSent",
	// offramp: emitted when a CCIP message execution state changes (in progress, success, failure).
	hash.CRC32("ExecutionStateChanged"): "ExecutionStateChanged",
	// offramp: emitted when a commit report is accepted, containing merkle roots and/or price updates.
	hash.CRC32("CommitReportAccepted"): "CommitReportAccepted",
	// ocr3base: emitted when a config is set.
	hash.CRC32("OCR3Base_ConfigSet"): "OCR3Base_ConfigSet",
	// ocr3base: emitted when a report is transmitted.
	hash.CRC32("OCR3Base_Transmitted"): "OCR3Base_Transmitted",
}

// getEventName returns the event name for a given topic, or a formatted unknown topic string.
func getEventName(topic uint32) string {
	if name, ok := CCIPEventTopics[topic]; ok {
		return name
	}
	return fmt.Sprintf("Unknown_0x%08x", topic)
}

// TODO(@jadepark-dev): clean up after verifying EVM2TON
func ConfirmCommitWithExpectedSeqNumRangeTON(t *testing.T, srcSelector uint64, tonChain cldf_ton.Chain, offRampContract address.Address, expectedSeqNumRange ccipocr3common.SeqNumRange) (bool, error) {
	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)

	// Create prefixed logger for TON event assertion
	lggr := logger.Named(logger.Test(t), "TON_EVENT_ASSERTION:COMMIT")
	lggr.Infof("waiting for commit report from srcSelector=%d, expectedSeqNumRange=[%d, %d], timeout=%v",
		srcSelector, expectedSeqNumRange.Start(), expectedSeqNumRange.End(), tests.WaitTimeout(t))

	tonBlockTicker := time.NewTicker(tonTickerInterval)
	// Start from a bit earlier to catch any events we might have missed
	// Use a lookback of about 50 blocks to ensure we don't miss commit events
	var startBlock uint32 = 0
	if currentBlock, err := tonChain.Client.CurrentMasterchainInfo(t.Context()); err == nil && currentBlock.SeqNo > safeLookbackBlocks {
		lggr.Infof("scan from block %d (50 blocks back from current %d)", startBlock, currentBlock.SeqNo)
	}

	eventCh, errCh := GetEvents[offramp.CommitReportAccepted](t, lggr, t.Context(), tonChain, &offRampContract, startBlock, tonBlockTicker)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	for {
		select {
		case commitEvent := <-eventCh:
			// Log event details with proper dereferencing
			if commitEvent.MerkleRoot != nil {
				lggr.Infof("received CommitReportAccepted event: MerkleRoot={SourceChainSelector:%d, MinSeqNr:%d, MaxSeqNr:%d, MerkleRoot:%x}, PriceUpdates=%+v",
					commitEvent.MerkleRoot.SourceChainSelector, commitEvent.MerkleRoot.MinSeqNr, commitEvent.MerkleRoot.MaxSeqNr,
					commitEvent.MerkleRoot.MerkleRoot, commitEvent.PriceUpdates)
			} else {
				lggr.Infof("received CommitReportAccepted event: MerkleRoot=<nil>, PriceUpdates=%+v", commitEvent.PriceUpdates)
			}

			// if merkle root is zero, it only contains price updates
			if commitEvent.MerkleRoot == nil {
				lggr.Infof("Skipping CommitReportAccepted with only price updates")
				continue
			}

			mr := commitEvent.MerkleRoot
			require.Equal(t, srcSelector, mr.SourceChainSelector)

			// TODO: this logic is duplicated with verifyCommitReport, share
			seenMessages.visitCommitReport(commitEvent.MerkleRoot.SourceChainSelector, mr.MinSeqNr, mr.MaxSeqNr)
			if mr.SourceChainSelector == srcSelector &&
				uint64(expectedSeqNumRange.Start()) >= mr.MinSeqNr &&
				uint64(expectedSeqNumRange.End()) <= mr.MaxSeqNr {
				t.Logf("All sequence numbers committed in a single report [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
				return true, nil
			}

			if seenMessages.allCommited(srcSelector) {
				t.Logf("All sequence numbers already committed from range [%d, %d]", expectedSeqNumRange.Start(), expectedSeqNumRange.End())
				return true, nil
			}
		case err := <-errCh:
			require.NoError(t, err)
		case <-timeout.C:
			return false, fmt.Errorf("timed out after waiting for commit report on chain selector %d from source selector %d expected seq nr range %s",
				tonChain.Selector, srcSelector, expectedSeqNumRange.String())
		}
	}
}

// ConfirmExecWithSeqNrsTON waits for execution state changes on TON for the given sequence numbers
// Returns a map of sequence number to execution state
func ConfirmExecWithSeqNrsTON(
	t *testing.T,
	srcSelector uint64,
	tonChain cldf_ton.Chain,
	offRampContract address.Address,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (map[uint64]int, error) {
	if len(expectedSeqNrs) == 0 {
		return nil, fmt.Errorf("no expected sequence numbers provided")
	}

	// Create prefixed logger for TON event assertion
	lggr := logger.Named(logger.Test(t), "TON_EVENT_ASSERTION:EXEC")
	lggr.Infof("waiting for execution state changes from srcSelector=%d, expectedSeqNrs=%v, timeout=%v",
		srcSelector, expectedSeqNrs, tests.WaitTimeout(t))

	executionStates := make(map[uint64]int)

	tonBlockTicker := time.NewTicker(tonTickerInterval)

	// Determine start block
	var scanStartBlock uint32 = 0
	if startBlock != nil {
		scanStartBlock = uint32(*startBlock)
	} else if currentBlock, err := tonChain.Client.CurrentMasterchainInfo(t.Context()); err == nil && currentBlock.SeqNo > safeLookbackBlocks {
		scanStartBlock = currentBlock.SeqNo - safeLookbackBlocks
		lggr.Infof("scan from block %d (50 blocks back from current %d)", scanStartBlock, currentBlock.SeqNo)
	}

	eventCh, errCh := GetEvents[offramp.ExecutionStateChanged](t, lggr, t.Context(), tonChain, &offRampContract, scanStartBlock, tonBlockTicker)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	// Track which sequence numbers we're waiting for
	expectedSeqNums := make(map[uint64]bool)
	for _, seqNum := range expectedSeqNrs {
		expectedSeqNums[seqNum] = true
	}

	for {
		select {
		case execEvent := <-eventCh:
			lggr.Infof("received ExecutionStateChanged event: SourceChainSelector=%d, SequenceNumber=%d, MessageID=%x, State=%d",
				execEvent.SourceChainSelector, execEvent.SequenceNumber, execEvent.MessageID, execEvent.State)

			// Check if this is for our source chain
			if execEvent.SourceChainSelector != srcSelector {
				lggr.Debugf("skipping event from different source chain: %d (expected %d)", execEvent.SourceChainSelector, srcSelector)
				continue
			}

			// Check if this is a sequence number we're waiting for
			if _, expected := expectedSeqNums[execEvent.SequenceNumber]; expected {
				// State 1 = IN_PROGRESS, State 2 = SUCCESS, State 3 = FAILURE
				// Only record and remove from expected if we got SUCCESS state
				if execEvent.State == 1 {
					lggr.Infof("received IN_PROGRESS state for seq num %d, waiting for SUCCESS state", execEvent.SequenceNumber)
					continue
				}

				if execEvent.State == 3 {
					lggr.Errorf("Execution failure detected for seq num %d, source chain %d, message id %x, state %d", execEvent.SequenceNumber, execEvent.SourceChainSelector, execEvent.MessageID, execEvent.State)
					// Don't delete from expected, don't record, just continue
					continue
				}

				executionStates[execEvent.SequenceNumber] = int(execEvent.State)
				delete(expectedSeqNums, execEvent.SequenceNumber)

				lggr.Infof("found execution state for seq num %d: state=%d, remaining=%d",
					execEvent.SequenceNumber, execEvent.State, len(expectedSeqNums))
			}

			// If we've found all expected sequence numbers, return
			if len(expectedSeqNums) == 0 {
				t.Logf("All sequence numbers executed: %v", expectedSeqNrs)
				return executionStates, nil
			}

		case err := <-errCh:
			return nil, fmt.Errorf("error while waiting for execution events: %w", err)

		case <-timeout.C:
			missing := make([]uint64, 0, len(expectedSeqNums))
			for seqNum := range expectedSeqNums {
				missing = append(missing, seqNum)
			}
			return executionStates, fmt.Errorf("timed out after waiting for execution on chain selector %d from source selector %d, missing seq nums: %v",
				tonChain.Selector, srcSelector, missing)
		}
	}
}

func GetEvents[T any](t *testing.T, lggr logger.Logger, ctx context.Context, tonChain cldf_ton.Chain, contractAddress *address.Address, startBlock uint32, ticker *time.Ticker) (<-chan T, <-chan error) {
	ch := make(chan T)
	errorCh := make(chan error)

	go func() {
		defer ticker.Stop()
		defer close(ch)
		defer close(errorCh)

		// lazy client provider
		clientProvider := func(ctx context.Context) (ton.APIClientWrapped, error) {
			return tonChain.Client.WithRetry(3), nil
		}
		// init loader
		loader := tonlploader.NewTxLoader(lggr, clientProvider, 100)

		lastProcessedBlock := startBlock

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 1. Get block range
				blockRange, newSeqNo, err := getBlockRange(ctx, tonChain, lastProcessedBlock)
				if err != nil {
					errorCh <- err
					return
				}

				// Skip if no new blocks
				if blockRange == nil {
					continue
				}

				// 2. Fetch transactions
				txs, err := loader.FetchTxsForAddress(ctx, blockRange, contractAddress)
				if err != nil {
					errorCh <- fmt.Errorf("failed to load transactions: %w", err)
					return
				}

				// 3. Get messages in transactions to see if there is event we're looking for
				events, err := extractEventMessage[T](txs)
				if err != nil {
					errorCh <- err
					return
				}

				// Send events to channel
				for _, event := range events {
					lggr.Infof("TON:FOUND EVENT: %+v", event)
					select {
					case ch <- event:
					case <-ctx.Done():
						return
					}
				}

				// Update last processed block
				lastProcessedBlock = newSeqNo
			}
		}
	}()

	return ch, errorCh
}

// getBlockRange creates a block range from lastProcessedBlock to current masterchain head
func getBlockRange(ctx context.Context, tonChain cldf_ton.Chain, lastProcessedBlock uint32) (*tonlptypes.BlockRange, uint32, error) {
	// lazy client provider
	clientProvider := func(_ context.Context) (ton.APIClientWrapped, error) {
		return tonChain.Client.WithRetry(3), nil
	}

	// 1. Get latest block number
	client, err := clientProvider(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get client: %w", err)
	}

	toBlock, err := client.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get current masterchain info: %w", err)
	}

	// Skip if no new blocks
	if toBlock.SeqNo <= lastProcessedBlock {
		return nil, toBlock.SeqNo, nil
	}

	// 2. Create a block range
	var prevBlock *ton.BlockIDExt
	if lastProcessedBlock > 0 {
		// Get the previous block reference
		prevBlock, err = client.LookupBlock(ctx, toBlock.Workchain, toBlock.Shard, lastProcessedBlock)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to lookup previous block %d: %w", lastProcessedBlock, err)
		}
	}

	blockRange := &tonlptypes.BlockRange{
		Prev: prevBlock,
		To:   toBlock,
	}

	return blockRange, toBlock.SeqNo, nil
}

// decodeEventTopic extracts the event topic from an external out message destination address
// Topic is encoded in the last 4 bytes of the address data
func decodeEventTopic(addr *address.Address) (uint32, error) {
	if addr == nil {
		return 0, fmt.Errorf("cannot decode from a nil address")
	}

	data := addr.Data()
	const eventTopicLength = 4

	if len(data) < eventTopicLength {
		return 0, fmt.Errorf("address data is too short to contain an event topic")
	}

	startIndex := len(data) - eventTopicLength
	return binary.BigEndian.Uint32(data[startIndex:]), nil
}

// extractEventMessage processes transactions to extract events of type T from external messages
func extractEventMessage[T any](txs []tonlptypes.TxWithBlock) ([]T, error) {
	var events []T

	for _, tx := range txs {
		if tx.Tx == nil {
			continue
		}
		if tx.Tx.IO.Out != nil {
			msgs, err := tx.Tx.IO.Out.ToSlice()
			if err != nil {
				// skip this tx
				continue
			}

			for _, msg := range msgs {
				if msg.MsgType == tlb.MsgTypeExternalOut {
					if extOut := msg.AsExternalOut(); extOut != nil {
						// Decode event topic from destination address
						topic, err := decodeEventTopic(extOut.DestAddr())
						if err != nil {
							fmt.Printf("\n[DEBUG] Failed to decode event topic: %v\n", err)
							continue
						}

						// First, try to identify what kind of event this is
						bodyCell := extOut.Payload()
						if bodyCell == nil {
							fmt.Printf("\n[DEBUG] No payload in external out message\n")
							continue
						}
						bodySlice := bodyCell.BeginParse()
						totalBits := bodySlice.BitsLeft()

						// Log ALL external out messages for debugging
						fmt.Printf("\n[DEBUG] Event: %s (topic=0x%08x), bits=%d, refs=%d\n",
							getEventName(topic), topic, totalBits, bodyCell.RefsNum())

						// Check if we're looking for ExecutionStateChanged events
						var eventType T
						_, isExecStateChanged := any(eventType).(offramp.ExecutionStateChanged)

						// ExecutionStateChanged topic (crc32("ExecutionStateChanged") = 0x4C94C360)
						const executionStateChangedTopic = 0x4C94C360

						if topic != executionStateChangedTopic {
							fmt.Printf("[DEBUG] Skipping %s (not ExecutionStateChanged)\n", getEventName(topic))
							if !isExecStateChanged {
								// Still try to parse for other event types
								var event T
								err := tlb.LoadFromCell(&event, bodyCell.BeginParse())
								if err == nil {
									events = append(events, event)
								}
							}
							continue
						}

						if isExecStateChanged {
							// Log cell structure (reuse variables from above)
							refsNum := bodyCell.RefsNum()

							fmt.Printf("\n=== DEBUG: ExecutionStateChanged Event Analysis ===\n")
							fmt.Printf("Cell bits available: %d\n", totalBits)
							fmt.Printf("Cell refs count: %d\n", refsNum)

							// Get the full BOC for complete analysis
							bocBytes := bodyCell.ToBOC()
							fmt.Printf("Full BOC (hex): %x\n", bocBytes)

							// Log raw bits in the main cell
							tempSlice := bodyCell.BeginParse()
							if tempSlice.BitsLeft() > 0 {
								bitsToLoad := tempSlice.BitsLeft()
								rawBits, err := tempSlice.LoadSlice(bitsToLoad)
								if err == nil {
									fmt.Printf("Main cell data (hex): %x\n", rawBits)
									fmt.Printf("Main cell data (bytes): %d\n", len(rawBits))

									// Manual parsing based on Go binding struct:
									// type ExecutionStateChanged struct {
									//     SourceChainSelector uint64 `tlb:"## 64"`
									//     SequenceNumber      uint64 `tlb:"## 64"`
									//     MessageID           []byte `tlb:"bits 256"`
									//     State               uint8  `tlb:"## 8"`
									// }
									// Total expected: 64 + 64 + 256 + 8 = 392 bits = 49 bytes

									fmt.Printf("\n--- Manual Field Extraction (Standard Order) ---\n")
									if len(rawBits) >= 8 {
										srcChain := uint64(rawBits[0])<<56 | uint64(rawBits[1])<<48 |
											uint64(rawBits[2])<<40 | uint64(rawBits[3])<<32 |
											uint64(rawBits[4])<<24 | uint64(rawBits[5])<<16 |
											uint64(rawBits[6])<<8 | uint64(rawBits[7])
										fmt.Printf("SourceChainSelector (bytes 0-7): %d (0x%016x)\n", srcChain, srcChain)
									}

									if len(rawBits) >= 16 {
										seqNum := uint64(rawBits[8])<<56 | uint64(rawBits[9])<<48 |
											uint64(rawBits[10])<<40 | uint64(rawBits[11])<<32 |
											uint64(rawBits[12])<<24 | uint64(rawBits[13])<<16 |
											uint64(rawBits[14])<<8 | uint64(rawBits[15])
										fmt.Printf("SequenceNumber (bytes 8-15): %d (0x%016x)\n", seqNum, seqNum)
									}

									if len(rawBits) >= 48 {
										msgID := rawBits[16:48] // 32 bytes = 256 bits
										fmt.Printf("MessageID (bytes 16-47): %x\n", msgID)
									} else if len(rawBits) > 16 {
										msgID := rawBits[16:]
										fmt.Printf("MessageID (partial, bytes 16-%d): %x\n", len(rawBits)-1, msgID)
										fmt.Printf("MessageID missing: %d bytes (%d bits)\n", 32-len(msgID), (32-len(msgID))*8)
									}

									if len(rawBits) >= 49 {
										state := rawBits[48]
										fmt.Printf("State (byte 48): %d\n", state)
									} else {
										fmt.Printf("State: MISSING (need byte 48, have only %d bytes)\n", len(rawBits))
									}

									// Try different offsets to find expected MessageID
									expectedMsgID := "bfe43d967f073320d3a6e38da2b949043085c4686de7fcc2d0b968659391da49"
									expectedSrcChain := uint64(909606746561742123) // 0x0ca09010777524fb
									expectedSeqNum := uint64(1)

									fmt.Printf("\n--- Searching for Expected Values ---\n")
									fmt.Printf("Expected MessageID: %s\n", expectedMsgID)
									fmt.Printf("Expected SourceChain: %d (0x%016x)\n", expectedSrcChain, expectedSrcChain)
									fmt.Printf("Expected SeqNum: %d\n", expectedSeqNum)

									// Try different offsets (0 to 10 bytes)
									fmt.Printf("\n--- Trying Different Offsets ---\n")
									for offset := 0; offset <= 10 && offset < len(rawBits)-32; offset++ {
										// Check if MessageID matches at this offset
										msgIDCandidate := rawBits[offset : offset+32]
										msgIDHex := fmt.Sprintf("%x", msgIDCandidate)

										if msgIDHex == expectedMsgID {
											fmt.Printf("*** FOUND MessageID at offset %d! ***\n", offset)

											// Try to extract other fields based on this offset
											if offset >= 16 {
												// MessageID should be at bytes 16-47, so calculate src/seq positions
												srcOffset := offset - 16
												seqOffset := offset - 8

												if srcOffset >= 0 && srcOffset+8 <= len(rawBits) {
													src := uint64(rawBits[srcOffset])<<56 | uint64(rawBits[srcOffset+1])<<48 |
														uint64(rawBits[srcOffset+2])<<40 | uint64(rawBits[srcOffset+3])<<32 |
														uint64(rawBits[srcOffset+4])<<24 | uint64(rawBits[srcOffset+5])<<16 |
														uint64(rawBits[srcOffset+6])<<8 | uint64(rawBits[srcOffset+7])
													fmt.Printf("  SourceChain at offset %d: %d (0x%016x)\n", srcOffset, src, src)
												}

												if seqOffset >= 0 && seqOffset+8 <= len(rawBits) {
													seq := uint64(rawBits[seqOffset])<<56 | uint64(rawBits[seqOffset+1])<<48 |
														uint64(rawBits[seqOffset+2])<<40 | uint64(rawBits[seqOffset+3])<<32 |
														uint64(rawBits[seqOffset+4])<<24 | uint64(rawBits[seqOffset+5])<<16 |
														uint64(rawBits[seqOffset+6])<<8 | uint64(rawBits[seqOffset+7])
													fmt.Printf("  SeqNum at offset %d: %d (0x%016x)\n", seqOffset, seq, seq)
												}
											}
										}
									}

									// Also search for the expected source chain selector
									fmt.Printf("\n--- Searching for SourceChain bytes ---\n")
									srcBytes := []byte{0x0c, 0xa0, 0x90, 0x10, 0x77, 0x75, 0x24, 0xfb}
									for i := 0; i <= len(rawBits)-8; i++ {
										match := true
										for j := 0; j < 8; j++ {
											if rawBits[i+j] != srcBytes[j] {
												match = false
												break
											}
										}
										if match {
											fmt.Printf("*** FOUND SourceChain at offset %d! ***\n", i)
										}
									}

									// Search for sequence number 1
									fmt.Printf("\n--- Searching for SeqNum=1 bytes ---\n")
									seqBytes := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
									for i := 0; i <= len(rawBits)-8; i++ {
										match := true
										for j := 0; j < 8; j++ {
											if rawBits[i+j] != seqBytes[j] {
												match = false
												break
											}
										}
										if match {
											fmt.Printf("*** FOUND SeqNum=1 at offset %d! ***\n", i)
										}
									}
								}
							}

							// Try TLB library decoding step-by-step
							fmt.Printf("\n--- TLB Library Step-by-Step Decode ---\n")
							manualSlice := bodyCell.BeginParse()

							if manualSlice.BitsLeft() >= 64 {
								sourceChain, err1 := manualSlice.LoadUInt(64)
								if err1 == nil {
									fmt.Printf("LoadUInt(64) SourceChain: %d, remaining bits: %d\n", sourceChain, manualSlice.BitsLeft())
								} else {
									fmt.Printf("LoadUInt(64) SourceChain FAILED: %v\n", err1)
								}
							}

							if manualSlice.BitsLeft() >= 64 {
								seqNum, err2 := manualSlice.LoadUInt(64)
								if err2 == nil {
									fmt.Printf("LoadUInt(64) SeqNum: %d, remaining bits: %d\n", seqNum, manualSlice.BitsLeft())
								} else {
									fmt.Printf("LoadUInt(64) SeqNum FAILED: %v\n", err2)
								}
							}

							if manualSlice.BitsLeft() >= 256 {
								msgID, err3 := manualSlice.LoadSlice(256)
								if err3 == nil {
									fmt.Printf("LoadSlice(256) MessageID: %x, remaining bits: %d\n", msgID, manualSlice.BitsLeft())
								} else {
									fmt.Printf("LoadSlice(256) MessageID FAILED: %v\n", err3)
								}
							} else {
								fmt.Printf("LoadSlice(256) MessageID: INSUFFICIENT BITS (have %d, need 256)\n", manualSlice.BitsLeft())
							}

							if manualSlice.BitsLeft() >= 8 {
								state, err4 := manualSlice.LoadUInt(8)
								if err4 == nil {
									fmt.Printf("LoadUInt(8) State: %d, remaining bits: %d\n", state, manualSlice.BitsLeft())
								} else {
									fmt.Printf("LoadUInt(8) State FAILED: %v\n", err4)
								}
							} else {
								fmt.Printf("LoadUInt(8) State: INSUFFICIENT BITS (have %d, need 8)\n", manualSlice.BitsLeft())
							}

							fmt.Printf("=== End Analysis ===\n\n")
						}

						// Try automatic TLB decoding
						var event T
						err = tlb.LoadFromCell(&event, bodyCell.BeginParse())
						if err == nil {
							if isExecStateChanged {
								fmt.Printf("DEBUG: TLB auto-decode succeeded: %+v\n", event)
							}
							events = append(events, event)
						} else {
							if isExecStateChanged {
								fmt.Printf("DEBUG: TLB auto-decode failed: %v\n", err)
							}
						}
					}
				}
			}
		}
	}

	return events, nil
}
