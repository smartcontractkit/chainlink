package solana_test

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/deployment"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestSolanaTokenOps(t *testing.T) {
	quarantine.Flaky(t, "DX-1728")
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(changeset_solana.SolanaContractV0_1_1))
	e := tenv.Env
	solChain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solChain1, true,
		changeset_solana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
			OffRamp:   true,
		})
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			// deployer creates token
			cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
			changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: shared.SPL2022Tokens,
				TokenDecimals:    9,
				TokenSymbol:      "TEST_TOKEN",
			},
		),
	)
	require.NoError(t, err)

	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			// deployer creates token
			cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
			changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: shared.SPLTokens,
				MintPrivateKey:   privKey,
				TokenDecimals:    9,
				TokenSymbol:      "SPL_TEST_TOKEN",
			},
		),
	)
	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(solChain1)
	require.NoError(t, err)
	tokenAddress := solanastateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Type:    shared.SPL2022Tokens,
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet("TEST_TOKEN"),
		},
		addresses,
	)

	deployerKey := e.BlockChains.SolanaChains()[solChain1].DeployerKey.PublicKey()

	testUser, _ := solana.NewRandomPrivateKey()
	testUserPubKey := testUser.PublicKey()

	e, err = commonchangeset.Apply(t, e, commonchangeset.Configure(
		// deployer creates ATA for itself and testUser
		cldf.CreateLegacyChangeSet(changeset_solana.CreateSolanaTokenATA),
		changeset_solana.CreateSolanaTokenATAConfig{
			ChainSelector: solChain1,
			TokenPubkey:   tokenAddress,
			ATAList:       []string{deployerKey.String(), testUserPubKey.String()},
		},
	), commonchangeset.Configure(
		// deployer mints token to itself and testUser
		cldf.CreateLegacyChangeSet(changeset_solana.MintSolanaToken),
		changeset_solana.MintSolanaTokenConfig{
			ChainSelector: solChain1,
			TokenPubkey:   tokenAddress.String(),
			AmountToAddress: map[string]uint64{
				deployerKey.String():    uint64(1000),
				testUserPubKey.String(): uint64(1000),
			},
		},
	))
	require.NoError(t, err)

	testUserATA, _, err := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenAddress, testUserPubKey)
	require.NoError(t, err)
	deployerATA, _, err := solTokenUtil.FindAssociatedTokenAddress(
		solana.Token2022ProgramID,
		tokenAddress,
		e.BlockChains.SolanaChains()[solChain1].DeployerKey.PublicKey(),
	)
	require.NoError(t, err)

	// test if minting was done correctly
	outDec, outVal, err := solTokenUtil.TokenBalance(context.Background(), e.BlockChains.SolanaChains()[solChain1].Client, deployerATA, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))

	outDec, outVal, err = solTokenUtil.TokenBalance(context.Background(), e.BlockChains.SolanaChains()[solChain1].Client, testUserATA, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))

	// now lets do it altogether
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			// deployer creates token
			cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
			changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: shared.SPLTokens,
				TokenDecimals:    9,
				TokenSymbol:      "TEST_TOKEN_2",
				ATAList:          []string{deployerKey.String(), testUserPubKey.String()},
				MintAmountToAddress: map[string]uint64{
					deployerKey.String():    uint64(1000),
					testUserPubKey.String(): uint64(1000),
				},
			},
		),
	)
	require.NoError(t, err)
	addresses, err = e.ExistingAddresses.AddressesForChain(solChain1)
	require.NoError(t, err)
	tokenAddress2 := solanastateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Type:    shared.SPLTokens,
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet("TEST_TOKEN_2"),
		},
		addresses,
	)
	testUserATA2, _, err := solTokenUtil.FindAssociatedTokenAddress(solana.TokenProgramID, tokenAddress2, testUserPubKey)
	require.NoError(t, err)
	deployerATA2, _, err := solTokenUtil.FindAssociatedTokenAddress(
		solana.TokenProgramID,
		tokenAddress2,
		e.BlockChains.SolanaChains()[solChain1].DeployerKey.PublicKey(),
	)
	require.NoError(t, err)
	// test if minting was done correctly
	outDec, outVal, err = solTokenUtil.TokenBalance(context.Background(), e.BlockChains.SolanaChains()[solChain1].Client, deployerATA2, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))

	outDec, outVal, err = solTokenUtil.TokenBalance(context.Background(), e.BlockChains.SolanaChains()[solChain1].Client, testUserATA2, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))
}

func TestDeployLinkToken(t *testing.T) {
	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithSolanaContainer(t, []uint64{selector}, t.TempDir(), map[string]string{}),
	))
	require.NoError(t, err)

	// solana test
	solLinkTokenPrivKey, _ := solana.NewRandomPrivateKey()

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeploySolanaLinkToken), commonchangeset.DeploySolanaLinkTokenConfig{
			ChainSelector: selector,
			TokenPrivKey:  solLinkTokenPrivKey,
			TokenDecimals: 9,
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	require.NotEmpty(t, addrs)
}
