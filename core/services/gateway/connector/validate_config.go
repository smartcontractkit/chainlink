package connector

import (
	"errors"
	"slices"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

func validateConnectorConfig(cfg *ConnectorConfig) error {
	if len(cfg.DonID) == 0 || len(cfg.DonID) > network.HandshakeDonIDLen {
		return errors.New("invalid DON ID")
	}

	multiDonMode := multiDonMode(cfg)
	if !multiDonMode {
		return nil
	}

	for _, g := range cfg.Gateways {
		if g.DonID == "" {
			return errors.New("all gateways must set DonID when multi-DON mode is enabled")
		}
	}
	return nil
}

func multiDonMode(cfg *ConnectorConfig) bool {
	return slices.ContainsFunc(cfg.Gateways, func(g ConnectorGatewayConfig) bool {
		return g.DonID != ""
	})
}
