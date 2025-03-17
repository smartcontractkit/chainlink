package types

import (
	"errors"
	"math/big"
	"time"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

type MCMSRole string

const (
	BypasserManyChainMultisig  deployment.ContractType = "BypasserManyChainMultiSig"
	CancellerManyChainMultisig deployment.ContractType = "CancellerManyChainMultiSig"
	ProposerManyChainMultisig  deployment.ContractType = "ProposerManyChainMultiSig"
	ManyChainMultisig          deployment.ContractType = "ManyChainMultiSig"
	RBACTimelock               deployment.ContractType = "RBACTimelock"
	CallProxy                  deployment.ContractType = "CallProxy"

	// roles
	ProposerRole  MCMSRole = "PROPOSER"
	BypasserRole  MCMSRole = "BYPASSER"
	CancellerRole MCMSRole = "CANCELLER"

	// LinkToken is the burn/mint link token. It should be used everywhere for
	// new deployments. Corresponds to
	// https://github.com/smartcontractkit/chainlink/blob/develop/core/gethwrappers/shared/generated/link_token/link_token.go#L34
	LinkToken deployment.ContractType = "LinkToken"
	// StaticLinkToken represents the (very old) non-burn/mint link token.
	// It is not used in new deployments, but still exists on some chains
	// and has a distinct ABI from the new LinkToken.
	// Corresponds to the ABI
	// https://github.com/smartcontractkit/chainlink/blob/develop/core/gethwrappers/generated/link_token_interface/link_token_interface.go#L34
	StaticLinkToken deployment.ContractType = "StaticLinkToken"
	// mcms Solana specific
	ManyChainMultisigProgram         deployment.ContractType = "ManyChainMultiSigProgram"
	RBACTimelockProgram              deployment.ContractType = "RBACTimelockProgram"
	AccessControllerProgram          deployment.ContractType = "AccessControllerProgram"
	ProposerAccessControllerAccount  deployment.ContractType = "ProposerAccessControllerAccount"
	ExecutorAccessControllerAccount  deployment.ContractType = "ExecutorAccessControllerAccount"
	CancellerAccessControllerAccount deployment.ContractType = "CancellerAccessControllerAccount"
	BypasserAccessControllerAccount  deployment.ContractType = "BypasserAccessControllerAccount"
)

func (role MCMSRole) String() string {
	return string(role)
}

type MCMSWithTimelockConfig struct {
	Canceller        config.Config
	Bypasser         config.Config
	Proposer         config.Config
	TimelockMinDelay *big.Int
	Label            *string
}

// MCMSWithTimelockConfigV2 holds the configuration for an MCMS with timelock.
// Note that this type already exists in types.go, but this one is using the new lib version.
type MCMSWithTimelockConfigV2 struct {
	Canceller        mcmstypes.Config
	Bypasser         mcmstypes.Config
	Proposer         mcmstypes.Config
	TimelockMinDelay *big.Int
	Label            *string
}

type OCRParameters struct {
	DeltaProgress                           time.Duration
	DeltaResend                             time.Duration
	DeltaInitial                            time.Duration
	DeltaRound                              time.Duration
	DeltaGrace                              time.Duration
	DeltaCertifiedCommitRequest             time.Duration
	DeltaStage                              time.Duration
	Rmax                                    uint64
	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationShouldAcceptAttestedReport   time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration
}

func (params OCRParameters) Validate() error {
	if params.DeltaProgress <= 0 {
		return errors.New("deltaProgress must be positive")
	}
	if params.DeltaResend <= 0 {
		return errors.New("deltaResend must be positive")
	}
	if params.DeltaInitial <= 0 {
		return errors.New("deltaInitial must be positive")
	}
	if params.DeltaRound <= 0 {
		return errors.New("deltaRound must be positive")
	}
	if params.DeltaGrace <= 0 {
		return errors.New("deltaGrace must be positive")
	}
	if params.DeltaCertifiedCommitRequest <= 0 {
		return errors.New("deltaCertifiedCommitRequest must be positive")
	}
	if params.DeltaStage < 0 {
		return errors.New("deltaStage must be positive or 0 for disabled")
	}
	if params.Rmax <= 0 {
		return errors.New("rmax must be positive")
	}
	if params.MaxDurationQuery <= 0 {
		return errors.New("maxDurationQuery must be positive")
	}
	if params.MaxDurationObservation <= 0 {
		return errors.New("maxDurationObservation must be positive")
	}
	if params.MaxDurationShouldAcceptAttestedReport <= 0 {
		return errors.New("maxDurationShouldAcceptAttestedReport must be positive")
	}
	if params.MaxDurationShouldTransmitAcceptedReport <= 0 {
		return errors.New("maxDurationShouldTransmitAcceptedReport must be positive")
	}
	return nil
}
