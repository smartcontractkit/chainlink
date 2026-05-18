package testconfig

import (
	"fmt"
	"math/big"
	"os"

	"github.com/AlekSi/pointer"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctfK8config "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"

	ccipcontracts "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
	testutils "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

const (
	CONTRACTS_OVERRIDE_CONFIG string = "BASE64_CCIP_CONFIG_OVERRIDE_CONTRACTS"
	TokenOnlyTransfer         string = "Token"
	DataOnlyTransfer          string = "Data"
	DataAndTokenTransfer      string = "DataWithToken"
)

type OffRampConfig struct {
	MaxDataBytes   *uint32          `toml:",omitempty"`
	BatchGasLimit  *uint32          `toml:",omitempty"`
	InflightExpiry *config.Duration `toml:",omitempty"`
	RootSnooze     *config.Duration `toml:",omitempty"`
}

type MsgDetails struct {
	MsgType        *string `toml:",omitempty"`
	DestGasLimit   *int64  `toml:",omitempty"`
	DataLength     *int64  `toml:",omitempty"`
	NoOfTokens     *int    `toml:",omitempty"`
	AmountPerToken *int64  `toml:",omitempty"`
}

func (m *MsgDetails) IsTokenTransfer() bool {
	return pointer.GetString(m.MsgType) == "Token" || pointer.GetString(m.MsgType) == "DataWithToken"
}

func (m *MsgDetails) IsDataTransfer() bool {
	return pointer.GetString(m.MsgType) == "Data" || pointer.GetString(m.MsgType) == "DataWithToken"
}

func (m *MsgDetails) TransferAmounts() []*big.Int {
	var transferAmounts []*big.Int
	if m.IsTokenTransfer() {
		for i := 0; i < pointer.GetInt(m.NoOfTokens); i++ {
			transferAmounts = append(transferAmounts, big.NewInt(pointer.GetInt64(m.AmountPerToken)))
		}
	}
	return transferAmounts
}

func (m *MsgDetails) Validate() error {
	if m == nil {
		return fmt.Errorf("msg details should be set")
	}
	if m.MsgType == nil {
		return fmt.Errorf("msg type should be set")
	}
	if m.IsDataTransfer() {
		if m.DataLength == nil || *m.DataLength == 0 {
			return fmt.Errorf("data length should be set and greater than 0")
		}
	}
	if m.DestGasLimit == nil {
		return fmt.Errorf("destination gas limit should be set")
	}
	if pointer.GetString(m.MsgType) != DataOnlyTransfer &&
		pointer.GetString(m.MsgType) != TokenOnlyTransfer &&
		pointer.GetString(m.MsgType) != DataAndTokenTransfer {
		return fmt.Errorf("msg type should be - %s/%s/%s", DataOnlyTransfer, TokenOnlyTransfer, DataAndTokenTransfer)
	}

	if m.IsTokenTransfer() {
		if pointer.GetInt64(m.AmountPerToken) == 0 {
			return fmt.Errorf("token amount should be greater than 0")
		}

		if pointer.GetInt(m.NoOfTokens) == 0 {
			return fmt.Errorf("number of tokens in msg should be greater than 0")
		}
	}

	return nil
}

// TokenConfig defines the configuration for tokens in a CCIP test group
type TokenConfig struct {
	NoOfTokensPerChain         *int             `toml:",omitempty"`
	WithPipeline               *bool            `toml:",omitempty"`
	TimeoutForPriceUpdate      *config.Duration `toml:",omitempty"`
	NoOfTokensWithDynamicPrice *int             `toml:",omitempty"`
	DynamicPriceUpdateInterval *config.Duration `toml:",omitempty"`
	// CCIPOwnerTokens dictates if tokens are deployed and controlled by the default CCIP owner account
	// By default, all tokens are deployed and owned by a separate address
	CCIPOwnerTokens *bool `toml:",omitempty"`
}

func (tc *TokenConfig) IsDynamicPriceUpdate() bool {
	return tc.NoOfTokensWithDynamicPrice != nil && *tc.NoOfTokensWithDynamicPrice > 0
}

func (tc *TokenConfig) IsPipelineSpec() bool {
	return pointer.GetBool(tc.WithPipeline)
}

func (tc *TokenConfig) Validate() error {
	if tc == nil {
		return fmt.Errorf("token config should be set")
	}
	if tc.TimeoutForPriceUpdate == nil || tc.TimeoutForPriceUpdate.Duration().Minutes() == 0 {
		return fmt.Errorf("timeout for price update should be set")
	}
	if tc.NoOfTokensWithDynamicPrice != nil && *tc.NoOfTokensWithDynamicPrice > 0 {
		if tc.DynamicPriceUpdateInterval == nil || tc.DynamicPriceUpdateInterval.Duration().Minutes() == 0 {
			return fmt.Errorf("dynamic price update interval should be set if NoOfTokensWithDynamicPrice is greater than 0")
		}
	}
	return nil
}

type MsgProfile struct {
	MsgDetails    *[]*MsgDetails `toml:",omitempty"`
	Frequencies   []int          `toml:",omitempty"`
	matrixByFreq  []int
	mapMsgDetails map[int]*MsgDetails
}

// msgDetailsIndexMatrixByFrequency creates a matrix of msg details index based on their frequency
// This matrix is used to select a msg detail based on the iteration number
// For example, if we have 3 msg details (msg1,msg2,msg3)  with frequencies 2, 3, 5 respectively,
// the matrixByFreq will be [0,0,1,1,1,2,2,2,2,2]
// and mapMsgDetails will be {0:msg1, 1:msg2, 2:msg3}
// So, for iteration 0, msg1 will be selected, for iteration 1, msg1 will be selected, for iteration 2, msg2 will be selected and so on
// This is useful to select a msg detail based on the iteration number
func (m *MsgProfile) msgDetailsIndexMatrixByFrequency() {
	m.mapMsgDetails = make(map[int]*MsgDetails)
	for i, msg := range *m.MsgDetails {
		m.mapMsgDetails[i] = msg
	}
	m.matrixByFreq = make([]int, 0)
	for i, freq := range m.Frequencies {
		for j := 0; j < freq; j++ {
			m.matrixByFreq = append(m.matrixByFreq, i)
		}
	}
	// we do not need frequencies and msg details after creating the matrix
	m.Frequencies = nil
	m.MsgDetails = nil
}

// MsgDetailsForIteration returns the msg details for the given iteration
// The iteration is used to select the msg details based on their frequency
// Refer to msgDetailsIndexMatrixByFrequency for more details
// If the iteration is greater than the number of matrixByFreq, it will loop back to the first msg detail
// if the final iteration in a load run is lesser than the number of matrixByFreq, there is a chance that some of the msg details might not be selected
func (m *MsgProfile) MsgDetailsForIteration(it int64) *MsgDetails {
	index := (it - 1) % int64(len(m.matrixByFreq))
	return m.mapMsgDetails[m.matrixByFreq[index]]
}

// MsgDetailWithMaxToken returns the msg details with the max no of tokens in the msg profile
func (m *MsgProfile) MsgDetailWithMaxToken() *MsgDetails {
	allDetails := *m.MsgDetails
	msgDetails := allDetails[0]
	for _, msg := range allDetails {
		if msg.NoOfTokens != nil && pointer.GetInt(msg.NoOfTokens) > pointer.GetInt(msgDetails.NoOfTokens) {
			msgDetails = msg
		}
	}
	return msgDetails
}

func (m *MsgProfile) Validate() error {
	if m == nil {
		return fmt.Errorf("msg profile should be set")
	}
	if m.MsgDetails == nil {
		return fmt.Errorf("msg details should be set")
	}
	allDetails := *m.MsgDetails
	if len(allDetails) == 0 {
		return fmt.Errorf("msg details should be set")
	}
	if len(m.Frequencies) == 0 {
		return fmt.Errorf("frequencies should be set")
	}
	if len(allDetails) != len(m.Frequencies) {
		return fmt.Errorf("number of msg details %d and frequencies %d should be same", len(allDetails), len(m.Frequencies))
	}
	for _, msg := range allDetails {
		if err := msg.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type LoadFrequency struct {
	RequestPerUnitTime []int64            `toml:",omitempty"`
	TimeUnit           *config.Duration   `toml:",omitempty"`
	StepDuration       []*config.Duration `toml:",omitempty"`
}

type LoadProfile struct {
	MsgProfile                 *MsgProfile               `toml:",omitempty"`
	FrequencyByDestination     map[string]*LoadFrequency `toml:",omitempty"`
	RequestPerUnitTime         []int64                   `toml:",omitempty"`
	TimeUnit                   *config.Duration          `toml:",omitempty"`
	StepDuration               []*config.Duration        `toml:",omitempty"`
	TestDuration               *config.Duration          `toml:",omitempty"`
	NetworkChaosDelay          *config.Duration          `toml:",omitempty"`
	WaitBetweenChaosDuringLoad *config.Duration          `toml:",omitempty"`
	OptimizeSpace              *bool                     `toml:",omitempty"`
	FailOnFirstErrorInLoad     *bool                     `toml:",omitempty"`
	SendMaxDataInEveryMsgCount *int64                    `toml:",omitempty"`
	TestRunName                string                    `toml:",omitempty"`
}

func (l *LoadProfile) Validate() error {
	if l == nil {
		return fmt.Errorf("load profile should be set")
	}
	if err := l.MsgProfile.Validate(); err != nil {
		return err
	}
	if len(l.RequestPerUnitTime) == 0 {
		return fmt.Errorf("request per unit time should be set")
	}
	if l.TimeUnit == nil || l.TimeUnit.Duration().Minutes() == 0 {
		return fmt.Errorf("time unit should be set")
	}
	if l.TestDuration == nil || l.TestDuration.Duration().Minutes() == 0 {
		return fmt.Errorf("test duration should be set")
	}
	return nil
}

func (l *LoadProfile) SetTestRunName(name string) {
	if l.TestRunName == "" && name != "" {
		l.TestRunName = name
	}
}

type ReorgProfile struct {
	FinalityDelta int `toml:",omitempty"`
}

func (gp *ReorgProfile) Validate() error {
	// FinalityDelta can be validated only relatively to CL nodes settings, see setupReorgSuite method
	return nil
}

// CCIPTestGroupConfig defines configuration input to change how a particular CCIP test group should run
type CCIPTestGroupConfig struct {
	Type                                       string                                `toml:",omitempty"`
	KeepEnvAlive                               *bool                                 `toml:",omitempty"`
	BiDirectionalLane                          *bool                                 `toml:",omitempty"`
	CommitAndExecuteOnSameDON                  *bool                                 `toml:",omitempty"`
	AllowOutOfOrder                            *bool                                 `toml:",omitempty"` // To set out of order execution globally
	NoOfCommitNodes                            int                                   `toml:",omitempty"`
	MsgDetails                                 *MsgDetails                           `toml:",omitempty"`
	TokenConfig                                *TokenConfig                          `toml:",omitempty"`
	MulticallInOneTx                           *bool                                 `toml:",omitempty"`
	NoOfSendsInMulticall                       int                                   `toml:",omitempty"`
	PhaseTimeout                               *config.Duration                      `toml:",omitempty"`
	LocalCluster                               *bool                                 `toml:",omitempty"`
	ExistingDeployment                         *bool                                 `toml:",omitempty"`
	ReuseContracts                             *bool                                 `toml:",omitempty"`
	NodeFunding                                float64                               `toml:",omitempty"`
	NetworkPairs                               []string                              `toml:",omitempty"`
	DenselyConnectedNetworkChainIds            []string                              `toml:",omitempty"`
	NoOfNetworks                               int                                   `toml:",omitempty"`
	NoOfRoutersPerPair                         int                                   `toml:",omitempty"`
	MaxNoOfLanes                               int                                   `toml:",omitempty"`
	ChaosDuration                              *config.Duration                      `toml:",omitempty"`
	USDCMockDeployment                         *bool                                 `toml:",omitempty"`
	LBTCMockDeployment                         *bool                                 `toml:",omitempty"`
	LBTCDestPoolDataAs32Bytes                  *bool                                 `toml:",omitempty"`
	CommitOCRParams                            *contracts.OffChainAggregatorV2Config `toml:",omitempty"`
	ExecOCRParams                              *contracts.OffChainAggregatorV2Config `toml:",omitempty"`
	OffRampConfig                              *OffRampConfig                        `toml:",omitempty"`
	CommitInflightExpiry                       *config.Duration                      `toml:",omitempty"`
	StoreLaneConfig                            *bool                                 `toml:",omitempty"`
	LoadProfile                                *LoadProfile                          `toml:",omitempty"`
	ReorgProfile                               *ReorgProfile                         `toml:",omitempty"`
	SkipRequestIfAnotherRequestTriggeredWithin *config.Duration                      `toml:",omitempty"`
}

func (c *CCIPTestGroupConfig) Validate() error {
	if c.Type == Load {
		if err := c.LoadProfile.Validate(); err != nil {
			return err
		}
		if c.MsgDetails == nil {
			c.MsgDetails = c.LoadProfile.MsgProfile.MsgDetailWithMaxToken()
		}
		c.LoadProfile.MsgProfile.msgDetailsIndexMatrixByFrequency()
		if c.ExistingDeployment != nil && *c.ExistingDeployment {
			if c.LoadProfile.TestRunName == "" && os.Getenv(ctfK8config.EnvVarJobImage) != "" {
				return fmt.Errorf("test run name should be set if existing deployment is true and test is running in k8s")
			}
		}
		if c.ReorgProfile != nil {
			if err := c.ReorgProfile.Validate(); err != nil {
				return err
			}
		}
	}
	err := c.MsgDetails.Validate()
	if err != nil {
		return err
	}
	if c.PhaseTimeout != nil && (c.PhaseTimeout.Duration().Minutes() < 1 || c.PhaseTimeout.Duration().Minutes() > 50) {
		return fmt.Errorf("phase timeout should be between 1 and 50 minutes")
	}

	if c.NoOfCommitNodes < 4 {
		return fmt.Errorf("insuffcient number of commit nodes provided")
	}
	if err := c.TokenConfig.Validate(); err != nil {
		return err
	}

	if c.MsgDetails.IsTokenTransfer() {
		if pointer.GetInt(c.TokenConfig.NoOfTokensPerChain) == 0 {
			return fmt.Errorf("number of tokens per chain should be greater than 0")
		}
	}
	if c.MulticallInOneTx != nil {
		if c.NoOfSendsInMulticall == 0 {
			return fmt.Errorf("number of sends in multisend should be greater than 0 if multisend is true")
		}
	}
	if c.SkipRequestIfAnotherRequestTriggeredWithin != nil && c.LoadProfile != nil &&
		c.LoadProfile.TimeUnit.Duration() < c.SkipRequestIfAnotherRequestTriggeredWithin.Duration() {
		return fmt.Errorf("SkipRequestIfAnotherRequestTriggeredWithin should be set below the load TimeUnit duration")
	}

	return nil
}

type CCIPContractConfig struct {
	DataFile *string `toml:",omitempty"`
	Data     string  `toml:",omitempty"`
}

func (c *CCIPContractConfig) DataFilePath() string {
	return pointer.GetString(c.DataFile)
}

// ContractsData reads the contract config passed in TOML
// CCIPContractConfig can accept contract config in string mentioned in Data field
// It also accepts DataFile. Data takes precedence over DataFile
// If you are providing contract config in DataFile, this will read the content of the file
// and set it to CONTRACTS_OVERRIDE_CONFIG env var in base 64 encoded format.
// This comes handy while running tests in remote runner. It ensures that you won't have to deal with copying the
// DataFile to remote runner pod. Instead, you can pass the base64ed content of the file with the help of
// an env var.
func (c *CCIPContractConfig) ContractsData() ([]byte, error) {
	// check if CONTRACTS_OVERRIDE_CONFIG is provided
	// load config from env var if specified for contracts
	rawConfig := os.Getenv(CONTRACTS_OVERRIDE_CONFIG)
	if rawConfig != "" {
		err := DecodeConfig(rawConfig, &c)
		if err != nil {
			return nil, err
		}
	}
	if c == nil {
		return nil, nil
	}
	if c.Data != "" {
		return []byte(c.Data), nil
	}
	// if DataFilePath is given, update c.Data with the content of file so that we can set CONTRACTS_OVERRIDE_CONFIG
	// to pass the file content to remote runner with override config var
	if c.DataFilePath() != "" {
		// if there is regex provided in filepath, reformat the filepath with actual filepath matching the regex
		filePath, err := testutils.FirstFileFromMatchingPath(c.DataFilePath())
		if err != nil {
			return nil, fmt.Errorf("error finding contract config file %s: %w", c.DataFilePath(), err)
		}
		dataContent, err := os.ReadFile(filePath)
		if err != nil {
			return dataContent, fmt.Errorf("error reading contract config file %s : %w", filePath, err)
		}
		c.Data = string(dataContent)
		// encode it to base64 and set to CONTRACTS_OVERRIDE_CONFIG so that the same content can be passed to remote runner
		// we add TEST_ prefix to CONTRACTS_OVERRIDE_CONFIG to ensure the env var is ported to remote runner.
		_, err = EncodeConfigAndSetEnv(c, fmt.Sprintf("TEST_%s", CONTRACTS_OVERRIDE_CONFIG))
		return dataContent, err
	}
	return nil, nil
}

type CCIP struct {
	Env              *Common                                      `toml:",omitempty"`
	ContractVersions map[ccipcontracts.Name]ccipcontracts.Version `toml:",omitempty"`
	Deployments      *CCIPContractConfig                          `toml:",omitempty"`
	Groups           map[string]*CCIPTestGroupConfig              `toml:",omitempty"`
}

// LoadFromEnv loads selected env vars into the CCIP config
func (c *CCIP) LoadFromEnv() error {
	if c.Env == nil {
		c.Env = &Common{}
	}
	err := c.Env.ReadFromEnvVar()
	if err != nil {
		return err
	}
	return nil
}

func (c *CCIP) Validate() error {
	if c.Env != nil {
		err := c.Env.Validate()
		if err != nil {
			return err
		}
	}

	for name, grp := range c.Groups {
		grp.Type = name
		if err := grp.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CCIP) ApplyOverrides(fromCfg *CCIP) error {
	if fromCfg == nil {
		return nil
	}
	logBytes, err := toml.Marshal(fromCfg)
	if err != nil {
		return err
	}
	return ctfconfig.BytesToAnyTomlStruct(zerolog.Logger{}, "", "", c, logBytes)
}
