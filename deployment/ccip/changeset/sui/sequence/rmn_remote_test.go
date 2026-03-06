package sequence

import (
	"context"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	rmn_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
)

func TestSuiCurseSequence_Success(t *testing.T) {
	restore := overrideSuiOps()
	t.Cleanup(restore)

	chain := cldf_sui.Chain{ChainMetadata: cldf_sui.ChainMetadata{Selector: 101}}
	blockChains := cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{chain})

	var subj fastcurse.Subject
	copy(subj[:], []byte("hello"))
	subjects := []fastcurse.Subject{subj}

	var gotInput rmn_ops.CurseUncurseChainInput
	stubCall := sui_ops.TransactionCall{
		PackageID:  "0xpkg",
		Module:     "rmn_remote",
		Function:   "curse_multiple",
		Data:       []byte("stub"),
		StateObjID: "0xstate",
		TypeArgs:   []string{},
	}

	rmn_ops.CurseChainOp = cldf_ops.NewOperation(
		"stub-curse-op",
		semver.MustParse("1.0.0"),
		"stub",
		func(_ cldf_ops.Bundle, _ sui_ops.OpTxDeps, in rmn_ops.CurseUncurseChainInput) (sui_ops.OpTxResult[rmn_ops.NoObjects], error) {
			gotInput = in
			return sui_ops.OpTxResult[rmn_ops.NoObjects]{
				Call: stubCall,
			}, nil
		},
	)

	bundle := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())

	rep, err := cldf_ops.ExecuteSequence(bundle, SuiCurseSequence, blockChains, SuiCurseUncurseInput{
		CCIPAddress:          "0xccip",
		CCIPObjectRef:        "0xref",
		CCIPOwnerCapObjectID: "0xcap",
		ChainSelector:        101,
		Subjects:             subjects,
	})
	require.NoError(t, err)

	require.Equal(t, "0xccip", gotInput.CCIPPackageId)
	require.Equal(t, "0xref", gotInput.StateObjectId)
	require.Equal(t, "0xcap", gotInput.OwnerCapObjectId)
	require.Len(t, gotInput.Subjects, 1)
	require.Equal(t, subjects[0][:], gotInput.Subjects[0])

	require.Len(t, rep.Output.BatchOps, 1)
	require.Equal(t, mcmstypes.ChainSelector(101), rep.Output.BatchOps[0].ChainSelector)
	require.Len(t, rep.Output.BatchOps[0].Transactions, 1)
}

func TestSuiUncurseSequence_Success(t *testing.T) {
	restore := overrideSuiOps()
	t.Cleanup(restore)

	chain := cldf_sui.Chain{ChainMetadata: cldf_sui.ChainMetadata{Selector: 202}}
	blockChains := cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{chain})

	var subj fastcurse.Subject
	copy(subj[:], []byte("world"))
	subjects := []fastcurse.Subject{subj}

	var gotInput rmn_ops.CurseUncurseChainInput
	stubCall := sui_ops.TransactionCall{
		PackageID:  "0xpkg",
		Module:     "rmn_remote",
		Function:   "uncurse_multiple",
		Data:       []byte("stub"),
		StateObjID: "0xstate",
		TypeArgs:   []string{},
	}

	rmn_ops.UncurseChainOp = cldf_ops.NewOperation(
		"stub-uncurse-op",
		semver.MustParse("1.0.0"),
		"stub",
		func(_ cldf_ops.Bundle, _ sui_ops.OpTxDeps, in rmn_ops.CurseUncurseChainInput) (sui_ops.OpTxResult[rmn_ops.NoObjects], error) {
			gotInput = in
			return sui_ops.OpTxResult[rmn_ops.NoObjects]{
				Call: stubCall,
			}, nil
		},
	)

	bundle := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())

	rep, err := cldf_ops.ExecuteSequence(bundle, SuiUncurseSequence, blockChains, SuiCurseUncurseInput{
		CCIPAddress:          "0xccip",
		CCIPObjectRef:        "0xref",
		CCIPOwnerCapObjectID: "0xcap",
		ChainSelector:        202,
		Subjects:             subjects,
	})
	require.NoError(t, err)

	require.Equal(t, "0xccip", gotInput.CCIPPackageId)
	require.Equal(t, "0xref", gotInput.StateObjectId)
	require.Equal(t, "0xcap", gotInput.OwnerCapObjectId)
	require.Len(t, gotInput.Subjects, 1)
	require.Equal(t, subjects[0][:], gotInput.Subjects[0])

	require.Len(t, rep.Output.BatchOps, 1)
	require.Equal(t, mcmstypes.ChainSelector(202), rep.Output.BatchOps[0].ChainSelector)
	require.Len(t, rep.Output.BatchOps[0].Transactions, 1)
}

func TestSuiCurseSequence_ChainNotFound(t *testing.T) {
	restore := overrideSuiOps()
	t.Cleanup(restore)

	bundle := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())
	_, err := cldf_ops.ExecuteSequence(bundle, SuiCurseSequence, cldf_chain.NewBlockChainsFromSlice(nil), SuiCurseUncurseInput{
		ChainSelector: 999,
	})
	require.Error(t, err)
}

func TestSuiUncurseSequence_ChainNotFound(t *testing.T) {
	restore := overrideSuiOps()
	t.Cleanup(restore)

	bundle := cldf_ops.NewBundle(context.Background, logger.Nop(), cldf_ops.NewMemoryReporter())
	_, err := cldf_ops.ExecuteSequence(bundle, SuiUncurseSequence, cldf_chain.NewBlockChainsFromSlice(nil), SuiCurseUncurseInput{
		ChainSelector: 888,
	})
	require.Error(t, err)
}

func overrideSuiOps() func() {
	origCurse := rmn_ops.CurseChainOp
	origUncurse := rmn_ops.UncurseChainOp
	return func() {
		rmn_ops.CurseChainOp = origCurse
		rmn_ops.UncurseChainOp = origUncurse
	}
}
