package testhelpers

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_1/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func addEVMDestChangesets(e *DeployedEnv, to, from uint64, isTestRouter bool) []commoncs.ConfiguredChangeSet {
	evmDstChangesets := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
			v1_6.UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
					to: {
						from: {
							IsEnabled:                 true,
							TestRouter:                isTestRouter,
							IsRMNVerificationDisabled: !e.RmnEnabledSourceChains[from],
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				TestRouter: isTestRouter,
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					// offramp update on dest chain
					to: {
						OffRampUpdates: map[uint64]bool{
							from: true,
						},
					},
				},
			},
		),
	}
	return evmDstChangesets
}

func addEVMSrcChangesets(from, to uint64, isTestRouter bool, gasprice map[uint64]*big.Int, tokenPrices map[common.Address]*big.Int, fqCfg fee_quoter.FeeQuoterDestChainConfig) []commoncs.ConfiguredChangeSet {
	evmSrcChangesets := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
			v1_6.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
					from: {
						to: {
							IsEnabled:        true,
							TestRouter:       isTestRouter,
							AllowListEnabled: false,
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterPricesChangeset),
			v1_6.UpdateFeeQuoterPricesConfig{
				PricesByChain: map[uint64]v1_6.FeeQuoterPriceUpdatePerSource{
					from: {
						TokenPrices: tokenPrices,
						GasPrices:   gasprice,
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
			v1_6.UpdateFeeQuoterDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
					from: {
						to: fqCfg,
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				TestRouter: isTestRouter,
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					// onRamp update on source chain
					from: {
						OnRampUpdates: map[uint64]bool{
							to: true,
						},
					},
				},
			},
		),
	}
	return evmSrcChangesets
}

func SendRequestEVM(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *CCIPSendReqConfig,
) (*onramp.OnRampCCIPMessageSent, error) {
	// Set default sender if not provided
	if cfg.Sender == nil {
		cfg.Sender = e.BlockChains.EVMChains()[cfg.SourceChain].DeployerKey
	}

	e.Logger.Infof("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, cfg.Sender.From.String())

	tx, blockNum, err := CCIPSendRequest(e, state, cfg)
	if err != nil {
		return nil, err
	}

	it, err := state.MustGetEVMChainState(cfg.SourceChain).OnRamp.FilterCCIPMessageSent(&bind.FilterOpts{
		Start:   blockNum,
		End:     &blockNum,
		Context: context.Background(),
	}, []uint64{cfg.DestChain}, []uint64{})
	if err != nil {
		return nil, err
	}

	if !it.Next() {
		return nil, errors.New("no CCIP message sent event found")
	}

	e.Logger.Infof("CCIP message (id %s) sent from chain selector %d to chain selector %d tx %s seqNum %d nonce %d sender %s testRouterEnabled %t",
		common.Bytes2Hex(it.Event.Message.Header.MessageId[:]),
		cfg.SourceChain,
		cfg.DestChain,
		tx.Hash().String(),
		it.Event.SequenceNumber,
		it.Event.Message.Header.Nonce,
		it.Event.Message.Sender.String(),
		cfg.IsTestRouter,
	)
	return it.Event, nil
}
