package reportingplugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/NethermindEth/juno/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	bigmath "github.com/smartcontractkit/chainlink/core/utils/big_math"
)

var (
	errNoEligibleUpkeepFound = errors.New("no eligible upkeep found")

	observationArguments = abi.Arguments{
		{Type: mustType(abi.NewType("uint256", "", nil))},
		{Type: mustType(abi.NewType("bytes", "", nil))},
	}
)

type upkeepFailureReason uint8

const (
	upkeepFailureReasonNone upkeepFailureReason = iota
	upkeepFailureReasonUpkeepCanceled
	upkeepFailureReasonUpkeepPaused
	upkeepFailureReasonTargetCheckReverted
	upkeepFailureReasonUpkeepNotNeeded
	upkeepFailureReasonPerformDataExceedsLimit
	upkeepFailureReasonInsufficientBalance
)

type ORM interface {
	RegistryByContractAddress(registryAddress ethkey.EIP55Address) (keeper.Registry, error)
	NewEligibleUpkeepsForRegistry(registryAddress ethkey.EIP55Address, blockNumber int64, gracePeriod int64, binaryHash string) (upkeeps []keeper.UpkeepRegistration, err error)
	EligibleUpkeepsForRegistry(registryAddress ethkey.EIP55Address, blockNumber, gracePeriod int64) (upkeeps []keeper.UpkeepRegistration, err error)
	SetLastRunInfoForUpkeepOnJob(jobID int32, upkeepID *utils.Big, height int64, fromAddress ethkey.EIP55Address, qopts ...pg.QOpt) (int64, error)
}

type Config interface {
	EvmEIP1559DynamicFees() bool
	KeySpecificMaxGasPriceWei(addr gethcommon.Address) *big.Int
	KeeperGasPriceBufferPercent() uint32
	KeeperGasTipCapBufferPercent() uint32
	KeeperBaseFeeBufferPercent() uint32
	KeeperMaximumGracePeriod() int64
	KeeperRegistryCheckGasOverhead() uint32
	KeeperRegistryPerformGasOverhead() uint32
	KeeperCheckUpkeepGasPriceFeatureEnabled() bool
	KeeperTurnLookBack() int64
	KeeperTurnFlagEnabled() bool
}

type observationData struct {
	Head        *evmtypes.Head             `json:"head"`
	Upkeep      *keeper.UpkeepRegistration `json:"upkeepID"`
	PerformData []byte                     `json:"performData"`
}

func newObservationDataFromRaw(data []byte) (*observationData, error) {
	var od observationData
	if err := json.Unmarshal(data, &od); err != nil {
		return nil, err
	}
	return &od, nil
}

func (qd *observationData) raw() ([]byte, error) {
	return json.Marshal(*qd)
}

type checkUpkeepOutput struct {
	GasUsed             *utils.Big          `json:"gasUsed"`
	PerformData         []byte              `json:"performData"`
	UpkeepFailureReason upkeepFailureReason `json:"upkeepFailureReason"`
	UpkeepNeeded        bool                `json:"upkeepNeeded"`
}

// plugin implements types.ReportingPlugin interface with the keepers-specific logic.
type plugin struct {
	logger          logger.Logger
	jobID           int32
	chainID         int64
	cfg             Config
	orm             ORM
	ethClient       evmclient.Client
	headsMngr       *headsMngr
	contractAddress string
	pr              pipeline.Runner
	gasEstimator    gas.Estimator
	chStop          chan struct{}
	// TODO: Keepers ORM
}

// NewPlugin is the constructor of plugin
func NewPlugin(
	logger logger.Logger,
	jobID int32,
	chainID int64,
	cfg Config,
	orm ORM,
	ethClient evmclient.Client,
	headBroadcaster httypes.HeadBroadcaster,
	contractAddress string,
	pr pipeline.Runner,
	gasEstimator gas.Estimator,
) types.ReportingPlugin {
	hm := newHeadsMngr(logger, headBroadcaster)
	hm.start()

	return &plugin{
		logger:          logger,
		jobID:           jobID,
		chainID:         chainID,
		cfg:             cfg,
		orm:             orm,
		ethClient:       ethClient,
		headsMngr:       hm,
		contractAddress: contractAddress,
		pr:              pr,
		gasEstimator:    gasEstimator,
		chStop:          make(chan struct{}),
	}
}

func (p *plugin) Query(_ context.Context, _ types.ReportTimestamp) (types.Query, error) {
	return nil, nil
}

func (p *plugin) Observation(ctx context.Context, _ types.ReportTimestamp, _ types.Query) (types.Observation, error) {
	currentHead := p.headsMngr.getCurrentHead()
	if currentHead == nil {
		p.logger.Warn("current head is nil")
		return nil, nil
	}

	upkeep, err := p.getEligibleUpkeep(currentHead)
	if err != nil {
		if err == errNoEligibleUpkeepFound {
			p.logger.Debugf("no eligible upkeeps found for the head %d", currentHead.Number)
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get eligible upkeep")
	}

	res, err := p.checkUpkeep(ctx, currentHead, upkeep)
	if err != nil {
		return nil, err
	}

	if res.UpkeepFailureReason == upkeepFailureReasonUpkeepNotNeeded || !res.UpkeepNeeded {
		return nil, nil
	} else if res.UpkeepFailureReason > 0 {
		return nil, fmt.Errorf("upkeep reverted with reason %d", res.UpkeepFailureReason)
	}

	return (&observationData{
		Head:        currentHead,
		Upkeep:      upkeep,
		PerformData: res.PerformData,
	}).raw()
}

func (p *plugin) Report(ctx context.Context, _ types.ReportTimestamp, _ types.Query, obs []types.AttributedObservation) (bool, types.Report, error) {
	var observations []*observationData
	for _, ob := range obs {
		if len(ob.Observation) == 0 {
			continue
		}

		p.logger.Info("Report:", string(ob.Observation))

		od, err := newObservationDataFromRaw(ob.Observation)
		if err != nil {
			return false, nil, err
		}

		observations = append(observations, od)
	}

	if len(observations) == 0 {
		p.logger.Info("no observations provided")
		return false, nil, nil
	}

	sort.SliceStable(observations, func(i, j int) bool {
		return observations[i].Head.GreaterThan(observations[j].Head)
	})

	var observation *observationData
	for _, od := range observations {
		checkUpkeep, err := p.checkUpkeep(ctx, od.Head, od.Upkeep)
		if err != nil {
			return false, nil, errors.Wrapf(err, "failed to check upkeep %s for head %d", od.Upkeep.UpkeepID, od.Head.Number)
		}

		if !checkUpkeep.UpkeepNeeded {
			return false, nil, fmt.Errorf("received observation for non-elligible upkeep")
		}

		if bytes.Equal(od.PerformData, checkUpkeep.PerformData) {
			observation = od
			break
		} else {
			p.logger.Warn("observed performa data does not match with the given data")
		}
	}

	if observation == nil {
		p.logger.Info("No observation data found")
		return false, nil, nil
	}

	payload, err := observationArguments.Pack(observation.Upkeep.UpkeepID.ToInt(), observation.PerformData)
	if err != nil {
		return false, nil, errors.Wrap(err, "failed to ABI encode observation results")
	}

	return true, payload, nil
}

func (p *plugin) ShouldAcceptFinalizedReport(_ context.Context, _ types.ReportTimestamp, r types.Report) (bool, error) {
	p.logger.Info("ShouldAcceptFinalizedReport()", string(r))
	return true, nil
}

func (p *plugin) ShouldTransmitAcceptedReport(_ context.Context, _ types.ReportTimestamp, r types.Report) (bool, error) {
	p.logger.Info("ShouldTransmitAcceptedReport()", string(r))
	return true, nil
}

func (p *plugin) Close() error {
	close(p.chStop)
	p.headsMngr.stop()
	return nil
}

func (p *plugin) turnBlockHashBinary(registry keeper.Registry, head *evmtypes.Head, lookback int64) (string, error) {
	turnBlock := head.Number - (head.Number % int64(registry.BlockCountPerTurn)) - lookback
	block, err := p.ethClient.BlockByNumber(context.Background(), big.NewInt(turnBlock))
	if err != nil {
		return "", err
	}
	hashAtHeight := block.Hash()
	binaryString := fmt.Sprintf("%b", hashAtHeight.Big())
	return binaryString, nil
}

func (p *plugin) estimateGasPrice(upkeep *keeper.UpkeepRegistration) (gasPrice *big.Int, fee gas.DynamicFee, err error) {
	var performTxData []byte
	performTxData, err = keeper.Registry1_1ABI.Pack(
		"performUpkeep", // performUpkeep is same across registry ABI versions
		upkeep.UpkeepID.ToInt(),
		common.Hex2Bytes("1234"), // placeholder
	)
	if err != nil {
		return nil, fee, errors.Wrap(err, "unable to construct performUpkeep data")
	}

	keySpecificGasPriceWei := p.cfg.KeySpecificMaxGasPriceWei(upkeep.Registry.FromAddress.Address())
	if p.cfg.EvmEIP1559DynamicFees() {
		fee, _, err = p.gasEstimator.GetDynamicFee(upkeep.ExecuteGas, keySpecificGasPriceWei)
		fee.TipCap = addBuffer(fee.TipCap, p.cfg.KeeperGasTipCapBufferPercent())
	} else {
		gasPrice, _, err = p.gasEstimator.GetLegacyGas(performTxData, upkeep.ExecuteGas, keySpecificGasPriceWei)
		gasPrice = addBuffer(gasPrice, p.cfg.KeeperGasPriceBufferPercent())
	}
	if err != nil {
		return nil, fee, errors.Wrap(err, "unable to estimate gas")
	}

	return gasPrice, fee, nil
}

func (p *plugin) checkUpkeep(ctx context.Context, head *evmtypes.Head, upkeep *keeper.UpkeepRegistration) (*checkUpkeepOutput, error) {
	svcLogger := p.logger.With("jobID", p.jobID, "blockNum", head.Number, "upkeepID", upkeep.UpkeepID)
	svcLogger.Debugw("checking upkeep", "lastRunBlockHeight", upkeep.LastRunBlockHeight, "lastKeeperIndex", upkeep.LastKeeperIndex)

	var gasPrice, gasTipCap, gasFeeCap *big.Int
	if p.cfg.KeeperCheckUpkeepGasPriceFeatureEnabled() {
		price, fee, err := p.estimateGasPrice(upkeep)
		if err != nil {
			return nil, errors.Wrap(err, "estimating gas price")
		}
		gasPrice, gasTipCap, gasFeeCap = price, fee.TipCap, fee.FeeCap

		// Make sure the gas price is at least as large as the basefee to avoid ErrFeeCapTooLow error from geth during eth call.
		// If head.BaseFeePerGas, we assume it is a EIP-1559 chain.
		// Note: gasPrice will be nil if EvmEIP1559DynamicFees is enabled.
		if head.BaseFeePerGas != nil && head.BaseFeePerGas.ToInt().BitLen() > 0 {
			baseFee := addBuffer(head.BaseFeePerGas.ToInt(), p.cfg.KeeperBaseFeeBufferPercent())
			if gasPrice == nil || gasPrice.Cmp(baseFee) < 0 {
				gasPrice = baseFee
			}
		}
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jobSpec": map[string]interface{}{
			"jobID":           p.jobID,
			"contractAddress": upkeep.Registry.ContractAddress.String(),
			"upkeepID":        upkeep.UpkeepID.ToInt(),
			"prettyID":        upkeep.PrettyID(),
			"checkUpkeepGasLimit": p.cfg.KeeperRegistryCheckGasOverhead() + upkeep.Registry.CheckGas +
				p.cfg.KeeperRegistryPerformGasOverhead() + upkeep.ExecuteGas,
			"gasPrice":   gasPrice,
			"gasTipCap":  gasTipCap,
			"gasFeeCap":  gasFeeCap,
			"evmChainID": fmt.Sprintf("%d", p.chainID),
		},
	})

	// DotDagSource in database is empty because all the Keeper pipeline runs make use of the same observation source
	run := pipeline.NewRun(pipeline.Spec{
		ID:              p.jobID,
		DotDagSource:    queryObservationSource,
		CreatedAt:       time.Now(),
		MaxTaskDuration: *models.NewInterval(time.Minute),
	}, vars)

	ctxService, cancel := utils.WithCloseChan(ctx, p.chStop)
	defer cancel()

	if _, err := p.pr.Run(ctxService, &run, svcLogger, false, nil); err != nil {
		return nil, errors.Wrap(err, "failed executing run")
	}

	// Only after task runs where a tx was broadcast
	if run.State == pipeline.RunStatusCompleted {
		_, err := p.orm.SetLastRunInfoForUpkeepOnJob(p.jobID, upkeep.UpkeepID, head.Number, upkeep.Registry.FromAddress, pg.WithParentCtx(ctxService))
		if err != nil {
			svcLogger.Error(errors.Wrap(err, "failed to set last run height for upkeep"))
		}
		svcLogger.Debugw("execute pipeline status completed", "fromAddr", upkeep.Registry.FromAddress)
	} else if run.State == pipeline.RunStatusErrored {
		return nil, fmt.Errorf("failed to check upkeep: %v", run.StringAllErrors())
	}

	runRaw, err := run.Outputs.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal run output")
	}

	var outputs []checkUpkeepOutput
	if err = json.Unmarshal(runRaw, &outputs); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal output")
	}

	if len(outputs) == 0 {
		return nil, fmt.Errorf("expected outputs: %s", string(runRaw))
	}

	return &outputs[0], nil
}

func (p *plugin) getEligibleUpkeep(head *evmtypes.Head) (*keeper.UpkeepRegistration, error) {
	registry, err := p.orm.RegistryByContractAddress(ethkey.MustEIP55Address(p.contractAddress))
	if err != nil {
		return nil, errors.Wrap(err, "unable to load registry")
	}

	var activeUpkeeps []keeper.UpkeepRegistration
	if p.cfg.KeeperTurnFlagEnabled() {
		var turnBinary string
		if turnBinary, err = p.turnBlockHashBinary(registry, head, p.cfg.KeeperTurnLookBack()); err != nil {
			return nil, errors.Wrap(err, "unable to get turn block number hash")
		}

		activeUpkeeps, err = p.orm.NewEligibleUpkeepsForRegistry(
			registry.ContractAddress,
			head.Number,
			p.cfg.KeeperMaximumGracePeriod(),
			turnBinary,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load active registrations")
		}
	} else {
		activeUpkeeps, err = p.orm.EligibleUpkeepsForRegistry(
			registry.ContractAddress,
			head.Number,
			p.cfg.KeeperMaximumGracePeriod(),
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load active registrations")
		}
	}

	if len(activeUpkeeps) == 0 {
		return nil, errNoEligibleUpkeepFound
	}

	// TODO: Implement a correct turn taking logic
	upkeep := activeUpkeeps[0]

	return &upkeep, nil
}

func addBuffer(val *big.Int, prct uint32) *big.Int {
	return bigmath.Div(
		bigmath.Mul(val, 100+prct),
		100,
	)
}

func mustType(tp abi.Type, err error) abi.Type {
	if err != nil {
		panic(err)
	}
	return tp
}
