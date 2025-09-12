package testhelpers

import (
	"math/big"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	aptosBind "github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/regulated_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink-aptos/bindings/regulated_token"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/deployment"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func DeployChainContractsToAptosCS(t *testing.T, e DeployedEnv, chainSelector uint64) commoncs.ConfiguredChangeSet {
	ccipConfig := config.DeployAptosChainConfig{
		ContractParamsPerChain: map[uint64]config.ChainContractParams{
			chainSelector: {
				FeeQuoterParams: config.FeeQuoterParams{
					MaxFeeJuelsPerMsg:            new(big.Int).Mul(big.NewInt(100_000_000), big.NewInt(1e18)), // 100M LINK @ 18 decimals
					TokenPriceStalenessThreshold: 24 * 60 * 60,
					FeeTokens:                    []aptos.AccountAddress{aptoscs.MustParseAddress(t, shared.AptosAPTAddress)}, // LINK token will be deployed and added here automatically
					PremiumMultiplierWeiPerEthByFeeToken: map[shared.TokenSymbol]uint64{
						shared.APTSymbol:  11e17,
						shared.LinkSymbol: 9e18,
					},
				},
				OffRampParams: config.OffRampParams{
					ChainSelector:                    chainSelector,
					PermissionlessExecutionThreshold: uint32(globals.PermissionLessExecutionThreshold.Seconds()),
					IsRMNVerificationDisabled:        nil,
					SourceChainSelectors:             nil,
					SourceChainIsEnabled:             nil,
					SourceChainsOnRamp:               nil,
				},
				OnRampParams: config.OnRampParams{
					ChainSelector:  chainSelector,
					AllowlistAdmin: e.Env.BlockChains.AptosChains()[chainSelector].DeployerSigner.AccountAddress(),
					FeeAggregator:  e.Env.BlockChains.AptosChains()[chainSelector].DeployerSigner.AccountAddress(),
				},
			},
		},
		MCMSDeployConfigPerChain: map[uint64]commontypes.MCMSWithTimelockConfigV2{
			chainSelector: {
				Canceller:        proposalutils.SingleGroupMCMSV2(t),
				Proposer:         proposalutils.SingleGroupMCMSV2(t),
				Bypasser:         proposalutils.SingleGroupMCMSV2(t),
				TimelockMinDelay: big.NewInt(1),
			},
		},
		MCMSTimelockConfigPerChain: map[uint64]proposalutils.TimelockConfig{
			chainSelector: {
				MinDelay:     time.Duration(1) * time.Second,
				MCMSAction:   mcmstypes.TimelockActionSchedule,
				OverrideRoot: false,
			},
		},
	}

	return commoncs.Configure(aptoscs.DeployAptosChain{}, ccipConfig)
}

// MakeBCSEVMExtraArgsV2 makes the BCS encoded extra args for a message sent from a Move based chain that is destined for an EVM chain.
// The extra args are used to specify the gas limit and allow out of order flag for the message.
func MakeBCSEVMExtraArgsV2(gasLimit *big.Int, allowOOO bool) []byte {
	s := &bcs.Serializer{}
	s.U256(*gasLimit)
	s.Bool(allowOOO)
	return append(hexutil.MustDecode(GenericExtraArgsV2Tag), s.ToBytes()...)
}

func DeployTransferableTokenAptos(
	t *testing.T,
	lggr logger.Logger,
	e cldf.Environment,
	evmChainSel, aptosChainSel uint64,
	tokenName string,
	mintAmount *config.TokenMint,
) (
	*burn_mint_erc677.BurnMintERC677,
	*burn_mint_token_pool.BurnMintTokenPool,
	aptos.AccountAddress,
	managed_token_pool.ManagedTokenPool,
	error,
) {
	selectorFamily, err := chainsel.GetSelectorFamily(evmChainSel)
	require.NoError(t, err)
	require.Equal(t, chainsel.FamilyEVM, selectorFamily)
	selectorFamily, err = chainsel.GetSelectorFamily(aptosChainSel)
	require.NoError(t, err)
	require.Equal(t, chainsel.FamilyAptos, selectorFamily)

	// EVM
	evmDeployerKey := e.BlockChains.EVMChains()[evmChainSel].DeployerKey
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	evmToken, evmPool, err := deployTransferTokenOneEnd(lggr, e.BlockChains.EVMChains()[evmChainSel], evmDeployerKey, e.ExistingAddresses, tokenName)
	require.NoError(t, err)
	err = attachTokenToTheRegistry(e.BlockChains.EVMChains()[evmChainSel], state.MustGetEVMChainState(evmChainSel), evmDeployerKey, evmToken.Address(), evmPool.Address())
	require.NoError(t, err)

	// Aptos
	e, err = commoncs.Apply(t, e,
		commoncs.Configure(aptoscs.AddTokenPool{},
			config.AddTokenPoolConfig{
				ChainSelector:                       aptosChainSel,
				TokenAddress:                        aptos.AccountAddress{}, // Will be deployed
				TokenCodeObjAddress:                 aptos.AccountAddress{}, // Will be deployed
				TokenPoolAddress:                    aptos.AccountAddress{}, // Will be deployed
				PoolType:                            shared.AptosManagedTokenPoolType,
				TokenTransferFeeByRemoteChainConfig: nil, // TODO - not needed?
				EVMRemoteConfigs: map[uint64]config.EVMRemoteConfig{
					evmChainSel: {
						TokenAddress:     evmToken.Address(),
						TokenPoolAddress: evmPool.Address(),
						RateLimiterConfig: config.RateLimiterConfig{
							RemoteChainSelector: evmChainSel,
							OutboundIsEnabled:   false,
							OutboundCapacity:    0,
							OutboundRate:        0,
							InboundIsEnabled:    false,
							InboundCapacity:     0,
							InboundRate:         0,
						},
					},
				},
				TokenParams: config.TokenParams{
					Name:     tokenName,
					Symbol:   "TKN",
					Decimals: 8,
				},
				TokenMint: mintAmount,
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Second, // TODO
				},
			},
		),
	)
	require.NoError(t, err)

	aptosAddresses, err := e.ExistingAddresses.AddressesForChain(aptosChainSel)
	require.NoError(t, err)
	tokenMetadataAddress := aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    "TKN",
			Version: deployment.Version1_6_0,
			Labels:  nil,
		},
		aptosAddresses,
	)
	lggr.Debugf("Deployed Token on Aptos: %v", tokenMetadataAddress.StringLong())
	tokenPoolAddress := aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    shared.AptosManagedTokenPoolType,
			Version: deployment.Version1_6_0,
			Labels:  cldf.NewLabelSet(tokenMetadataAddress.StringLong()),
		},
		aptosAddresses,
	)
	aptosTokenPool := managed_token_pool.Bind(tokenPoolAddress, e.BlockChains.AptosChains()[aptosChainSel].Client)
	lggr.Debugf("Deployed Token Pool for %v to %v", tokenMetadataAddress.StringLong(), tokenPoolAddress.StringLong())

	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChainSel], evmPool, evmDeployerKey, aptosChainSel, tokenMetadataAddress[:], tokenPoolAddress[:])
	require.NoError(t, err)

	err = grantMintBurnPermissions(lggr, e.BlockChains.EVMChains()[evmChainSel], evmToken, evmDeployerKey, evmPool.Address())
	require.NoError(t, err)

	return evmToken, evmPool, tokenMetadataAddress, aptosTokenPool, nil
}

func DeployRegulatedTransferableTokenAptos(
	t *testing.T,
	lggr logger.Logger,
	e cldf.Environment,
	evmChainSel,
	aptosChainSel uint64,
	tokenName string,
	mintAmount *config.TokenMint,
) (
	*burn_mint_erc677.BurnMintERC677,
	*burn_mint_token_pool.BurnMintTokenPool,
	aptos.AccountAddress,
	regulated_token_pool.RegulatedTokenPool,
	error,
) {
	selectorFamily, err := chainsel.GetSelectorFamily(evmChainSel)
	require.NoError(t, err)
	require.Equal(t, chainsel.FamilyEVM, selectorFamily)
	selectorFamily, err = chainsel.GetSelectorFamily(aptosChainSel)
	require.NoError(t, err)
	require.Equal(t, chainsel.FamilyAptos, selectorFamily)
	// EVM
	evmDeployerKey := e.BlockChains.EVMChains()[evmChainSel].DeployerKey
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	evmToken, evmPool, err := deployTransferTokenOneEnd(lggr, e.BlockChains.EVMChains()[evmChainSel], evmDeployerKey, e.ExistingAddresses, tokenName)
	require.NoError(t, err)
	err = attachTokenToTheRegistry(e.BlockChains.EVMChains()[evmChainSel], state.MustGetEVMChainState(evmChainSel), evmDeployerKey, evmToken.Address(), evmPool.Address())
	require.NoError(t, err)

	// Regulated token must be initialized via EOA, not mcms
	signer := e.BlockChains.AptosChains()[aptosChainSel].DeployerSigner
	client := e.BlockChains.AptosChains()[aptosChainSel].Client
	opts := &aptosBind.TransactOpts{Signer: signer}
	aptosAddresses, err := e.ExistingAddresses.AddressesForChain(aptosChainSel)
	require.NoError(t, err)
	mcmsAddress := aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    shared.AptosMCMSType,
			Version: deployment.Version1_6_0,
		},
		aptosAddresses,
	)
	require.NotEqualf(t, aptos.AccountAddress{}, mcmsAddress, "Aptos mcms address not found")

	// Only admin can grant roles
	adminAddress := signer.AccountAddress()
	tokenAddress, tx, token, err := regulated_token.DeployToObject(signer, client, adminAddress)
	require.NoError(t, err)
	data, err := client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy regulated token: %v", data.VmStatus)

	tx, _, err = regulated_token.DeployMCMSRegistrarToExistingObject(signer, client, tokenAddress, adminAddress, mcmsAddress, true)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy regulated token MCMS registrar: %v", data.VmStatus)

	tx, err = token.RegulatedToken().Initialize(opts, nil, tokenName, "TKN", 8, "", "")
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to initialize regulated token: %v", data.VmStatus)

	if mintAmount != nil {
		lggr.Infof("Minting %v tokens to %v...", mintAmount.Amount, mintAmount.To)
		tx, err = token.RegulatedToken().GrantRole(opts, 4, adminAddress)
		require.NoError(t, err)
		data, err = client.WaitForTransaction(tx.Hash)
		require.NoError(t, err)
		require.True(t, data.Success, "failed to grant mint role to deployer: %v", data.VmStatus)

		tx, err = token.RegulatedToken().Mint(opts, mintAmount.To, mintAmount.Amount)
		require.NoError(t, err)
		data, err = client.WaitForTransaction(tx.Hash)
		require.NoError(t, err)
		require.True(t, data.Success, "failed to mint %d token to %s: %v", mintAmount.Amount, mintAmount.To, data.VmStatus)
	}

	tokenMetadata, err := token.RegulatedToken().TokenMetadata(nil)
	require.NoError(t, err)

	// Save addresses in address book
	typeAndVersion := cldf.NewTypeAndVersion(shared.AptosRegulatedTokenType, deployment.Version1_6_0)
	typeAndVersion.AddLabel("TKN")
	err = e.ExistingAddresses.Save(aptosChainSel, tokenAddress.StringLong(), typeAndVersion)
	require.NoError(t, err)
	typeAndVersion = cldf.NewTypeAndVersion("TKN", deployment.Version1_6_0)
	err = e.ExistingAddresses.Save(aptosChainSel, tokenMetadata.StringLong(), typeAndVersion)
	require.NoError(t, err)

	// Transfer ownership to mcms
	mcmsContract := mcms.Bind(mcmsAddress, client)
	tokenOwnerAddress, err := mcmsContract.MCMSRegistry().GetPreexistingCodeObjectOwnerAddress(nil, tokenAddress)
	require.NoError(t, err)

	tx, err = token.RegulatedToken().TransferOwnership(opts, tokenOwnerAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to propose ownership transfer to mcms %v: %v", tokenOwnerAddress, data.VmStatus)

	tx, err = token.RegulatedToken().TransferAdmin(opts, tokenOwnerAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to propose admin transfer to mcms %v: %v", tokenOwnerAddress, data.VmStatus)

	_, err = commoncs.Apply(t, e,
		commoncs.Configure(aptoscs.AcceptTokenOwnership{},
			config.AcceptTokenOwnershipInput{
				ChainSelector: aptosChainSel,
				Accepts: []config.TokenAcceptInput{
					{
						TokenCodeObjectAddress: tokenAddress,
						TokenType:              shared.AptosRegulatedTokenType,
					},
				},
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Second,
				},
			},
		),
	)
	require.NoError(t, err)

	tx, err = token.RegulatedToken().ExecuteOwnershipTransfer(opts, tokenOwnerAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to execute ownership transfer to mcms %v: %v", tokenOwnerAddress, data.VmStatus)

	// Transfer admin role to mcms
	tx, err = token.RegulatedToken().TransferAdmin(opts, tokenOwnerAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to propose admin transfer to mcms %v: %v", tokenOwnerAddress, data.VmStatus)

	_, err = commoncs.Apply(t, e,
		commoncs.Configure(aptoscs.AcceptTokenAdmin{},
			config.AcceptTokenAdminInput{
				ChainSelector: aptosChainSel,
				Accepts: []config.TokenAcceptInput{
					{
						TokenCodeObjectAddress: tokenAddress,
						TokenType:              shared.AptosRegulatedTokenType,
					},
				},
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Second,
				},
			},
		),
	)
	require.NoError(t, err)

	// Deploy lane
	e, err = commoncs.Apply(t, e,
		commoncs.Configure(aptoscs.AddTokenPool{},
			config.AddTokenPoolConfig{
				ChainSelector:                       aptosChainSel,
				TokenAddress:                        tokenMetadata,
				TokenCodeObjAddress:                 tokenAddress,
				TokenPoolAddress:                    aptos.AccountAddress{},             // Will be deployed
				PoolType:                            shared.AptosRegulatedTokenPoolType, // Use regulated token pool type
				TokenTransferFeeByRemoteChainConfig: nil,                                // TODO - not needed?
				EVMRemoteConfigs: map[uint64]config.EVMRemoteConfig{
					evmChainSel: {
						TokenAddress:     evmToken.Address(),
						TokenPoolAddress: evmPool.Address(),
						RateLimiterConfig: config.RateLimiterConfig{
							RemoteChainSelector: evmChainSel,
							OutboundIsEnabled:   false,
							OutboundCapacity:    0,
							OutboundRate:        0,
							InboundIsEnabled:    false,
							InboundCapacity:     0,
							InboundRate:         0,
						},
					},
				},
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Second, // TODO
				},
			},
		),
	)
	require.NoError(t, err)
	aptosAddresses, err = e.ExistingAddresses.AddressesForChain(aptosChainSel)
	require.NoError(t, err)
	tokenMetadataAddress := aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    "TKN", // Regulated Token symbol
			Version: deployment.Version1_6_0,
			Labels:  nil,
		},
		aptosAddresses,
	)
	lggr.Debugf("Deployed Regulated Token on Aptos: %v", tokenMetadataAddress.StringLong())
	tokenPoolAddress := aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    shared.AptosRegulatedTokenPoolType,
			Version: deployment.Version1_6_0,
			Labels:  cldf.NewLabelSet(tokenMetadataAddress.StringLong()),
		},
		aptosAddresses,
	)
	aptosTokenPool := regulated_token_pool.Bind(tokenPoolAddress, e.BlockChains.AptosChains()[aptosChainSel].Client)
	lggr.Debugf("Deployed Regulated Token Pool for %v to %v", tokenMetadataAddress.StringLong(), tokenPoolAddress.StringLong())
	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChainSel], evmPool, evmDeployerKey, aptosChainSel, tokenMetadataAddress[:], tokenPoolAddress[:])
	require.NoError(t, err)
	err = grantMintBurnPermissions(lggr, e.BlockChains.EVMChains()[evmChainSel], evmToken, evmDeployerKey, evmPool.Address())
	require.NoError(t, err)
	return evmToken, evmPool, tokenMetadataAddress, aptosTokenPool, nil
}
