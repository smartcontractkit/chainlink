package gateway_connector

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"reflect"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	gwConfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gwCommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

type serviceWrapper struct {
	services.StateMachine

	config    *config.GatewayConnector
	keystore  keystore.Eth
	connector connector.GatewayConnector
	lggr      logger.Logger
}

type workflowConnectorSigner struct {
	services.StateMachine

	connector   connector.GatewayConnector
	signerKey   *ecdsa.PrivateKey
	nodeAddress string
	lggr        logger.Logger
}

var (
	_ connector.Signer = &workflowConnectorSigner{}
)

func NewWorkflowConnectorSigner(config *config.GatewayConnector, signerKey *ecdsa.PrivateKey, lggr logger.Logger) (*workflowConnectorSigner, error) {
	return &workflowConnectorSigner{
		nodeAddress: (*config).NodeAddress(),
		signerKey:   signerKey,
		lggr:        lggr.Named("WorkflowConnectorSigner"),
	}, nil
}

func (h *workflowConnectorSigner) Sign(data ...[]byte) ([]byte, error) {
	return gwCommon.SignData(h.signerKey, data...)
}

func (h *workflowConnectorSigner) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
}
func (h *workflowConnectorSigner) Start(ctx context.Context) error {
	return h.StartOnce("WorkflowConnectorHandler", func() error {
		return nil
	})
}
func (h *workflowConnectorSigner) Close() error {
	return h.StopOnce("WorkflowConnectorHandler", func() (err error) {
		return nil
	})
}

func (h *workflowConnectorSigner) SetConnector(connector connector.GatewayConnector) {
	h.connector = connector
}

func translateConfigs(f config.GatewayConnector) connector.ConnectorConfig {
	r := connector.ConnectorConfig{}
	if f.NodeAddress() != "" {
		r.NodeAddress = f.NodeAddress()
	}

	if f.DonID() != "" {
		r.DonId = f.DonID()
	}

	if f.Gateways != nil {
		r.Gateways = make([]connector.ConnectorGatewayConfig, len(r.Gateways))
		for index, element := range f.Gateways() {
			r.Gateways[index] = connector.ConnectorGatewayConfig{Id: element.ID(), URL: element.URL()}
		}
	}

	if !reflect.ValueOf(f.WSHandshakeTimeoutMillis).IsZero() {
		r.WsClientConfig = network.WebSocketClientConfig{HandshakeTimeoutMillis: f.WSHandshakeTimeoutMillis()}
	}

	// 0 are valid values
	r.AuthMinChallengeLen = f.AuthMinChallengeLen()
	r.AuthTimestampToleranceSec = f.AuthTimestampToleranceSec()
	return r
}

// NOTE: this wrapper is needed to make sure that our services are started after Keystore.
func NewGatewayConnectorServiceWrapper(config *gwConfig.GatewayConnector, keystore keystore.Eth, lggr logger.Logger) *serviceWrapper {
	return &serviceWrapper{
		config:   config,
		keystore: keystore,
		lggr:     lggr,
	}
}

func (e *serviceWrapper) Start(ctx context.Context) error {
	return e.StartOnce("GatewayConnectorServiceWrapper", func() error {
		conf := *e.config
		e.lggr.Infow("Starting GatewayConnectorServiceWrapper", "chainID", conf.ChainIDForNodeKey())
		chainId, _ := new(big.Int).SetString(conf.ChainIDForNodeKey(), 0)
		enabledKeys, err := e.keystore.EnabledKeysForChain(ctx, chainId)
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

		signer, err := NewWorkflowConnectorSigner(e.config, signerKey, e.lggr)
		if err != nil {
			return err
		}
		translated := translateConfigs(conf)
		e.connector, err = connector.NewGatewayConnector(&translated, signer, signer, clockwork.NewRealClock(), e.lggr)
		if err != nil {
			return err
		}
		return e.connector.Start(ctx)
	})
}

func (e *serviceWrapper) Close() error {
	return e.StopOnce("GatewayConnectorServiceWrapper", func() (err error) {
		return e.connector.Close()
	})
}

func (e *serviceWrapper) Ready() error {
	return nil
}

func (e *serviceWrapper) HealthReport() map[string]error {
	return nil
}

func (e *serviceWrapper) Name() string {
	return "GatewayConnectorServiceWrapper"
}
