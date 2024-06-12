package evm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipcommit"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipexec"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// CCIPCommitProvider provides all components needed for a CCIP Relay OCR2 plugin.
type CCIPCommitProvider interface {
	commontypes.Plugin
}

// CCIPExecutionProvider provides all components needed for a CCIP Execution OCR2 plugin.
type CCIPExecutionProvider interface {
	commontypes.Plugin
}

type ccipCommitProvider struct {
	*configWatcher
	contractTransmitter *contractTransmitter
}

func chainToUUID(chainID *big.Int) uuid.UUID {
	// See https://www.rfc-editor.org/rfc/rfc4122.html#section-4.1.3 for the list of supported versions.
	const VersionSHA1 = 5
	var buf bytes.Buffer
	buf.WriteString("CCIP:")
	buf.Write(chainID.Bytes())
	// We use SHA-256 instead of SHA-1 because the former has better collision resistance.
	// The UUID will contain only the first 16 bytes of the hash.
	// You can't say which algorithms was used just by looking at the UUID bytes.
	return uuid.NewHash(sha256.New(), uuid.NameSpaceOID, buf.Bytes(), VersionSHA1)
}

func NewCCIPCommitProvider(ctx context.Context, lggr logger.Logger, chainSet legacyevm.Chain, rargs commontypes.RelayArgs, transmitterID string, ks keystore.Eth) (CCIPCommitProvider, error) {
	relayOpts := types.NewRelayOpts(rargs)
	configWatcher, err := newStandardConfigProvider(ctx, lggr, chainSet, relayOpts)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(relayOpts.ContractID)
	typ, ver, err := ccipconfig.TypeAndVersion(address, chainSet.Client())
	if err != nil {
		return nil, err
	}
	fn, err := ccipcommit.CommitReportToEthTxMeta(typ, ver)
	if err != nil {
		return nil, err
	}
	subjectID := chainToUUID(configWatcher.chain.ID())
	contractTransmitter, err := newOnChainContractTransmitter(ctx, lggr, rargs, transmitterID, ks, configWatcher, configTransmitterOpts{
		subjectID: &subjectID,
	}, OCR2AggregatorTransmissionContractABI, fn, 0)
	if err != nil {
		return nil, err
	}
	return &ccipCommitProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

func (c *ccipCommitProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}

func (c *ccipCommitProvider) ChainReader() commontypes.ContractReader {
	return nil
}

func (c *ccipCommitProvider) Codec() commontypes.Codec {
	return nil
}

type ccipExecutionProvider struct {
	*configWatcher
	contractTransmitter *contractTransmitter
}

var _ commontypes.Plugin = (*ccipExecutionProvider)(nil)

func NewCCIPExecutionProvider(ctx context.Context, lggr logger.Logger, chainSet legacyevm.Chain, rargs commontypes.RelayArgs, transmitterID string, ks keystore.Eth) (CCIPExecutionProvider, error) {
	relayOpts := types.NewRelayOpts(rargs)

	configWatcher, err := newStandardConfigProvider(ctx, lggr, chainSet, relayOpts)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(relayOpts.ContractID)
	typ, ver, err := ccipconfig.TypeAndVersion(address, chainSet.Client())
	if err != nil {
		return nil, err
	}
	fn, err := ccipexec.ExecReportToEthTxMeta(ctx, typ, ver)
	if err != nil {
		return nil, err
	}
	subjectID := chainToUUID(configWatcher.chain.ID())
	contractTransmitter, err := newOnChainContractTransmitter(ctx, lggr, rargs, transmitterID, ks, configWatcher, configTransmitterOpts{
		subjectID: &subjectID,
	}, OCR2AggregatorTransmissionContractABI, fn, 0)
	if err != nil {
		return nil, err
	}
	return &ccipExecutionProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

func (c *ccipExecutionProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}

func (c *ccipExecutionProvider) ChainReader() commontypes.ContractReader {
	return nil
}

func (c *ccipExecutionProvider) Codec() commontypes.Codec {
	return nil
}
