package ccip

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/environment/crib"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	toncodec "github.com/smartcontractkit/chainlink-ton/pkg/ccip/codec"
)

func TestStaging_CCIP_Load(t *testing.T) {
	lggr := logger.Test(t)

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(lggr, *userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(crib.DeployerKeys{
		EVMKey:   *userOverrides.TestnetConfig.EVMPrivateKey,
		SolKey:   *userOverrides.TestnetConfig.SolanaPrivateKey,
		AptosKey: *userOverrides.TestnetConfig.AptosPrivateKey,
	})
	require.NoError(t, err)
	env, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, env)
	userOverrides.Validate(t, env)
	state, err := stateview.LoadOnchainState(*env)
	require.NoError(t, err)

	// Create context for the test duration
	ctx, cancel := context.WithTimeout(context.Background(), userOverrides.GetLoadDuration()+5*time.Minute)
	defer cancel()

	// Calculate block times for EVM chains (needed for metrics)
	blockTimes := make(map[uint64]uint64)
	evmChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	for _, cs := range evmChains {
		client := env.BlockChains.EVMChains()[cs].Client
		block1, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			lggr.Warnw("Failed to get block header", "chainSelector", cs, "error", err)
			blockTimes[cs] = 12 // default block time
			continue
		}
		time.Sleep(2 * time.Second)
		block2, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			lggr.Warnw("Failed to get second block header", "chainSelector", cs, "error", err)
			blockTimes[cs] = 12 // default block time
			continue
		}
		time1 := time.Unix(int64(block1.Time), 0)
		time2 := time.Unix(int64(block2.Time), 0)
		blockTimeDiff := int64(time2.Sub(time1))
		blockNumberDiff := new(big.Int).Sub(block2.Number, block1.Number).Int64()
		if blockNumberDiff > 0 {
			blockTime := blockTimeDiff / blockNumberDiff / int64(time.Second)
			blockTimes[cs] = uint64(blockTime) //nolint:gosec // G115
		} else {
			blockTimes[cs] = 12 // default block time
		}
		lggr.Infow("Chain block time", "chainSelector", cs, "blockTime", blockTimes[cs])
	}

	// Set block times for TON chains (TON has ~5 second block time)
	tonChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyTon))
	for _, cs := range tonChains {
		blockTimes[cs] = 5 // TON testnet block time ~5 seconds
		lggr.Infow("TON chain block time", "chainSelector", cs, "blockTime", blockTimes[cs])
	}

	// initialize additional accounts on other chains
	// Fund NumDestinationChains + len(tonChains) keys per source so each gun gets its own key (no nonce conflicts)
	numGuns := *userOverrides.NumDestinationChains + len(tonChains)
	gunSlice := make([]uint64, numGuns)
	copy(gunSlice, evmChains[:*userOverrides.NumDestinationChains])
	copy(gunSlice[*userOverrides.NumDestinationChains:], tonChains)
	transmitKeys, err := fundAdditionalKeys(lggr, *env, gunSlice, *userOverrides.TestnetConfig.FundingAmountEth, userOverrides.TestnetConfig.ChainFundingOverrides)
	require.NoError(t, err)

	// Build TON endpoint map from crib chain configs (private liteserver URLs)
	tonEndpoints := make(map[uint64]string)
	for _, chainCfg := range cribDeployOutput.Chains {
		if chainCfg.ChainType == "TON" && len(chainCfg.HTTPRPCs) > 0 && chainCfg.HTTPRPCs[0].External != "" {
			chainDetails, detailErr := chain_selectors.GetChainDetailsByChainIDAndFamily(chainCfg.ChainID, "ton")
			if detailErr != nil {
				lggr.Warnw("Failed to get TON chain selector", "chainID", chainCfg.ChainID, "error", detailErr)
				continue
			}
			tonEndpoints[chainDetails.ChainSelector] = chainCfg.HTTPRPCs[0].External
			lggr.Infow("Found TON endpoint from crib config", "chainSelector", chainDetails.ChainSelector, "endpoint", chainCfg.HTTPRPCs[0].External)
		}
	}

	// Initialize TON source keys before lane discovery so we can reuse the client
	// to query the TON router for TON→EVM lanes
	tonSourceKeys, initErr := initializeTonSourceKeys(env.GetContext(), lggr, env, *userOverrides.TestnetConfig.TonMnemonic, tonEndpoints)
	if initErr != nil {
		lggr.Warnw("Failed to initialize TON source keys", "error", initErr)
	}

	// Discover lanes from deployed state (EVM→TON discovered via EVM OnRamp, TON→EVM skipped due to nil client in env)
	laneConfig := &crib.LaneConfiguration{}
	err = laneConfig.DiscoverLanesFromDeployedState(lggr, *env, &state)
	require.NoError(t, err)

	// Discover TON→EVM lanes using the initialized TON client (reuses existing connection)
	for chainSel, manager := range tonSourceKeys {
		dests, queryErr := manager.GetEnabledDestinations(ctx)
		if queryErr != nil {
			lggr.Warnw("Failed to query TON router for destinations", "chainSelector", chainSel, "error", queryErr)
			continue
		}
		var tonLanes []crib.LaneConfig
		for _, dst := range dests {
			tonLanes = append(tonLanes, crib.LaneConfig{SourceChain: chainSel, DestinationChain: dst})
		}
		if len(tonLanes) > 0 {
			laneConfig.AppendLanes(tonLanes)
			lggr.Infow("Discovered TON→EVM lanes from router", "chainSelector", chainSel, "destinations", dests)
		}
	}

	// Log discovered lanes for operator verification
	laneConfig.LogLaneConfigInfo(lggr)

	// Fail fast if no lanes discovered — indicates deployment misconfiguration
	discoveredLanes, err := laneConfig.GetLanes()
	require.NoError(t, err)
	require.NotEmpty(t, discoveredLanes,
		"Lane discovery found zero lanes - verify staging deployment has configured TON↔EVM routes")

	// Initialize TON destination managers for execution event tracking
	// Reuses clients from tonSourceKeys to use the same private endpoint
	tonDestManagers, initErr := initializeTonDestinationManagers(ctx, lggr, env, tonSourceKeys)
	if initErr != nil {
		lggr.Warnw("Failed to initialize TON destination managers", "error", initErr)
	}

	// Initialize MetricsManager for Loki integration
	mm := NewMetricsManager(t, env.Logger, userOverrides, blockTimes)
	go mm.Start(ctx)

	// Track start blocks for event subscriptions
	startBlocks := make(map[uint64]*uint64)
	finalSeqNrCommitChannels := make(map[uint64]chan finalSeqNrReport)
	finalSeqNrExecChannels := make(map[uint64]chan finalSeqNrReport)
	var wg sync.WaitGroup

	// gunMap holds a destinationGun for every enabled destination chain
	gunMap := make(map[uint64]*DestinationGun)
	p := wasp.NewProfile()
	for ind := range *userOverrides.NumDestinationChains {
		cs := evmChains[ind]

		// Get start block for event subscriptions
		latesthdr, err := env.BlockChains.EVMChains()[cs].Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[cs] = &block

		messageKeys := make(map[uint64]*bind.TransactOpts)
		other := env.BlockChains.ListChainSelectors(
			cldf_chain.WithFamily(chain_selectors.FamilyEVM),
			cldf_chain.WithChainSelectorsExclusion([]uint64{cs}),
		)
		var mu sync.Mutex
		var wg2 sync.WaitGroup
		wg2.Add(len(other))
		for _, src := range other {
			go func(src uint64) {
				defer wg2.Done()
				mu.Lock()
				messageKeys[src] = transmitKeys[src][ind]
				mu.Unlock()
			}(src)
		}
		wg2.Wait()
		srcChains := laneConfig.GetSourceChainsForDestination(cs)

		// Initialize channels for this destination
		finalSeqNrCommitChannels[cs] = make(chan finalSeqNrReport)
		finalSeqNrExecChannels[cs] = make(chan finalSeqNrReport)

		gunMap[cs], err = NewDestinationGun(
			env.Logger,
			cs,
			*env,
			&state,
			state.MustGetEVMChainState(cs).Receiver.Address().Bytes(),
			userOverrides,
			messageKeys,
			nil,
			tonSourceKeys,
			mm.InputChan, // Pass metrics pipe for Loki integration
			srcChains,
		)
		if err != nil {
			lggr.Errorw("Failed to initialize DestinationGun for", "chainSelector", cs, "error", err)
			t.Fatal(err)
		}

		// Subscribe to commit and execution events for this destination
		wg.Add(2)
		go subscribeCommitEvents(
			ctx,
			lggr,
			state.Chains[cs].OffRamp,
			srcChains,
			startBlocks[cs],
			cs,
			env.BlockChains.EVMChains()[cs].Client,
			finalSeqNrCommitChannels[cs],
			&wg,
			mm.InputChan)
		go subscribeExecutionEvents(
			ctx,
			lggr,
			state.Chains[cs].OffRamp,
			srcChains,
			startBlocks[cs],
			cs,
			env.BlockChains.EVMChains()[cs].Client,
			finalSeqNrExecChannels[cs],
			&wg,
			mm.InputChan)
	}

	// Create DestinationGuns for TON destination chains (EVM → TON)
	for _, tonChainSel := range tonChains {
		srcChains := laneConfig.GetSourceChainsForDestination(tonChainSel)
		if len(srcChains) == 0 {
			lggr.Warnw("No source chains found for TON destination", "chainSelector", tonChainSel)
			continue
		}

		tonState, exists := state.TonChains[tonChainSel]
		if !exists {
			lggr.Warnw("No TON state found for chain", "chainSelector", tonChainSel)
			continue
		}

		// Get EVM message keys for all EVM source chains
		messageKeys := make(map[uint64]*bind.TransactOpts)
		tonKeyIndex := *userOverrides.NumDestinationChains // TON gun gets its own key index after EVM guns
		for _, evmChain := range evmChains {
			if tonKeyIndex < len(transmitKeys[evmChain]) {
				messageKeys[evmChain] = transmitKeys[evmChain][tonKeyIndex]
			}
		}

		// Convert TON receiver address to bytes for the gun
		addrCodec := toncodec.NewAddressCodec()
		receiverBytes, err := addrCodec.AddressStringToBytes(tonState.ReceiverAddress.String())
		if err != nil {
			lggr.Errorw("Failed to encode TON receiver address", "error", err, "chainSelector", tonChainSel)
			continue
		}

		gunMap[tonChainSel], err = NewDestinationGun(
			env.Logger,
			tonChainSel,
			*env,
			&state,
			receiverBytes,
			userOverrides,
			messageKeys,
			nil,
			tonSourceKeys,
			mm.InputChan,
			srcChains,
		)
		if err != nil {
			lggr.Errorw("Failed to initialize DestinationGun for TON chain", "chainSelector", tonChainSel, "error", err)
			continue
		}
		lggr.Infow("Created DestinationGun for TON chain", "chainSelector", tonChainSel, "sourceChains", srcChains)
	}

	// Subscribe to TON execution events for TON destination chains
	for tonChainSel, tonDestManager := range tonDestManagers {
		srcChains := laneConfig.GetSourceChainsForDestination(tonChainSel)
		if len(srcChains) == 0 {
			lggr.Warnw("No source chains found for TON destination", "chainSelector", tonChainSel)
			continue
		}

		lggr.Infow("Setting up TON destination event subscriber",
			"tonChainSelector", tonChainSel,
			"sourceChains", srcChains)

		wg.Add(1)
		go subscribeTonExecutionEvents(
			ctx,
			lggr,
			tonDestManager,
			srcChains,
			tonChainSel,
			&wg,
			mm.InputChan)
	}

	requestFrequency, err := time.ParseDuration(*userOverrides.RequestFrequency)
	require.NoError(t, err)

	// Only pass LokiConfig if URL is configured
	var lokiConfig *wasp.LokiConfig
	if cfg := wasp.NewEnvLokiConfig(); cfg.URL != "" {
		lokiConfig = cfg
	}

	for _, gun := range gunMap {
		p.Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			GenName:     "ccipLoad",
			LoadType:    wasp.RPS,
			CallTimeout: userOverrides.GetLoadDuration(),
			// 1 request per second for n seconds
			Schedule: wasp.Plain(1, userOverrides.GetLoadDuration()),
			// limit requests to 1 per duration
			RateLimitUnitDuration: requestFrequency,
			// will need to be divided by number of chains
			// this schedule is per generator
			// in this example, it would be 1 request per 5seconds per generator (dest chain)
			// so if there are 3 generators, it would be 3 requests per 5 seconds over the network
			Gun:        gun,
			Labels:     CommonTestLabels,
			LokiConfig: lokiConfig,
		}))
	}

	// DEBUG BREAKPOINT: Setup complete
	lggr.Info("=== SETUP COMPLETE ===")
	numTonDestGuns := 0
	for _, cs := range tonChains {
		if _, exists := gunMap[cs]; exists {
			numTonDestGuns++
		}
	}
	lggr.Infow("Discovered lanes summary",
		"totalDestinationGuns", len(gunMap),
		"numEVMDestinations", len(gunMap)-numTonDestGuns,
		"numTONDestinations", numTonDestGuns,
		"tonSourceKeys", len(tonSourceKeys),
	)

	_, err = p.Run(true)
	require.NoError(t, err)

	lggr.Info("Load test complete, shutting down goroutines")
	cancel() // Signal all goroutines to stop
	wg.Wait() // Wait for all goroutines to finish before test cleanup

	lggr.Info("All goroutines stopped, returning funds")
	// return funds to source address at the end of the test
	sourcePk, err := crypto.HexToECDSA(*userOverrides.TestnetConfig.EVMPrivateKey)
	if err != nil {
		lggr.Errorw("could not return funds to source address")
	}
	// Derive the public key
	publicKey := sourcePk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		lggr.Errorw("could not return funds to source address")
	}

	// Get the address from the public key
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	err = reclaimFunds(lggr, *env, transmitKeys, common.HexToAddress(address))
	if err != nil {
		lggr.Errorw(err.Error())
	}
}
