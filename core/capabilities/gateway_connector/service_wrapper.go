package gatewayconnector

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gwCommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

type ServiceWrapper struct {
	services.StateMachine

	config    config.GatewayConnector
	keystore  keystore.Eth
	connector connector.GatewayConnector
	lggr      logger.Logger
}

type connectorSigner struct {
	services.StateMachine

	connector   connector.GatewayConnector
	signerKey   *ecdsa.PrivateKey
	nodeAddress string
	lggr        logger.Logger
}

var _ connector.Signer = &connectorSigner{}

func NewConnectorSigner(config config.GatewayConnector, signerKey *ecdsa.PrivateKey, lggr logger.Logger) (*connectorSigner, error) {
	return &connectorSigner{
		nodeAddress: config.NodeAddress(),
		signerKey:   signerKey,
		lggr:        lggr.Named("ConnectorSigner"),
	}, nil
}

func (h *connectorSigner) Sign(data ...[]byte) ([]byte, error) {
	return gwCommon.SignData(h.signerKey, data...)
}

func (h *connectorSigner) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
}
func (h *connectorSigner) Start(ctx context.Context) error {
	return h.StartOnce("ConnectorSigner", func() error {
		return nil
	})
}
func (h *connectorSigner) Close() error {
	return h.StopOnce("ConnectorSigner", func() (err error) {
		return nil
	})
}

func (h *connectorSigner) SetConnector(connector connector.GatewayConnector) {
	h.connector = connector
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
func NewGatewayConnectorServiceWrapper(config config.GatewayConnector, keystore keystore.Eth, lggr logger.Logger) *ServiceWrapper {
	return &ServiceWrapper{
		config:   config,
		keystore: keystore,
		lggr:     lggr,
	}
}

func (e *ServiceWrapper) Start(ctx context.Context) error {
	return e.StartOnce("GatewayConnectorServiceWrapper", func() error {
		conf := e.config
		e.lggr.Infow("Starting GatewayConnectorServiceWrapper", "chainID")

		//     logger.go:146: 2024-09-10T07:52:35.248-0700	ERROR	zap@v1.27.0/sugar.go:257	Ignored key without a value.	{"version": "unset@unset", "ignored": "chainID"}
		// go.uber.org/zap.(*SugaredLogger).Infow
		///Users/davidorchard/go/pkg/mod/go.uber.org/zap@v1.27.0/sugar.go:257

		e.lggr.Infow("Starting GatewayConnectorServiceWrapper2", "chainID", conf.ChainIDForNodeKey())
		chainID, _ := new(big.Int).SetString(conf.ChainIDForNodeKey(), 0)
		enabledKeys, err := e.keystore.EnabledKeysForChain(ctx, chainID)
		if err != nil {
			return err
		}
		if len(enabledKeys) == 0 {
			return errors.New("no available keys found")
		}
		configuredNodeAddress := common.HexToAddress(conf.NodeAddress())
		idx := slices.IndexFunc(enabledKeys, func(key ethkey.KeyV2) bool { return key.Address == configuredNodeAddress })
		if idx == -1 {
			return errors.New("key for configured node address not found")
		}
		signerKey := enabledKeys[idx].ToEcdsaPrivKey()
		if enabledKeys[idx].ID() != conf.NodeAddress() {
			return errors.New("node address mismatch")
		}

		signer, err := NewConnectorSigner(e.config, signerKey, e.lggr)
		if err != nil {
			return err
		}
		translated := translateConfigs(conf)
		e.connector, err = connector.NewGatewayConnector(&translated, signer, clockwork.NewRealClock(), e.lggr)
		if err != nil {
			return err
		}
		return e.connector.Start(ctx)
	})
}

func (e *ServiceWrapper) Close() error {
	return e.StopOnce("GatewayConnectorServiceWrapper", func() (err error) {
		return e.connector.Close()
	})
}

func (e *ServiceWrapper) Ready() error {
	return nil
}

func (e *ServiceWrapper) HealthReport() map[string]error {
	return nil
}

func (e *ServiceWrapper) Name() string {
	return "GatewayConnectorServiceWrapper"
}
