package v0_5

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	verifier "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

type VerifierView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}

// GenerateVerifierView generates a VerifierView from a Verifier contract.
func GenerateVerifierView(verifier *verifier.Verifier) (VerifierView, error) {
	if verifier == nil {
		return VerifierView{}, errors.New("cannot generate view for nil Verifier")
	}

	owner, err := verifier.Owner(nil)
	if err != nil {
		return VerifierView{}, fmt.Errorf("failed to get owner for Verifier: %w", err)
	}

	return VerifierView{
		Address:        verifier.Address(),
		Owner:          owner,
		TypeAndVersion: "Verifier 0.5.0",
	}, nil
}
