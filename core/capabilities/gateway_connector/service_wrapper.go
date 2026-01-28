package gatewayconnector

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

type Keystore interface {
	keys.AddressChecker
	keys.MessageSigner
}

// DeterministicKeyProvider is an interface for keystores that support
// deterministic key discovery. This allows auto-discovery of node addresses.
type DeterministicKeyProvider interface {
	EnabledKeysWithDeterminism(ctx context.Context, chainID *big.Int) (keys []ethkey.KeyV2, err error)
}

type ServiceWrapper struct {
	services.StateMachine
	stopCh services.StopChan

	config                config.GatewayConnector
	keystore              Keystore
	deterministicProvider DeterministicKeyProvider
	chainID               *big.Int
	connector             connector.GatewayConnector
	lggr                  logger.Logger
	clock                 clockwork.Clock
	discoveredNodeAddress string // Stores auto-discovered node address if not configured
}

func translateConfigs(f config.GatewayConnector) connector.ConnectorConfig {
	r := connector.ConnectorConfig{}
	r.NodeAddress = f.NodeAddress()
	r.DonId = f.DonID()

	if len(f.Gateways()) != 0 {
		r.Gateways = make([]connector.ConnectorGatewayConfig, len(f.Gateways()))
		for index, element := range f.Gateways() {
			r.Gateways[index] = connector.ConnectorGatewayConfig{Id: element.ID(), URL: element.URL()}
		}
	}

	r.WsClientConfig = network.WebSocketClientConfig{HandshakeTimeoutMillis: f.WSHandshakeTimeoutMillis()}
	r.AuthMinChallengeLen = f.AuthMinChallengeLen()
	r.AuthTimestampToleranceSec = f.AuthTimestampToleranceSec()
	return r
}

// NOTE: this wrapper is needed to make sure that our services are started after Keystore.
// keystore is used for signing operations.
// deterministicProvider is used for auto-discovering the node address if not configured.
// chainID is the chain ID for which keys should be discovered.
func NewGatewayConnectorServiceWrapper(
	config config.GatewayConnector,
	keystore Keystore,
	deterministicProvider DeterministicKeyProvider,
	chainID *big.Int,
	clock clockwork.Clock,
	lggr logger.Logger,
) *ServiceWrapper {
	return &ServiceWrapper{
		stopCh:                make(services.StopChan),
		config:                config,
		keystore:              keystore,
		deterministicProvider: deterministicProvider,
		chainID:               chainID,
		clock:                 clock,
		lggr:                  logger.Named(lggr, "GatewayConnectorServiceWrapper"),
	}
}

func (e *ServiceWrapper) Start(ctx context.Context) error {
	return e.StartOnce("GatewayConnectorServiceWrapper", func() error {
		conf := e.config
		nodeAddress := conf.NodeAddress()

		// Auto-discover node address if not configured
		if nodeAddress == "" {
			if e.deterministicProvider == nil {
				return errors.New("NodeAddress must be configured when deterministic key provider is not available")
			}
			keys, err := e.deterministicProvider.EnabledKeysWithDeterminism(ctx, e.chainID)
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				return errors.New("no enabled keys found for auto-discovery")
			}
			// Use the first account (lowest State.ID) as the node address
			nodeAddress = keys[0].Address.String()
			e.discoveredNodeAddress = nodeAddress
			e.lggr.Infow("Auto-discovered node address", "address", nodeAddress)
		}

		configuredNodeAddress := common.HexToAddress(nodeAddress)
		err := e.keystore.CheckEnabled(ctx, configuredNodeAddress)
		if err != nil {
			return err
		}

		// Update config with discovered address for translateConfigs
		translated := translateConfigs(conf)
		// Override NodeAddress in translated config if we auto-discovered it
		if translated.NodeAddress == "" {
			translated.NodeAddress = nodeAddress
		}

		e.connector, err = connector.NewGatewayConnector(&translated, e, e.clock, e.lggr)
		if err != nil {
			return err
		}
		return e.connector.Start(ctx)
	})
}

func (e *ServiceWrapper) Sign(ctx context.Context, data ...[]byte) ([]byte, error) {
	nodeAddress := e.config.NodeAddress()
	if nodeAddress == "" {
		nodeAddress = e.discoveredNodeAddress
	}
	account := common.HexToAddress(nodeAddress)
	return e.keystore.SignMessage(ctx, account, gwcommon.Flatten(data...))
}

func (e *ServiceWrapper) Close() error {
	return e.StopOnce("GatewayConnectorServiceWrapper", func() (err error) {
		close(e.stopCh)
		return e.connector.Close()
	})
}

func (e *ServiceWrapper) HealthReport() map[string]error {
	return map[string]error{e.Name(): e.Healthy()}
}

func (e *ServiceWrapper) Name() string {
	return e.lggr.Name()
}

func (e *ServiceWrapper) GetGatewayConnector() connector.GatewayConnector {
	return e.connector
}
