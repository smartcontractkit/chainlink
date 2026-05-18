package config

import (
	"errors"

	"github.com/aptos-labs/aptos-go-sdk"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
)

// DeployRegulatedTokenConfig drives the regulated token publish + initialize +
// transfer_ownership + transfer_admin to the MCMS registry owner, and emits an MCMS
// proposal with accept_ownership and accept_admin. Note: regulated_token cannot be
// deployed via MCMS due to DFA re-entrancy; those steps are signed directly by the
// deployer.
type DeployRegulatedTokenConfig struct {
	ChainSelector uint64
	TokenParams   TokenParams
	TokenMint     *TokenMint
	// RegistrarPreregister is passed to DeployMCMSRegistrarToExistingObject (default true).
	RegistrarPreregister *bool
	MCMSConfig           *cldfproposalutils.TimelockConfig
}

func (c DeployRegulatedTokenConfig) Validate() error {
	var errs []error
	if c.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMSConfig is required"))
	}
	if err := c.TokenParams.Validate(); err != nil {
		errs = append(errs, err)
	}
	// The deploy sequence grants MINTER_ROLE to the deployer and calls Mint via direct EOA
	// transactions before the MCMS proposal is built; reject invalid TokenMint up front so
	// we don't fail after partial on-chain side effects.
	if c.TokenMint != nil {
		if (c.TokenMint.To == aptos.AccountAddress{}) {
			errs = append(errs, errors.New("TokenMint.To address is empty"))
		}
		if c.TokenMint.Amount == 0 {
			errs = append(errs, errors.New("TokenMint.Amount is 0"))
		}
	}
	return errors.Join(errs...)
}

// RegistrarPreregisterOrDefault returns true when unset (matches test helper default).
func (c DeployRegulatedTokenConfig) RegistrarPreregisterOrDefault() bool {
	if c.RegistrarPreregister == nil {
		return true
	}
	return *c.RegistrarPreregister
}
