package ccip

import (
	"context"
	"encoding/hex"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipevents"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	MaxQueryLength       = 0       // empty for both plugins
	MaxObservationLength = 250_000 // plugins's Observation should make sure to cap to this limit
	CommitPluginLabel    = "commit"
	ExecPluginLabel      = "exec"
)

var zeroAddress = common.HexToAddress("0")

var ErrCommitStoreIsDown = errors.New("commitStore is down")

func LoadOnRamp(onRampAddress common.Address, pluginName string, client client.Client) (evm_2_evm_onramp.EVM2EVMOnRampInterface, error) {
	err := ccipconfig.VerifyTypeAndVersion(onRampAddress, client, ccipconfig.EVM2EVMOnRamp)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid onRamp contract")
	}
	return observability.NewObservedEvm2EvmOnRamp(onRampAddress, pluginName, client)
}

func LoadOnRampDynamicConfig(onRamp evm_2_evm_onramp.EVM2EVMOnRampInterface, client client.Client) (evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig, error) {
	versionString, err := onRamp.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
	}

	_, version, err := ccipconfig.ParseTypeAndVersion(versionString)
	if err != nil {
		return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
	}

	opts := &bind.CallOpts{}

	switch version {
	case "1.0.0":
		legacyOnramp, err := evm_2_evm_onramp_1_0_0.NewEVM2EVMOnRamp(onRamp.Address(), client)
		if err != nil {
			return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
		}
		legacyDynamicConfig, err := legacyOnramp.GetDynamicConfig(opts)
		if err != nil {
			return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
		}
		return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:          legacyDynamicConfig.Router,
			MaxTokensLength: legacyDynamicConfig.MaxTokensLength,
			PriceRegistry:   legacyDynamicConfig.PriceRegistry,
			MaxDataSize:     legacyDynamicConfig.MaxDataSize,
			MaxGasLimit:     uint32(legacyDynamicConfig.MaxGasLimit),
		}, nil
	case "1.1.0":
		legacyOnramp, err := evm_2_evm_onramp_1_1_0.NewEVM2EVMOnRamp(onRamp.Address(), client)
		if err != nil {
			return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
		}
		legacyDynamicConfig, err := legacyOnramp.GetDynamicConfig(opts)
		if err != nil {
			return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, err
		}
		return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                legacyDynamicConfig.Router,
			MaxTokensLength:       legacyDynamicConfig.MaxTokensLength,
			DestGasOverhead:       legacyDynamicConfig.DestGasOverhead,
			DestGasPerPayloadByte: legacyDynamicConfig.DestGasPerPayloadByte,
			PriceRegistry:         legacyDynamicConfig.PriceRegistry,
			MaxDataSize:           legacyDynamicConfig.MaxDataSize,
			MaxGasLimit:           uint32(legacyDynamicConfig.MaxGasLimit),
		}, nil
	case "1.2.0":
		return onRamp.GetDynamicConfig(opts)
	default:
		return evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{}, errors.Errorf("Invalid onramp version: %s", version)
	}
}

func LoadOffRamp(offRampAddress common.Address, pluginName string, client client.Client) (evm_2_evm_offramp.EVM2EVMOffRampInterface, error) {
	err := ccipconfig.VerifyTypeAndVersion(offRampAddress, client, ccipconfig.EVM2EVMOffRamp)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid offRamp contract")
	}
	return observability.NewObservedEvm2EvmOffRamp(offRampAddress, pluginName, client)
}

func LoadCommitStore(commitStoreAddress common.Address, pluginName string, client client.Client) (commit_store.CommitStoreInterface, error) {
	err := ccipconfig.VerifyTypeAndVersion(commitStoreAddress, client, ccipconfig.CommitStore)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid commitStore contract")
	}
	return observability.NewObservedCommitStore(commitStoreAddress, pluginName, client)
}

func contiguousReqs(lggr logger.Logger, min, max uint64, seqNrs []uint64) bool {
	if int(max-min+1) != len(seqNrs) {
		return false
	}
	for i, j := min, 0; i <= max && j < len(seqNrs); i, j = i+1, j+1 {
		if seqNrs[j] != i {
			lggr.Errorw("unexpected gap in seq nums", "seqNr", i, "minSeqNr", min, "maxSeqNr", max)
			return false
		}
	}
	return true
}

func calculateUsdPerUnitGas(sourceGasPrice *big.Int, usdPerFeeCoin *big.Int) *big.Int {
	// (wei / gas) * (usd / eth) * (1 eth / 1e18 wei)  = usd/gas
	tmp := new(big.Int).Mul(sourceGasPrice, usdPerFeeCoin)
	return tmp.Div(tmp, big.NewInt(1e18))
}

// Extracts the hashed leaves from a given set of logs
func leavesFromIntervals(
	lggr logger.Logger,
	interval commit_store.CommitStoreInterval,
	hasher hashlib.LeafHasherInterface[[32]byte],
	sendReqs []ccipevents.Event[evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested],
) ([][32]byte, error) {
	var seqNrs []uint64
	for _, req := range sendReqs {
		seqNrs = append(seqNrs, req.Data.Message.SequenceNumber)
	}

	if !contiguousReqs(lggr, interval.Min, interval.Max, seqNrs) {
		return nil, errors.Errorf("do not have full range [%v, %v] have %v", interval.Min, interval.Max, seqNrs)
	}
	var leaves [][32]byte

	for _, sendReq := range sendReqs {
		hash, err2 := hasher.HashLeaf(sendReq.Data.Raw)
		if err2 != nil {
			return nil, err2
		}
		leaves = append(leaves, hash)
	}

	return leaves, nil
}

// Checks whether the commit store is down by doing an onchain check for Paused and ARM status
func isCommitStoreDownNow(ctx context.Context, lggr logger.Logger, commitStore commit_store.CommitStoreInterface) bool {
	unPausedAndHealthy, err := commitStore.IsUnpausedAndARMHealthy(&bind.CallOpts{Context: ctx})
	if err != nil {
		// If we cannot read the state, assume the worst
		lggr.Errorw("Unable to read CommitStore IsUnpausedAndARMHealthy", "err", err)
		return true
	}
	return !unPausedAndHealthy
}

func filterContainsZeroAddress(addrs []common.Address) bool {
	for _, addr := range addrs {
		if addr == zeroAddress {
			return true
		}
	}
	return false
}

func registerLpFilters(lp logpoller.LogPoller, filters []logpoller.Filter, qopts ...pg.QOpt) error {
	for _, lpFilter := range filters {
		if filterContainsZeroAddress(lpFilter.Addresses) {
			continue
		}
		if err := lp.RegisterFilter(lpFilter, qopts...); err != nil {
			return err
		}
	}
	return nil
}

func unregisterLpFilters(lp logpoller.LogPoller, filters []logpoller.Filter, qopts ...pg.QOpt) error {
	for _, lpFilter := range filters {
		if filterContainsZeroAddress(lpFilter.Addresses) {
			continue
		}
		if err := lp.UnregisterFilter(lpFilter.Name, qopts...); err != nil {
			return err
		}
	}
	return nil
}

func containsFilter(filters []logpoller.Filter, f logpoller.Filter) bool {
	for _, existing := range filters {
		if existing.Name == f.Name {
			return true
		}
	}
	return false
}

func filtersDiff(filtersBefore, filtersNow []logpoller.Filter) (created, deleted []logpoller.Filter) {
	created = make([]logpoller.Filter, 0, len(filtersNow))
	deleted = make([]logpoller.Filter, 0, len(filtersBefore))

	for _, f := range filtersNow {
		if !containsFilter(filtersBefore, f) {
			created = append(created, f)
		}
	}

	for _, f := range filtersBefore {
		if !containsFilter(filtersNow, f) {
			deleted = append(deleted, f)
		}
	}

	return created, deleted
}

func max[T constraints.Ordered](first T, rest ...T) T {
	max := first
	for _, v := range rest {
		if v > max {
			max = v
		}
	}
	return max
}

func median(vals []*big.Int) *big.Int {
	valsCopy := make([]*big.Int, len(vals))
	copy(valsCopy[:], vals[:])
	sort.Slice(valsCopy, func(i, j int) bool {
		return valsCopy[i].Cmp(valsCopy[j]) == -1
	})
	return valsCopy[len(valsCopy)/2]
}

// deviation_parts_per_billion = ((x2 - x1) / x1) * 1e9
func deviates(x1, x2 *big.Int, ppb int64) bool {
	// if x1 == 0, deviates if x2 != x1, to avoid the relative division by 0 error
	if x1.BitLen() == 0 {
		return x1.Cmp(x2) != 0
	}
	diff := big.NewInt(0).Sub(x1, x2)
	diff.Mul(diff, big.NewInt(1e9))
	diff.Div(diff, x1)
	return diff.CmpAbs(big.NewInt(ppb)) > 0
}

func bytesOfBytesKeccak(b [][]byte) ([32]byte, error) {
	if len(b) == 0 {
		return [32]byte{}, nil
	}

	h := utils.Keccak256Fixed(b[0])
	for _, v := range b[1:] {
		h = utils.Keccak256Fixed(append(h[:], v...))
	}
	return h, nil
}

func mergeEpochAndRound(epoch uint32, round uint8) uint64 {
	return uint64(epoch)<<8 + uint64(round)
}

func getMessageIDsAsHexString(messages []evm_2_evm_offramp.InternalEVM2EVMMessage) []string {
	messageIDs := make([]string, 0, len(messages))
	for _, m := range messages {
		messageIDs = append(messageIDs, "0x"+hex.EncodeToString(m.MessageId[:]))
	}
	return messageIDs
}
