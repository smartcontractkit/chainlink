package solana

import (
	"context"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink/deployment"
)

// https://github.com/anza-xyz/agave/blob/master/programs/address-lookup-table/src/processor.rs
// https://github.com/anza-xyz/agave/blob/489f483e1d7b30ef114e0123994818b2accfa389/sdk/program/src/address_lookup_table/instruction.rs
const (
	InstructionCreateLookupTable uint32 = iota
	InstructionFreezeLookupTable
	InstructionExtendLookupTable
	InstructionDeactivateLookupTable
	InstructionCloseLookupTable
)

func NewCreateLookupTableInstruction(
	programId,
	authority,
	funder solana.PublicKey,
	slot uint64,
) (solana.PublicKey, solana.Instruction, error) {
	// https://github.com/solana-labs/solana-web3.js/blob/c1c98715b0c7900ce37c59bffd2056fa0037213d/src/programs/address-lookup-table/index.ts#L274
	slotLE := make([]byte, 8)
	binary.LittleEndian.PutUint64(slotLE, slot)
	account, bumpSeed, err := solana.FindProgramAddress([][]byte{authority.Bytes(), slotLE}, programId)
	if err != nil {
		return solana.PublicKey{}, nil, err
	}

	data := binary.LittleEndian.AppendUint32([]byte{}, InstructionCreateLookupTable)
	data = binary.LittleEndian.AppendUint64(data, slot)
	data = append(data, bumpSeed)
	return account, solana.NewInstruction(
		programId,
		solana.AccountMetaSlice{
			solana.Meta(account).WRITE(),
			solana.Meta(authority).SIGNER(),
			solana.Meta(funder).SIGNER().WRITE(),
			solana.Meta(solana.SystemProgramID),
		},
		data,
	), nil
}

func NewExtendLookupTableInstruction(
	programID, table, authority, funder solana.PublicKey,
	accounts []solana.PublicKey,
) solana.Instruction {
	// https://github.com/solana-labs/solana-web3.js/blob/c1c98715b0c7900ce37c59bffd2056fa0037213d/src/programs/address-lookup-table/index.ts#L113

	data := binary.LittleEndian.AppendUint32([]byte{}, InstructionExtendLookupTable)
	data = binary.LittleEndian.AppendUint64(data, uint64(len(accounts))) // note: this is usually u32 + 8 byte buffer
	for _, a := range accounts {
		data = append(data, a.Bytes()...)
	}

	return solana.NewInstruction(
		programID,
		solana.AccountMetaSlice{
			solana.Meta(table).WRITE(),
			solana.Meta(authority).SIGNER(),
			solana.Meta(funder).SIGNER().WRITE(),
			solana.Meta(solana.SystemProgramID),
		},
		data,
	)
}

func CreateLookupTable(ctx context.Context, chain deployment.SolChain, programID solana.PublicKey) (solana.PublicKey, error) {
	slot, err := chain.GetSlot(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.PublicKey{}, err
	}

	table, instruction, ierr := NewCreateLookupTableInstruction(
		programID,
		chain.DeployerKey.PublicKey(),
		chain.DeployerKey.PublicKey(),
		slot-1, // Using the most recent slot sometimes results in errors when submitting the transaction.
	)
	if ierr != nil {
		return solana.PublicKey{}, ierr
	}

	chain.Confirm([]solana.Instruction{instruction})

	return table, nil
}

func ExtendLookupTable(ctx context.Context, chain deployment.SolChain, programID, table solana.PublicKey, entries []solana.PublicKey) {
	chain.Confirm([]solana.Instruction{
		NewExtendLookupTableInstruction(
			programID,
			table,
			chain.DeployerKey.PublicKey(),
			chain.DeployerKey.PublicKey(),
			entries,
		),
	})
}

func SetupLookupTable(ctx context.Context, chain deployment.SolChain, programID solana.PublicKey, entries []solana.PublicKey) (solana.PublicKey, error) {
	table, err := CreateLookupTable(ctx, chain, programID)
	if err != nil {
		return solana.PublicKey{}, err
	}

	ExtendLookupTable(ctx, chain, programID, table, entries)

	// Address lookup tables have to "warm up" for at least 1 slot before they can be used.
	// So, we wait for a new slot to be produced before returning the table, so it's available
	// and warmed up as soon as this method returns it.
	err = chain.AwaitSlotChange(ctx)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return table, nil
}
