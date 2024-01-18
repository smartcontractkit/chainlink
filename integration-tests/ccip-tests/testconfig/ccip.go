package testconfig

import (
	"github.com/pkg/errors"

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
	ExistingEnv                string             `toml:",omitempty"`
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
}

func (c *CCIPTestConfig) ApplyOverrides(fromCfg *CCIPTestConfig) error {
	if fromCfg == nil {
		return nil
	}
	if fromCfg.BiDirectionalLane != nil {
		c.BiDirectionalLane = fromCfg.BiDirectionalLane
	}
	if fromCfg.KeepEnvAlive != nil {
		c.KeepEnvAlive = fromCfg.KeepEnvAlive
	}
	if fromCfg.NoOfCommitNodes > 0 && fromCfg.NoOfCommitNodes != c.NoOfCommitNodes {
		c.NoOfCommitNodes = fromCfg.NoOfCommitNodes
	}
	if fromCfg.MsgType != "" {
		c.MsgType = fromCfg.MsgType
	}
	if fromCfg.PhaseTimeout != nil {
		c.PhaseTimeout = fromCfg.PhaseTimeout
	}
	if fromCfg.TestDuration != nil {
		c.TestDuration = fromCfg.TestDuration
	}
	if fromCfg.LocalCluster != nil {
		c.LocalCluster = fromCfg.LocalCluster
	}
	if fromCfg.DestGasLimit != nil {
		c.DestGasLimit = fromCfg.DestGasLimit
	}
	if fromCfg.ExistingDeployment != nil {
		c.ExistingDeployment = fromCfg.ExistingDeployment
	}
	if fromCfg.ExistingEnv != "" {
		c.ExistingEnv = fromCfg.ExistingEnv
	}
	if fromCfg.ReuseContracts != nil {
		c.ReuseContracts = fromCfg.ReuseContracts
	}

	if fromCfg.NodeFunding != 0 {
		c.NodeFunding = fromCfg.NodeFunding
	}
	if len(fromCfg.RequestPerUnitTime) != 0 {
		c.RequestPerUnitTime = fromCfg.RequestPerUnitTime
	}
	if fromCfg.TimeUnit != nil {
		c.TimeUnit = fromCfg.TimeUnit
	}
	if len(fromCfg.StepDuration) != 0 {
		c.StepDuration = fromCfg.StepDuration
	}
	if fromCfg.WaitBetweenChaosDuringLoad != nil {
		c.WaitBetweenChaosDuringLoad = fromCfg.WaitBetweenChaosDuringLoad
	}
	if fromCfg.ChaosDuration != nil {
		c.ChaosDuration = fromCfg.ChaosDuration
	}
	if len(fromCfg.NetworkPairs) != 0 {
		c.NetworkPairs = fromCfg.NetworkPairs
	}
	if fromCfg.NoOfNetworks != 0 {
		c.NoOfNetworks = fromCfg.NoOfNetworks
	}
	if fromCfg.NoOfRoutersPerPair != 0 {
		c.NoOfRoutersPerPair = fromCfg.NoOfRoutersPerPair
	}
	if fromCfg.Blockscout {
		c.Blockscout = fromCfg.Blockscout
	}
	if fromCfg.NoOfTokensPerChain != 0 {
		c.NoOfTokensPerChain = fromCfg.NoOfTokensPerChain
	}
	if fromCfg.NoOfTokensInMsg != 0 {
		c.NoOfTokensInMsg = fromCfg.NoOfTokensInMsg
	}
	if fromCfg.MaxNoOfLanes != 0 {
		c.MaxNoOfLanes = fromCfg.MaxNoOfLanes
	}
	if fromCfg.AmountPerToken != 0 {
		c.AmountPerToken = fromCfg.AmountPerToken
	}
	if fromCfg.MulticallInOneTx != nil {
		c.MulticallInOneTx = fromCfg.MulticallInOneTx
	}
	if fromCfg.NoOfSendsInMulticall != 0 {
		c.NoOfSendsInMulticall = fromCfg.NoOfSendsInMulticall
	}

	return nil
}

func (c *CCIPTestConfig) ReadSecrets() error {
	// no secrets to read
	return nil
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

func (c *CCIPContractConfig) ApplyOverrides(from *CCIPContractConfig) error {
	if from == nil {
		return nil
	}
	if from.Data != "" {
		c.Data = from.Data
	}
	return nil
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

func (c *CCIP) ReadSecrets() error {
	err := c.Env.ReadSecrets()
	if err != nil {
		return err
	}
	for _, grp := range c.Groups {
		if err := grp.ReadSecrets(); err != nil {
			return err
		}
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

	for _, grp := range c.Groups {
		if err := grp.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CCIP) ApplyOverrides(fromCfg *CCIP) error {
	if c.Env == nil {
		if fromCfg.Env != nil {
			c.Env = fromCfg.Env
		}
	} else {
		if err := c.Env.ApplyOverrides(fromCfg.Env); err != nil {
			return err
		}
	}
	if c.Deployments == nil {
		if fromCfg.Deployments != nil {
			c.Deployments = fromCfg.Deployments
		}
	} else {
		if err := c.Deployments.ApplyOverrides(fromCfg.Deployments); err != nil {
			return err
		}
	}
	if len(fromCfg.Groups) != 0 {
		for name, grp := range fromCfg.Groups {
			if c.Groups == nil {
				c.Groups = map[string]*CCIPTestConfig{}
			}
			if _, ok := c.Groups[name]; !ok {
				c.Groups[name] = &CCIPTestConfig{}
			}
			if err := c.Groups[name].ApplyOverrides(grp); err != nil {
				return err
			}
		}
	}
	return nil
}
