package evm

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocrconfigurationstoreevmsimple"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// configPollerEVMSimple polls config from OCRConfigurationStoreEVMSimple contract
type configPollerEVMSimple struct {
	services.Service

	lggr                logger.Logger
	configStoreContract *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimple
	configStoreAddr     common.Address
	eventName           string
	abi                 *abi.ABI
	configCh            chan *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleNewConfiguration
	configDigestMutex   sync.Mutex
	latestConfigDigest  ocrtypes.ConfigDigest
}

func newConfigPollerEVMSimple(ctx context.Context,
	lggr logger.Logger,
	configStoreAddress common.Address,
	client client.Client) (
	evmRelayTypes.ConfigPoller, error) {

	configStoreContract, err := ocrconfigurationstoreevmsimple.NewOCRConfigurationStoreEVMSimple(configStoreAddress, client)
	if err != nil {
		return nil, err
	}

	abi, err := ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &configPollerEVMSimple{
		lggr:                lggr,
		configStoreAddr:     configStoreAddress,
		configStoreContract: configStoreContract,
		eventName:           "NewConfiguration",
		abi:                 abi,
	}, nil
}

func (cp *configPollerEVMSimple) Start() {
	cp.lggr.Infof("Starting config poller for contract %s", cp.configStoreAddr)
	cp.configCh = make(chan *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleNewConfiguration, 10)
	cp.configStoreContract.WatchNewConfiguration(nil, cp.configCh, nil)
	go func() {
		for config := range cp.configCh {
			cp.lggr.Infof("TRACE New configuration added to OCRConfigurationStoreEVMSimple: %s", fmt.Sprintf("0x%x", config.ConfigDigest))
			cp.configDigestMutex.Lock()
			defer cp.configDigestMutex.Unlock()
			// Update the latest config digest with the new configuration
			cp.latestConfigDigest = config.ConfigDigest
		}
	}()

}

func (cp *configPollerEVMSimple) Close() error {
	cp.lggr.Infof("Closing config poller for contract %s", cp.configStoreAddr)
	if cp.configCh != nil {
		close(cp.configCh)
	}
	return nil
}

func (cp *configPollerEVMSimple) Notify() <-chan struct{} {
	// TODO(gg): to be implemented if needed
	cp.lggr.Warn("Notify channel not implemented for configPollerEVMSimple")
	return nil
}

func (cp *configPollerEVMSimple) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfigDetails called")

	// TODO(gg) consider decoding config.Raw to get the block number
	cp.lggr.Warn("LatestConfigDetails always returning block 0")
	return 0, cp.latestConfigDigest, nil
}

func (cp *configPollerEVMSimple) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfig called for block %d", changedInBlock)

	if cp.latestConfigDigest == (ocrtypes.ConfigDigest{}) {
		cp.lggr.Warn("LatestConfig called but no config digest available, returning empty config")
		return ocrtypes.ContractConfig{}, nil
	}

	storedConfig, err := cp.configStoreContract.ReadConfig(&bind.CallOpts{}, cp.latestConfigDigest)
	if err != nil {
		cp.lggr.Errorf("Failed to read config for digest %s: %v", cp.latestConfigDigest, err)
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to read config for digest %s: %w", cp.latestConfigDigest, err)
	}

	signers := make([]ocrtypes.OnchainPublicKey, len(storedConfig.Signers))
	for i := range signers {
		signers[i] = storedConfig.Signers[i].Bytes()
	}
	transmitters := make([]ocrtypes.Account, len(storedConfig.Transmitters))
	for i := range transmitters {
		transmitters[i] = ocrtypes.Account(storedConfig.Transmitters[i].Hex())
	}

	ocrConfig := ocrtypes.ContractConfig{
		ConfigDigest:          cp.latestConfigDigest,
		ConfigCount:           uint64(storedConfig.ConfigCount),
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     storedConfig.F,
		OnchainConfig:         storedConfig.OnchainConfig,
		OffchainConfigVersion: storedConfig.OffchainConfigVersion,
		OffchainConfig:        storedConfig.OffchainConfig,
	}

	dgst, err := evmutil.EVMOffchainConfigDigester{}.ConfigDigest(ctx, ocrConfig)
	if err != nil {
		cp.lggr.Errorf("Failed to compute config digest: %v", err)
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to compute config digest: %w", err)
	}
	cp.lggr.Infof("TRACE configPollerEVMSimple LatestConfig found config: %s", dgst.Hex())

	return ocrConfig, nil
}

func (cp *configPollerEVMSimple) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	cp.lggr.Infof("TRACE configPollerEVMSimple LatestBlockHeight called")
	cp.lggr.Warn("LatestBlockHeight always returning block 0, as this is not implemented for configPollerEVMSimple")
	// TODO(gg): implement this if needed
	return uint64(0), nil
}

func (cp *configPollerEVMSimple) Replay(ctx context.Context, fromBlock int64) error {
	cp.lggr.Infof("TRACE configPollerEVMSimple Replay called from block %d", fromBlock)

	cp.lggr.Warn("Replay not implemented for configPollerEVMSimple, returning nil")
	// TODO(gg): implement this if needed
	return nil
}
