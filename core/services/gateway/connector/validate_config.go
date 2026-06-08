package connector

import (
	"errors"
	"slices"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

func validateConnectorConfig(cfg *ConnectorConfig) error {
	if len(cfg.DonId) == 0 || len(cfg.DonId) > network.HandshakeDonIdLen {
		return errors.New("invalid DON ID")
	}

	multiDonMode := slices.ContainsFunc(cfg.Gateways, func(g ConnectorGatewayConfig) bool {
		return g.DonId != ""
	})
	if !multiDonMode {
		return nil
	}

	for _, g := range cfg.Gateways {
		if g.DonId == "" {
			return errors.New("all gateways must set DonID when multi-DON mode is enabled")
		}
	}
	return nil
}
