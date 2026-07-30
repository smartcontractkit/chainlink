package nodeconfig

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pelletier/go-toml/v2"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// Inputs is the role/topology-dependent data that determines a node's 30-cre layer.
type Inputs struct {
	CapRegAddress      string
	CapRegVersion      string
	WorkflowRegAddress string
	WorkflowRegVersion string
	RegistryChainID    uint64
	Allowlist          []string

	IsBootstrapNode bool
	IsGatewayNode   bool

	BootstrapPeerID string
	BootstrapHost   string
	P2PPort         int

	GatewayConnector *GatewayConnectorConfig
}

type GatewayConnectorConfig struct {
	DonID             string
	ChainIDForNodeKey string
	NodeAddress       string
	Gateways          []ConnectorGateway
}

type ConnectorGateway struct {
	ID    string
	URL   string
	DonID string
}

// Generate builds the typed config for a node and marshals it to the 30-cre TOML layer.
func Generate(in Inputs) (string, error) {
	port := in.P2PPort
	if port == 0 {
		port = 5001
	}

	cfg := corechainlink.Config{Core: coretoml.Core{
		Feature: coretoml.Feature{LogPoller: new(true)},
		OCR2: coretoml.OCR2{
			Enabled:              new(true),
			DatabaseTimeout:      commonconfig.MustNewDuration(time.Second),
			ContractPollInterval: commonconfig.MustNewDuration(time.Second),
		},
		CRE: coretoml.CreConfig{
			EnableDKGRecipient:   new(true),
			UseLocalTimeProvider: new(false),
			DebugMode:            new(true),
			WorkflowFetcher:      &coretoml.WorkflowFetcherConfig{URL: new("file:///home/chainlink/workflows")},
		},
	}}

	// P2P.V2 + rage P2P
	cfg.P2P.EnableExperimentalRageP2P = new(true)
	cfg.P2P.V2.Enabled = new(true)
	cfg.P2P.V2.ListenAddresses = new([]string{fmt.Sprintf("0.0.0.0:%d", port)})
	switch {
	case in.IsBootstrapNode:
		bl, err := bootstrapper(in.BootstrapPeerID, fmt.Sprintf("localhost:%d", port))
		if err != nil {
			return "", err
		}
		cfg.P2P.V2.DefaultBootstrappers = new([]ocrcommontypes.BootstrapperLocator{bl})
	case in.BootstrapPeerID != "" && in.BootstrapHost != "":
		bl, err := bootstrapper(in.BootstrapPeerID, fmt.Sprintf("%s:%d", in.BootstrapHost, port))
		if err != nil {
			return "", err
		}
		cfg.P2P.V2.DefaultBootstrappers = new([]ocrcommontypes.BootstrapperLocator{bl})
	}

	// Capabilities peering/dispatch
	cfg.Capabilities.Peering.V2.Enabled = new(false)
	cfg.Capabilities.SharedPeering.Enabled = new(true)
	cfg.Capabilities.Dispatcher.SendToSharedPeer = new(true)

	// ExternalRegistry (only if address set)
	if in.CapRegAddress != "" {
		v := in.CapRegVersion
		if v == "" {
			v = "2.0.0"
		}
		cfg.Capabilities.ExternalRegistry = coretoml.ExternalRegistry{
			Address:         new(in.CapRegAddress),
			NetworkID:       new("evm"),
			ChainID:         new(strconv.FormatUint(in.RegistryChainID, 10)),
			ContractVersion: new(v),
		}
	}

	// WorkflowRegistry (non-bootstrap only, if address set)
	if !in.IsBootstrapNode && in.WorkflowRegAddress != "" {
		v := in.WorkflowRegVersion
		if v == "" {
			v = "2.0.0"
		}
		cfg.Capabilities.WorkflowRegistry = coretoml.WorkflowRegistry{
			Address:         new(in.WorkflowRegAddress),
			NetworkID:       new("evm"),
			ChainID:         new(strconv.FormatUint(in.RegistryChainID, 10)),
			ContractVersion: new(v),
			SyncStrategy:    new("reconciliation"),
		}
	}

	// Allowlist
	if len(in.Allowlist) > 0 {
		cfg.Capabilities.Local.RegistryBasedLaunchAllowlist = append([]string(nil), in.Allowlist...)
	}

	// Gateway connector
	switch {
	case in.IsGatewayNode:
		cfg.Capabilities.GatewayConnector = coretoml.GatewayConnector{
			DonID:             new("doesnt-matter-for-gateway-node"),
			ChainIDForNodeKey: new(strconv.FormatUint(in.RegistryChainID, 10)),
		}
		// WorkflowStorage stub (TODO: remove once not required) — nested under WorkflowRegistry.
		cfg.Capabilities.WorkflowRegistry.WorkflowStorage = coretoml.WorkflowStorage{
			URL:                 new("localhost"),
			TLSEnabled:          new(false),
			ArtifactStorageHost: new("localhost"),
		}
	case in.GatewayConnector != nil:
		gc := in.GatewayConnector
		conn := coretoml.GatewayConnector{
			DonID:             new(gc.DonID),
			ChainIDForNodeKey: new(gc.ChainIDForNodeKey),
			NodeAddress:       new(gc.NodeAddress),
		}
		for _, g := range gc.Gateways {
			cg := coretoml.ConnectorGateway{ID: new(g.ID), URL: new(g.URL)}
			if g.DonID != "" {
				cg.DonID = new(g.DonID)
			}
			conn.Gateways = append(conn.Gateways, cg)
		}
		cfg.Capabilities.GatewayConnector = conn
	}

	out, err := toml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal node config: %w", err)
	}
	return string(out), nil
}

// bootstrapper creates a libocr bootstrapper locator from peer ID and address strings.
func bootstrapper(peerID, addr string) (ocrcommontypes.BootstrapperLocator, error) {
	bl, err := ocrcommontypes.NewBootstrapperLocator(peerID, []string{addr})
	if err != nil {
		return ocrcommontypes.BootstrapperLocator{}, fmt.Errorf("invalid bootstrap peer id %q: %w", peerID, err)
	}
	return *bl, nil
}
