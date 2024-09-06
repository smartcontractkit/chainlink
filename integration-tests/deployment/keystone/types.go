package keystone

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	v1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/node/v1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var (
	CapabilitiesRegistry deployment.ContractType = "CapabilitiesRegistry"
	KeystoneForwarder    deployment.ContractType = "KeystoneForwarder"
	OCR3Capability       deployment.ContractType = "OCR3Capability"
)

type deployResponse struct {
	Address common.Address
	Tx      common.Hash // todo: chain agnostic
	Tv      deployment.TypeAndVersion
}

type deployRequest struct {
	Chain deployment.Chain
}

type DonNode struct {
	Don  string
	Node string // not unique across environments
}

type CapabilityHost struct {
	NodeID       string // globally unique
	Capabilities []capabilities_registry.CapabilitiesRegistryCapability
}

type Nop struct {
	capabilities_registry.CapabilitiesRegistryNodeOperator
	NodeIDs []string // nodes run by this operator
}

// ocr2Node is a subset of the node configuration that is needed to register a node
// with the capabilities registry. Signer and P2PKey are chain agnostic.
type ocr2Node struct {
	ID         string
	Signer     [32]byte
	P2PKey     p2pkey.PeerID
	IsBoostrap bool
}

func newOcr2Node(id string, cfg *v1.OCR2Config) (*ocr2Node, error) {
	if cfg == nil {
		return nil, errors.New("nil ocr2config")
	}
	p := p2pkey.PeerID{}
	if err := p.UnmarshalString(cfg.P2PKeyBundle.PeerId); err != nil {
		return nil, fmt.Errorf("failed to unmarshal peer id %s: %w", cfg.P2PKeyBundle.PeerId, err)
	}

	signer := cfg.OcrKeyBundle.OnchainSigningAddress
	if len(signer) != 40 {
		return nil, fmt.Errorf("invalid onchain signing address %s", cfg.OcrKeyBundle.OnchainSigningAddress)
	}
	signerB, err := hex.DecodeString(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signer %s: %w", signer, err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return &ocr2Node{
		ID:         id,
		Signer:     sigb,
		P2PKey:     p,
		IsBoostrap: cfg.IsBootstrap,
	}, nil
}
