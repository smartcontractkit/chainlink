package evm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type MultichainMeta interface {
	GetDestinationChainID() string
}

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
	// return a deterministic account based on the list of accounts
	// it has to be 20 bytes long because the transmitter is of type address in
	// the solidity contracts.
	h := gethcrypto.Keccak256Hash([]byte(strings.Join(accounts, ",")))
	return types.Account(hexutil.Encode(h[:20])), nil
}

// Transmit implements ocr3types.ContractTransmitter.
func (m *multichainTransmitterOCR3[RI]) Transmit(ctx context.Context, configDigest types.ConfigDigest, seqNr uint64, rwi ocr3types.ReportWithInfo[RI], sigs []types.AttributedOnchainSignature) error {
	transmitter, ok := m.transmitters[rwi.Info.GetDestinationChainID()]
	if !ok {
		return fmt.Errorf("no transmitter for chain %s", rwi.Info.GetDestinationChainID())
	}
	return transmitter.Transmit(ctx, configDigest, seqNr, rwi, sigs)
}
