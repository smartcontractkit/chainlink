package testconfig

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/config"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

type CCIPTestConfig struct {
	KeepEnvAlive               *bool              `toml:",omitempty"`
	BiDirectionalLane          *bool              `toml:",omitempty"`
	CommitAndExecuteOnSameDON  *bool              `toml:",omitempty"`
	NoOfCommitNodes            int                `toml:",omitempty"`
	MsgType                    string             `toml:",omitempty"`
	DestGasLimit               *int64             `toml:",omitempty"`
	MulticallInOneTx           *bool              `toml:",omitempty"`
	NoOfSendsInMulticall       int                `toml:",omitempty"`
	PhaseTimeout               *config.Duration   `toml:",omitempty"`
	TestDuration               *config.Duration   `toml:",omitempty"`
	LocalCluster               *bool              `toml:",omitempty"`
	ExistingDeployment         *bool              `toml:",omitempty"`
	TestRunName                string             `toml:",omitempty"`
	ReuseContracts             *bool              `toml:",omitempty"`
	NodeFunding                float64            `toml:",omitempty"`
	RequestPerUnitTime         []int64            `toml:",omitempty"`
	TimeUnit                   *config.Duration   `toml:",omitempty"`
	StepDuration               []*config.Duration `toml:",omitempty"`
	WaitBetweenChaosDuringLoad *config.Duration   `toml:",omitempty"`
	NetworkPairs               []string           `toml:",omitempty"`
	NoOfNetworks               int                `toml:",omitempty"`
	NoOfRoutersPerPair         int                `toml:",omitempty"`
	Blockscout                 bool               `toml:",omitempty"`
	NoOfTokensPerChain         int                `toml:",omitempty"`
	NoOfTokensInMsg            int                `toml:",omitempty"`
	AmountPerToken             int64              `toml:",omitempty"`
	MaxNoOfLanes               int                `toml:",omitempty"`
	ChaosDuration              *config.Duration   `toml:",omitempty"`
	USDCMockDeployment         *bool              `toml:",omitempty"`
	TimeoutForPriceUpdate      *config.Duration   `toml:",omitempty"`
}

func (c *CCIPTestConfig) SetTestRunName(name string) {
	if c.TestRunName == "" && name != "" {
		c.TestRunName = name
	}
}

func (c *CCIPTestConfig) Validate() error {
	if c.PhaseTimeout != nil && (c.PhaseTimeout.Duration().Minutes() < 1 || c.PhaseTimeout.Duration().Minutes() > 50) {
		return errors.Errorf("phase timeout should be between 1 and 50 minutes")
	}
	if c.TestDuration != nil && c.TestDuration.Duration().Minutes() < 1 {
		return errors.Errorf("test duration should be greater than 1 minute")
	}
	if c.MsgType != "WithoutToken" && c.MsgType != "WithToken" {
		return errors.Errorf("msg type should be either WithoutToken or WithToken")
	}

	if c.MsgType == "WithToken" {
		if c.AmountPerToken == 0 {
			return errors.Errorf("token amount should be greater than 0")
		}
		if c.NoOfTokensPerChain == 0 {
			return errors.Errorf("number of tokens per chain should be greater than 0")
		}
		if c.NoOfTokensInMsg == 0 {
			return errors.Errorf("number of tokens in msg should be greater than 0")
		}
	}

	if c.MulticallInOneTx != nil {
		if c.NoOfSendsInMulticall == 0 {
			return errors.Errorf("number of sends in multisend should be greater than 0 if multisend is true")
		}
	}

	if c.DestGasLimit == nil {
		return errors.Errorf("destination gas limit should be set")
	}

	return nil
}

type CCIPContractConfig struct {
	Data string `toml:",omitempty"`
}

func (c *CCIPContractConfig) ContractsData() []byte {
	if c == nil || c.Data == "" {
		return nil
	}
	return []byte(c.Data)
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

	for _, grp := range c.Groups {
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
