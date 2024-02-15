package ccipexec

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/merklemulti"
)

func getProofData(
	ctx context.Context,
	sourceReader ccipdata.OnRampReader,
	interval cciptypes.CommitStoreInterval,
) (sendReqsInRoot []cciptypes.EVM2EVMMessageWithTxMeta, leaves [][32]byte, tree *merklemulti.Tree[[32]byte], err error) {
	// We don't need to double-check if logs are finalized because we already checked that in the Commit phase.
	sendReqs, err := sourceReader.GetSendRequestsBetweenSeqNums(ctx, interval.Min, interval.Max, false)
	if err != nil {
		return nil, nil, nil, err
	}

	if err1 := validateSendRequests(sendReqs, interval); err1 != nil {
		return nil, nil, nil, err1
	}

	leaves = make([][32]byte, 0, len(sendReqs))
	for _, req := range sendReqs {
		leaves = append(leaves, req.Hash)
	}
	tree, err = merklemulti.NewTree(hashlib.NewKeccakCtx(), leaves)
	if err != nil {
		return nil, nil, nil, err
	}
	return sendReqs, leaves, tree, nil
}

func validateSendRequests(sendReqs []cciptypes.EVM2EVMMessageWithTxMeta, interval cciptypes.CommitStoreInterval) error {
	if len(sendReqs) == 0 {
		return fmt.Errorf("could not find any requests in the provided interval %v", interval)
	}

	gotInterval := cciptypes.CommitStoreInterval{
		Min: sendReqs[0].SequenceNumber,
		Max: sendReqs[0].SequenceNumber,
	}

	for _, req := range sendReqs[1:] {
		if req.SequenceNumber < gotInterval.Min {
			gotInterval.Min = req.SequenceNumber
		}
		if req.SequenceNumber > gotInterval.Max {
			gotInterval.Max = req.SequenceNumber
		}
	}

	if (gotInterval.Min != interval.Min) || (gotInterval.Max != interval.Max) {
		return fmt.Errorf("interval %v is not the expected %v", gotInterval, interval)
	}
	return nil
}

func buildExecutionReportForMessages(
	msgsInRoot []cciptypes.EVM2EVMMessageWithTxMeta,
	tree *merklemulti.Tree[[32]byte],
	commitInterval cciptypes.CommitStoreInterval,
	observedMessages []ccip.ObservedMessage,
) (cciptypes.ExecReport, error) {
	innerIdxs := make([]int, 0, len(observedMessages))
	var messages []cciptypes.EVM2EVMMessage
	var offchainTokenData [][][]byte
	for _, observedMessage := range observedMessages {
		if observedMessage.SeqNr < commitInterval.Min || observedMessage.SeqNr > commitInterval.Max {
			// We only return messages from a single root (the root of the first message).
			continue
		}
		innerIdx := int(observedMessage.SeqNr - commitInterval.Min)
		if innerIdx >= len(msgsInRoot) || innerIdx < 0 {
			return cciptypes.ExecReport{}, fmt.Errorf("invalid inneridx SeqNr=%d IntervalMin=%d msgsInRoot=%d",
				observedMessage.SeqNr, commitInterval.Min, len(msgsInRoot))
		}
		messages = append(messages, msgsInRoot[innerIdx].EVM2EVMMessage)
		offchainTokenData = append(offchainTokenData, observedMessage.TokenData)
		innerIdxs = append(innerIdxs, innerIdx)
	}

	merkleProof, err := tree.Prove(innerIdxs)
	if err != nil {
		return cciptypes.ExecReport{}, err
	}

	// any capped proof will have length <= this one, so we reuse it to avoid proving inside loop, and update later if changed
	return cciptypes.ExecReport{
		Messages:          messages,
		Proofs:            merkleProof.Hashes,
		ProofFlagBits:     abihelpers.ProofFlagsToBits(merkleProof.SourceFlags),
		OffchainTokenData: offchainTokenData,
	}, nil
}

// Validates the given message observations do not exceed the committed sequence numbers
// in the commitStoreReader.
func validateSeqNumbers(serviceCtx context.Context, commitStore ccipdata.CommitStoreReader, observedMessages []ccip.ObservedMessage) error {
	nextMin, err := commitStore.GetExpectedNextSequenceNumber(serviceCtx)
	if err != nil {
		return err
	}
	// observedMessages are always sorted by SeqNr and never empty, so it's safe to take last element
	maxSeqNumInBatch := observedMessages[len(observedMessages)-1].SeqNr

	if maxSeqNumInBatch >= nextMin {
		return errors.Errorf("Cannot execute uncommitted seq num. nextMin %v, seqNums %v", nextMin, observedMessages)
	}
	return nil
}

// Gets the commit report from the saved logs for a given sequence number.
func getCommitReportForSeqNum(ctx context.Context, commitStoreReader ccipdata.CommitStoreReader, seqNum uint64) (cciptypes.CommitStoreReport, error) {
	acceptedReports, err := commitStoreReader.GetCommitReportMatchingSeqNum(ctx, seqNum, 0)
	if err != nil {
		return cciptypes.CommitStoreReport{}, err
	}

	if len(acceptedReports) == 0 {
		return cciptypes.CommitStoreReport{}, errors.Errorf("seq number not committed")
	}

	return acceptedReports[0].CommitStoreReport, nil
}
