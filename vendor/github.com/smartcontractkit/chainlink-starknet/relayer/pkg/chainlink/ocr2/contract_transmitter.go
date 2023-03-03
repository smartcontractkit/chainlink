package ocr2

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"

	caigotypes "github.com/dontpanicdao/caigo/types"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2/medianreport"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ types.ContractTransmitter = (*contractTransmitter)(nil)

type contractTransmitter struct {
	reader *transmissionsCache

	contractAddress caigotypes.Hash
	senderAddress   caigotypes.Hash

	txm txm.TxManager
}

func NewContractTransmitter(
	reader *transmissionsCache,
	contract string,
	sender string,
	txm txm.TxManager,
) *contractTransmitter {
	return &contractTransmitter{
		reader:          reader,
		contractAddress: caigotypes.HexToHash(contract),
		senderAddress:   caigotypes.HexToHash(sender),
		txm:             txm,
	}
}

func (c *contractTransmitter) Transmit(
	ctx context.Context,
	reportCtx types.ReportContext,
	report types.Report,
	sigs []types.AttributedOnchainSignature,
) error {
	// flat array of arguments
	// convert everything to hex string -> caigo internally converts into big.int
	var transmitPayload []string

	// ReportContext:
	//    config_digest
	//    epoch_and_round
	//    extra_hash
	reportContext := medianreport.RawReportContext(reportCtx)

	for _, r := range reportContext {
		transmitPayload = append(transmitPayload, "0x"+hex.EncodeToString(r[:]))
	}

	slices, err := medianreport.SplitReport(report)
	if err != nil {
		return err
	}
	for i := 0; i < len(slices); i++ {
		hexStr := hex.EncodeToString(slices[i])
		transmitPayload = append(transmitPayload, "0x"+hexStr)
	}

	transmitPayload = append(transmitPayload, "0x"+fmt.Sprintf("%x", len(sigs))) // signatures_len
	for _, sig := range sigs {
		// signature: 32 byte public key + 32 byte R + 32 byte S
		signature := sig.Signature
		if len(signature) != 32+32+32 {
			return errors.New("invalid length of the signature")
		}
		transmitPayload = append(transmitPayload, "0x"+hex.EncodeToString(signature[32:64])) // r
		transmitPayload = append(transmitPayload, "0x"+hex.EncodeToString(signature[64:]))   // s
		transmitPayload = append(transmitPayload, "0x"+hex.EncodeToString(signature[:32]))   // public key
	}

	err = c.txm.Enqueue(c.senderAddress, caigotypes.FunctionCall{
		ContractAddress:    c.contractAddress,
		EntryPointSelector: "transmit",
		Calldata:           transmitPayload,
	})

	return err
}

func (c *contractTransmitter) LatestConfigDigestAndEpoch(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	err error,
) {
	configDigest, epoch, _, _, _, err = c.reader.LatestTransmissionDetails(ctx)
	if err != nil {
		err = errors.Wrap(err, "couldn't fetch latest transmission details")
	}

	return
}

func (c *contractTransmitter) FromAccount() types.Account {
	return types.Account(c.senderAddress.String())
}
