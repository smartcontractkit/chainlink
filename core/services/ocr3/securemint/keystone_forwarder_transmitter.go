package securemint

import (
	"context"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// keystoneForwarderContractTransmitter implements the ocr3types.ContractTransmitter interface
// by submitting transactions to the KeystoneForwarder's report function via the ContractWriter
type keystoneForwarderContractTransmitter struct {
	logger           logger.Logger
	contractWriter   commontypes.ContractWriter
	fromAccount      ocr2types.Account
	forwarderAddress string
	receiverAddress  string
	chainSelector    por.ChainSelector
	contractName     string
}

// Ensure keystoneForwarderContractTransmitter implements the ContractTransmitter interface
var _ ocr3types.ContractTransmitter[por.ChainSelector] = (*keystoneForwarderContractTransmitter)(nil)

// createKeystoneForwarderContractTransmitter creates a contract transmitter that submits transactions to the KeystoneForwarder
func createKeystoneForwarderContractTransmitter(
	lggr logger.Logger,
	contractWriter commontypes.ContractWriter,
	fromAccount ocr2types.Account,
	forwarderAddress string,
	receiverAddress string,
	chainSelector por.ChainSelector,
) (ocr3types.ContractTransmitter[por.ChainSelector], error) {
	return &keystoneForwarderContractTransmitter{
		logger:           lggr,
		contractWriter:   contractWriter,
		fromAccount:      fromAccount,
		forwarderAddress: forwarderAddress,
		receiverAddress:  receiverAddress,
		chainSelector:    chainSelector,
		contractName:     "keystoneforwarder",
	}, nil
}

// Transmit submits a report to the KeystoneForwarder contract via the ContractWriter
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
		"receiverAddress", s.receiverAddress)

	// Convert signatures to the format expected by KeystoneForwarder
	// KeystoneForwarder expects signatures as bytes[] where each signature is 65 bytes
	signatures := make([][]byte, len(sigs))
	for i, sig := range sigs {
		signatures[i] = sig.Signature
	}

	// Create a dummy report context (96 bytes as expected by KeystoneForwarder)
	reportContext := make([]byte, 96)

	// Prepare the arguments for the report function
	// The KeystoneForwarder.report function signature is:
	// function report(address receiver, bytes calldata report, bytes calldata reportContext, bytes[] calldata signatures)
	args := map[string]any{
		"receiver":      s.receiverAddress,
		"report":        reportWithInfo.Report,
		"reportContext": reportContext,
		"signatures":    signatures,
	}

	// Generate a unique transaction ID
	txID, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Create transaction metadata
	meta := commontypes.TxMeta{
		// Add any relevant metadata here
	}

	// Submit the transaction via the ContractWriter
	// The ContractWriter will handle the transaction creation, signing, and broadcasting
	zero := big.NewInt(0)
	s.logger.Infow("Submitting transaction to contractWriter",
		"txID", txID.String(),
		"contractName", s.contractName,
		"method", "report",
		"chainSelector", s.chainSelector,
		"args", args)

	if err := s.contractWriter.SubmitTransaction(
		ctx,
		s.contractName,
		"report",
		args,
		fmt.Sprintf("%s-%s-%s", s.contractName, s.forwarderAddress, txID.String()),
		s.forwarderAddress,
		&meta,
		zero,
	); err != nil {
		s.logger.Errorw("Failed to submit transaction to contractWriter",
			"error", err,
			"chainSelector", s.chainSelector,
			"forwarderAddress", s.forwarderAddress)
		return fmt.Errorf("failed to submit transaction to contractWriter: %w", err)
	}

	s.logger.Infow("Transaction submitted successfully to contractWriter",
		"txID", txID.String(),
		"chainSelector", s.chainSelector,
		"forwarderAddress", s.forwarderAddress)

	return nil
}

// FromAccount returns the account that will be used to sign transactions
func (s *keystoneForwarderContractTransmitter) FromAccount(_ context.Context) (ocr2types.Account, error) {
	s.logger.Debugw("FromAccount called", "account", string(s.fromAccount))
	return s.fromAccount, nil
}
