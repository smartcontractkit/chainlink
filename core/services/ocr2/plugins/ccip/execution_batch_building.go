package ccip

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/merklemulti"
)

func getProofData(
	ctx context.Context,
	sourceReader ccipdata.OnRampReader,
	interval commit_store.CommitStoreInterval,
) (sendReqsInRoot []ccipdata.Event[ccipdata.EVM2EVMMessage], leaves [][32]byte, tree *merklemulti.Tree[[32]byte], err error) {
	sendReqs, err := sourceReader.GetSendRequestsBetweenSeqNums(
		ctx,
		interval.Min,
		interval.Max,
		0, // no need for confirmations, commitReport was already confirmed and we need all msgs in it
	)
	if err != nil {
		return nil, nil, nil, err
	}
	leaves = make([][32]byte, 0, len(sendReqs))
	for _, req := range sendReqs {
		leaves = append(leaves, req.Data.Hash)
	}
	tree, err = merklemulti.NewTree(hashlib.NewKeccakCtx(), leaves)
	if err != nil {
		return nil, nil, nil, err
	}
	return sendReqs, leaves, tree, nil
}

func buildExecutionReportForMessages(
	msgsInRoot []*evm_2_evm_offramp.InternalEVM2EVMMessage,
	leaves [][32]byte,
	tree *merklemulti.Tree[[32]byte],
	commitInterval commit_store.CommitStoreInterval,
	observedMessages []ObservedMessage,
) (report evm_2_evm_offramp.InternalExecutionReport, hashes [][32]byte, err error) {
	innerIdxs := make([]int, 0, len(observedMessages))
	report.Messages = []evm_2_evm_offramp.InternalEVM2EVMMessage{}
	for _, observedMessage := range observedMessages {
		if observedMessage.SeqNr < commitInterval.Min || observedMessage.SeqNr > commitInterval.Max {
			// We only return messages from a single root (the root of the first message).
			continue
		}
		innerIdx := int(observedMessage.SeqNr - commitInterval.Min)
		report.Messages = append(report.Messages, *msgsInRoot[innerIdx])
		report.OffchainTokenData = append(report.OffchainTokenData, observedMessage.TokenData)

		innerIdxs = append(innerIdxs, innerIdx)
		hashes = append(hashes, leaves[innerIdx])
	}

	merkleProof, err := tree.Prove(innerIdxs)
	if err != nil {
		return evm_2_evm_offramp.InternalExecutionReport{}, nil, err
	}

	// any capped proof will have length <= this one, so we reuse it to avoid proving inside loop, and update later if changed
	report.Proofs = merkleProof.Hashes
	report.ProofFlagBits = abihelpers.ProofFlagsToBits(merkleProof.SourceFlags)

	return report, hashes, nil
}

// Validates the given message observations do not exceed the committed sequence numbers
// in the commitStore.
func validateSeqNumbers(serviceCtx context.Context, commitStore commit_store.CommitStoreInterface, observedMessages []ObservedMessage) error {
	nextMin, err := commitStore.GetExpectedNextSequenceNumber(&bind.CallOpts{Context: serviceCtx})
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
func getCommitReportForSeqNum(ctx context.Context, destReader ccipdata.Reader, commitStore commit_store.CommitStoreInterface, seqNum uint64) (commit_store.CommitStoreCommitReport, error) {
	acceptedReports, err := destReader.GetAcceptedCommitReportsGteSeqNum(ctx, commitStore.Address(), seqNum, 0)
	if err != nil {
		return commit_store.CommitStoreCommitReport{}, err
	}

	for _, acceptedReport := range acceptedReports {
		reportInterval := acceptedReport.Data.Report.Interval
		if reportInterval.Min <= seqNum && seqNum <= reportInterval.Max {
			return acceptedReport.Data.Report, nil
		}
	}

	return commit_store.CommitStoreCommitReport{}, errors.Errorf("seq number not committed")
}
