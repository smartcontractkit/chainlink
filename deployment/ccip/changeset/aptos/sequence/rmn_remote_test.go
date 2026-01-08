package sequence

import (
	"context"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	aptosops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/aptos"
)

func TestAptosCurseSequence_Success(t *testing.T) {
	restore := overrideOps()
	t.Cleanup(restore)

	ccipAddr := aptos.AccountAddress{}
	chain := cldf_aptos.Chain{Selector: 101}
	blockChains := cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{chain})

	var subj fastcurse.Subject
	copy(subj[:], []byte("hello"))
	subjects := []fastcurse.Subject{subj}

	var gotDeps dependency.AptosDeps
	var gotInput aptosops.CurseMultipleInput
	stubTx := mcmstypes.Transaction{}

	aptosops.CurseMultipleOp = cldf_ops.NewOperation(
		"stub-curse-op",
		operation.Version1_0_0,
		"stub",
		func(_ cldf_ops.Bundle, deps dependency.AptosDeps, in aptosops.CurseMultipleInput) (mcmstypes.Transaction, error) {
			gotDeps = deps
			gotInput = in
			return stubTx, nil
		},
	)

	bundle := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())

	rep, err := cldf_ops.ExecuteSequence(bundle, AptosCurseSequence, blockChains, AptosCurseUncurseInput{
		CCIPAddress:   ccipAddr,
		ChainSelector: chain.Selector,
		Subjects:      subjects,
	})
	require.NoError(t, err)

	require.Equal(t, chain, gotDeps.AptosChain)
	require.Equal(t, ccipAddr, gotInput.CCIPAddress)
	require.Len(t, gotInput.Subjects, 1)
	require.Equal(t, subjects[0][:], gotInput.Subjects[0])

	require.Len(t, rep.Output.BatchOps, 1)
	require.Equal(t, mcmstypes.ChainSelector(chain.Selector), rep.Output.BatchOps[0].ChainSelector)
	require.Equal(t, []mcmstypes.Transaction{stubTx}, rep.Output.BatchOps[0].Transactions)
}

func TestAptosUncurseSequence_Success(t *testing.T) {
	restore := overrideOps()
	t.Cleanup(restore)

	ccipAddr := aptos.AccountAddress{}
	chain := cldf_aptos.Chain{Selector: 202}
	blockChains := cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{chain})

	var subj fastcurse.Subject
	copy(subj[:], []byte("world"))
	subjects := []fastcurse.Subject{subj}

	var gotDeps dependency.AptosDeps
	var gotInput aptosops.UncurseMultipleInput
	stubTx := mcmstypes.Transaction{}

	aptosops.UncurseMultipleOp = cldf_ops.NewOperation(
		"stub-uncurse-op",
		operation.Version1_0_0,
		"stub",
		func(_ cldf_ops.Bundle, deps dependency.AptosDeps, in aptosops.UncurseMultipleInput) (mcmstypes.Transaction, error) {
			gotDeps = deps
			gotInput = in
			return stubTx, nil
		},
	)

	bundle := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())

	rep, err := cldf_ops.ExecuteSequence(bundle, AptosUncurseSequence, blockChains, AptosCurseUncurseInput{
		CCIPAddress:   ccipAddr,
		ChainSelector: chain.Selector,
		Subjects:      subjects,
	})
	require.NoError(t, err)

	require.Equal(t, chain, gotDeps.AptosChain)
	require.Equal(t, ccipAddr, gotInput.CCIPAddress)
	require.Len(t, gotInput.Subjects, 1)
	require.Equal(t, subjects[0][:], gotInput.Subjects[0])

	require.Len(t, rep.Output.BatchOps, 1)
	require.Equal(t, mcmstypes.ChainSelector(chain.Selector), rep.Output.BatchOps[0].ChainSelector)
	require.Equal(t, []mcmstypes.Transaction{stubTx}, rep.Output.BatchOps[0].Transactions)
}

func TestAptosCurseSequence_ChainNotFound(t *testing.T) {
	restore := overrideOps()
	t.Cleanup(restore)

	bundle := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())
	_, err := cldf_ops.ExecuteSequence(bundle, AptosCurseSequence, cldf_chain.NewBlockChainsFromSlice(nil), AptosCurseUncurseInput{
		ChainSelector: 999,
	})
	require.Error(t, err)
}

func TestAptosUncurseSequence_ChainNotFound(t *testing.T) {
	restore := overrideOps()
	t.Cleanup(restore)

	bundle := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())
	_, err := cldf_ops.ExecuteSequence(bundle, AptosUncurseSequence, cldf_chain.NewBlockChainsFromSlice(nil), AptosCurseUncurseInput{
		ChainSelector: 888,
	})
	require.Error(t, err)
}

func overrideOps() func() {
	origCurse := aptosops.CurseMultipleOp
	origUncurse := aptosops.UncurseMultipleOp
	return func() {
		aptosops.CurseMultipleOp = origCurse
		aptosops.UncurseMultipleOp = origUncurse
	}
}
