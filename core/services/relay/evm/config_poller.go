package evm

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Common to all OCR2 evm based contracts: https://github.com/smartcontractkit/libocr/blob/master/contract2/OCR2Abstract.sol#L23
var ConfigSet = common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")

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

type FullConfigFromLog struct {
	ocrtypes.ContractConfig
	FeedId [32]byte
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

func ConfigFromLog(logData []byte, withFeedId bool) (FullConfigFromLog, error) {
	args := makeConfigSetMsgArgs(withFeedId)
	unpacked, err := args.Unpack(logData)
	if err != nil {
		return FullConfigFromLog{}, err
	}
	if len(unpacked) != len(args) {
		return FullConfigFromLog{}, errors.Errorf("invalid number of fields, got %v", len(unpacked))
	}
	initialIndex := 1
	if withFeedId {
		initialIndex = 2
	}
	configDigest, ok := unpacked[initialIndex].([32]byte)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid config digest, got %T", unpacked[initialIndex])
	}
	configCount, ok := unpacked[initialIndex+1].(uint64)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid config count, got %T", unpacked[initialIndex+1])
	}
	signersAddresses, ok := unpacked[initialIndex+2].([]common.Address)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid signers, got %T", unpacked[initialIndex+2])
	}
	transmitters, ok := unpacked[initialIndex+3].([]common.Address)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid transmitters, got %T", unpacked[initialIndex+3])
	}
	f, ok := unpacked[initialIndex+4].(uint8)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid f, got %T", unpacked[initialIndex+4])
	}
	onchainConfig, ok := unpacked[initialIndex+5].([]byte)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid onchain config, got %T", unpacked[initialIndex+5])
	}
	offchainConfigVersion, ok := unpacked[initialIndex+6].(uint64)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid config digest, got %T", unpacked[initialIndex+6])
	}
	offchainConfig, ok := unpacked[initialIndex+7].([]byte)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid offchainConfig, got %T", unpacked[initialIndex+7])
	}
	var feedId [32]byte
	if withFeedId {
		feedId, ok = unpacked[0].([32]byte)
		if !ok {
			return FullConfigFromLog{}, errors.Errorf("invalid feed ID, got %T", unpacked[0])
		}
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
	return FullConfigFromLog{
		FeedId: feedId,
		ContractConfig: ocrtypes.ContractConfig{
			ConfigDigest:          configDigest,
			ConfigCount:           configCount,
			Signers:               signers,
			Transmitters:          transmitAccounts,
			F:                     f,
			OnchainConfig:         onchainConfig,
			OffchainConfigVersion: offchainConfigVersion,
			OffchainConfig:        offchainConfig,
		},
	}, nil
}

type ConfigPoller struct {
	lggr               logger.Logger
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
	_, err := destChainPoller.RegisterFilter(logpoller.Filter{EventSigs: []common.Hash{ConfigSet}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}
	cp := &ConfigPoller{
		lggr:               lggr,
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
	latestConfigSet, err := ConfigFromLog(latest.Data, lp.WithFeedId())
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
	latestConfigSet, err := ConfigFromLog(lgs[len(lgs)-1].Data, lp.WithFeedId())
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
