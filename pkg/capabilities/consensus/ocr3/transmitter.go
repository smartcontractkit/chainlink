package ocr3

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	pbtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

var _ (ocr3types.ContractTransmitter[[]byte]) = (*ContractTransmitter)(nil)

// ContractTransmitter is a custom transmitter for the OCR3 capability.
// When called it will forward the report + its signatures back to the
// OCR3 capability by making a call to Execute with a special "method"
// parameter.
type ContractTransmitter struct {
	lggr        logger.Logger
	registry    core.CapabilitiesRegistry
	capability  capabilities.CallbackCapability
	fromAccount string
}

func (c *ContractTransmitter) Transmit(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte], signatures []types.AttributedOnchainSignature) error {
	info := &pbtypes.ReportInfo{}
	err := proto.Unmarshal(rwi.Info, info)
	if err != nil {
		c.lggr.Error("could not unmarshal info")
		return err
	}

	resp := map[string]any{
		methodHeader: methodSendResponse,
	}
	if info.ShouldReport {
		resp["report"] = []byte(rwi.Report)

		sigs := [][]byte{}
		for _, s := range signatures {
			sigs = append(sigs, s.Signature)
		}
		resp["signatures"] = sigs
	} else {
		resp["report"] = nil
		resp["signatures"] = [][]byte{}
	}

	inputs, err := values.Wrap(resp)
	if err != nil {
		c.lggr.Error("could not wrap report", "payload", resp)
		return err
	}

	c.lggr.Debugw("ContractTransmitter transmitting", "shouldReport", info.ShouldReport, "len", len(rwi.Report))
	if c.capability == nil {
		cp, innerErr := c.registry.Get(ctx, ocrCapabilityID)
		if innerErr != nil {
			return fmt.Errorf("failed to fetch ocr3 capability from registry: %w", innerErr)
		}

		c.capability = cp.(capabilities.CallbackCapability)
	}

	_, err = c.capability.Execute(ctx, capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowExecutionID: info.Id.WorkflowExecutionId,
			WorkflowID:          info.Id.WorkflowId,
		},
		Inputs: inputs.(*values.Map),
	})
	if err != nil {
		c.lggr.Errorw("could not transmit response", "error", err, "weid", info.Id.WorkflowExecutionId)
	}
	return err
}

func (c *ContractTransmitter) FromAccount() (types.Account, error) {
	return types.Account(c.fromAccount), nil
}

func NewContractTransmitter(lggr logger.Logger, registry core.CapabilitiesRegistry, fromAccount string) *ContractTransmitter {
	return &ContractTransmitter{lggr: lggr, registry: registry, fromAccount: fromAccount}
}
