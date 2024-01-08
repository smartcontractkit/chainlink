package ocr3

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type MultichainMeta interface {
	GetDestinationChainID() string
}

// multichainTransmitterOCR3 is a transmitter that can transmit to multiple chains.
// It uses the information in the MultichainMeta to determine which chain to transmit to.
// Note that this would only work with the appropriate multi-chain config tracker implementation.
type multichainTransmitterOCR3[RI MultichainMeta] struct {
	transmitters map[string]ocr3types.ContractTransmitter[RI]
	lp           logpoller.LogPoller
	lggr         logger.Logger
}

func NewMultichainTransmitterOCR3[RI MultichainMeta](
	transmitters map[string]ocr3types.ContractTransmitter[RI],
	lp logpoller.LogPoller,
	lggr logger.Logger,
) (*multichainTransmitterOCR3[RI], error) {
	return &multichainTransmitterOCR3[RI]{
		transmitters: transmitters,
		lp:           lp,
		lggr:         lggr,
	}, nil
}

// FromAccount implements ocr3types.ContractTransmitter.
func (m *multichainTransmitterOCR3[RI]) FromAccount() (types.Account, error) {
	var accounts []string
	for _, t := range m.transmitters {
		account, err := t.FromAccount()
		if err != nil {
			return "", err
		}
		accounts = append(accounts, string(account))
	}
	return types.Account(JoinTransmitters(accounts)), nil
}

// Transmit implements ocr3types.ContractTransmitter.
func (m *multichainTransmitterOCR3[RI]) Transmit(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64, rwi ocr3types.ReportWithInfo[RI], sigs []types.AttributedOnchainSignature) error {
	transmitter, ok := m.transmitters[rwi.Info.GetDestinationChainID()]
	if !ok {
		return fmt.Errorf("no transmitter for chain %s", rwi.Info.GetDestinationChainID())
	}
	return transmitter.Transmit(ctx, configDigest, seqNr, rwi, sigs)
}
