package transmitter

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/report"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"
)

var _ ocr3types.ContractTransmitter[[]byte] = (*ContractTransmitter)(nil)

type SendResponse func(ctx context.Context, response oracle.ConsensusResponse)

type ContractTransmitter struct {
	lggr         logger.Logger
	sendResponse SendResponse
}

func (c *ContractTransmitter) Transmit(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte], signatures []types.AttributedOnchainSignature) error {
	unmarshalledInfo := new(structpb.Struct)
	err := proto.Unmarshal(rwi.Info, unmarshalledInfo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal report info: %w", err)
	}
	infoMap := unmarshalledInfo.AsMap()
	requestID, ok := infoMap[plugin.InfoRequestID]
	if !ok {
		return errors.New("infoRequestID not found in report info")
	}

	requestIDStr, ok := requestID.(string)
	if !ok {
		return errors.New("infoRequestID is not a string")
	}

	if _, exists := infoMap[plugin.InfoConsensusFailureCode]; exists {
		failureCode, failureMessageStr, err := getFailureCodeAndMessageFromReportInfo(infoMap)
		if err != nil {
			return fmt.Errorf("failed to get failure code and message from report info: %w", err)
		}

		c.lggr.Debugw("received consensus failure message report", "requestID", requestIDStr,
			"failureMessage", failureMessageStr, "failureCode", failureCode)

		var failureErr caperrors.Error
		switch failureCode {
		case oracletypes.ConsensusFailureCode_RECEIVED_FPLUS1_ERRORS:
			// This is considered to be a user error as the caller of the consensus capability has sent too many errors and
			// so consensus cannot be reached.
			failureErr = caperrors.NewPublicUserError(errors.New(failureMessageStr), caperrors.ConsensusFailed)
		case oracletypes.ConsensusFailureCode_MORE_THAN_ONE_VALID_OUTCOME_FOR_IDENTICAL_CONSENSUS:
			// This is considered to be a user error as the caller of the consensus capability is attempting to achieve
			// identical consensus on a value which has multiple valid (>=f+1) observations sets.  For example this
			// could occur if the caller is attempting to achieve identical consensus on a value which is relatively volatile
			// resulting in f+1 nodes seeing value A, and f+1 nodes seeing value B.
			failureErr = caperrors.NewPublicUserError(errors.New(failureMessageStr), caperrors.ConsensusFailed)
		default:
			failureErr = caperrors.NewPublicSystemError(errors.New(failureMessageStr), caperrors.ConsensusFailed)
		}

		response := oracle.ConsensusResponse{
			ReqID: requestIDStr,
			SeqNr: seqNr,
			Err:   failureErr,
		}

		c.sendResponse(ctx, response)
	} else {
		c.lggr.Debugw("received consensus success report", "requestID", requestIDStr)
		// report context is the config digest + the sequence number padded with zeros
		repContext := report.GenerateReportContext(seqNr, configDigest)

		response := oracle.ConsensusResponse{
			ReqID:         requestIDStr,
			ConfigDigest:  configDigest,
			SeqNr:         seqNr,
			ReportContext: repContext,
			RawReport:     rwi.Report,
			Sigs:          signatures,
		}

		c.sendResponse(ctx, response)
	}

	return nil
}

func getFailureCodeAndMessageFromReportInfo(infoMap map[string]any) (oracletypes.ConsensusFailureCode, string, error) {
	failureCodeEntry, failureCodeExists := infoMap[plugin.InfoConsensusFailureCode]
	if !failureCodeExists {
		return 0, "", errors.New("failure code not found in report info")
	}

	failureCodeStr, ok := failureCodeEntry.(string)
	if !ok {
		return 0, "", errors.New("failure code is not a string")
	}

	codeInt, ok := oracletypes.ConsensusFailureCode_value[failureCodeStr]
	if !ok {
		return 0, "", fmt.Errorf("invalid failure code: %s", failureCodeStr)
	}

	failureCode := oracletypes.ConsensusFailureCode(codeInt)

	failureMessage, failureMsgExists := infoMap[plugin.InfoConsensusFailureMessage]
	if !failureMsgExists {
		return 0, "", errors.New("failure message not found in report info")
	}

	failureMessageStr, ok := failureMessage.(string)
	if !ok {
		return 0, "", errors.New("message is not a string")
	}
	return failureCode, failureMessageStr, nil
}

func (c *ContractTransmitter) FromAccount(_ context.Context) (types.Account, error) {
	return "", nil
}

func NewContractTransmitter(lggr logger.Logger, sendResponse SendResponse) *ContractTransmitter {
	return &ContractTransmitter{lggr: lggr, sendResponse: sendResponse}
}
