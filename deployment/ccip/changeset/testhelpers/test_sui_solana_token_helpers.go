package testhelpers

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	tokenpoolops "github.com/smartcontractkit/chainlink-ccip/chains/solana/deployment/v1_6_0/operations/token_pools"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeSetSolanaV0_1_1 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// solSuiBnMTokenSymbol is the symbol/label used for the Solana burn-mint token deployed by
// DeployAndCrossRegisterBurnMintPoolSolanaSui. It doubles as the AddressRef label used to look the
// mint back up, so it must be unique enough not to collide with other Solana tokens in the env.
const solSuiBnMTokenSymbol = "SOLSUIBNT"

// DeployAndCrossRegisterBurnMintPoolSolanaSui deploys a burn-mint token pool pair for a
// Solana<->Sui lane and cross-registers each side to the other. It exists because the legacy
// helpers are unusable for this pair:
//   - HandleTokenAndBurnMintTokenPoolDeploymentForSUI (test_sui_helpers.go:737) is EVM-dest
//     hardcoded and deploys an EVM token+pool leg.
//   - DeployTransferableTokenSolanaV0_1_1 (test_helpers_solana_v0_1_1.go:115) requires an EVM leg.
//   - The legacy Solana cross-reg changeset SetupTokenPoolForRemoteChain keys on EVMRemoteConfigs
//     and hard-rejects non-EVM remotes.
//
// Wiring is hybrid (proven changesets + one direct op), mirroring the EVM<->Sui pair but with the
// Solana leg swapped in:
//
//  1. Solana token:  ccipChangeSetSolanaV0_1_1.DeploySolanaToken (SPL-2022, 9 decimals, ATA + mint
//     to the Solana deployer).
//  2. Solana pool:   ccipChangeSetSolanaV0_1_1.E2ETokenPool (deploy program + init + register
//     token-admin-registry + accept + setPool). Its RemoteChainTokenPool section is intentionally
//     omitted (EVM-only); the Sui remote is added in step 4.
//  3. Sui pool:      sui_cs.DeployTPAndConfigure for the existing Sui LINK token, cross-registered
//     to the Solana remote. The Sui remote-pool address is the Solana pool CONFIG PDA (Solana pools
//     are identified to remotes by their pool-config PDA, see DeployTransferableTokenSolanaV0_1_1
//     line 292) and the remote-token address is the Solana mint pubkey; both passed as hex strings
//     (StrToBytes / StrTo32 inside the Sui op).
//  4. Solana->Sui cross-reg: the v1_6_0 token_pools.UpsertRemoteChainConfigBurnMint op, executed
//     directly with raw []byte remotes (no chain-family check; on-chain RemoteAddress.Address is a
//     Vec<u8> so a 32-byte Sui object id fits). The Sui remote-token is the Sui LINK CoinMetadata
//     object id and the remote-pool is the Sui BnM pool package id, matching what
//     HandleTokenAndBurnMintTokenPoolDeploymentForSUI stores on the EVM counterparty (lines 827/837).
//
// Rate limiters are disabled on both sides. The env (e.Env) is updated in place with every
// changeset's output. Returns the Solana mint + pool pubkeys and the Sui LINK package id (the
// package id is what WaitForTokenBalanceSui keys its coin-type filter on:
// 0x2::coin::Coin<<suiLinkPkgID>::link::LINK>).
//
// Runtime unknowns to confirm at nix (flagged in the plan):
//   - The Sui object id that the Solana pool must store as RemoteTokenAddress so the Sui OffRamp's
//     token_admin_registry::get_pool(dest_token_address) resolves. We use LinkTokenCoinMetadataId
//     by analogy with the EVM<->Sui helper; a wrong value => Sui OffRamp EUnsupportedToken at exec.
func DeployAndCrossRegisterBurnMintPoolSolanaSui(
	t *testing.T,
	e *DeployedEnv,
	solChainSel, suiChainSel uint64,
) (solTokenMint, solPool solana.PublicKey, suiLinkPkgID string, err error) {
	t.Helper()
	env := e.Env

	// --- 1. Solana token -----------------------------------------------------------
	solDeployerKey := env.BlockChains.SolanaChains()[solChainSel].DeployerKey.PublicKey()
	env, err = commoncs.Apply(nil, env,
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.DeploySolanaToken),
			ccipChangeSetSolanaV0_1_1.DeploySolanaTokenConfig{
				ChainSelector:    solChainSel,
				TokenProgramName: shared.SPL2022Tokens,
				TokenDecimals:    9,
				TokenSymbol:      solSuiBnMTokenSymbol,
				ATAList:          []string{solDeployerKey.String()},
				MintAmountToAddress: map[string]uint64{
					solDeployerKey.String(): uint64(1000e9),
				},
			},
		),
	)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("deploy solana token: %w", err)
	}

	solAddresses, err := env.ExistingAddresses.AddressesForChain(solChainSel)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("read solana addresses: %w", err)
	}
	solTokenMint = solanastateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Type:    shared.SPL2022Tokens,
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet(solSuiBnMTokenSymbol),
		},
		solAddresses,
	)
	bnm := shared.BurnMintTokenPool

	// --- 2. Solana burn-mint pool (deploy + init + register + accept + setPool) ----
	// RemoteChainTokenPool omitted: the legacy Solana cross-reg is EVM-only. The Sui remote is
	// wired in step 4 via the v1_6_0 UpsertRemoteChainConfigBurnMint op, which accepts raw []byte.
	env, err = commoncs.Apply(nil, env,
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.E2ETokenPool),
			ccipChangeSetSolanaV0_1_1.E2ETokenPoolConfig{
				InitializeGlobalTokenPoolConfig: []ccipChangeSetSolanaV0_1_1.TokenPoolConfigWithMCM{
					{
						ChainSelector: solChainSel,
						TokenPoolConfigs: []ccipChangeSetSolanaV0_1_1.TokenPoolConfig{
							{TokenPubKey: solTokenMint, PoolType: bnm, Metadata: shared.CLLMetadata},
						},
					},
				},
				AddTokenPoolAndLookupTable: []ccipChangeSetSolanaV0_1_1.AddTokenPoolAndLookupTableConfig{
					{
						ChainSelector: solChainSel,
						TokenPoolConfigs: []ccipChangeSetSolanaV0_1_1.TokenPoolConfig{
							{TokenPubKey: solTokenMint, PoolType: bnm, Metadata: shared.CLLMetadata},
						},
					},
				},
				RegisterTokenAdminRegistry: []ccipChangeSetSolanaV0_1_1.RegisterTokenAdminRegistryConfig{
					{
						ChainSelector: solChainSel,
						RegisterTokenConfigs: []ccipChangeSetSolanaV0_1_1.RegisterTokenConfig{
							{
								TokenPubKey:             solTokenMint,
								TokenAdminRegistryAdmin: solDeployerKey,
								RegisterType:            ccipChangeSetSolanaV0_1_1.ViaGetCcipAdminInstruction,
							},
						},
					},
				},
				AcceptAdminRoleTokenAdminRegistry: []ccipChangeSetSolanaV0_1_1.AcceptAdminRoleTokenAdminRegistryConfig{
					{
						ChainSelector: solChainSel,
						AcceptAdminRoleTokenConfigs: []ccipChangeSetSolanaV0_1_1.AcceptAdminRoleTokenConfig{
							{TokenPubKey: solTokenMint},
						},
					},
				},
				SetPool: []ccipChangeSetSolanaV0_1_1.SetPoolConfig{
					{
						ChainSelector: solChainSel,
						SetPoolTokenConfigs: []ccipChangeSetSolanaV0_1_1.SetPoolTokenConfig{
							{TokenPubKey: solTokenMint, PoolType: bnm, Metadata: shared.CLLMetadata},
						},
					},
				},
			},
		),
	)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("deploy/configure solana burn-mint pool: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("load onchain state after solana pool: %w", err)
	}
	solPool, ok := state.SolChains[solChainSel].BurnMintTokenPools[shared.CLLMetadata]
	if !ok {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("solana BurnMintTokenPool not found under metadata %s", shared.CLLMetadata)
	}
	// Solana pools are identified to remotes by their pool-config PDA (see
	// DeployTransferableTokenSolanaV0_1_1 line 292).
	solPoolConfigPDA, err := soltokens.TokenPoolConfigAddress(solTokenMint, solPool)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("derive solana pool config PDA: %w", err)
	}

	// --- 3. Sui burn-mint pool (LINK), cross-registered to the Solana remote --------
	suiState := state.SuiChains[suiChainSel]
	suiLinkPkgID = suiState.LinkTokenAddress
	linkTokenMetadataID := suiState.LinkTokenCoinMetadataId
	linkTokenTreasuryCapID := suiState.LinkTokenTreasuryCapId

	env, _, err = commoncs.ApplyChangesets(t, env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployTPAndConfigure{}, sui_cs.DeployTPAndConfigureConfig{
			SuiChainSelector: suiChainSel,
			TokenPoolTypes:   []sui_deployment.TokenPoolType{sui_deployment.TokenPoolTypeBurnMint},
			BurnMintTpInput: burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
				CoinObjectTypeArg:    suiLinkPkgID + "::link::LINK",
				CoinMetadataObjectId: linkTokenMetadataID,
				TreasuryCapObjectId:  linkTokenTreasuryCapID,

				// Apply dest chain updates: register the Solana remote.
				RemoteChainSelectorsToRemove: []uint64{},
				RemoteChainSelectorsToAdd:    []uint64{solChainSel},
				// Pool = Solana pool-config PDA (StrToBytes, raw); Token = Solana mint (StrTo32, 32B).
				RemotePoolAddressesToAdd: [][]string{{fmt.Sprintf("0x%x", solPoolConfigPDA.Bytes())}},
				RemoteTokenAddressesToAdd: []string{
					fmt.Sprintf("0x%x", solTokenMint.Bytes()),
				},

				// Rate limiter config for the Solana remote (disabled).
				RemoteChainSelectors: []uint64{solChainSel},
				OutboundIsEnableds:   []bool{false},
				OutboundCapacities:   []uint64{0},
				OutboundRates:        []uint64{0},
				InboundIsEnableds:    []bool{false},
				InboundCapacities:    []uint64{0},
				InboundRates:         []uint64{0},
			},
		}),
	})
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("deploy/configure sui burn-mint pool: %w", err)
	}

	// Reload to read the freshly deployed Sui BnM pool package id.
	state, err = stateview.LoadOnchainState(env)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("load onchain state after sui pool: %w", err)
	}
	suiBnMPool, ok := state.SuiChains[suiChainSel].BnMTokenPools[TokenSymbolLINK]
	if !ok {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("sui BurnMintTokenPool not found for token %s", TokenSymbolLINK)
	}

	// --- 4. Solana -> Sui cross-registration via the v1_6_0 op (raw []byte remote) --
	suiRemoteTokenBytes, err := hex.DecodeString(strings.TrimPrefix(linkTokenMetadataID, "0x"))
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("decode sui link metadata id: %w", err)
	}
	suiRemotePoolBytes, err := hex.DecodeString(strings.TrimPrefix(suiBnMPool.PackageID, "0x"))
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("decode sui bnm pool package id: %w", err)
	}

	solChain := env.BlockChains.SolanaChains()[solChainSel]
	if _, err = operations.ExecuteOperation(env.OperationsBundle, tokenpoolops.UpsertRemoteChainConfigBurnMint, solChain, tokenpoolops.RemoteChainConfig{
		TokenPool:          solPool,
		TokenMint:          solTokenMint,
		TokenProgramID:     solana.Token2022ProgramID,
		RemoteSelector:     suiChainSel,
		RemoteTokenAddress: suiRemoteTokenBytes,
		RemotePoolAddress:  suiRemotePoolBytes,
		RemoteDecimals:     9,
		// UpsertRemoteChainConfigBurnMint hardcodes disabled rate limits when initializing a new
		// remote, so the Inbound/Outbound rate-limiter fields are intentionally left zero.
	}); err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, "", fmt.Errorf("upsert solana remote chain config for sui: %w", err)
	}

	e.Env = env
	require.NoError(t, err)
	return solTokenMint, solPool, suiLinkPkgID, nil
}
