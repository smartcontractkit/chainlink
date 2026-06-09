package types

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type MCMSRole string

func (role MCMSRole) String() string {
	return string(role)
}

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
