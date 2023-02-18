package evm

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"

	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
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

func makeConfigSetMsgArgs() abi.Arguments {
	return []abi.Argument{
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
}

func ConfigFromLog(logData []byte) (ocrtypes.ContractConfig, error) {
	unpacked, err := makeConfigSetMsgArgs().Unpack(logData)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(unpacked) != 9 {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid number of fields, got %v", len(unpacked))
	}
	configDigest, ok := unpacked[1].([32]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config digest, got %T", unpacked[1])
	}
	configCount, ok := unpacked[2].(uint64)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config count, got %T", unpacked[2])
	}
	signersAddresses, ok := unpacked[3].([]common.Address)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid signers, got %T", unpacked[3])
	}
	transmitters, ok := unpacked[4].([]common.Address)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid transmitters, got %T", unpacked[4])
	}
	f, ok := unpacked[5].(uint8)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid f, got %T", unpacked[5])
	}
	onchainConfig, ok := unpacked[6].([]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid onchain config, got %T", unpacked[6])
	}
	offchainConfigVersion, ok := unpacked[7].(uint64)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid config digest, got %T", unpacked[7])
	}
	offchainConfig, ok := unpacked[8].([]byte)
	if !ok {
		return ocrtypes.ContractConfig{}, errors.Errorf("invalid offchainConfig, got %T", unpacked[8])
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

type ConfigPoller struct {
	lggr               logger.Logger
	filterName         string
	destChainLogPoller logpoller.LogPoller
	addr               common.Address
}

func NewConfigPoller(lggr logger.Logger, destChainPoller logpoller.LogPoller, addr common.Address) (*ConfigPoller, error) {
	configFilterName := logpoller.FilterName("OCR2ConfigPoller", addr.String())
	err := destChainPoller.RegisterFilter(logpoller.Filter{Name: configFilterName, EventSigs: []common.Hash{ConfigSet}, Addresses: []common.Address{addr}})
	if err != nil {
		return nil, err
	}
	return &ConfigPoller{
		lggr:               lggr,
		filterName:         configFilterName,
		destChainLogPoller: destChainPoller,
		addr:               addr,
	}, nil
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
	latestConfigSet, err := ConfigFromLog(latest.Data)
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
	latestConfigSet, err := ConfigFromLog(lgs[len(lgs)-1].Data)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	lp.lggr.Infof("LatestConfig %+v\n", latestConfigSet)
	return latestConfigSet, nil
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
