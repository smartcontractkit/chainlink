package agentcommerce

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type CCIPSubmission struct {
	RequestHash              string
	SourceChainSelector      ccipocr3.ChainSelector
	DestinationChainSelector ccipocr3.ChainSelector
	MessageID                ccipocr3.Bytes32
	SourceTransactionHash    ccipocr3.Bytes32
	Phase                    SettlementPhase
}

type CCIPClient interface {
	SubmitCCIP(context.Context, RailSettlementRequest, ccipocr3.ChainSelector, ccipocr3.ChainSelector) (CCIPSubmission, error)
	ObserveCCIP(context.Context, RailSettlementRequest, CCIPSubmission) (RailSettlementReceipt, error)
}

func (a *CCIPSettlementAdapter) Observe(ctx context.Context, req RailSettlementRequest, submission CCIPSubmission) (RailSettlementReceipt, error) {
	wantHash, err := SettlementRequestHash(req)
	if err != nil {
		return RailSettlementReceipt{}, err
	}
	if submission.RequestHash != wantHash || submission.MessageID == (ccipocr3.Bytes32{}) || submission.SourceTransactionHash == (ccipocr3.Bytes32{}) {
		return RailSettlementReceipt{}, errors.New("invalid CCIP submission identity")
	}
	receipt, err := a.client.ObserveCCIP(ctx, req, submission)
	if err != nil {
		return RailSettlementReceipt{}, err
	}
	if receipt.RequestHash != wantHash || receipt.Rail != "ccip" || receipt.ExternalID != hex.EncodeToString(submission.MessageID[:]) {
		return RailSettlementReceipt{}, errors.New("CCIP destination observation identity mismatch")
	}
	return SealSettlementReceipt(receipt)
}

type CCIPSettlementAdapter struct{ client CCIPClient }

func NewCCIPSettlementAdapter(client CCIPClient) (*CCIPSettlementAdapter, error) {
	if client == nil {
		return nil, errors.New("CCIP client is required")
	}
	return &CCIPSettlementAdapter{client: client}, nil
}

func (a *CCIPSettlementAdapter) Submit(ctx context.Context, req RailSettlementRequest, source, destination ccipocr3.ChainSelector) (CCIPSubmission, error) {
	if req.Rail != "ccip" {
		return CCIPSubmission{}, errors.New("CCIP adapter requires ccip rail")
	}
	if source == 0 || destination == 0 || source == destination {
		return CCIPSubmission{}, errors.New("distinct non-zero chain selectors are required")
	}
	wantHash, err := SettlementRequestHash(req)
	if err != nil {
		return CCIPSubmission{}, err
	}
	result, err := a.client.SubmitCCIP(ctx, req, source, destination)
	if err != nil {
		return CCIPSubmission{}, fmt.Errorf("submit CCIP message: %w", err)
	}
	if result.RequestHash != wantHash || result.SourceChainSelector != source || result.DestinationChainSelector != destination {
		return CCIPSubmission{}, errors.New("CCIP submission identity mismatch")
	}
	if result.MessageID == (ccipocr3.Bytes32{}) || result.SourceTransactionHash == (ccipocr3.Bytes32{}) {
		return CCIPSubmission{}, errors.New("CCIP submission evidence is incomplete")
	}
	if result.Phase != SettlementSubmitted && result.Phase != SettlementBroadcast {
		return CCIPSubmission{}, errors.New("CCIP submission returned an unsupported phase")
	}
	return result, nil
}
