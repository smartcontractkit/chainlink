package ocr3impls

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type MultichainMeta interface {
	GetDestinationChain() commontypes.RelayID
	GetDestinationConfigDigest() types.ConfigDigest
}

// multichainTransmitterOCR3 is a transmitter that can transmit to multiple chains.
// It uses the information in the MultichainMeta to determine which chain to transmit to.
// Note that this would only work with the appropriate multi-chain config tracker implementation.
type multichainTransmitterOCR3[RI MultichainMeta] struct {
	transmitters map[commontypes.RelayID]ocr3types.ContractTransmitter[RI]
	lggr         logger.Logger
}

func NewMultichainTransmitterOCR3[RI MultichainMeta](
	transmitters map[commontypes.RelayID]ocr3types.ContractTransmitter[RI],
	lggr logger.Logger,
) (*multichainTransmitterOCR3[RI], error) {
	return &multichainTransmitterOCR3[RI]{
		transmitters: transmitters,
		lggr:         lggr,
	}, nil
}

// FromAccount implements ocr3types.ContractTransmitter.
func (m *multichainTransmitterOCR3[RI]) FromAccount() (types.Account, error) {
	var accounts []string
	for relayID, t := range m.transmitters {
		account, err := t.FromAccount()
		if err != nil {
			return "", err
		}
		accounts = append(accounts, EncodeTransmitter(relayID, account))
	}
	return types.Account(JoinTransmitters(accounts)), nil
}

// Transmit implements ocr3types.ContractTransmitter.
func (m *multichainTransmitterOCR3[RI]) Transmit(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64, rwi ocr3types.ReportWithInfo[RI], sigs []types.AttributedOnchainSignature) error {
	destChain := rwi.Info.GetDestinationChain()
	transmitter, ok := m.transmitters[destChain]
	if !ok {
		return fmt.Errorf("no transmitter for chain %s", destChain)
	}
	m.lggr.Infow("multichain transmitter: transmitting to chain",
		"destChain", destChain.String(),
		"configDigest", rwi.Info.GetDestinationConfigDigest().Hex(),
	)
	return transmitter.Transmit(ctx, rwi.Info.GetDestinationConfigDigest(), seqNr, rwi, sigs)
}
