package config

import (
	"errors"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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
	MCMSConfig           *proposalutils.TimelockConfig
}

func (c DeployRegulatedTokenConfig) Validate() error {
	var errs []error
	if c.MCMSConfig == nil {
		errs = append(errs, errors.New("MCMSConfig is required"))
	}
	if err := c.TokenParams.Validate(); err != nil {
		errs = append(errs, err)
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
