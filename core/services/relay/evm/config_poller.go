package evm

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
)

// Common to all OCR2 evm based contracts: https://github.com/smartcontractkit/libocr/blob/master/contract2/OCR2Abstract.sol#L23
var ConfigSet = common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")

const (
	firstIndexWithoutFeedId = 1
	firstIndexWithFeedId    = 2
)

type OCR2AbstractConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigDigest              [32]byte
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	F                         uint8
	OnchainConfig             []byte
	OffchainConfigVersion     uint64
	OffchainConfig            []byte
}

var configSetFeedIDArg = abi.Argument{
	Name: "feedId",
	Type: utils.MustAbiType("bytes32", nil),
}

func makeConfigSetMsgArgs(withFeedId bool) abi.Arguments {
	args := []abi.Argument{
		{
			Name: "previousConfigBlockNumber",
			Type: utils.MustAbiType("uint32", nil),
		},
		{
			Name: "configDigest",
			Type: utils.MustAbiType("bytes32", nil),
		},
		{
			Name: "configCount",
			Type: utils.MustAbiType("uint64", nil),
		},
		{
			Name: "signers",
			Type: utils.MustAbiType("address[]", nil),
		},
		{
			Name: "transmitters",
			Type: utils.MustAbiType("address[]", nil),
		},
		{
			Name: "f",
			Type: utils.MustAbiType("uint8", nil),
		},
		{
			Name: "onchainConfig",
			Type: utils.MustAbiType("bytes", nil),
		},
		{
			Name: "offchainConfigVersion",
			Type: utils.MustAbiType("uint64", nil),
		},
		{
			Name: "offchainConfig",
			Type: utils.MustAbiType("bytes", nil),
		},
	}

	if withFeedId {
		args = append([]abi.Argument{configSetFeedIDArg}, args...)
	}

	return args
}

type FullConfigFromLog struct {
	ocrtypes.ContractConfig
	FeedId [32]byte
}

func NewContractConfigFromLog(unpacked []interface{}, fromIndex int) (ocrtypes.ContractConfig, error) {
	configDigest, ok := unpacked[fromIndex].([32]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config digest, got %T", unpacked[fromIndex])
	}
	configCount, ok := unpacked[fromIndex+1].(uint64)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config count, got %T", unpacked[fromIndex+1])
	}
	signersAddresses, ok := unpacked[fromIndex+2].([]common.Address)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid signers, got %T", unpacked[fromIndex+2])
	}
	transmitters, ok := unpacked[fromIndex+3].([]common.Address)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid transmitters, got %T", unpacked[fromIndex+3])
	}
	f, ok := unpacked[fromIndex+4].(uint8)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid f, got %T", unpacked[fromIndex+4])
	}
	onchainConfig, ok := unpacked[fromIndex+5].([]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid onchain config, got %T", unpacked[fromIndex+5])
	}
	offchainConfigVersion, ok := unpacked[fromIndex+6].(uint64)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config digest, got %T", unpacked[fromIndex+6])
	}
	offchainConfig, ok := unpacked[fromIndex+7].([]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid offchainConfig, got %T", unpacked[fromIndex+7])
	}
	var transmitAccounts []ocrtypes.Account
	for _, addr := range transmitters {
		transmitAccounts = append(transmitAccounts, ocrtypes.Account(addr.Hex()))
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

func unpackLogData(d []byte, withFeedId bool) ([]interface{}, error) {
	args := makeConfigSetMsgArgs(withFeedId)
	unpacked, err := args.Unpack(d)
	if err != nil {
		return nil, err
	}
	if len(unpacked) != len(args) {
		return nil, errors.Errorf("invalid number of fields, got %v", len(unpacked))
	}
	return unpacked, nil
}

func ConfigFromLog(logData []byte) (FullConfigFromLog, error) {
	unpacked, err := unpackLogData(logData, false)
	if err != nil {
		return FullConfigFromLog{}, err
	}
	contractConfig, err := NewContractConfigFromLog(unpacked, firstIndexWithoutFeedId)
	if err != nil {
		return FullConfigFromLog{}, err
	}
	return FullConfigFromLog{
		FeedId:         [32]byte{},
		ContractConfig: contractConfig,
	}, nil
}

func ConfigFromLogWithFeedId(logData []byte) (FullConfigFromLog, error) {
	unpacked, err := unpackLogData(logData, true)
	if err != nil {
		return FullConfigFromLog{}, err
	}
	feedId, ok := unpacked[0].([32]byte)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid feed ID, got %T", unpacked[0])
	}
	contractConfig, err := NewContractConfigFromLog(unpacked, firstIndexWithFeedId)
	if err != nil {
		return FullConfigFromLog{}, err
	}

	return FullConfigFromLog{
		FeedId:         feedId,
		ContractConfig: contractConfig,
	}, nil
}

type ConfigPoller struct {
	lggr               logger.Logger
	filterName         string
	destChainLogPoller logpoller.LogPoller
	addr               common.Address
	feedId             common.Hash
}

type ConfigPollerOption func(cp *ConfigPoller)

func WithFeedId(feedId common.Hash) ConfigPollerOption {
	return func(cp *ConfigPoller) {
		cp.feedId = feedId
	}
}

func NewConfigPoller(lggr logger.Logger, destChainPoller logpoller.LogPoller, addr common.Address, opts ...ConfigPollerOption) (*ConfigPoller, error) {
	configFilterName := logpoller.FilterName("OCR2ConfigPoller", addr.String())
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: configFilterName, EventSigs: []common.Hash{ConfigSet}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}

	cp := &ConfigPoller{
		lggr:               lggr,
		filterName:         configFilterName,
		destChainLogPoller: destChainPoller,
		addr:               addr,
	}

	for _, opt := range opts {
		opt(cp)
	}

	return cp, nil
}

func (lp *ConfigPoller) WithFeedId() bool {
	return lp.feedId != (common.Hash{})
}

func (lp *ConfigPoller) Notify() <-chan struct{} {
	return nil
}

func (lp *ConfigPoller) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	latest, err := lp.destChainLogPoller.LatestLogByEventSigWithConfs(ConfigSet, lp.addr, 1, pg.WithParentCtx(ctx))
	if err != nil {
		// If contract is not configured, we will not have the log.
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ocrtypes.ConfigDigest{}, nil
		}
		return 0, ocrtypes.ConfigDigest{}, err
	}
	var latestConfigSet FullConfigFromLog
	if lp.WithFeedId() {
		latestConfigSet, err = ConfigFromLogWithFeedId(latest.Data)
	} else {
		latestConfigSet, err = ConfigFromLog(latest.Data)
	}
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}
	return uint64(latest.BlockNumber), latestConfigSet.ConfigDigest, nil
}

func (lp *ConfigPoller) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	lgs, err := lp.destChainLogPoller.Logs(int64(changedInBlock), int64(changedInBlock), ConfigSet, lp.addr, pg.WithParentCtx(ctx))
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	var latestConfigSet FullConfigFromLog
	if lp.WithFeedId() {
		latestConfigSet, err = ConfigFromLogWithFeedId(lgs[len(lgs)-1].Data)
	} else {
		latestConfigSet, err = ConfigFromLog(lgs[len(lgs)-1].Data)
	}
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	lp.lggr.Infof("LatestConfig %+v\n", latestConfigSet)
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
