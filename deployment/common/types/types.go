package types

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"

	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	BypasserManyChainMultisig  deployment.ContractType = "BypasserManyChainMultiSig"
	CancellerManyChainMultisig deployment.ContractType = "CancellerManyChainMultiSig"
	ProposerManyChainMultisig  deployment.ContractType = "ProposerManyChainMultiSig"
	RBACTimelock               deployment.ContractType = "RBACTimelock"
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
)

type MCMSWithTimelockConfig struct {
	Canceller         config.Config
	Bypasser          config.Config
	Proposer          config.Config
	TimelockExecutors []common.Address
	TimelockMinDelay  *big.Int
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
		return fmt.Errorf("deltaProgress must be positive")
	}
	if params.DeltaResend <= 0 {
		return fmt.Errorf("deltaResend must be positive")
	}
	if params.DeltaInitial <= 0 {
		return fmt.Errorf("deltaInitial must be positive")
	}
	if params.DeltaRound <= 0 {
		return fmt.Errorf("deltaRound must be positive")
	}
	if params.DeltaGrace <= 0 {
		return fmt.Errorf("deltaGrace must be positive")
	}
	if params.DeltaCertifiedCommitRequest <= 0 {
		return fmt.Errorf("deltaCertifiedCommitRequest must be positive")
	}
	if params.DeltaStage <= 0 {
		return fmt.Errorf("deltaStage must be positive")
	}
	if params.Rmax <= 0 {
		return fmt.Errorf("rmax must be positive")
	}
	if params.MaxDurationQuery <= 0 {
		return fmt.Errorf("maxDurationQuery must be positive")
	}
	if params.MaxDurationObservation <= 0 {
		return fmt.Errorf("maxDurationObservation must be positive")
	}
	if params.MaxDurationShouldAcceptAttestedReport <= 0 {
		return fmt.Errorf("maxDurationShouldAcceptAttestedReport must be positive")
	}
	if params.MaxDurationShouldTransmitAcceptedReport <= 0 {
		return fmt.Errorf("maxDurationShouldTransmitAcceptedReport must be positive")
	}
	return nil
}
