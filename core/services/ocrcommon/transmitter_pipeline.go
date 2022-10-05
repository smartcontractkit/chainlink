package ocrcommon

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

const txObservationSource = `
    transmit_tx [type=ethtx
                 minConfirmations=0
                 to="$(jobSpec.contractAddress)"
                 from="[$(jobSpec.fromAddress)]"
                 evmChainID="$(jobSpec.evmChainID)"
                 data="$(jobSpec.data)"
                 gasLimit="$(jobSpec.gasLimit)"]
    transmit_tx
`

type pipelineTransmitter struct {
	lgr                         logger.Logger
	fromAddress                 common.Address
	gasLimit                    uint32
	effectiveTransmitterAddress common.Address
	strategy                    txmgr.TxStrategy
	checker                     txmgr.TransmitCheckerSpec
	pr                          pipeline.Runner
	chainID                     string
}

// NewPipelineTransmitter creates a new eth transmitter using the job pipeline mechanism
func NewPipelineTransmitter(
	lgr logger.Logger,
	fromAddress common.Address,
	gasLimit uint32,
	effectiveTransmitterAddress common.Address,
	strategy txmgr.TxStrategy,
	checker txmgr.TransmitCheckerSpec,
	pr pipeline.Runner,
	chainID string,
) Transmitter {
	return &pipelineTransmitter{
		lgr:                         lgr,
		fromAddress:                 fromAddress,
		gasLimit:                    gasLimit,
		effectiveTransmitterAddress: effectiveTransmitterAddress,
		strategy:                    strategy,
		checker:                     checker,
		pr:                          pr,
		chainID:                     chainID,
	}
}

func (t *pipelineTransmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error {
	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"contractAddress": toAddress.String(),
			"fromAddress":     t.fromAddress.String(),
			"gasLimit":        t.gasLimit,
			"evmChainID":      t.chainID,
			"data":            payload,
		},
	})

	run := pipeline.NewRun(pipeline.Spec{
		DotDagSource: txObservationSource,
		GasLimit:     &t.gasLimit,
	}, vars)

	if _, err := t.pr.Run(ctx, &run, t.lgr, true, nil); err != nil {
		return errors.Wrap(err, "Skipped OCR transmission")
	}

	return nil
}

func (t *pipelineTransmitter) FromAddress() common.Address {
	return t.effectiveTransmitterAddress
}

func (t *pipelineTransmitter) forwarderAddress() common.Address {
	if t.effectiveTransmitterAddress != t.fromAddress {
		return t.effectiveTransmitterAddress
	}
	return common.Address{}
}
