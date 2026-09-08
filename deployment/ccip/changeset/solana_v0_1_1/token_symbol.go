package solana

import (
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

// validateTokenSymbol checks the two rules for a token qualifier: it must be supplied, and it
// must name the token rather than copying the token's address.
func validateTokenSymbol(tokenSymbol, mint string) error {
	if tokenSymbol == "" {
		return errors.New("token symbol is required: it qualifies the token's datastore ref")
	}
	if strings.Contains(strings.ToLower(tokenSymbol), strings.ToLower(mint)) {
		return fmt.Errorf("token symbol %q contains the mint address: a qualifier is never derived from an address", tokenSymbol)
	}

	return nil
}

// recordOnboardedTokenMint adds the mint to this changeset's address-book and datastore outputs.
// The symbol is the datastore key; the labels retain the metadata needed to describe the pool.
// If the mint is already in the environment, its existing address-book metadata is reused so
// the returned address-book row is an idempotent merge rather than a conflicting rewrite.
func recordOnboardedTokenMint(
	chainSelector uint64,
	ab cldf.AddressBook,
	ds datastore.MutableDataStore,
	envAddresses map[string]cldf.TypeAndVersion,
	cfg OnboardTokenPoolConfig,
	poolProgramID string,
) error {
	if err := validateTokenSymbol(cfg.TokenSymbol, cfg.TokenMint.String()); err != nil {
		return err
	}

	tv := cldf.NewTypeAndVersion(cfg.TokenProgramName, deployment.Version1_0_0)
	tv.AddLabel(cfg.Metadata)
	tv.AddLabel(cfg.PoolType.String())
	tv.AddLabel(poolProgramID)
	for knownAddr := range envAddresses {
		if shared.AddressesEqual(chainSelector, knownAddr, cfg.TokenMint.String()) {
			// Reuse address book metadata for the emitted row.
			// The datastore still gets the current symbol as its key.
			tv = envAddresses[knownAddr]
			break
		}
	}

	return shared.RecordAddress(ab, ds, chainSelector, cfg.TokenMint.String(), tv, cfg.TokenSymbol)
}
