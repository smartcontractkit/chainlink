package targets

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	_ capabilities.ActionCapability = &WriteTarget{}
)

// required field of target's config in the workflow spec
const signedReportField = "signed_report"

type WriteTarget struct {
	cr               commontypes.ContractReader
	cw               commontypes.ChainWriter
	forwarderAddress string
	capabilities.CapabilityInfo
	lggr logger.Logger

	// options for configuring local node info retrieval
	registry         core.CapabilitiesRegistry
	localNodeRetryMs int // how often to retry fetching local node info
	localNode        *capabilities.Node
}

func NewWriteTarget(lggr logger.Logger, id string, cr commontypes.ContractReader, cw commontypes.ChainWriter, forwarderAddress string, registry core.CapabilitiesRegistry, localNodeInfoIntervalMs int) *WriteTarget {
	info := capabilities.MustNewCapabilityInfo(
		id,
		capabilities.CapabilityTypeTarget,
		"Write target.",
	)

	logger := lggr.Named("WriteTarget")
	if localNodeInfoIntervalMs == 0 {
		localNodeInfoIntervalMs = 5000
	}

	wt := &WriteTarget{
		cr,
		cw,
		forwarderAddress,
		info,
		logger,
		registry,
		localNodeInfoIntervalMs,
		nil,
	}

	wt.ResolveLocalNodeInfo(context.Background())
	return wt
}

type EvmConfig struct {
	Address string
}

func parseConfig(rawConfig *values.Map) (config EvmConfig, err error) {
	if err := rawConfig.UnwrapTo(&config); err != nil {
		return config, err
	}
	if !common.IsHexAddress(config.Address) {
		return config, fmt.Errorf("'%v' is not a valid address", config.Address)
	}
	return config, nil
}

func success() <-chan capabilities.CapabilityResponse {
	callback := make(chan capabilities.CapabilityResponse)
	go func() {
		callback <- capabilities.CapabilityResponse{}
		close(callback)
	}()
	return callback
}

// ResolveLocalNodeInfo fetches the local node info and updates the logger with the peerID and workflowDONID
func (cap *WriteTarget) ResolveLocalNodeInfo(ctx context.Context) {
	if cap.registry == nil {
		cap.lggr.Warn("Capabilities registry not set, skipping ResolveLocalNodeInfo")
		return
	}

	go func() {
		err := utils.Retryable(ctx, cap.lggr.Errorf, cap.localNodeRetryMs, 0, func() error {
			node, iErr := cap.registry.GetLocalNode(ctx)
			if iErr != nil {
				return iErr
			}
			if node.PeerID == nil {
				return fmt.Errorf("local node not found or not part of DON")
			}

			cap.lggr = cap.lggr.With("peerID", node.PeerID, "workflowDONID", node.WorkflowDON.ID, "workflowDONConfigVersion", node.WorkflowDON.ConfigVersion)
			cap.lggr.Debug("Resolved local node info")
			cap.localNode = &node
			return nil
		})

		if err != nil {
			cap.lggr.Errorf("Failed to resolve local node info: %v", err)
		}
	}()
}

func (cap *WriteTarget) Execute(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	l := cap.lggr.With("workflowID", request.Metadata.WorkflowID, "executionID", request.Metadata.WorkflowExecutionID)
	// both donid and configversion start at 1, so if either is 0, this node is purely a capability node and not part of a workflow DON
	if cap.localNode == nil || cap.localNode.WorkflowDON.ID == 0 || cap.localNode.WorkflowDON.ConfigVersion == 0 {
		l = l.With("workflowDONID", request.Metadata.WorkflowDonID, "workflowDONConfigVersion", request.Metadata.WorkflowDonConfigVersion)
	}
	l.Debugw("Execute", "request", request)

	reqConfig, err := parseConfig(request.Config)
	if err != nil {
		return nil, err
	}

	signedReport, ok := request.Inputs.Underlying[signedReportField]
	if !ok {
		return nil, fmt.Errorf("missing required field %s", signedReportField)
	}

	inputs := types.SignedReport{}
	if err = signedReport.UnwrapTo(&inputs); err != nil {
		return nil, err
	}

	if len(inputs.Report) == 0 {
		// We received any empty report -- this means we should skip transmission.
		l.Debugw("Skipping empty report")
		return success(), nil
	}
	// TODO: validate encoded report is prefixed with workflowID and executionID that match the request meta

	rawExecutionID, err := hex.DecodeString(request.Metadata.WorkflowExecutionID)
	if err != nil {
		return nil, err
	}
	// Check whether value was already transmitted on chain
	queryInputs := struct {
		Receiver            string
		WorkflowExecutionID []byte
		ReportId            []byte
	}{
		Receiver:            reqConfig.Address,
		WorkflowExecutionID: rawExecutionID,
		ReportId:            inputs.ID,
	}
	l = l.With("receiver", queryInputs.Receiver, "reportID", queryInputs.ReportId)
	var transmitter common.Address
	if err = cap.cr.GetLatestValue(ctx, "forwarder", "getTransmitter", queryInputs, &transmitter); err != nil {
		return nil, err
	}
	if transmitter != common.HexToAddress("0x0") {
		l.Infow("WriteTarget report already onchain - returning without a tranmission attempt")
		return success(), nil
	}

	l.Infow("WriteTarget non-empty report - attempting to push to txmgr", "reportLen", len(inputs.Report), "reportContextLen", len(inputs.Context), "nSignatures", len(inputs.Signatures))
	txID, err := uuid.NewUUID() // NOTE: CW expects us to generate an ID, rather than return one
	if err != nil {
		return nil, err
	}

	// Note: The codec that ChainWriter uses to encode the parameters for the contract ABI cannot handle
	// `nil` values, including for slices. Until the bug is fixed we need to ensure that there are no
	// `nil` values passed in the request.
	req := struct {
		Receiver      string
		RawReport     []byte
		ReportContext []byte
		Signatures    [][]byte
	}{reqConfig.Address, inputs.Report, inputs.Context, inputs.Signatures}

	if req.RawReport == nil {
		req.RawReport = make([]byte, 0)
	}

	if req.ReportContext == nil {
		req.ReportContext = make([]byte, 0)
	}

	if req.Signatures == nil {
		req.Signatures = make([][]byte, 0)
	}
	l.Debugw("Transaction raw report", "report", hex.EncodeToString(req.RawReport))

	meta := commontypes.TxMeta{WorkflowExecutionID: &request.Metadata.WorkflowExecutionID}
	value := big.NewInt(0)
	if err := cap.cw.SubmitTransaction(ctx, "forwarder", "report", req, txID.String(), cap.forwarderAddress, &meta, value); err != nil {
		return nil, err
	}
	l.Debugw("Transaction submitted", "transactionID", txID)
	return success(), nil
}

func (cap *WriteTarget) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *WriteTarget) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
