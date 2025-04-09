package v0_5

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
)

type VerifierProxyView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}

// GenerateVerifierProxyView generates a VerifierProxyView from a VerifierProxy contract.
func GenerateVerifierProxyView(verifierProxy *verifier_proxy_v0_5_0.VerifierProxy) (VerifierProxyView, error) {
	if verifierProxy == nil {
		return VerifierProxyView{}, errors.New("cannot generate view for nil VerifierProxy")
	}

	owner, err := verifierProxy.Owner(nil)
	if err != nil {
		return VerifierProxyView{}, fmt.Errorf("failed to get owner for VerifierProxy: %w", err)
	}

	return VerifierProxyView{
		Address:        verifierProxy.Address(),
		Owner:          owner,
		TypeAndVersion: "VerifierProxy 0.5.0",
	}, nil
}
