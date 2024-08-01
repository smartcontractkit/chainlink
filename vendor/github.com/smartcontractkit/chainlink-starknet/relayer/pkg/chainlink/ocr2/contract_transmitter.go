package ocr2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	starknetrpc "github.com/NethermindEth/starknet.go/rpc"
	starknetutils "github.com/NethermindEth/starknet.go/utils"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2/medianreport"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ types.ContractTransmitter = (*contractTransmitter)(nil)

type contractTransmitter struct {
	reader *transmissionsCache

	contractAddress *felt.Felt
	senderAddress   *felt.Felt // account.publicKey
	accountAddress  *felt.Felt

	txm txm.TxManager
}

func NewContractTransmitter(
	reader *transmissionsCache,
	contractAddress string,
	senderAddress string,
	accountAddress string,
	txm txm.TxManager,
) *contractTransmitter {
	contractAddr, _ := starknetutils.HexToFelt(contractAddress)
	senderAddr, _ := starknetutils.HexToFelt(senderAddress)
	accountAddr, _ := starknetutils.HexToFelt(accountAddress)

	return &contractTransmitter{
		reader:          reader,
		contractAddress: contractAddr,
		senderAddress:   senderAddr,
		accountAddress:  accountAddr,
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

	// TODO: build felts directly rather than afterwards
	calldata, err := starknetutils.HexArrToFelt(transmitPayload)
	if err != nil {
		return err
	}

	err = c.txm.Enqueue(c.accountAddress, c.senderAddress, starknetrpc.FunctionCall{
		ContractAddress:    c.contractAddress,
		EntryPointSelector: starknetutils.GetSelectorFromNameFelt("transmit"),
		Calldata:           calldata,
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
		err = fmt.Errorf("couldn't fetch latest transmission details: %w", err)
	}

	return
}

func (c *contractTransmitter) FromAccount() (types.Account, error) {
	return types.Account(c.accountAddress.String()), nil
}
