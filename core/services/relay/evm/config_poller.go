package evm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/mercury_verifier"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Common to all OCR2 evm based contracts: https://github.com/smartcontractkit/libocr/blob/master/contract2/dev/OCR2Abstract.sol
// event ConfigSet(
//
//	bytes32 feedId,
//	uint32 previousConfigBlockNumber,
//	bytes32 configDigest,
//	uint64 configCount,
//	address[] signers,
//	uint8 f,
//	bytes onchainConfig,
//	uint64 offchainConfigVersion,
//	bytes offchainConfig
//
// );
var ConfigSet = common.HexToHash("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")

// FeedScopedConfigSet ConfigSet with FeedID for use with mercury (and multi-config DON)
var FeedScopedConfigSet common.Hash

func init() {
	abi, err := abi.JSON(strings.NewReader(mercury_verifier.MercuryVerifierABI))
	if err != nil {
		panic(err)
	}
	FeedScopedConfigSet = abi.Events["ConfigSet"].ID
	fmt.Printf("BALLS FeedScopedConfigSet: 0x%x\n", FeedScopedConfigSet)
}

const (
	firstIndexWithoutFeedID = 1
	firstIndexWithFeedID    = 2
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

func makeConfigSetMsgArgs(withFeedID bool) abi.Arguments {
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
	}

	if withFeedID {
		args = append([]abi.Argument{configSetFeedIDArg}, args...)
	} else {
		// We only support `transmitters` when not using feedId
		transmittersArg := abi.Argument{
			Name: "transmitters",
			Type: utils.MustAbiType("address[]", nil),
		}
		args = append(args, transmittersArg)
	}

	lastArgs := []abi.Argument{
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

	return append(args, lastArgs...)
}

type FullConfigFromLog struct {
	ocrtypes.ContractConfig
	feedID [32]byte
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
	var transmitters []common.Address
	// Mercury does not support transmitters
	if fromIndex == firstIndexWithoutFeedID {
		transmitters, ok = unpacked[fromIndex+3].([]common.Address)
		if !ok {
			return ocrtypes.ContractConfig{}, errors.Errorf("invalid transmitters, got %T", unpacked[fromIndex+3])
		}
	} else {
		// We decrease the `fromIndex` by one since we don't have `transmitters` to unpack
		fromIndex -= 1
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

func unpackLogData(d []byte, withFeedID bool) ([]interface{}, error) {
	args := makeConfigSetMsgArgs(withFeedID)
	unpacked, err := args.Unpack(d)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack log data")
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
	contractConfig, err := NewContractConfigFromLog(unpacked, firstIndexWithoutFeedID)
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
	feedID, ok := unpacked[0].([32]byte)
	if !ok {
		return FullConfigFromLog{}, errors.Errorf("invalid feed ID, got %T", unpacked[0])
	}
	contractConfig, err := NewContractConfigFromLog(unpacked, firstIndexWithFeedID)
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
	eventSig           common.Hash
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

	cp.eventSig = ConfigSet
	if cp.WithFeedID() {
		cp.eventSig = FeedScopedConfigSet
	}
	fmt.Println("BALLS listen on address", addr.Hex())
	fmt.Printf("BALLS listen to event sigs %#v %#v\n", ConfigSet, FeedScopedConfigSet)
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: configFilterName, EventSigs: []common.Hash{cp.eventSig}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}

	lggr.Infow("BALLS feed ID", "feedID", cp.feedID)

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
	latest, err := lp.destChainLogPoller.LatestLogByEventSigWithConfs(lp.eventSig, lp.addr, 1, pg.WithParentCtx(ctx))
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
	lgs, err := lp.destChainLogPoller.Logs(int64(changedInBlock), int64(changedInBlock), lp.eventSig, lp.addr, pg.WithParentCtx(ctx))
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
