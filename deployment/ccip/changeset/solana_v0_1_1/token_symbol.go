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

// recordOnboardedTokenMint adds the mint to the datastore output and, when it is new to the
// environment, to the legacy AddressBook output. The symbol is the datastore key; the labels
// retain the metadata needed to describe the pool. A known mint is omitted from the returned
// AddressBook because ApplyChangesets merges that output with the environment's existing
// AddressBook, whose duplicate-address check rejects re-emitting an existing row. The datastore
// write is still returned so reruns and older onboardings can backfill the token ref.
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
	addressBook := ab
	for knownAddr := range envAddresses {
		if shared.AddressesEqual(chainSelector, knownAddr, cfg.TokenMint.String()) {
			addressBook = nil
			break
		}
	}

	return shared.RecordAddress(addressBook, ds, chainSelector, cfg.TokenMint.String(), tv, cfg.TokenSymbol)
}
