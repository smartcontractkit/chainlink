package txmeta

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/decode"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/version"
)

const (
	V1_2_0          = "1.2.0"
	V1_5_0          = "1.5.0"
	ManuallyExecute = "manuallyExecute"
	ReportAccepted  = "ReportAccepted"
)

// ExecReportToEthTxMeta generates a txmgr.EthTxMeta from the given report.
// Only MessageIDs will be populated in the TxMeta.
func ExecReportToEthTxMeta(ctx context.Context, typ version.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	if typ != version.EVM2EVMOffRamp {
		return nil, errors.Errorf("expected %v got %v", version.EVM2EVMOffRamp, typ)
	}
	switch ver.String() {
	case V1_2_0, V1_5_0:
		offRampABI := abihelpers.MustParseABI(evm_2_evm_offramp.EVM2EVMOffRampABI)
		return func(report []byte) (*txmgr.TxMeta, error) {
			execReport, err := decode.DecodeExecReport(ctx, abihelpers.MustGetMethodInputs(ManuallyExecute, offRampABI)[:1], report)
			if err != nil {
				return nil, err
			}
			return execReportToEthTxMeta(execReport)
		}, nil
	default:
		return nil, errors.Errorf("got unexpected version %v", ver.String())
	}
}

func execReportToEthTxMeta(execReport ccip.ExecReport) (*txmgr.TxMeta, error) {
	msgIDs := make([]string, len(execReport.Messages))
	for i, msg := range execReport.Messages {
		msgIDs[i] = hexutil.Encode(msg.MessageID[:])
	}

	return &txmgr.TxMeta{
		MessageIDs: msgIDs,
	}, nil
}

func CommitReportToEthTxMeta(typ version.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	if typ != version.CommitStore {
		return nil, errors.Errorf("expected %v got %v", version.CommitStore, typ)
	}
	switch ver.String() {
	case V1_2_0, V1_5_0:
		commitStoreABI := abihelpers.MustParseABI(commit_store.CommitStoreABI)
		return func(report []byte) (*txmgr.TxMeta, error) {
			commitReport, err := decode.DecodeCommitReport(abihelpers.MustGetEventInputs(ReportAccepted, commitStoreABI), report)
			if err != nil {
				return nil, err
			}
			return commitReportToEthTxMeta(commitReport)
		}, nil
	default:
		return nil, errors.Errorf("got unexpected version %v", ver.String())
	}
}

// CommitReportToEthTxMeta generates a txmgr.EthTxMeta from the given commit report.
// sequence numbers of the committed messages will be added to tx metadata
func commitReportToEthTxMeta(commitReport ccip.CommitStoreReport) (*txmgr.TxMeta, error) {
	n := (commitReport.Interval.Max - commitReport.Interval.Min) + 1
	seqRange := make([]uint64, n)
	for i := range n {
		seqRange[i] = i + commitReport.Interval.Min
	}
	return &txmgr.TxMeta{
		SeqNumbers: seqRange,
	}, nil
}
