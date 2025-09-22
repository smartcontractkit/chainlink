package testhelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/offramp"
	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tl"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func ConfirmCommitWithExpectedSeqNumRangeTON(t *testing.T, srcSelector uint64, tonChain cldf_ton.Chain, offRampContract address.Address, expectedSeqNumRange ccipocr3common.SeqNumRange) (bool, error) {
	seenMessages := NewCommitReportTracker(srcSelector, expectedSeqNumRange)
	tonBlockTicker := time.NewTicker(2500 * time.Millisecond)
	eventCh, errCh := TONEventEmitter[offramp.CommitReportAccepted](t.Context(), tonChain, offRampContract, 0, tonBlockTicker)

	timeout := time.NewTimer(tests.WaitTimeout(t))
	defer timeout.Stop()

	for {
		select {
		case commitEvent := <-eventCh:
			// if merkle root is zero, it only contains price updates
			if commitEvent.MerkleRoot == nil {
				t.Logf("Skipping CommitReportAccepted with only price updates")
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

func TONEventEmitter[T any](ctx context.Context, tonChain cldf_ton.Chain, contractAddress address.Address, startBlock uint32, ticker *time.Ticker) (<-chan T, <-chan error) {
	ch := make(chan T)
	errorCh := make(chan error)

	go func() {
		defer ticker.Stop()
		for {
			toBlock, err := tonChain.Client.CurrentMasterchainInfo(ctx)
			if err != nil {
				errorCh <- err
				return
			}

			res, err := tonChain.Client.WaitForBlock(toBlock.SeqNo).GetAccount(ctx, toBlock, &contractAddress)
			if err != nil {
				errorCh <- err
				return
			}

			txHash := res.LastTxHash

			var resp tl.Serializable
			err = tonChain.Client.Client().QueryLiteserver(ctx, ton.GetTransactions{
				Limit: int32(10),
				AccID: &ton.AccountID{
					Workchain: contractAddress.Workchain(),
					ID:        contractAddress.Data(),
				},
				LT:     int64(res.LastTxLT),
				TxHash: txHash,
			}, &resp)

			if err != nil {
				errorCh <- err
				return
			}

			var out T
			var txList []*cell.Cell
			var msgs []tlb.Message
			switch t := resp.(type) {
			case ton.TransactionList:
				if len(t.Transactions) == 0 {
					errorCh <- ton.ErrNoTransactionsWereFound
					return
				}

				txList, err = cell.FromBOCMultiRoot(t.Transactions)
				if err != nil {
					errorCh <- fmt.Errorf("failed to parse cell from transaction bytes: %w", err)
				}

				for i := 0; i < len(txList); i++ {
					loader := txList[i].BeginParse()

					var tx tlb.Transaction
					err = tlb.LoadFromCell(&tx, loader)
					if err != nil {
						errorCh <- fmt.Errorf("failed to load transaction from cell: %w", err)
						return
					}

					if tx.IO.Out != nil {
						msgs, err = tx.IO.Out.ToSlice()
						if err != nil {
							// skip this tx
							continue
						}

						for _, msg := range msgs {
							if msg.MsgType != tlb.MsgTypeExternalOut {
								continue
							}

							var event T
							c := msg.AsExternalOut().Body
							err = tlb.LoadFromCell(&event, c.BeginParse())
							if err == nil {
								ch <- out
							}
						}
					}
				}

			case ton.LSError:
				if t.Code == 0 {
					errorCh <- ton.ErrNoTransactionsWereFound
					return
				}
				errorCh <- t
				return
			}
		}
	}()

	return ch, errorCh
}
