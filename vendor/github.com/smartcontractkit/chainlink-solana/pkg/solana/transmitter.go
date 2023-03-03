package solana

import (
	"bytes"
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/logger"
)

var _ types.ContractTransmitter = (*Transmitter)(nil)

type Transmitter struct {
	stateID, programID, storeProgramID, transmissionsID, transmissionSigner solana.PublicKey
	reader                                                                  client.Reader
	stateCache                                                              *StateCache
	lggr                                                                    logger.Logger
	txManager                                                               TxManager
}

// Transmit sends the report to the on-chain OCR2Aggregator smart contract's Transmit method
func (c *Transmitter) Transmit(
	ctx context.Context,
	reportCtx types.ReportContext,
	report types.Report,
	sigs []types.AttributedOnchainSignature,
) error {
	blockhash, err := c.reader.LatestBlockhash()
	if err != nil {
		return errors.Wrap(err, "error on Transmit.GetRecentBlockhash")
	}
	if blockhash == nil || blockhash.Value == nil {
		return errors.New("nil pointer returned from Transmit.GetRecentBlockhash")
	}

	// Determine store authority
	seeds := [][]byte{[]byte("store"), c.stateID.Bytes()}
	storeAuthority, storeNonce, err := solana.FindProgramAddress(seeds, c.programID)
	if err != nil {
		return errors.Wrap(err, "error on Transmit.FindProgramAddress")
	}

	accounts := []*solana.AccountMeta{
		// state, transmitter, transmissions, store_program, store, store_authority, instructions_sysvar
		{PublicKey: c.stateID, IsWritable: true, IsSigner: false},
		{PublicKey: c.transmissionSigner, IsWritable: false, IsSigner: true},
		{PublicKey: c.transmissionsID, IsWritable: true, IsSigner: false},
		{PublicKey: c.storeProgramID, IsWritable: false, IsSigner: false},
		{PublicKey: storeAuthority, IsWritable: false, IsSigner: false},
		{PublicKey: solana.SysVarInstructionsPubkey, IsWritable: false, IsSigner: false},
	}

	reportContext := utils.RawReportContext(reportCtx)

	// Construct the instruction payload
	data := new(bytes.Buffer) // store_nonce || report_context || raw_report || raw_signatures
	data.WriteByte(storeNonce)
	data.Write(reportContext[0][:])
	data.Write(reportContext[1][:])
	data.Write(reportContext[2][:])
	data.Write([]byte(report))
	for _, sig := range sigs {
		// Signature = 64 bytes + 1 byte recovery id
		data.Write(sig.Signature)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			solana.NewInstruction(c.programID, accounts, data.Bytes()),
		},
		blockhash.Value.Blockhash,
		solana.TransactionPayer(c.transmissionSigner),
	)
	if err != nil {
		return errors.Wrap(err, "error on Transmit.NewTransaction")
	}

	// pass transmit payload to tx manager queue
	c.lggr.Debugf("Queuing transmit tx: state (%s) + transmissions (%s)", c.stateID.String(), c.transmissionsID.String())
	err = c.txManager.Enqueue(c.stateID.String(), tx)
	return errors.Wrap(err, "error on Transmit.txManager.Enqueue")
}

func (c *Transmitter) LatestConfigDigestAndEpoch(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	err error,
) {
	state, err := c.stateCache.ReadState()
	return state.Config.LatestConfigDigest, state.Config.Epoch, err
}

func (c *Transmitter) FromAccount() types.Account {
	return types.Account(c.transmissionSigner.String())
}
