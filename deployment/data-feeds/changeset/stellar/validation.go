package stellar

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stellar/go-stellar-sdk/strkey"
)

// validateContract restricts a request's target to the two DF contract types.
func validateContract(t datastore.ContractType) error {
	if t != CacheContract && t != ProxyContract {
		return fmt.Errorf("unsupported contract type %q: must be %q or %q", t, CacheContract, ProxyContract)
	}
	return nil
}

// validateAddress accepts a G... account or C... contract strkey.
func validateAddress(s string) error {
	if strkey.IsValidEd25519PublicKey(s) {
		return nil
	}
	if _, err := strkey.Decode(strkey.VersionByteContract, s); err == nil {
		return nil
	}
	return fmt.Errorf("%q is not a valid Stellar account or contract address", s)
}

// verifyFeedPreconditions checks the preconditions shared by the cache feed
// changesets: cache address ref, valid admin, non-empty parseable DataIDs.
func verifyFeedPreconditions(env cldf.Environment, chainSel uint64, qualifier, version, admin string, dataIDs []string) error {
	if err := verifyContractRef(env, chainSel, CacheContract, qualifier, version); err != nil {
		return err
	}
	if err := validateAddress(admin); err != nil {
		return err
	}
	if len(dataIDs) == 0 {
		return errors.New("DataIDs cannot be empty")
	}
	_, err := dataIDsToBytes(dataIDs)
	return err
}
