package securemint

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// secureMintContractTransmitter implements the ocr3types.ContractTransmitter interface
// by using a chain writer to submit transactions to the blockchain
type secureMintContractTransmitter struct {
	logger          logger.Logger
	chainWriter     types.ContractWriter
	fromAccount     ocrtypes.Account
	contractAddress string
	chainSelector   por.ChainSelector
}

// Ensure secureMintContractTransmitter implements the ContractTransmitter interface
var _ ocr3types.ContractTransmitter[por.ChainSelector] = (*secureMintContractTransmitter)(nil)

// Transmit submits a report to the blockchain using the chain writer
func (s *secureMintContractTransmitter) Transmit(
	ctx context.Context,
	configDigest ocrtypes.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[por.ChainSelector],
	sigs []ocrtypes.AttributedOnchainSignature,
) error {
	s.logger.Infow("Transmitting report",
		"configDigest", fmt.Sprintf("%x", configDigest),
		"sequenceNumber", seqNr,
		"reportLength", len(reportWithInfo.Report),
		"reportInfo", reportWithInfo.Info,
		"signaturesCount", len(sigs),
		"chainSelector", s.chainSelector)

	// Create a unique transaction ID for this transmission
	transactionID := fmt.Sprintf("secure-mint-%d-%d-%x", s.chainSelector, seqNr, configDigest[:8])

	// Prepare the transaction arguments
	// This would typically include the report data, signatures, and other metadata
	args := map[string]interface{}{
		"configDigest": configDigest,
		"seqNr":        seqNr,
		"report":       reportWithInfo.Report,
		"signatures":   sigs,
		"info":         reportWithInfo.Info,
	}

	// Submit the transaction using the chain writer
	err := s.chainWriter.SubmitTransaction(
		ctx,
		"SecureMint",      // contract name
		"transmit",        // method name
		args,              // arguments
		transactionID,     // transaction ID
		s.contractAddress, // contract address
		nil,               // metadata
		big.NewInt(0),     // value (0 for regular transactions)
	)
	if err != nil {
		s.logger.Errorw("Failed to submit transaction",
			"error", err,
			"transactionID", transactionID,
			"chainSelector", s.chainSelector)
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	s.logger.Infow("Transaction submitted successfully",
		"transactionID", transactionID,
		"chainSelector", s.chainSelector)

	return nil
}

// FromAccount returns the account that will be used to sign transactions
func (s *secureMintContractTransmitter) FromAccount(_ context.Context) (ocrtypes.Account, error) {
	s.logger.Debugw("FromAccount called", "account", string(s.fromAccount))
	return s.fromAccount, nil
}
