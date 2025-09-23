package testhelpers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	aptosBind "github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_dummy_receiver"
	module_onramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_onramp/onramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	aptos_burn_mint_token_pool "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/regulated_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/helpers"
	"github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink-aptos/bindings/regulated_token"
	module_regulated_token "github.com/smartcontractkit/chainlink-aptos/bindings/regulated_token/regulated_token"
	"github.com/smartcontractkit/chainlink-aptos/bindings/test_token/bnm_registrar"
	"github.com/smartcontractkit/chainlink-aptos/bindings/test_token/lnr_registrar"
	"github.com/smartcontractkit/chainlink-aptos/bindings/test_token/test_token"
	"github.com/smartcontractkit/chainlink-aptos/relayer/codec"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/deployment"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
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

// Aptos doesn't provide any struct that we could reuse here

type AptosSendRequest struct {
	Receiver      []byte
	Data          []byte
	ExtraArgs     []byte
	FeeToken      aptos.AccountAddress
	FeeTokenStore aptos.AccountAddress
	TokenAmounts  []AptosTokenAmount
}

type AptosTokenAmount struct {
	Token  aptos.AccountAddress
	Amount uint64
}

func SendRequestAptos(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *ccipclient.CCIPSendReqConfig,
) (*ccipclient.AnyMsgSentEvent, error) {
	sender := e.BlockChains.AptosChains()[cfg.SourceChain].DeployerSigner
	senderAddress := sender.AccountAddress()
	client := e.BlockChains.AptosChains()[cfg.SourceChain].Client

	e.Logger.Infof("(Aptos) Sending CCIP request from chain selector %d to chain selector %d using sender %s",
		cfg.SourceChain, cfg.DestChain, senderAddress.StringLong())

	msg := cfg.Message.(AptosSendRequest)
	router := state.AptosChains[cfg.SourceChain].CCIPAddress
	if cfg.IsTestRouter {
		router = state.AptosChains[cfg.DestChain].TestRouterAddress
	}

	tokenAddresses := make([]aptos.AccountAddress, len(msg.TokenAmounts))
	tokenAmounts := make([]uint64, len(msg.TokenAmounts))
	tokenStoreAddresses := make([]aptos.AccountAddress, len(msg.TokenAmounts))
	for i, v := range msg.TokenAmounts {
		tokenAddresses[i] = v.Token
		tokenAmounts[i] = v.Amount
		tokenStoreAddresses[i] = aptos.AccountAddress{}
	}

	// Debug information
	var (
		tokenAddressStrings []string
		tokenStoreStrings   []string
	)
	feeTokenBalance, err := helpers.GetFungibleAssetBalance(client, senderAddress, msg.FeeToken)
	if err != nil {
		return nil, err
	}
	e.Logger.Debugw("Fungible Asset balance", "feeToken", feeTokenBalance)
	for _, address := range tokenAddresses {
		tokenAddressStrings = append(tokenAddressStrings, address.StringLong())
		transferTokenBalance, err := helpers.GetFungibleAssetBalance(client, senderAddress, address)
		if err != nil {
			return nil, err
		}
		e.Logger.Debugw("Fungible Asset balance", "transferToken", transferTokenBalance)
	}
	for _, address := range tokenStoreAddresses {
		tokenStoreStrings = append(tokenStoreStrings, address.StringLong())
	}
	e.Logger.Debugw("(Aptos) Sending message: ",
		"destChainSelector", cfg.DestChain,
		"routerAddress", router.StringLong(),
		"receiver", hex.EncodeToString(msg.Receiver),
		"data", hex.EncodeToString(msg.Data),
		"tokenAddresses", tokenAddressStrings,
		"tokenAmounts", tokenAmounts,
		"tokenStoreAddresses", tokenStoreStrings,
		"feeToken", msg.FeeToken.StringLong(),
		"feeTokenStore", msg.FeeTokenStore.StringLong(),
		"extraArgs", hex.EncodeToString(msg.ExtraArgs),
	)

	routerContract := ccip_router.Bind(router, client)
	fee, err := routerContract.Router().GetFee(
		nil,
		cfg.DestChain,
		msg.Receiver,
		msg.Data,
		tokenAddresses,
		tokenAmounts,
		tokenStoreAddresses,
		msg.FeeToken,
		msg.FeeTokenStore,
		msg.ExtraArgs,
	)
	if err != nil {
		e.Logger.Errorf("Estimating fee: %v", err)
	}
	e.Logger.Infof("Estimated fee: %v", fee)

	opts := &aptosBind.TransactOpts{
		Signer: sender,
	}
	tx, err := routerContract.Router().CCIPSend(
		opts,
		cfg.DestChain,
		msg.Receiver,
		msg.Data,
		tokenAddresses,
		tokenAmounts,
		tokenStoreAddresses,
		msg.FeeToken,
		msg.FeeTokenStore,
		msg.ExtraArgs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send CCIP message: %w", err)
	}
	data, err := client.WaitForTransaction(tx.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}
	if !data.Success {
		return nil, fmt.Errorf("transaction reverted: %v", data.VmStatus)
	}
	e.Logger.Infof("(Aptos) CCIP message sent (tx %s) from chain selector %d to chain selector %d", tx.Hash, cfg.SourceChain, cfg.DestChain)

	for _, event := range data.Events {
		e.Logger.Debugf("(Aptos) Message contains event type: %v", event.Type)
		// The RPC strips all leading zeroes from the event type
		if event.Type == fmt.Sprintf("0x%s::onramp::CCIPMessageSent", strings.TrimLeft(strings.TrimPrefix(router.String(), "0x"), "0")) {
			var msgSentEvent module_onramp.CCIPMessageSent
			if err := codec.DecodeAptosJsonValue(event.Data, &msgSentEvent); err != nil {
				return nil, fmt.Errorf("failed to decode CCIPMessageSentEvent: %w", err)
			}
			e.Logger.Debugf("CCIPMessageSentEvent: %v", msgSentEvent)
			return &ccipclient.AnyMsgSentEvent{
				SequenceNumber: msgSentEvent.SequenceNumber,
				RawEvent:       msgSentEvent,
			}, nil
		}
	}
	return nil, errors.New("sent message but didn't receive CCIPMessageSent event")
}

// DeployTransferableTokenAptos deploys two tokens onto the EVM and Aptos chain respectively, setting up a lane between them.
// For the aptos token the managed_token package will be used, alongside the managed_token_pool package for the token pool
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

// DeployRegulatedTransferableTokenAptos deploys two tokens onto the EVM and Aptos chain and sets up a lane between them
// For Aptos, the regulated_token will be used along with the regulated_token_pool token pool.
// Since the regulated_token must be initialized from an EOA, not mcms, it will be deployed from the deployer account
// and then transferred over to mcms
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
	// helper function to wait for a transaction to be mined
	assertTxSuccess := func(err error, tx *api.PendingTransaction, msg string, args ...any) {
		require.NoError(t, err)
		data, err := client.WaitForTransaction(tx.Hash)
		require.NoError(t, err)
		require.True(t, data.Success, "%s: %s", fmt.Sprintf(msg, args...), data.VmStatus)
	}
	mcmsAddress := aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    shared.AptosMCMSType,
			Version: deployment.Version1_6_0,
		},
		aptosAddresses,
	)
	require.NotEqualf(t, aptos.AccountAddress{}, mcmsAddress, "Aptos mcms address not found")

	// Deploy the token and token registrar, setting the deployer as the administrator
	adminAddress := signer.AccountAddress()
	tokenAddress, tx, token, err := regulated_token.DeployToObject(signer, client, adminAddress)
	assertTxSuccess(err, tx, "failed to deploy regulated token")
	tx, _, err = regulated_token.DeployMCMSRegistrarToExistingObject(signer, client, tokenAddress, adminAddress, mcmsAddress, true)
	assertTxSuccess(err, tx, "failed to deploy regulated token MCMS registrar")

	// Initialize the token
	tx, err = token.RegulatedToken().Initialize(opts, nil, tokenName, "TKN", 8, "", "")
	assertTxSuccess(err, tx, "failed to initialize regulated token")

	// If requested, set the deployer as an allowed minter and mint the requested tokens
	if mintAmount != nil {
		lggr.Infof("Minting %v tokens to %v...", mintAmount.Amount, mintAmount.To)
		tx, err = token.RegulatedToken().GrantRole(opts, module_regulated_token.MINTER_ROLE, adminAddress)
		assertTxSuccess(err, tx, "failed to grant mint role to deployer")
		tx, err = token.RegulatedToken().Mint(opts, mintAmount.To, mintAmount.Amount)
		assertTxSuccess(err, tx, "failed to mint %d token to %s", mintAmount.Amount, mintAmount.To)
	}

	// Save token addresses in address book
	tokenMetadata, err := token.RegulatedToken().TokenMetadata(nil)
	require.NoError(t, err)
	typeAndVersion := cldf.NewTypeAndVersion(shared.AptosRegulatedTokenType, deployment.Version1_6_0)
	typeAndVersion.AddLabel("TKN")
	err = e.ExistingAddresses.Save(aptosChainSel, tokenAddress.StringLong(), typeAndVersion)
	require.NoError(t, err)
	typeAndVersion = cldf.NewTypeAndVersion("TKN", deployment.Version1_6_0)
	err = e.ExistingAddresses.Save(aptosChainSel, tokenMetadata.StringLong(), typeAndVersion)
	require.NoError(t, err)

	// Transfer token ownership to mcms
	mcmsContract := mcms.Bind(mcmsAddress, client)
	tokenOwnerAddress, err := mcmsContract.MCMSRegistry().GetPreexistingCodeObjectOwnerAddress(nil, tokenAddress)
	require.NoError(t, err)
	tx, err = token.RegulatedToken().TransferOwnership(opts, tokenOwnerAddress)
	assertTxSuccess(err, tx, "failed to propose ownership transfer to mcms %v", tokenOwnerAddress)
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
	assertTxSuccess(err, tx, "failed to execute ownership transfer to mcms %v", tokenOwnerAddress)

	// Transfer admin role to mcms
	tx, err = token.RegulatedToken().TransferAdmin(opts, tokenOwnerAddress)
	assertTxSuccess(err, tx, "failed to propose admin transfer to mcms %v", tokenOwnerAddress)
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
				TokenTransferFeeByRemoteChainConfig: nil,
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
					MinDelay: time.Second,
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

// DeployAptosCCIPReceiver deploys the ccip_dummy_receiver package to all Aptos chains, saving the resulting address in the address book for future use
func DeployAptosCCIPReceiver(t *testing.T, e cldf.Environment) {
	state, err := aptosstate.LoadOnchainStateAptos(e)
	require.NoError(t, err)
	for selector, onchainState := range state {
		addr, tx, _, err := ccip_dummy_receiver.DeployToObject(e.BlockChains.AptosChains()[selector].DeployerSigner, e.BlockChains.AptosChains()[selector].Client, onchainState.CCIPAddress, onchainState.MCMSAddress)
		require.NoError(t, err)
		t.Logf("(Aptos) CCIPDummyReceiver(ccip: %s, mcms: %s) deployed to %s in tx %s", onchainState.CCIPAddress.StringLong(), onchainState.MCMSAddress.StringLong(), addr.StringLong(), tx.Hash)
		require.NoError(t, e.BlockChains.AptosChains()[selector].Confirm(tx.Hash))
		err = e.ExistingAddresses.Save(selector, addr.StringLong(), cldf.NewTypeAndVersion(shared.AptosReceiverType, deployment.Version1_0_0))
		require.NoError(t, err)
	}
}

// DeployBnMTokenAptos deploys two tokens on to the EVM and Aptos chain and sets up a lane between them.
// For Aptos, the test_token will be used along with the burn_mint_token_pool token pool.
func DeployBnMTokenAptos(
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
	aptos_burn_mint_token_pool.BurnMintTokenPool,
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

	signer := e.BlockChains.AptosChains()[aptosChainSel].DeployerSigner
	signerAddress := signer.AccountAddress()
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
	require.Falsef(t, (mcmsAddress == aptos.AccountAddress{}), "Aptos mcms address not found")
	ccipAddress := aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    shared.AptosCCIPType,
			Version: deployment.Version1_6_0,
		},
		aptosAddresses,
	)
	require.Falsef(t, (ccipAddress == aptos.AccountAddress{}), "Aptos CCIP address not found")

	// Deploy test token
	tokenObjectAddress, tx, testToken, err := test_token.DeployToObject(signer, client)
	require.NoError(t, err)
	data, err := client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy test_token: %v", data.VmStatus)

	tx, err = testToken.TestToken().Initialize(opts, nil, "Test Token", "TKN", 8, "", "", true)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to initialize test_token: %v", data.VmStatus)

	if mintAmount != nil {
		lggr.Infof("Minting %v tokens to %v...", mintAmount.Amount, mintAmount.To.StringLong())
		tx, err = testToken.TestToken().Mint(opts, mintAmount.To, mintAmount.Amount)
		require.NoError(t, err)
		data, err = client.WaitForTransaction(tx.Hash)
		require.NoError(t, err)
		require.True(t, data.Success, "failed to mint %d tokens to %v: %v", mintAmount.Amount, mintAmount.To.StringLong(), data.VmStatus)
	}

	tokenAddress, err := testToken.TestToken().TokenMetadata(nil)
	require.NoError(t, err)

	// Save addresses in address book
	typeAndVersion := cldf.NewTypeAndVersion(shared.AptosTestTokenType, deployment.Version1_6_0)
	typeAndVersion.AddLabel("TKN")
	err = e.ExistingAddresses.Save(aptosChainSel, tokenObjectAddress.StringLong(), typeAndVersion)
	require.NoError(t, err)
	typeAndVersion = cldf.NewTypeAndVersion(cldf.ContractType("TKN"), deployment.Version1_6_0)
	err = e.ExistingAddresses.Save(aptosChainSel, tokenAddress.StringLong(), typeAndVersion)
	require.NoError(t, err)

	// Deploy BnM Token Pool
	tokenPoolAddress, tx, _, err := token_pool.DeployToObject(signer, client, ccipAddress, mcmsAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy token_pool package: %v", data.VmStatus)

	tx, burnMintTokenPool, err := aptos_burn_mint_token_pool.DeployToExistingObject(signer, client, ccipAddress, mcmsAddress, tokenPoolAddress, tokenAddress, true)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy burn mint token pool: %v", data.VmStatus)

	typeAndVersion = cldf.NewTypeAndVersion(shared.BurnMintTokenPool, deployment.Version1_6_0)
	typeAndVersion.AddLabel(tokenAddress.StringLong())
	err = e.ExistingAddresses.Save(aptosChainSel, tokenPoolAddress.StringLong(), typeAndVersion)
	require.NoError(t, err)

	// Deploy BnM registrar
	tx, bnmRegistrar, err := bnm_registrar.DeployToExistingObject(signer, client, tokenObjectAddress, tokenPoolAddress, ccipAddress, tokenPoolAddress, mcmsAddress, tokenAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy BnM Registrar: %v", data.VmStatus)

	// Initialize token pool
	tx, err = bnmRegistrar.BnMRegistrar().Initialize(opts)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to initialize BnM token pool: %v", data.VmStatus)

	ccipContract := ccip.Bind(ccipAddress, client)
	tx, err = ccipContract.TokenAdminRegistry().ProposeAdministrator(opts, tokenAddress, signer.AccountAddress())
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.Truef(t, data.Success, "failed to propose %v as an administrator for token %v: %v", signerAddress.StringLong(), tokenAddress.StringLong(), data.VmStatus)

	tx, err = ccipContract.TokenAdminRegistry().AcceptAdminRole(opts, tokenAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.Truef(t, data.Success, "failed to accept administrator role for token %v: %v", tokenAddress.StringLong(), data.VmStatus)

	tx, err = ccipContract.TokenAdminRegistry().SetPool(opts, tokenAddress, tokenPoolAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.Truef(t, data.Success, "failed to call set_pool for token %v and token pool %v: %v", tokenAddress.StringLong(), tokenPoolAddress.StringLong(), data.VmStatus)

	// Transfer token pool to mcms
	mcmsContract := mcms.Bind(mcmsAddress, client)
	tokenPoolOwnerAddress, err := mcmsContract.MCMSRegistry().GetPreexistingCodeObjectOwnerAddress(nil, burnMintTokenPool.Address())
	require.NoError(t, err)
	tx, err = burnMintTokenPool.BurnMintTokenPool().TransferOwnership(opts, tokenPoolOwnerAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to initiate ownership transfer of BnM token pool to %v: %v", tokenPoolOwnerAddress, data.VmStatus)

	_, err = commoncs.Apply(t, e,
		commoncs.Configure(aptoscs.AcceptTokenPoolOwnership{},
			config.AcceptTokenPoolOwnershipInput{
				ChainSelector: aptosChainSel,
				Accepts: []config.TokenPoolAccept{
					{
						TokenPoolAddress: tokenPoolAddress,
						TokenPoolType:    shared.BurnMintTokenPool,
					},
				},
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Second,
				},
			},
		),
	)
	require.NoError(t, err)

	tx, err = burnMintTokenPool.BurnMintTokenPool().ExecuteOwnershipTransfer(opts, tokenPoolOwnerAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to execute ownership transfer of BnM token pool to %v: %", tokenPoolOwnerAddress, data.VmStatus)

	e, err = commoncs.Apply(t, e,
		commoncs.Configure(aptoscs.AddTokenPool{},
			config.AddTokenPoolConfig{
				ChainSelector:                       aptosChainSel,
				TokenAddress:                        tokenAddress,
				TokenCodeObjAddress:                 tokenObjectAddress,
				TokenPoolAddress:                    tokenPoolAddress,
				PoolType:                            shared.BurnMintTokenPool,
				TokenTransferFeeByRemoteChainConfig: nil,
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
					MinDelay: time.Second,
				},
			},
		),
	)
	require.NoError(t, err)

	aptosAddresses, err = e.ExistingAddresses.AddressesForChain(aptosChainSel)
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
	tokenPoolAddress = aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    shared.BurnMintTokenPool,
			Version: deployment.Version1_6_0,
			Labels:  cldf.NewLabelSet(tokenMetadataAddress.StringLong()),
		},
		aptosAddresses,
	)
	aptosTokenPool := aptos_burn_mint_token_pool.Bind(tokenPoolAddress, e.BlockChains.AptosChains()[aptosChainSel].Client)
	lggr.Debugf("Deployed Burn Mint Token Pool for %v to %v", tokenMetadataAddress.StringLong(), tokenPoolAddress.StringLong())

	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChainSel], evmPool, evmDeployerKey, aptosChainSel, tokenMetadataAddress[:], tokenPoolAddress[:])
	require.NoError(t, err)

	err = grantMintBurnPermissions(lggr, e.BlockChains.EVMChains()[evmChainSel], evmToken, evmDeployerKey, evmPool.Address())
	require.NoError(t, err)

	return evmToken, evmPool, tokenMetadataAddress, aptosTokenPool, nil
}

// DeployLnRTokenAptos deploys two tokens onto the EVM and Aptos chain and sets up a lane between them
// For Aptos, the test_token will be used along with the lock_release_token_pool token pool.
//
// The `withDispatchHooks` parameter decides whether the token pool will be initialized with a TransferRef or not:
//   - If set to true, the token will be initialized as a dispatchable fungible asset with a withdrawal/deposit hook.
//     Since this requires the LnR pool to have access to a TransferRef, the pool will be initialized with one.
//   - If set to false, the token will be initialized as a fungible asset and the token pool will be initialized without a TransferRef
func DeployLnRTokenAptos(
	t *testing.T,
	lggr logger.Logger,
	e cldf.Environment,
	evmChainSel, aptosChainSel uint64,
	tokenName string,
	mintAmount *config.TokenMint,
	withDispatchHooks bool,
) (
	*burn_mint_erc677.BurnMintERC677,
	*burn_mint_token_pool.BurnMintTokenPool,
	aptos.AccountAddress,
	lock_release_token_pool.LockReleaseTokenPool,
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

	signer := e.BlockChains.AptosChains()[aptosChainSel].DeployerSigner
	signerAddress := signer.AccountAddress()
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
	require.Falsef(t, (mcmsAddress == aptos.AccountAddress{}), "Aptos mcms address not found")
	ccipAddress := aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    shared.AptosCCIPType,
			Version: deployment.Version1_6_0,
		},
		aptosAddresses,
	)
	require.Falsef(t, (ccipAddress == aptos.AccountAddress{}), "Aptos CCIP address not found")

	// Deploy test token
	tokenObjectAddress, tx, testToken, err := test_token.DeployToObject(signer, client)
	require.NoError(t, err)
	data, err := client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy test_token: %v", data.VmStatus)

	tx, err = testToken.TestToken().Initialize(opts, nil, "Test Token", "TKN", 8, "", "", withDispatchHooks)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to initialize test_token: %v", data.VmStatus)

	if mintAmount != nil {
		lggr.Infof("Minting %v tokens to %v...", mintAmount.Amount, mintAmount.To.StringLong())
		tx, err = testToken.TestToken().Mint(opts, mintAmount.To, mintAmount.Amount)
		require.NoError(t, err)
		data, err = client.WaitForTransaction(tx.Hash)
		require.NoError(t, err)
		require.True(t, data.Success, "failed to mint %d tokens to %v: %v", mintAmount.Amount, mintAmount.To.StringLong(), data.VmStatus)
	}

	tokenAddress, err := testToken.TestToken().TokenMetadata(nil)
	require.NoError(t, err)

	// Save addresses in address book
	typeAndVersion := cldf.NewTypeAndVersion(shared.AptosTestTokenType, deployment.Version1_6_0)
	typeAndVersion.AddLabel("TKN")
	err = e.ExistingAddresses.Save(aptosChainSel, tokenObjectAddress.StringLong(), typeAndVersion)
	require.NoError(t, err)
	typeAndVersion = cldf.NewTypeAndVersion(cldf.ContractType("TKN"), deployment.Version1_6_0)
	err = e.ExistingAddresses.Save(aptosChainSel, tokenAddress.StringLong(), typeAndVersion)
	require.NoError(t, err)

	// Deploy LnR Token Pool
	tokenPoolAddress, tx, _, err := token_pool.DeployToObject(signer, client, ccipAddress, mcmsAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy token_pool package: %v", data.VmStatus)

	tx, lockReleaseTokenPool, err := lock_release_token_pool.DeployToExistingObject(signer, client, ccipAddress, mcmsAddress, tokenPoolAddress, tokenAddress, true)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy lock release token pool: %v", data.VmStatus)

	typeAndVersion = cldf.NewTypeAndVersion(shared.LockReleaseTokenPool, deployment.Version1_6_0)
	typeAndVersion.AddLabel(tokenAddress.StringLong())
	err = e.ExistingAddresses.Save(aptosChainSel, tokenPoolAddress.StringLong(), typeAndVersion)
	require.NoError(t, err)

	// Deploy LnR registrar
	tx, lnrRegistrar, err := lnr_registrar.DeployToExistingObject(signer, client, tokenObjectAddress, tokenPoolAddress, ccipAddress, tokenPoolAddress, mcmsAddress, tokenAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to deploy LnR Registrar: %v", data.VmStatus)

	// Initialize token pool
	if withDispatchHooks {
		tx, err = lnrRegistrar.LnRRegistrar().Initialize(opts)
		require.NoError(t, err)
		data, err = client.WaitForTransaction(tx.Hash)
		require.NoError(t, err)
		require.True(t, data.Success, "failed to initialize LnR token pool: %v", data.VmStatus)
	} else {
		tx, err = lnrRegistrar.LnRRegistrar().InitializeWithoutTransferRef(opts)
		require.NoError(t, err)
		data, err = client.WaitForTransaction(tx.Hash)
		require.NoError(t, err)
		require.True(t, data.Success, "failed to initialize LnR token pool without TransferRef: %v", data.VmStatus)
	}

	ccipContract := ccip.Bind(ccipAddress, client)
	tx, err = ccipContract.TokenAdminRegistry().ProposeAdministrator(opts, tokenAddress, signer.AccountAddress())
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.Truef(t, data.Success, "failed to propose %v as an administrator for token %v: %v", signerAddress.StringLong(), tokenAddress.StringLong(), data.VmStatus)

	tx, err = ccipContract.TokenAdminRegistry().AcceptAdminRole(opts, tokenAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.Truef(t, data.Success, "failed to accept administrator role for token %v: %v", tokenAddress.StringLong(), data.VmStatus)

	tx, err = ccipContract.TokenAdminRegistry().SetPool(opts, tokenAddress, tokenPoolAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.Truef(t, data.Success, "failed to call set_pool for token %v and token pool %v: %v", tokenAddress.StringLong(), tokenPoolAddress.StringLong(), data.VmStatus)

	// Transfer token pool to mcms
	mcmsContract := mcms.Bind(mcmsAddress, client)
	tokenPoolOwnerAddress, err := mcmsContract.MCMSRegistry().GetPreexistingCodeObjectOwnerAddress(nil, lockReleaseTokenPool.Address())
	require.NoError(t, err)
	tx, err = lockReleaseTokenPool.LockReleaseTokenPool().TransferOwnership(opts, tokenPoolOwnerAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to initiate ownership transfer of BnM token pool to %v: %v", tokenPoolOwnerAddress, data.VmStatus)

	_, err = commoncs.Apply(t, e,
		commoncs.Configure(aptoscs.AcceptTokenPoolOwnership{},
			config.AcceptTokenPoolOwnershipInput{
				ChainSelector: aptosChainSel,
				Accepts: []config.TokenPoolAccept{
					{
						TokenPoolAddress: tokenPoolAddress,
						TokenPoolType:    shared.LockReleaseTokenPool,
					},
				},
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Second,
				},
			},
		),
	)
	require.NoError(t, err)

	tx, err = lockReleaseTokenPool.LockReleaseTokenPool().ExecuteOwnershipTransfer(opts, tokenPoolOwnerAddress)
	require.NoError(t, err)
	data, err = client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to execute ownership transfer of LnR token pool to %v: %", tokenPoolOwnerAddress, data.VmStatus)

	e, err = commoncs.Apply(t, e,
		commoncs.Configure(aptoscs.AddTokenPool{},
			config.AddTokenPoolConfig{
				ChainSelector:                       aptosChainSel,
				TokenAddress:                        tokenAddress,
				TokenCodeObjAddress:                 tokenObjectAddress,
				TokenPoolAddress:                    tokenPoolAddress,
				PoolType:                            shared.LockReleaseTokenPool,
				TokenTransferFeeByRemoteChainConfig: nil,
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
					MinDelay: time.Second,
				},
			},
		),
	)
	require.NoError(t, err)

	aptosAddresses, err = e.ExistingAddresses.AddressesForChain(aptosChainSel)
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
	tokenPoolAddress = aptosstate.FindAptosAddress(
		cldf.TypeAndVersion{
			Type:    shared.LockReleaseTokenPool,
			Version: deployment.Version1_6_0,
			Labels:  cldf.NewLabelSet(tokenMetadataAddress.StringLong()),
		},
		aptosAddresses,
	)
	aptosTokenPool := lock_release_token_pool.Bind(tokenPoolAddress, e.BlockChains.AptosChains()[aptosChainSel].Client)
	lggr.Debugf("Deployed Lock Release Token Pool for %v to %v", tokenMetadataAddress.StringLong(), tokenPoolAddress.StringLong())

	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChainSel], evmPool, evmDeployerKey, aptosChainSel, tokenMetadataAddress[:], tokenPoolAddress[:])
	require.NoError(t, err)

	err = grantMintBurnPermissions(lggr, e.BlockChains.EVMChains()[evmChainSel], evmToken, evmDeployerKey, evmPool.Address())
	require.NoError(t, err)

	return evmToken, evmPool, tokenMetadataAddress, aptosTokenPool, nil
}
