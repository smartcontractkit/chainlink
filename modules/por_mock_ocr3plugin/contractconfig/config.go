package contractconfig

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

var N = 4

var configDigest types.ConfigDigest = types.ConfigDigest{0x13, 0x37}

var contractConfig types.ContractConfig = mustMakeContractConfig()

func mustMakeContractConfig() types.ContractConfig {
	reportingPluginConfig, err := json.Marshal(por.PorOffchainConfig{
		MaxChains: 100,
	})

	if err != nil {
		panic(err)
	}

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		10*time.Second,
		10*time.Second,
		3*time.Second,
		time.Second,
		time.Second,
		time.Second,
		time.Second,
		10,
		[]int{31},
		OracleIdentities(N),
		reportingPluginConfig,
		nil,
		time.Second,
		time.Second,
		time.Second,
		time.Second,
		1,
		nil,
	)
	if err != nil {
		panic(err)
	}

	return types.ContractConfig{
		ConfigDigest:          configDigest,
		ConfigCount:           1,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

var _ types.ContractConfigTracker = &FakeContractConfigTracker{}

type FakeContractConfigTracker struct{}

func (f *FakeContractConfigTracker) Notify() <-chan struct{} {
	return nil
}

func (f *FakeContractConfigTracker) LatestConfigDetails(ctx context.Context) (uint64, types.ConfigDigest, error) {
	return 0, configDigest, nil
}

func (f *FakeContractConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (types.ContractConfig, error) {
	return contractConfig, nil
}

func (f *FakeContractConfigTracker) LatestBlockHeight(ctx context.Context) (uint64, error) {
	return 0, nil
}

var _ types.OffchainConfigDigester = &FakeOffchainConfigDigester{}

type FakeOffchainConfigDigester struct{}

func (f *FakeOffchainConfigDigester) ConfigDigest(ctx context.Context, config types.ContractConfig) (types.ConfigDigest, error) {
	return configDigest, nil
}

func (f *FakeOffchainConfigDigester) ConfigDigestPrefix(ctx context.Context) (types.ConfigDigestPrefix, error) {
	return types.ConfigDigestPrefix(binary.BigEndian.Uint16(configDigest[0:2])), nil
}
