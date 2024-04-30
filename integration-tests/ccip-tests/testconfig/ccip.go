package testconfig

import (
	"fmt"
	"math/big"
	"os"

	"github.com/AlekSi/pointer"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"

	testutils "github.com/smartcontractkit/ccip/integration-tests/ccip-tests/utils"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/config"
	ctfK8config "github.com/smartcontractkit/chainlink-testing-framework/k8s/config"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

const (
	CONTRACTS_OVERRIDE_CONFIG        = "BASE64_CCIP_CONFIG_OVERRIDE_CONTRACTS"
	TokenOnlyTransfer         string = "Token"

	DataOnlyTransfer string = "Data"

	DataAndTokenTransfer string = "DataWithToken"
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

type TokenConfig struct {
	NoOfTokensPerChain         *int             `toml:",omitempty"`
	WithPipeline               *bool            `toml:",omitempty"`
	TimeoutForPriceUpdate      *config.Duration `toml:",omitempty"`
	NoOfTokensWithDynamicPrice *int             `toml:",omitempty"`
	DynamicPriceUpdateInterval *config.Duration `toml:",omitempty"`
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

type CCIPTestConfig struct {
	Type                                       string                                `toml:",omitempty"`
	KeepEnvAlive                               *bool                                 `toml:",omitempty"`
	BiDirectionalLane                          *bool                                 `toml:",omitempty"`
	CommitAndExecuteOnSameDON                  *bool                                 `toml:",omitempty"`
	NoOfCommitNodes                            int                                   `toml:",omitempty"`
	MsgDetails                                 *MsgDetails                           `toml:",omitempty"`
	TokenConfig                                *TokenConfig                          `toml:",omitempty"`
	MulticallInOneTx                           *bool                                 `toml:",omitempty"`
	NoOfSendsInMulticall                       int                                   `toml:",omitempty"`
	PhaseTimeout                               *config.Duration                      `toml:",omitempty"`
	TestDuration                               *config.Duration                      `toml:",omitempty"`
	LocalCluster                               *bool                                 `toml:",omitempty"`
	ExistingDeployment                         *bool                                 `toml:",omitempty"`
	TestRunName                                string                                `toml:",omitempty"`
	ReuseContracts                             *bool                                 `toml:",omitempty"`
	NodeFunding                                float64                               `toml:",omitempty"`
	RequestPerUnitTime                         []int64                               `toml:",omitempty"`
	TimeUnit                                   *config.Duration                      `toml:",omitempty"`
	StepDuration                               []*config.Duration                    `toml:",omitempty"`
	WaitBetweenChaosDuringLoad                 *config.Duration                      `toml:",omitempty"`
	NetworkPairs                               []string                              `toml:",omitempty"`
	NoOfNetworks                               int                                   `toml:",omitempty"`
	NoOfRoutersPerPair                         int                                   `toml:",omitempty"`
	MaxNoOfLanes                               int                                   `toml:",omitempty"`
	ChaosDuration                              *config.Duration                      `toml:",omitempty"`
	USDCMockDeployment                         *bool                                 `toml:",omitempty"`
	FailOnFirstErrorInLoad                     *bool                                 `toml:",omitempty"`
	SendMaxDataInEveryMsgCount                 *int64                                `toml:",omitempty"`
	CommitOCRParams                            *contracts.OffChainAggregatorV2Config `toml:",omitempty"`
	ExecOCRParams                              *contracts.OffChainAggregatorV2Config `toml:",omitempty"`
	OffRampConfig                              *OffRampConfig                        `toml:",omitempty"`
	CommitInflightExpiry                       *config.Duration                      `toml:",omitempty"`
	OptimizeSpace                              *bool                                 `toml:",omitempty"`
	SkipRequestIfAnotherRequestTriggeredWithin *config.Duration                      `toml:",omitempty"`
	StoreLaneConfig                            *bool                                 `toml:",omitempty"`
}

func (c *CCIPTestConfig) SetTestRunName(name string) {
	if c.TestRunName == "" && name != "" {
		c.TestRunName = name
	}
}

func (c *CCIPTestConfig) Validate() error {
	err := c.MsgDetails.Validate()
	if err != nil {
		return err
	}
	if c.PhaseTimeout != nil && (c.PhaseTimeout.Duration().Minutes() < 1 || c.PhaseTimeout.Duration().Minutes() > 50) {
		return fmt.Errorf("phase timeout should be between 1 and 50 minutes")
	}
	if c.TestDuration != nil && c.TestDuration.Duration().Minutes() < 1 {
		return fmt.Errorf("test duration should be greater than 1 minute")
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
	if c.ExistingDeployment != nil && *c.ExistingDeployment {
		if c.Type != Smoke {
			if c.TestRunName == "" && os.Getenv(ctfK8config.EnvVarJobImage) != "" {
				return fmt.Errorf("test run name should be set if existing deployment is true and test is running in k8s")
			}
		}
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
	Env         *Common                    `toml:",omitempty"`
	Deployments *CCIPContractConfig        `toml:",omitempty"`
	Groups      map[string]*CCIPTestConfig `toml:",omitempty"`
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
	lggr := zerolog.Logger{}
	return ctfconfig.BytesToAnyTomlStruct(lggr, "", "", c, logBytes)
}
