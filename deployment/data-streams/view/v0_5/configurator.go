package v0_5

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
)

type ConfiguratorView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}

// GenerateConfiguratorView generates a ConfiguratorView from a Configurator contract.
func GenerateConfiguratorView(configurator *configurator.Configurator) (ConfiguratorView, error) {
	if configurator == nil {
		return ConfiguratorView{}, errors.New("cannot generate view for nil Configurator")
	}

	owner, err := configurator.Owner(nil)
	if err != nil {
		return ConfiguratorView{}, fmt.Errorf("failed to get owner for Configurator: %w", err)
	}

	return ConfiguratorView{
		Address:        configurator.Address(),
		Owner:          owner,
		TypeAndVersion: "Configurator 0.5.0",
	}, nil
}
