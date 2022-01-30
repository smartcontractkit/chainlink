package plugins

import (
	relay "github.com/smartcontractkit/chainlink-relay/pkg/plugin"

	"github.com/smartcontractkit/chainlink/core/services"
)

var (
	Solana SolanaRelayer = &solana{}
)

type SolanaRelayer interface {
	services.Service
	relay.Solana
}
