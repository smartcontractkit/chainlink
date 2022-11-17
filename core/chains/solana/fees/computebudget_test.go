package fees

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetComputeUnitPrice(t *testing.T) {
	key, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	t.Run("noAccount_nofee", func(t *testing.T) {
		// build base tx (no fee)
		tx, err := solana.NewTransaction([]solana.Instruction{
			system.NewTransferInstruction(
				0,
				key.PublicKey(),
				key.PublicKey(),
			).Build(),
		}, solana.Hash{})
		require.NoError(t, err)
		instructionCount := len(tx.Message.Instructions)

		// add fee
		require.NoError(t, SetComputeUnitPrice(tx, 1))

		// evaluate
		currentCount := len(tx.Message.Instructions)
		assert.Greater(t, currentCount, instructionCount)
		assert.Equal(t, 2, currentCount)
		assert.Equal(t, COMPUTE_BUDGET_PROGRAM, tx.Message.AccountKeys[tx.Message.Instructions[0].ProgramIDIndex].String())
		data, err := ComputeUnitPrice(1).Data()
		assert.NoError(t, err)
		assert.Equal(t, data, []byte(tx.Message.Instructions[0].Data))
	})

	t.Run("accountExists_noFee", func(t *testing.T) {
		// build base tx (no fee)
		tx, err := solana.NewTransaction([]solana.Instruction{
			system.NewTransferInstruction(
				0,
				key.PublicKey(),
				key.PublicKey(),
			).Build(),
		}, solana.Hash{})
		require.NoError(t, err)
		accountCount := len(tx.Message.AccountKeys)
		tx.Message.AccountKeys = append(tx.Message.AccountKeys, ComputeUnitPrice(0).ProgramID())
		accountCount++

		// add fee
		require.NoError(t, SetComputeUnitPrice(tx, 1))

		// accounts should not have changed
		assert.Equal(t, accountCount, len(tx.Message.AccountKeys))
		assert.Equal(t, 2, len(tx.Message.Instructions))
		assert.Equal(t, COMPUTE_BUDGET_PROGRAM, tx.Message.AccountKeys[tx.Message.Instructions[0].ProgramIDIndex].String())
		data, err := ComputeUnitPrice(1).Data()
		assert.NoError(t, err)
		assert.Equal(t, data, []byte(tx.Message.Instructions[0].Data))

	})

	// // not a valid test, account must exist for tx to be added
	// t.Run("noAccount_feeExists", func(t *testing.T) {})

	t.Run("exists_notFirst", func(t *testing.T) {
		// build base tx (no fee)
		tx, err := solana.NewTransaction([]solana.Instruction{
			system.NewTransferInstruction(
				0,
				key.PublicKey(),
				key.PublicKey(),
			).Build(),
		}, solana.Hash{})
		require.NoError(t, err)
		transferInstruction := tx.Message.Instructions[0]

		// add fee
		require.NoError(t, SetComputeUnitPrice(tx, 0))

		// swap order of instructions
		tx.Message.Instructions[0], tx.Message.Instructions[1] = tx.Message.Instructions[1], tx.Message.Instructions[0]
		require.Equal(t, transferInstruction, tx.Message.Instructions[0])
		oldFeeInstruction := tx.Message.Instructions[1]
		accountCount := len(tx.Message.AccountKeys)

		// set fee with existing fee instruction
		require.NoError(t, SetComputeUnitPrice(tx, 100))
		require.Equal(t, transferInstruction, tx.Message.Instructions[0]) // transfer should not have been touched
		assert.NotEqual(t, oldFeeInstruction, tx.Message.Instructions[1])
		assert.Equal(t, accountCount, len(tx.Message.AccountKeys))
		assert.Equal(t, 2, len(tx.Message.Instructions)) // instruction count did not change
		data, err := ComputeUnitPrice(100).Data()
		assert.NoError(t, err)
		assert.Equal(t, data, []byte(tx.Message.Instructions[1].Data))
	})

}
