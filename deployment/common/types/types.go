package types

import (
	"errors"
	"math/big"
	"time"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type MCMSRole string

const (
	BypasserManyChainMultisig  cldf.ContractType = "BypasserManyChainMultiSig"
	CancellerManyChainMultisig cldf.ContractType = "CancellerManyChainMultiSig"
	ProposerManyChainMultisig  cldf.ContractType = "ProposerManyChainMultiSig"
	ManyChainMultisig          cldf.ContractType = "ManyChainMultiSig"
	RBACTimelock               cldf.ContractType = "RBACTimelock"
	CallProxy                  cldf.ContractType = "CallProxy"

	// roles
	ProposerRole  MCMSRole = "PROPOSER"
	BypasserRole  MCMSRole = "BYPASSER"
	CancellerRole MCMSRole = "CANCELLER"

	// LinkToken is the burn/mint link token. It should be used everywhere for
	// new deployments. Corresponds to
	// https://github.com/smartcontractkit/chainlink/blob/develop/core/gethwrappers/shared/generated/link_token/link_token.go#L34
	LinkToken cldf.ContractType = "LinkToken"
	// StaticLinkToken represents the (very old) non-burn/mint link token.
	// It is not used in new deployments, but still exists on some chains
	// and has a distinct ABI from the new LinkToken.
	// Corresponds to the ABI
	// https://github.com/smartcontractkit/chainlink/blob/develop/core/gethwrappers/generated/link_token_interface/link_token_interface.go#L34
	StaticLinkToken cldf.ContractType = "StaticLinkToken"
	// mcms Solana specific
	ManyChainMultisigProgram         cldf.ContractType = "ManyChainMultiSigProgram"
	RBACTimelockProgram              cldf.ContractType = "RBACTimelockProgram"
	AccessControllerProgram          cldf.ContractType = "AccessControllerProgram"
	ProposerAccessControllerAccount  cldf.ContractType = "ProposerAccessControllerAccount"
	ExecutorAccessControllerAccount  cldf.ContractType = "ExecutorAccessControllerAccount"
	CancellerAccessControllerAccount cldf.ContractType = "CancellerAccessControllerAccount"
	BypasserAccessControllerAccount  cldf.ContractType = "BypasserAccessControllerAccount"
)

func (role MCMSRole) String() string {
	return string(role)
}

type MCMSWithTimelockConfig struct {
	Canceller        config.Config `json:"canceller"`
	Bypasser         config.Config `json:"bypasser"`
	Proposer         config.Config `json:"proposer"`
	TimelockMinDelay *big.Int      `json:"timelockMinDelay"`
	Label            *string       `json:"label"`
}

// MCMSWithTimelockConfigV2 holds the configuration for an MCMS with timelock.
// Note that this type already exists in types.go, but this one is using the new lib version.
type MCMSWithTimelockConfigV2 struct {
	Canceller        mcmstypes.Config `json:"canceller"`
	Bypasser         mcmstypes.Config `json:"bypasser"`
	Proposer         mcmstypes.Config `json:"proposer"`
	TimelockMinDelay *big.Int         `json:"timelockMinDelay"`
	Label            *string          `json:"label"`
}

type OCRParameters struct {
	DeltaProgress                           time.Duration `json:"deltaProgress"`
	DeltaResend                             time.Duration `json:"deltaResend"`
	DeltaInitial                            time.Duration `json:"deltaInitial"`
	DeltaRound                              time.Duration `json:"deltaRound"`
	DeltaGrace                              time.Duration `json:"deltaGrace"`
	DeltaCertifiedCommitRequest             time.Duration `json:"deltaCertifiedCommitRequest"`
	DeltaStage                              time.Duration `json:"deltaStage"`
	Rmax                                    uint64        `json:"rmax"`
	MaxDurationQuery                        time.Duration `json:"maxDurationQuery"`
	MaxDurationObservation                  time.Duration `json:"maxDurationObservation"`
	MaxDurationShouldAcceptAttestedReport   time.Duration `json:"maxDurationShouldAcceptAttestedReport"`
	MaxDurationShouldTransmitAcceptedReport time.Duration `json:"maxDurationShouldTransmitAcceptedReport"`
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
