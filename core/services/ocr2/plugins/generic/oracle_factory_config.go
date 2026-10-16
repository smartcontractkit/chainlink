package generic

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"math/big"
	"slices"

	gethCommon "github.com/ethereum/go-ethereum/common"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
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
func ResolveOracleFactoryConfig(ctx context.Context, params ResolveOracleFactoryConfigParams) (job.OracleFactoryConfig, job.OnchainSigningStrategy, error) {
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
		for family, kb := range params.OCRKeyBundles {
			signing.Config[family] = kb.ID()
		}
	}

	// Prefer the transmitter extracted from the on-chain OCR config by the caller.
	if cfg.TransmitterID == "" && params.Transmitter != "" {
		cfg.TransmitterID = params.Transmitter
	}

	// Fall back to a round-robin keystore address when no transmitter was resolved
	// from the on-chain OCR config.
	if cfg.TransmitterID == "" && params.EthKeystore != nil && cfg.ChainID != "" {
		transmitter, err := DefaultTransmitterForChain(ctx, params.EthKeystore, cfg.ChainID)
		if err != nil {
			return cfg, signing, fmt.Errorf("failed to resolve transmitter: %w", err)
		}
		cfg.TransmitterID = transmitter
	}

	return cfg, signing, nil
}

// SelectOCRKeyBundlesForConfig locates this node's signer in the on-chain OCR config and
// returns the key bundles, keyed by chain family name, that make up that signer, along
// with the signer's index in cc.Signers. A signer is either a raw EVM public key or a
// multichain public key covering one or more chain families (e.g. EVM + Solana + Aptos);
// the node matches only if it holds a bundle for every family in the signer. Returns
// false when no config is provided or no signer matches.
func SelectOCRKeyBundlesForConfig(bundles []ocr2key.KeyBundle, cc *ocrtypes.ContractConfig) (map[string]ocr2key.KeyBundle, int, bool) {
	if cc == nil {
		return nil, -1, false
	}
	for i, s := range cc.Signers {
		if matched, ok := matchSigner(bundles, s); ok {
			return matched, i, true
		}
	}
	return nil, -1, false
}

func matchSigner(bundles []ocr2key.KeyBundle, signer ocrtypes.OnchainPublicKey) (map[string]ocr2key.KeyBundle, bool) {
	if len(signer) == 0 {
		return nil, false
	}
	for _, kb := range bundles {
		if kb.ChainType() == corekeys.EVM && bytes.Equal(signer, kb.PublicKey()) {
			return map[string]ocr2key.KeyBundle{string(corekeys.EVM): kb}, true
		}
	}

	pubKeys, err := ocrcommon.UnmarshalMultichainPublicKey(signer)
	if err != nil || len(pubKeys) == 0 {
		return nil, false
	}
	matched := make(map[string]ocr2key.KeyBundle, len(pubKeys))
	for family, pub := range pubKeys {
		for _, kb := range bundles {
			if string(kb.ChainType()) == family && bytes.Equal(pub, kb.PublicKey()) {
				matched[family] = kb
				break
			}
		}
		if _, ok := matched[family]; !ok {
			return nil, false
		}
	}
	// Unmarshalling skips unknown families, so re-encode and compare to make sure the
	// matched bundles produce exactly the signer in the config.
	encoded, err := ocrcommon.MarshalMultichainKeyBundle(matched)
	if err != nil || !bytes.Equal(encoded, signer) {
		return nil, false
	}
	return matched, true
}

// OffchainKeyBundle picks the bundle whose offchain keys are used for OCR from the
// bundles matching this node's signer. The EVM bundle is preferred; otherwise the
// bundle of the lexicographically first family is used for determinism.
func OffchainKeyBundle(kbs map[string]ocr2key.KeyBundle) ocr2key.KeyBundle {
	if kb, ok := kbs[string(corekeys.EVM)]; ok {
		return kb
	}
	families := slices.Sorted(maps.Keys(kbs))
	if len(families) == 0 {
		return nil
	}
	return kbs[families[0]]
}

// TransmitterAt returns the transmitter account at index i of the on-chain OCR config,
// i.e. the transmitter paired with Signers[i]. EVM addresses are checksummed.
func TransmitterAt(cc ocrtypes.ContractConfig, i int) (string, bool) {
	if i < 0 || i >= len(cc.Transmitters) {
		return "", false
	}
	transmitter := string(cc.Transmitters[i])
	if gethCommon.IsHexAddress(transmitter) {
		transmitter = gethCommon.HexToAddress(transmitter).Hex()
	}
	return transmitter, true
}

func DefaultTransmitterForChain(ctx context.Context, ethKS keystore.Eth, chainID string) (string, error) {
	chainIDBig, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		return "", fmt.Errorf("invalid chain_id %q", chainID)
	}

	addr, err := ethKS.GetRoundRobinAddress(ctx, chainIDBig)
	if err != nil {
		return "", fmt.Errorf("failed to get transmitter for chain_id %s: %w", chainID, err)
	}

	return addr.String(), nil
}
