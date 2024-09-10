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
	Signer     [32]byte // note that in capabilities registry we need a [32]byte, but in the forwarder we need a common.Address [20]byte
	P2PKey     p2pkey.PeerID
	IsBoostrap bool
	// useful when have to register the ocr3 contract config
	p2pKeyBundle   *v1.OCR2Config_P2PKeyBundle
	ocrKeyBundle   *v1.OCR2Config_OCRKeyBundle
	accountAddress string
}

func (o *ocr2Node) signerAddress() common.Address {
	return common.BytesToAddress(o.Signer[:])
}

func (o *ocr2Node) p2pKeyBundleConfig() *v1.OCR2Config_P2PKeyBundle {
	return o.p2pKeyBundle
}

func (o *ocr2Node) ocrKeyBundleConfig() *v1.OCR2Config_OCRKeyBundle {
	return o.ocrKeyBundle
}

func (o *ocr2Node) toNodeKeys() NodeKeys {
	return NodeKeys{
		EthAddress:            o.accountAddress,
		P2PPeerID:             o.p2pKeyBundle.PeerId,
		OCR2BundleID:          o.ocrKeyBundle.BundleId,
		OCR2OnchainPublicKey:  o.ocrKeyBundle.OnchainSigningAddress,
		OCR2OffchainPublicKey: o.ocrKeyBundle.OffchainPublicKey,
		OCR2ConfigPublicKey:   o.ocrKeyBundle.ConfigPublicKey,
	}
}

func newOcr2Node(id string, ccfg *v1.ChainConfig) (*ocr2Node, error) {
	if ccfg == nil {
		return nil, errors.New("nil ocr2config")
	}
	ocfg := ccfg.Ocr2Config
	p := p2pkey.PeerID{}
	if err := p.UnmarshalString(ocfg.P2PKeyBundle.PeerId); err != nil {
		return nil, fmt.Errorf("failed to unmarshal peer id %s: %w", ocfg.P2PKeyBundle.PeerId, err)
	}

	signer := ocfg.OcrKeyBundle.OnchainSigningAddress
	if len(signer) != 40 {
		return nil, fmt.Errorf("invalid onchain signing address %s", ocfg.OcrKeyBundle.OnchainSigningAddress)
	}
	signerB, err := hex.DecodeString(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signer %s: %w", signer, err)
	}

	var sigb [32]byte
	copy(sigb[:], signerB)

	return &ocr2Node{
		ID:             id,
		Signer:         sigb,
		P2PKey:         p,
		IsBoostrap:     ocfg.IsBootstrap,
		p2pKeyBundle:   ocfg.P2PKeyBundle,
		ocrKeyBundle:   ocfg.OcrKeyBundle,
		accountAddress: ccfg.AccountAddress,
	}, nil
}

func makeNodeKeysSlice(nodes []*ocr2Node) []NodeKeys {
	var out []NodeKeys
	for _, n := range nodes {
		out = append(out, n.toNodeKeys())
	}
	return out
}
