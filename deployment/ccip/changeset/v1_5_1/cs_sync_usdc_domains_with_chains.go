package v1_5_1

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var _ deployment.ChangeSet[SyncUSDCDomainsWithChainsConfig] = SyncUSDCDomainsWithChainsChangeset

// SyncUSDCDomainsWithChainsConfig defines the chain selector -> USDC domain mappings.
type SyncUSDCDomainsWithChainsConfig struct {
	// USDCVersionByChain defines the USDC domain and pool version for each chain selector.
	USDCVersionByChain map[uint64]semver.Version
	// ChainSelectorToUSDCDomain maps chains selectors to their USDC domain identifiers.
	ChainSelectorToUSDCDomain map[uint64]uint32
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
}

func (c SyncUSDCDomainsWithChainsConfig) Validate(env deployment.Environment, state changeset.CCIPOnChainState) error {
	ctx := env.GetContext()

	if c.ChainSelectorToUSDCDomain == nil {
		return errors.New("chain selector to usdc domain must be defined")
	}

	// Validate that all USDC configs inputted are for valid chains that define USDC pools.
	for chainSelector, version := range c.USDCVersionByChain {
		err := deployment.IsValidChainSelector(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in state", chainSelector)
		}
		if chainState.USDCTokenPools == nil {
			return fmt.Errorf("%s does not define any USDC token pools, config should be removed", chain)
		}
		if chainState.Timelock == nil {
			return fmt.Errorf("missing timelock on %s", chain.String())
		}
		if chainState.ProposerMcm == nil {
			return fmt.Errorf("missing proposerMcm on %s", chain.String())
		}
		usdcTokenPool, ok := chainState.USDCTokenPools[version]
		if !ok {
			return fmt.Errorf("no USDC token pool found on %s with version %s", chain, version)
		}
		if len(usdcTokenPool.Address().Bytes()) > 32 {
			// Will never be true for EVM
			return fmt.Errorf("expected USDC token pool address on %s (%s) to be less than 32 bytes", chain, usdcTokenPool.Address())
		}
		// Validate that the USDC token pool is owned by the address that will be actioning the transactions (i.e. Timelock or deployer key)
		if err := commoncs.ValidateOwnership(ctx, c.MCMS != nil, chain.DeployerKey.From, chainState.Timelock.Address(), usdcTokenPool); err != nil {
			return fmt.Errorf("token pool with address %s on %s failed ownership validation: %w", usdcTokenPool.Address(), chain, err)
		}
		// Validate that each supported chain has a domain ID defined
		supportedChains, err := usdcTokenPool.GetSupportedChains(&bind.CallOpts{Context: ctx})
		if err != nil {
			return fmt.Errorf("failed to get supported chains from USDC token pool on %s with address %s: %w", chain, usdcTokenPool.Address(), err)
		}
		for _, supportedChain := range supportedChains {
			if _, ok := c.ChainSelectorToUSDCDomain[supportedChain]; !ok {
				return fmt.Errorf("no USDC domain ID defined for chain with selector %d", supportedChain)
			}
		}
	}
	// Check that our input covers all chains that define USDC pools.
	for chainSelector, chainState := range state.Chains {
		if _, ok := c.USDCVersionByChain[chainSelector]; !ok && chainState.USDCTokenPools != nil {
			return fmt.Errorf("no USDC chain config defined for %s, which does support USDC", env.Chains[chainSelector])
		}
	}
	return nil
}

// SyncUSDCDomainsWithChainsChangeset syncs domain support on specified USDC token pools with its chain support.
// As such, it is expected that ConfigureTokenPoolContractsChangeset is executed before running this changeset.
func SyncUSDCDomainsWithChainsChangeset(env deployment.Environment, c SyncUSDCDomainsWithChainsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if err := c.Validate(env, state); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid SyncUSDCDomainsWithChainsConfig: %w", err)
	}
	readOpts := &bind.CallOpts{Context: env.GetContext()}

	deployerGroup := changeset.NewDeployerGroup(env, state, c.MCMS).WithDeploymentContext("sync domain support with chain support on USDC token pools")

	for chainSelector, version := range c.USDCVersionByChain {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]
		writeOpts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get transaction opts for %s", chain)
		}

		usdcTokenPool := chainState.USDCTokenPools[version]
		supportedChains, err := usdcTokenPool.GetSupportedChains(readOpts)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch supported chains from USDC token pool with address %s on %s: %w", usdcTokenPool.Address(), chain, err)
		}

		domainUpdates := make([]usdc_token_pool.USDCTokenPoolDomainUpdate, 0)
		for _, remoteChainSelector := range supportedChains {
			remoteChainState := state.Chains[remoteChainSelector]
			remoteUSDCVersion := c.USDCVersionByChain[remoteChainSelector]
			remoteUSDCTokenPool := remoteChainState.USDCTokenPools[remoteUSDCVersion]

			var desiredAllowedCaller [32]byte
			copy(desiredAllowedCaller[:], common.LeftPadBytes(remoteUSDCTokenPool.Address().Bytes(), 32))

			desiredDomainIdentifier := c.ChainSelectorToUSDCDomain[remoteChainSelector]

			currentDomain, err := usdcTokenPool.GetDomain(readOpts, remoteChainSelector)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch domain for %d from USDC token pool with address %s on %s: %w", remoteChainSelector, usdcTokenPool.Address(), chain, err)
			}
			// If any parameters are different, we need to add a setDomains call
			if currentDomain.AllowedCaller != desiredAllowedCaller ||
				currentDomain.DomainIdentifier != desiredDomainIdentifier {
				domainUpdates = append(domainUpdates, usdc_token_pool.USDCTokenPoolDomainUpdate{
					AllowedCaller:     desiredAllowedCaller,
					Enabled:           true,
					DomainIdentifier:  desiredDomainIdentifier,
					DestChainSelector: remoteChainSelector,
				})
			}
		}

		if len(domainUpdates) > 0 {
			_, err := usdcTokenPool.SetDomains(writeOpts, domainUpdates)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create set domains operation on %s: %w", chain, err)
			}
		}
	}

	return deployerGroup.Enact()
}
