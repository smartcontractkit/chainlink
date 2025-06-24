package securemint

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// keystoneForwarderContractTransmitter implements the ocr3types.ContractTransmitter interface
// by calling the KeystoneForwarder's report function
type keystoneForwarderContractTransmitter struct {
	logger          logger.Logger
	forwarder       *forwarder.KeystoneForwarder
	fromAccount     ocr2types.Account
	receiverAddress common.Address
	chainSelector   por.ChainSelector
	ethKey          ethkey.KeyV2
	backend         bind.ContractBackend
}

// Ensure keystoneForwarderContractTransmitter implements the ContractTransmitter interface
var _ ocr3types.ContractTransmitter[por.ChainSelector] = (*keystoneForwarderContractTransmitter)(nil)

// createKeystoneForwarderContractTransmitter creates a contract transmitter that calls the KeystoneForwarder
func createKeystoneForwarderContractTransmitter(
	lggr logger.Logger,
	forwarderAddress common.Address,
	receiverAddress common.Address,
	fromAccount ocr2types.Account,
	chainSelector por.ChainSelector,
	ethKey ethkey.KeyV2,
	backend bind.ContractBackend,
) (ocr3types.ContractTransmitter[por.ChainSelector], error) {
	forwarder, err := forwarder.NewKeystoneForwarder(forwarderAddress, backend)
	if err != nil {
		return nil, fmt.Errorf("failed to create forwarder contract instance: %w", err)
	}

	return &keystoneForwarderContractTransmitter{
		logger:          lggr,
		forwarder:       forwarder,
		fromAccount:     fromAccount,
		receiverAddress: receiverAddress,
		chainSelector:   chainSelector,
		ethKey:          ethKey,
		backend:         backend,
	}, nil
}

// Transmit submits a report to the KeystoneForwarder contract
func (s *keystoneForwarderContractTransmitter) Transmit(
	ctx context.Context,
	configDigest ocr2types.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[por.ChainSelector],
	sigs []ocr2types.AttributedOnchainSignature,
) error {
	s.logger.Infow("Transmitting report to KeystoneForwarder",
		"configDigest", fmt.Sprintf("%x", configDigest),
		"sequenceNumber", seqNr,
		"reportLength", len(reportWithInfo.Report),
		"reportInfo", reportWithInfo.Info,
		"signaturesCount", len(sigs),
		"chainSelector", s.chainSelector,
		"receiverAddress", s.receiverAddress.Hex())

	// Convert signatures to the format expected by KeystoneForwarder
	// KeystoneForwarder expects signatures as bytes[] where each signature is 65 bytes
	signatures := make([][]byte, len(sigs))
	for i, sig := range sigs {
		signatures[i] = sig.Signature
	}

	// Create transaction options
	// For now, we'll use a simple approach since this is just for testing
	// In a real implementation, we'd need to properly extract the private key
	auth := &bind.TransactOpts{
		From: s.ethKey.Address,
		Signer: func(address common.Address, tx *gethtypes.Transaction) (*gethtypes.Transaction, error) {
			// This is a simplified approach for testing
			// In a real implementation, we'd properly sign the transaction
			return tx, nil
		},
	}

	// Call the KeystoneForwarder's report function
	// For now, we'll use a dummy receiver address (address(0)) since we're not forwarding to a real receiver
	// In a real implementation, this would be the DataFeedsCache address
	receiver := common.Address{} // address(0) for now

	// Create a dummy report context (96 bytes as expected by KeystoneForwarder)
	reportContext := make([]byte, 96)

	// Call the report function
	tx, err := s.forwarder.Report(auth, receiver, reportWithInfo.Report, reportContext, signatures)
	if err != nil {
		s.logger.Errorw("Failed to submit report to KeystoneForwarder",
			"error", err,
			"chainSelector", s.chainSelector,
			"receiverAddress", receiver.Hex())
		return fmt.Errorf("failed to submit report to KeystoneForwarder: %w", err)
	}

	s.logger.Infow("Report submitted successfully to KeystoneForwarder",
		"transactionHash", tx.Hash().Hex(),
		"chainSelector", s.chainSelector,
		"receiverAddress", receiver.Hex())

	return nil
}

// FromAccount returns the account that will be used to sign transactions
func (s *keystoneForwarderContractTransmitter) FromAccount(_ context.Context) (ocr2types.Account, error) {
	s.logger.Debugw("FromAccount called", "account", string(s.fromAccount))
	return s.fromAccount, nil
}
