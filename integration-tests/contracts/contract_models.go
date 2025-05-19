// Package contracts handles deployment, management, and interactions of smart contracts on various chains
package contracts

import (
	"context"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrConfigHelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
)

type FluxAggregatorOptions struct {
	PaymentAmount *big.Int       // The amount of LINK paid to each oracle per submission, in wei (units of 10⁻¹⁸ LINK)
	Timeout       uint32         // The number of seconds after the previous round that are allowed to lapse before allowing an oracle to skip an unfinished round
	Validator     common.Address // An optional contract address for validating external validation of answers
	MinSubValue   *big.Int       // An immutable check for a lower bound of what submission values are accepted from an oracle
	MaxSubValue   *big.Int       // An immutable check for an upper bound of what submission values are accepted from an oracle
	Decimals      uint8          // The number of decimals to offset the answer by
	Description   string         // A short description of what is being reported
}

type FluxAggregatorData struct {
	AllocatedFunds  *big.Int         // The amount of payment yet to be withdrawn by oracles
	AvailableFunds  *big.Int         // The amount of future funding available to oracles
	LatestRoundData RoundData        // Data about the latest round
	Oracles         []common.Address // Addresses of oracles on the contract
}

type FluxAggregatorSetOraclesOptions struct {
	AddList            []common.Address // oracle addresses to add
	RemoveList         []common.Address // oracle addresses to remove
	AdminList          []common.Address // oracle addresses to become admin
	MinSubmissions     uint32           // min amount of submissions in round
	MaxSubmissions     uint32           // max amount of submissions in round
	RestartDelayRounds uint32           // rounds to wait after oracles has changed
}

type SubmissionEvent struct {
	Contract    common.Address
	Submission  *big.Int
	Round       uint32
	BlockNumber uint64
	Oracle      common.Address
}

type FluxAggregator interface {
	Address() string
	Fund(ethAmount *big.Float) error
	LatestRoundID(ctx context.Context) (*big.Int, error)
	LatestRoundData(ctx context.Context) (flux_aggregator_wrapper.LatestRoundData, error)
	GetContractData(ctxt context.Context) (*FluxAggregatorData, error)
	UpdateAvailableFunds() error
	PaymentAmount(ctx context.Context) (*big.Int, error)
	RequestNewRound(ctx context.Context) error
	WithdrawPayment(ctx context.Context, from common.Address, to common.Address, amount *big.Int) error
	WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error)
	GetOracles(ctx context.Context) ([]string, error)
	SetOracles(opts FluxAggregatorSetOraclesOptions) error
	Description(ctxt context.Context) (string, error)
	SetRequesterPermissions(ctx context.Context, addr common.Address, authorized bool, roundsDelay uint32) error
	WatchSubmissionReceived(ctx context.Context, eventChan chan<- *SubmissionEvent) error
}

type LinkToken interface {
	Address() string
	Approve(to string, amount *big.Int) error
	Transfer(to string, amount *big.Int) error
	BalanceOf(ctx context.Context, addr string) (*big.Int, error)
	TransferAndCall(to string, amount *big.Int, data []byte) (*types.Transaction, error)
	TransferAndCallFromKey(to string, amount *big.Int, data []byte, keyNum int) (*types.Transaction, error)
	Name(context.Context) (string, error)
	Decimals() uint
}

type WETHToken interface {
	Address() string
	Approve(to string, amount *big.Int) error
	Transfer(to string, amount *big.Int) error
	BalanceOf(ctx context.Context, addr string) (*big.Int, error)
	Name(context.Context) (string, error)
	Decimals() uint
}

type OffchainOptions struct {
	MaximumGasPrice           uint32         // The highest gas price for which transmitter will be compensated
	ReasonableGasPrice        uint32         // The transmitter will receive reward for gas prices under this value
	MicroLinkPerEth           uint32         // The reimbursement per ETH of gas cost, in 1e-6LINK units
	LinkGweiPerObservation    uint32         // The reward to the oracle for contributing an observation to a successfully transmitted report, in 1e-9LINK units
	LinkGweiPerTransmission   uint32         // The reward to the transmitter of a successful report, in 1e-9LINK units
	MinimumAnswer             *big.Int       // The lowest answer the median of a report is allowed to be
	MaximumAnswer             *big.Int       // The highest answer the median of a report is allowed to be
	BillingAccessController   common.Address // The access controller for billing admin functions
	RequesterAccessController common.Address // The access controller for requesting new rounds
	Decimals                  uint8          // Answers are stored in fixed-point format, with this many digits of precision
	Description               string         // A short description of what is being reported
}

// https://uploads-ssl.webflow.com/5f6b7190899f41fb70882d08/603651a1101106649eef6a53_chainlink-ocr-protocol-paper-02-24-20.pdf
type OffChainAggregatorConfig struct {
	DeltaProgress    time.Duration // The duration in which a leader must achieve progress or be replaced
	DeltaResend      time.Duration // The interval at which nodes resend NEWEPOCH messages
	DeltaRound       time.Duration // The duration after which a new round is started
	DeltaGrace       time.Duration // The duration of the grace period during which delayed oracles can still submit observations
	DeltaC           time.Duration // Limits how often updates are transmitted to the contract as long as the median isn’t changing by more then AlphaPPB
	AlphaPPB         uint64        // Allows larger changes of the median to be reported immediately, bypassing DeltaC
	DeltaStage       time.Duration // Used to stagger stages of the transmission protocol. Multiple Ethereum blocks must be mineable in this period
	RMax             uint8         // The maximum number of rounds in an epoch
	S                []int         // Transmission Schedule
	F                int           // The allowed number of "bad" oracles
	N                int           // The number of oracles
	OracleIdentities []ocrConfigHelper.OracleIdentityExtra
}

type OffChainAggregatorV2Config struct {
	DeltaProgress                           *config.Duration                       `toml:",omitempty"`
	DeltaResend                             *config.Duration                       `toml:",omitempty"`
	DeltaRound                              *config.Duration                       `toml:",omitempty"`
	DeltaGrace                              *config.Duration                       `toml:",omitempty"`
	DeltaStage                              *config.Duration                       `toml:",omitempty"`
	RMax                                    uint8                                  `toml:"-"`
	S                                       []int                                  `toml:"-"`
	Oracles                                 []ocrConfigHelper2.OracleIdentityExtra `toml:"-"`
	ReportingPluginConfig                   []byte                                 `toml:"-"`
	MaxDurationQuery                        *config.Duration                       `toml:",omitempty"`
	MaxDurationObservation                  *config.Duration                       `toml:",omitempty"`
	MaxDurationReport                       *config.Duration                       `toml:",omitempty"`
	MaxDurationShouldAcceptFinalizedReport  *config.Duration                       `toml:",omitempty"`
	MaxDurationShouldTransmitAcceptedReport *config.Duration                       `toml:",omitempty"`
	F                                       int                                    `toml:"-"`
	OnchainConfig                           []byte                                 `toml:"-"`
}

type OffchainAggregatorData struct {
	LatestRoundData RoundData // Data about the latest round
}

type ChainlinkNodeWithKeysAndAddress interface {
	MustReadOCRKeys() (*nodeclient.OCRKeys, error)
	MustReadP2PKeys() (*nodeclient.P2PKeys, error)
	PrimaryEthAddress() (string, error)
	EthAddresses() ([]string, error)
	ChainlinkKeyExporter
}

type ChainlinkKeyExporter interface {
	ExportEVMKeysForChain(string) ([]*nodeclient.ExportedEVMKey, error)
}

type ChainlinkNodeWithForwarder interface {
	TrackForwarder(chainID *big.Int, address common.Address) (*nodeclient.Forwarder, *http.Response, error)
	GetConfig() nodeclient.ChainlinkConfig
}

type OffChainAggregatorWithRounds interface {
	Address() string
	GetLatestRound(ctx context.Context) (*RoundData, error)
	RequestNewRound() error
}

type OffchainAggregator interface {
	Address() string
	SetConfig(chainlinkNodes []ChainlinkNodeWithKeysAndAddress, ocrConfig OffChainAggregatorConfig, transmitters []common.Address) error
	SetPayees([]string, []string) error
	RequestNewRound() error
	GetLatestAnswer(ctx context.Context) (*big.Int, error)
	GetLatestRound(ctx context.Context) (*RoundData, error)
	GetRound(ctx context.Context, roundID *big.Int) (*RoundData, error)
	ParseEventAnswerUpdated(log types.Log) (*offchainaggregator.OffchainAggregatorAnswerUpdated, error)
	LatestRoundDataUpdatedAt() (*big.Int, error)
}

type OffchainAggregatorV2 interface {
	Address() string
	RequestNewRound() error
	SetConfig(ocrConfig *OCRv2Config) error
	SetPayees(transmitters, payees []string) error
	GetLatestAnswer(ctx context.Context) (*big.Int, error)
	GetLatestRound(ctx context.Context) (*RoundData, error)
	GetRound(ctx context.Context, roundID *big.Int) (*RoundData, error)
	ParseEventAnswerUpdated(log types.Log) (*ocr2aggregator.OCR2AggregatorAnswerUpdated, error)
}

type KeeperRegistryCheckUpkeepGasUsageWrapper interface {
	Address() string
}

type Oracle interface {
	Address() string
	Fund(ethAmount *big.Float) error
	SetFulfillmentPermission(address string, allowed bool) error
}

type APIConsumer interface {
	Address() string
	RoundID(ctx context.Context) (*big.Int, error)
	Fund(ethAmount *big.Float) error
	Data(ctx context.Context) (*big.Int, error)
	CreateRequestTo(
		oracleAddr string,
		jobID [32]byte,
		payment *big.Int,
		url string,
		path string,
		times *big.Int,
	) error
}

type Storage interface {
	Get(ctxt context.Context) (*big.Int, error)
	Set(*big.Int) error
}

// JobByInstance helper struct to match job + instance ID
type JobByInstance struct {
	ID       string
	Instance string
}

type MockLINKETHFeed interface {
	Address() string
	LatestRoundData() (*big.Int, error)
	LatestRoundDataUpdatedAt() (*big.Int, error)
}

type MockETHLINKFeed interface {
	Address() string
	LatestRoundData() (*big.Int, error)
	LatestRoundDataUpdatedAt() (*big.Int, error)
}

type MockETHUSDFeed interface {
	Address() string
	LatestRoundData() (*big.Int, error)
	LatestRoundDataUpdatedAt() (*big.Int, error)
	Decimals() uint
}

type MockGasFeed interface {
	Address() string
}

type BlockHashStore interface {
	Address() string
	GetBlockHash(ctx context.Context, blockNumber *big.Int) ([32]byte, error)
	StoreVerifyHeader(blockNumber *big.Int, blockHeader []byte) error
}

type Staking interface {
	Address() string
	Fund(ethAmount *big.Float) error
	AddOperators(operators []common.Address) error
	RemoveOperators(operators []common.Address) error
	SetFeedOperators(operators []common.Address) error
	RaiseAlert() error
	Start(amount *big.Int, initialRewardRate *big.Int) error
	SetMerkleRoot(newMerkleRoot [32]byte) error
}

type StakingEventsMock interface {
	Address() string
	PoolSizeIncreased(maxPoolSize *big.Int) error
	MaxCommunityStakeAmountIncreased(maxStakeAmount *big.Int) error
	MaxOperatorStakeAmountIncreased(maxStakeAmount *big.Int) error
	RewardInitialized(rate *big.Int, available *big.Int, startTimestamp *big.Int, endTimestamp *big.Int) error
	AlertRaised(alerter common.Address, roundId *big.Int, rewardAmount *big.Int) error
	Staked(staker common.Address, newStake *big.Int, totalStake *big.Int) error
	OperatorAdded(operator common.Address) error
	OperatorRemoved(operator common.Address, amount *big.Int) error
	FeedOperatorsSet(feedOperators []common.Address) error
}

type OffchainAggregatorEventsMock interface {
	Address() string
	ConfigSet(previousConfigBlockNumber uint32, configCount uint64, signers []common.Address, transmitters []common.Address, threshold uint8, encodedConfigVersion uint64, encoded []byte) error
	NewTransmission(aggregatorRoundId uint32, answer *big.Int, transmitter common.Address, observations []*big.Int, observers []byte, rawReportContext [32]byte) error
}

type KeeperRegistry11Mock interface {
	Address() string
	EmitUpkeepPerformed(id *big.Int, success bool, from common.Address, payment *big.Int, performData []byte) error
	EmitUpkeepCanceled(id *big.Int, atBlockHeight uint64) error
	EmitFundsWithdrawn(id *big.Int, amount *big.Int, to common.Address) error
	EmitKeepersUpdated(keepers []common.Address, payees []common.Address) error
	EmitUpkeepRegistered(id *big.Int, executeGas uint32, admin common.Address) error
	EmitFundsAdded(id *big.Int, from common.Address, amount *big.Int) error
	SetUpkeepCount(_upkeepCount *big.Int) error
	SetCanceledUpkeepList(_canceledUpkeepList []*big.Int) error
	SetKeeperList(_keepers []common.Address) error
	SetConfig(_paymentPremiumPPB uint32, _flatFeeMicroLink uint32, _blockCountPerTurn *big.Int, _checkGasLimit uint32, _stalenessSeconds *big.Int, _gasCeilingMultiplier uint16, _fallbackGasPrice *big.Int, _fallbackLinkPrice *big.Int) error
	SetUpkeep(id *big.Int, _target common.Address, _executeGas uint32, _balance *big.Int, _admin common.Address, _maxValidBlocknumber uint64, _lastKeeper common.Address, _checkData []byte) error
	SetMinBalance(id *big.Int, minBalance *big.Int) error
	SetCheckUpkeepData(id *big.Int, performData []byte, maxLinkPayment *big.Int, gasLimit *big.Int, adjustedGasWei *big.Int, linkEth *big.Int) error
	SetPerformUpkeepSuccess(id *big.Int, success bool) error
}

type KeeperRegistrar12Mock interface {
	Address() string
	EmitRegistrationRequested(hash [32]byte, name string, encryptedEmail []byte, upkeepContract common.Address, gasLimit uint32, adminAddress common.Address, checkData []byte, amount *big.Int, source uint8) error
	EmitRegistrationApproved(hash [32]byte, displayName string, upkeepId *big.Int) error
	SetRegistrationConfig(_autoApproveConfigType uint8, _autoApproveMaxAllowed uint32, _approvedCount uint32, _keeperRegistry common.Address, _minLINKJuels *big.Int) error
}

type KeeperGasWrapperMock interface {
	Address() string
	SetMeasureCheckGasResult(result bool, payload []byte, gas *big.Int) error
}

type FunctionsV1EventsMock interface {
	Address() string
	EmitRequestProcessed(requestId [32]byte, subscriptionId uint64, totalCostJuels *big.Int, transmitter common.Address, resultCode uint8, response []byte, errByte []byte, callbackReturnData []byte) error
	EmitRequestStart(requestId [32]byte, donId [32]byte, subscriptionId uint64, subscriptionOwner common.Address, requestingContract common.Address, requestInitiator common.Address, data []byte, dataVersion uint16, callbackGasLimit uint32, estimatedTotalCostJuels *big.Int) error
	EmitSubscriptionCanceled(subscriptionId uint64, fundsRecipient common.Address, fundsAmount *big.Int) error
	EmitSubscriptionConsumerAdded(subscriptionId uint64, consumer common.Address) error
	EmitSubscriptionConsumerRemoved(subscriptionId uint64, consumer common.Address) error
	EmitSubscriptionCreated(subscriptionId uint64, owner common.Address) error
	EmitSubscriptionFunded(subscriptionId uint64, oldBalance *big.Int, newBalance *big.Int) error
	EmitSubscriptionOwnerTransferred(subscriptionId uint64, from common.Address, to common.Address) error
	EmitSubscriptionOwnerTransferRequested(subscriptionId uint64, from common.Address, to common.Address) error
	EmitRequestNotProcessed(requestId [32]byte, coordinator common.Address, transmitter common.Address, resultCode uint8) error
	EmitContractUpdated(id [32]byte, from common.Address, to common.Address) error
}

type MockAggregatorProxy interface {
	Address() string
	UpdateAggregator(aggregator common.Address) error
	Aggregator() (common.Address, error)
}

type RoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

// ReadAccessController is read/write access controller, just named by interface
type ReadAccessController interface {
	Address() string
	AddAccess(addr string) error
	DisableAccessCheck() error
}

// Flags flags contract interface
type Flags interface {
	Address() string
	GetFlag(ctx context.Context, addr string) (bool, error)
}

// OperatorFactory creates Operator contracts for node operators
type OperatorFactory interface {
	Address() string
	DeployNewOperatorAndForwarder() (*types.Transaction, error)
	ParseAuthorizedForwarderCreated(log types.Log) (*operator_factory.OperatorFactoryAuthorizedForwarderCreated, error)
	ParseOperatorCreated(log types.Log) (*operator_factory.OperatorFactoryOperatorCreated, error)
}

// Operator operates forwarders
type Operator interface {
	Address() string
	AcceptAuthorizedReceivers(forwarders []common.Address, eoa []common.Address) error
}

// AuthorizedForwarder forward requests from cll nodes eoa
type AuthorizedForwarder interface {
	Address() string
	Owner(ctx context.Context) (string, error)
	GetAuthorizedSenders(ctx context.Context) ([]string, error)
}

type FunctionsCoordinator interface {
	Address() string
	GetThresholdPublicKey() ([]byte, error)
	GetDONPublicKey() ([]byte, error)
}

type FunctionsRouter interface {
	Address() string
	CreateSubscriptionWithConsumer(consumer string) (uint64, error)
}

type FunctionsLoadTestClient interface {
	Address() string
	ResetStats() error
	GetStats() (*EthereumFunctionsLoadStats, error)
	SendRequest(times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) error
	SendRequestWithDONHostedSecrets(times uint32, source string, slotID uint8, slotVersion uint64, args []string, subscriptionId uint64, donID [32]byte) error
}

type MercuryVerifier interface {
	Address() common.Address
	Verify(signedReport []byte, sender common.Address) error
	SetConfig(feedId [32]byte, signers []common.Address, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte, recipientAddressesAndWeights []verifier.CommonAddressAndWeight) (*types.Transaction, error)
	LatestConfigDetails(ctx context.Context, feedId [32]byte) (verifier.LatestConfigDetails, error)
}

type MercuryVerifierProxy interface {
	Address() common.Address
	InitializeVerifier(verifierAddress common.Address) (*types.Transaction, error)
	Verify(signedReport []byte, parameterPayload []byte, value *big.Int) (*types.Transaction, error)
	VerifyBulk(signedReports [][]byte, parameterPayload []byte, value *big.Int) (*types.Transaction, error)
	SetFeeManager(feeManager common.Address) (*types.Transaction, error)
}

type MercuryFeeManager interface {
	Address() common.Address
	UpdateSubscriberDiscount(subscriber common.Address, feedId [32]byte, token common.Address, discount uint64) (*types.Transaction, error)
}

type MercuryRewardManager interface {
	Address() common.Address
	SetFeeManager(feeManager common.Address) (*types.Transaction, error)
}

type WERC20Mock interface {
	Address() common.Address
	BalanceOf(ctx context.Context, addr string) (*big.Int, error)
	Approve(to string, amount *big.Int) error
	Transfer(to string, amount *big.Int) error
	Mint(account common.Address, amount *big.Int) (*types.Transaction, error)
}

type LogEmitter interface {
	Address() common.Address
	EmitLogInts(ints []int) (*types.Transaction, error)
	EmitLogIntsIndexed(ints []int) (*types.Transaction, error)
	EmitLogIntMultiIndexed(ints int, ints2 int, count int) (*types.Transaction, error)
	EmitLogStrings(strings []string) (*types.Transaction, error)
	EmitLogIntsFromKey(ints []int, keyNum int) (*types.Transaction, error)
	EmitLogIntsIndexedFromKey(ints []int, keyNum int) (*types.Transaction, error)
	EmitLogIntMultiIndexedFromKey(ints int, ints2 int, count int, keyNum int) (*types.Transaction, error)
	EmitLogStringsFromKey(strings []string, keyNum int) (*types.Transaction, error)
	EmitLogInt(payload int) (*types.Transaction, error)
	EmitLogIntIndexed(payload int) (*types.Transaction, error)
	EmitLogString(strings string) (*types.Transaction, error)
}
