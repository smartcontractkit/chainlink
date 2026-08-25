package generic

import (
	"bytes"
	"context"

	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

type ResolveOracleFactoryConfigParams struct {
	Config job.OracleFactoryConfig
	// OnchainSigning is the signing strategy from the job spec. When empty, it is
	// populated from the node's OCR key bundles keyed by chain family name.
	OnchainSigning job.OnchainSigningStrategy
	// CapRegistryAddress and CapRegistryChainID are the Capabilities Registry
	// contract address and its chain ID. The oracle factory reads its OCR config
	// (signers/transmitters) from the registry, so the contract/chain default to
	// the registry's own address/chain when not set in the job spec.
	CapRegistryAddress string
	CapRegistryChainID string
	// OCRKeyBundles maps chain family name (e.g. "evm") to the node's OCR key bundle
	// for that family. Used to fill in the signing strategy config when the job spec
	// does not provide one.
	OCRKeyBundles map[string]ocr2key.KeyBundle
	// Transmitter is the transmitter address extracted from the on-chain OCR config,
	// paired with this node's signer. When non-empty, it is used as the transmitter
	// when the job spec does not set one. Empty when no registry config is available.
	Transmitter string
	EthKeystore keystore.Eth
	Logger      common.Logger
}

// ResolveOracleFactoryConfig fills missing oracle factory fields. Contract address
// and chain ID default to the Capabilities Registry's address/chain. The signing
// config defaults to the node's OCR key bundles. The transmitter is taken from the
// Transmitter param (extracted from the on-chain OCR config by the caller) when
// available, falling back to a round-robin keystore address otherwise. Job spec
// values take precedence when set.
func ResolveOracleFactoryConfig(_ context.Context, params ResolveOracleFactoryConfigParams) (job.OracleFactoryConfig, job.OnchainSigningStrategy, error) {
	cfg := params.Config
	signing := params.OnchainSigning

	// Nothing to resolve when this job does not use the oracle factory. Resolving would
	// otherwise force a transmitter/key lookup for jobs that never build an oracle.
	if !cfg.Enabled {
		return cfg, signing, nil
	}

	if cfg.OCRContractAddress == "" {
		cfg.OCRContractAddress = params.CapRegistryAddress
	}
	if cfg.ChainID == "" {
		cfg.ChainID = params.CapRegistryChainID
	}

	// TODO: support other chain families besides EVM.
	if cfg.OCRKeyBundleID == "" {
		if kb, ok := params.OCRKeyBundles["evm"]; ok {
			cfg.OCRKeyBundleID = kb.ID()
		}
	}

	if len(signing.Config) == 0 && len(params.OCRKeyBundles) > 0 {
		if signing.StrategyName == "" {
			signing.StrategyName = "multi-chain"
		}
		signing.Config = make(map[string]string)
		// TODO: support other chain families besides EVM.
		for family, kb := range params.OCRKeyBundles {
			signing.Config[family] = kb.ID()
		}
	}

	// Prefer the transmitter extracted from the on-chain OCR config by the caller.
	if cfg.TransmitterID == "" && params.Transmitter != "" {
		cfg.TransmitterID = params.Transmitter
	}

	return cfg, signing, nil
}

// SelectOCRKeyBundleForConfig returns the key bundle whose signer public key appears in
// the on-chain OCR config. It disambiguates when a node holds multiple EVM OCR key
// bundles by picking the one the registry registered as a signer. Returns false when no
// config is provided or none of the bundles match.
func SelectOCRKeyBundleForConfig(bundles []ocr2key.KeyBundle, cc *ocrtypes.ContractConfig) (ocr2key.KeyBundle, bool) {
	if cc == nil {
		return nil, false
	}
	for _, kb := range bundles {
		pub := kb.PublicKey()
		for _, s := range cc.Signers {
			if bytes.Equal(s, pub) {
				return kb, true
			}
		}
	}
	return nil, false
}

// TransmitterForSigner returns the transmitter account paired with the given signer in
// the on-chain OCR config. Signers[i] and Transmitters[i] describe the same oracle, so
// locating this node's signer yields the transmitter the OCR config expects for it.
func TransmitterForSigner(cc ocrtypes.ContractConfig, signer ocrtypes.OnchainPublicKey) (string, bool) {
	for i, s := range cc.Signers {
		if bytes.Equal(s, signer) {
			if i < len(cc.Transmitters) {
				return string(cc.Transmitters[i]), true
			}
			return "", false
		}
	}
	return "", false
}
