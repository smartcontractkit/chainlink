package evm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/mercury_verifier"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// TODO: This probably ought to be split into regular/mercury configpoller
// TODO: Consider using UnpackIntoInterface instead of UnpackIntoMap and take
// the log structs directly from contract e.g.
// mercury_verifier.MercuryVerifierConfigSet

// ConfigSet Common to all OCR2 evm based contracts: https://github.com/smartcontractkit/libocr/blob/master/contract2/dev/OCR2Abstract.sol
var ConfigSet common.Hash

// FeedScopedConfigSet ConfigSet with FeedID for use with mercury (and multi-config DON)
var FeedScopedConfigSet common.Hash

var defaultABI abi.ABI
var verifierABI abi.ABI

const configSetEventName = "ConfigSet"

func init() {
	var err error
	abiPointer, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	defaultABI = *abiPointer
	verifierABI, err = abi.JSON(strings.NewReader(mercury_verifier.MercuryVerifierABI))
	if err != nil {
		panic(err)
	}
	ConfigSet = defaultABI.Events[configSetEventName].ID
	FeedScopedConfigSet = verifierABI.Events[configSetEventName].ID
}

type FullConfigFromLog struct {
	ocrtypes.ContractConfig
	feedID [32]byte
}

func NewContractConfigFromLog(unpacked map[string]interface{}, withOffchainTransmitters bool) (ocrtypes.ContractConfig, error) {
	configDigest, ok := unpacked["configDigest"].([32]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config digest, got %T", unpacked["configDigest"])
	}
	configCount, ok := unpacked["configCount"].(uint64)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config count, got %T", unpacked["configCount"])
	}
	signersAddresses, ok := unpacked["signers"].([]common.Address)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid signers, got %T", unpacked["signers"])
	}
	var transmitters [][]byte
	if withOffchainTransmitters {
		offchainTransmitters, ok := unpacked["offchainTransmitters"].([][32]byte)
		if !ok {
			return ocrtypes.ContractConfig{}, errors.Errorf("invalid offchain transmitters, got %T", unpacked["offchainTransmitters"])
		}
		for _, d := range offchainTransmitters {
			c := d
			transmitters = append(transmitters, c[:])
		}
	} else {
		t, ok := unpacked["transmitters"].([]common.Address)
		if !ok {
			return ocrtypes.ContractConfig{}, errors.Errorf("invalid transmitters, got %T", unpacked["transmitters"])
		}
		for _, d := range t {
			c := d
			transmitters = append(transmitters, c[:])
		}
	}
	f, ok := unpacked["f"].(uint8)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid f, got %T", unpacked["f"])
	}
	onchainConfig, ok := unpacked["onchainConfig"].([]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid onchain config, got %T", unpacked["onchainConfig"])
	}
	offchainConfigVersion, ok := unpacked["offchainConfigVersion"].(uint64)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config digest, got %T", unpacked["offchainConfigVersion"])
	}
	offchainConfig, ok := unpacked["offchainConfig"].([]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid offchainConfig, got %T", unpacked["offchainConfig"])
	}
	var transmitAccounts []ocrtypes.Account
	for _, addr := range transmitters {
		if withOffchainTransmitters {
			transmitAccounts = append(transmitAccounts, ocrtypes.Account(fmt.Sprintf("%x", addr)))
		} else {
			transmitAccounts = append(transmitAccounts, ocrtypes.Account(fmt.Sprintf("0x%x", addr)))
		}
	}
	var signers []ocrtypes.OnchainPublicKey
	for _, addr := range signersAddresses {
		addr := addr
		signers = append(signers, addr[:])
	}
	return ocrtypes.ContractConfig{
		ConfigDigest:          configDigest,
		ConfigCount:           configCount,
		Signers:               signers,
		Transmitters:          transmitAccounts,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, nil
}

func unpackLogData(d []byte, withFeedID bool) (map[string]interface{}, error) {
	var err error
	unpacked := map[string]interface{}{}
	if withFeedID {
		err = verifierABI.UnpackIntoMap(unpacked, configSetEventName, d)
	} else {
		err = defaultABI.UnpackIntoMap(unpacked, configSetEventName, d)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack log data")
	}
	return unpacked, nil
}

func ConfigFromLog(logData []byte) (FullConfigFromLog, error) {
	unpacked, err := unpackLogData(logData, false)
	if err != nil {
		return FullConfigFromLog{}, err
	}
	contractConfig, err := NewContractConfigFromLog(unpacked, false)
	if err != nil {
		return FullConfigFromLog{}, err
	}
	return FullConfigFromLog{
		feedID:         [32]byte{},
		ContractConfig: contractConfig,
	}, nil
}

func ConfigFromLogWithFeedID(logData []byte) (FullConfigFromLog, error) {
	unpacked, err := unpackLogData(logData, true)
	if err != nil {
		return FullConfigFromLog{}, err
	}
	feedID, ok := unpacked["feedId"].([32]byte)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid feed ID, got %T", unpacked["feedId"])
	}
	contractConfig, err := NewContractConfigFromLog(unpacked, true)
	if err != nil {
		return FullConfigFromLog{}, err
	}

	return FullConfigFromLog{
		feedID:         feedID,
		ContractConfig: contractConfig,
	}, nil
}

type ConfigPoller struct {
	lggr               logger.Logger
	filterName         string
	destChainLogPoller logpoller.LogPoller
	addr               common.Address
	feedID             common.Hash
}

type ConfigPollerOption func(cp *ConfigPoller)

func WithFeedID(feedID *common.Hash) ConfigPollerOption {
	return func(cp *ConfigPoller) {
		if feedID != nil {
			cp.feedID = *feedID
		}
	}
}

func NewConfigPoller(lggr logger.Logger, destChainPoller logpoller.LogPoller, addr common.Address, opts ...ConfigPollerOption) (*ConfigPoller, error) {
	configFilterName := logpoller.FilterName("OCR2ConfigPoller", addr.String())

	cp := &ConfigPoller{
		lggr:               lggr,
		filterName:         configFilterName,
		destChainLogPoller: destChainPoller,
		addr:               addr,
	}

	for _, opt := range opts {
		opt(cp)
	}

	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: configFilterName, EventSigs: []common.Hash{cp.ConfigSetEventID()}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}

	return cp, nil
}

func (lp *ConfigPoller) WithFeedID() bool {
	return lp.feedID != (common.Hash{})
}

func (lp *ConfigPoller) ConfigSetEventID() common.Hash {
	if lp.WithFeedID() {
		return FeedScopedConfigSet
	} else {
		return ConfigSet
	}
}

func (lp *ConfigPoller) Notify() <-chan struct{} {
	return nil
}

func (lp *ConfigPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	latest, err := lp.destChainLogPoller.LatestLogByEventSigWithConfs(lp.ConfigSetEventID(), lp.addr, 1, pg.WithParentCtx(ctx))
	if err != nil {
		// If contract is not configured, we will not have the log.
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ocrtypes.ConfigDigest{}, nil
		}
		return 0, ocrtypes.ConfigDigest{}, err
	}
	var latestConfigSet FullConfigFromLog
	if lp.WithFeedID() {
		latestConfigSet, err = ConfigFromLogWithFeedID(latest.Data)
	} else {
		latestConfigSet, err = ConfigFromLog(latest.Data)
	}
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
}

func (lp *ConfigPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	lgs, err := lp.destChainLogPoller.Logs(int64(changedInBlock), int64(changedInBlock), lp.ConfigSetEventID(), lp.addr, pg.WithParentCtx(ctx))
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	var latestConfigSet FullConfigFromLog
	if lp.WithFeedID() {
		latestConfigSet, err = ConfigFromLogWithFeedID(lgs[len(lgs)-1].Data)
	} else {
		latestConfigSet, err = ConfigFromLog(lgs[len(lgs)-1].Data)
	}
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	lp.lggr.Infow("LatestConfig", "latestConfig", latestConfigSet)
	return latestConfigSet.ContractConfig, nil
}

func (lp *ConfigPoller) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latest, err := lp.destChainLogPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latest), nil
}
