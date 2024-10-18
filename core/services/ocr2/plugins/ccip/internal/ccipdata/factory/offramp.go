package factory

import (
	"context"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_5_0"
)

func NewOffRampReader(lggr logger.Logger, versionFinder VersionFinder, addr cciptypes.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, registerFilters bool) (ccipdata.OffRampReader, error) {
	return initOrCloseOffRampReader(lggr, versionFinder, addr, destClient, lp, estimator, destMaxGasPrice, false, registerFilters)
}

func CloseOffRampReader(lggr logger.Logger, versionFinder VersionFinder, addr cciptypes.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int) error {
	_, err := initOrCloseOffRampReader(lggr, versionFinder, addr, destClient, lp, estimator, destMaxGasPrice, true, false)
	return err
}

func initOrCloseOffRampReader(lggr logger.Logger, versionFinder VersionFinder, addr cciptypes.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, closeReader bool, registerFilters bool) (ccipdata.OffRampReader, error) {
	contractType, version, err := versionFinder.TypeAndVersion(addr, destClient)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read type and version")
	}
	if contractType != ccipconfig.EVM2EVMOffRamp {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.EVM2EVMOffRamp, contractType)
	}

	evmAddr, err := ccipcalc.GenericAddrToEvm(addr)
	if err != nil {
		return nil, err
	}

	lggr.Infow("Initializing OffRamp Reader", "version", version.String(), "destMaxGasPrice", destMaxGasPrice.String())

	switch version.String() {
	case ccipdata.V1_2_0:
		offRamp, err := v1_2_0.NewOffRamp(lggr, evmAddr, destClient, lp, estimator, destMaxGasPrice)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, offRamp.Close()
		}
		return offRamp, offRamp.RegisterFilters()
	case ccipdata.V1_5_0:
		offRamp, err := v1_5_0.NewOffRamp(lggr, evmAddr, destClient, lp, estimator, destMaxGasPrice)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, offRamp.Close()
		}
		return offRamp, offRamp.RegisterFilters()
	default:
		return nil, errors.Errorf("unsupported offramp version %v", version.String())
	}
	// TODO can validate it pointing to the correct version
}

func ExecReportToEthTxMeta(ctx context.Context, typ ccipconfig.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	if typ != ccipconfig.EVM2EVMOffRamp {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.EVM2EVMOffRamp, typ)
	}
	switch ver.String() {
	case ccipdata.V1_2_0, ccipdata.V1_5_0:
		offRampABI := abihelpers.MustParseABI(evm_2_evm_offramp.EVM2EVMOffRampABI)
		return func(report []byte) (*txmgr.TxMeta, error) {
			execReport, err := v1_2_0.DecodeExecReport(ctx, abihelpers.MustGetMethodInputs(ccipdata.ManuallyExecute, offRampABI)[:1], report)
			if err != nil {
				return nil, err
			}
			return execReportToEthTxMeta(execReport)
		}, nil
	default:
		return nil, errors.Errorf("got unexpected version %v", ver.String())
	}
}

func execReportToEthTxMeta(execReport cciptypes.ExecReport) (*txmgr.TxMeta, error) {
	msgIDs := make([]string, len(execReport.Messages))
	for i, msg := range execReport.Messages {
		msgIDs[i] = hexutil.Encode(msg.MessageID[:])
	}

	return &txmgr.TxMeta{
		MessageIDs: msgIDs,
	}, nil
}
