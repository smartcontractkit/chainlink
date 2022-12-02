package ocrcommon

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

const txObservationSource = `
    transmit_tx [type=ethtx
                 minConfirmations=0
                 to="$(jobSpec.contractAddress)"
                 from="[$(jobSpec.fromAddress)]"
                 evmChainID="$(jobSpec.evmChainID)"
                 data="$(jobSpec.data)"
                 gasLimit="$(jobSpec.gasLimit)"
                 forwardingAllowed="$(jobSpec.forwardingAllowed)"
                 transmitChecker="$(jobSpec.transmitChecker)"]
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
	spec                        job.Job
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
	spec job.Job,
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
		spec:                        spec,
		chainID:                     chainID,
	}
}

func (t *pipelineTransmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error {
	// t.strategy is ignored currently as pipeline does not support passing this (sc-55115)
	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"contractAddress":   toAddress.String(),
			"fromAddress":       t.fromAddress.String(),
			"gasLimit":          t.gasLimit,
			"evmChainID":        t.chainID,
			"forwardingAllowed": t.spec.ForwardingAllowed,
			"data":              payload,
			"transmitChecker":   t.checker,
		},
	})

	t.spec.PipelineSpec.DotDagSource = txObservationSource
	run := pipeline.NewRun(*t.spec.PipelineSpec, vars)

	if _, err := t.pr.Run(ctx, &run, t.lgr, true, nil); err != nil {
		return errors.Wrap(err, "Skipped OCR transmission")
	}

	if run.State != pipeline.RunStatusCompleted {
		return fmt.Errorf("unexpected pipeline run state: %s", run.State)
	}

	return nil
}

func (t *pipelineTransmitter) FromAddress() common.Address {
	return t.effectiveTransmitterAddress
}
