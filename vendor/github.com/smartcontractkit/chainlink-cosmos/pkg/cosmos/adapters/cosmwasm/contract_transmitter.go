package cosmwasm

import (
	"context"
	"encoding/json"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cosmosSDK "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"

	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ types.ContractTransmitter = (*ContractTransmitter)(nil)

type ContractTransmitter struct {
	*OCR2Reader
	msgEnqueuer adapters.MsgEnqueuer
	lggr        logger.Logger
	jobID       string
	contract    cosmosSDK.AccAddress
	sender      cosmosSDK.AccAddress
	cfg         config.Config
}

func NewContractTransmitter(
	reader *OCR2Reader,
	jobID string,
	contract cosmosSDK.AccAddress,
	sender cosmosSDK.AccAddress,
	msgEnqueuer adapters.MsgEnqueuer,
	lggr logger.Logger,
	cfg config.Config,
) *ContractTransmitter {
	return &ContractTransmitter{
		OCR2Reader:  reader,
		jobID:       jobID,
		contract:    contract,
		msgEnqueuer: msgEnqueuer,
		sender:      sender,
		lggr:        lggr,
		cfg:         cfg,
	}
}

// Transmit signs and sends the report
func (ct *ContractTransmitter) Transmit(
	ctx context.Context,
	reportCtx types.ReportContext,
	report types.Report,
	sigs []types.AttributedOnchainSignature,
) error {
	ct.lggr.Infof("[%s] Sending TX to %s", ct.jobID, ct.contract.String())
	msgStruct := TransmitMsg{}
	reportContext := evmutil.RawReportContext(reportCtx)
	for _, r := range reportContext {
		msgStruct.Transmit.ReportContext = append(msgStruct.Transmit.ReportContext, r[:]...)
	}
	msgStruct.Transmit.Report = []byte(report)
	for _, sig := range sigs {
		msgStruct.Transmit.Signatures = append(msgStruct.Transmit.Signatures, sig.Signature)
	}
	msgBytes, err := json.Marshal(msgStruct)
	if err != nil {
		return err
	}
	m := &wasmtypes.MsgExecuteContract{
		Sender:   ct.sender.String(),
		Contract: ct.contract.String(),
		Msg:      msgBytes,
		Funds:    cosmosSDK.Coins{},
	}
	_, err = ct.msgEnqueuer.Enqueue(ctx, ct.contract.String(), m)
	return err
}

func (ct *ContractTransmitter) FromAccount() (types.Account, error) {
	return types.Account(ct.sender.String()), nil
}
