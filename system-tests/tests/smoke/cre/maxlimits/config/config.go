package config

// LogTriggerConfig defines an EVM log trigger subscription.
type LogTriggerConfig struct {
	ChainSelector uint64   `yaml:"chainSelector"`
	Addresses     []string `yaml:"addresses"`
	TopicSlots    []TopicSlotConfig `yaml:"topicSlots"`
	Confidence    int32    `yaml:"confidence"`
}

// TopicSlotConfig represents topic values for a single slot.
type TopicSlotConfig struct {
	Values []string `yaml:"values"`
}

// AuthorizedKeyConfig holds an authorized key for the HTTP trigger.
type AuthorizedKeyConfig struct {
	Type      string `yaml:"type"`
	PublicKey string `yaml:"publicKey"`
}

// ChainReadCallConfig describes a single EVM read call.
type ChainReadCallConfig struct {
	Method          string `yaml:"method"` // BalanceAt, CallContract, FilterLogs, HeaderByNumber, GetTxReceipt, GetTxByHash, EstimateGas
	ChainSelector   uint64 `yaml:"chainSelector"`
	AccountAddress  []byte `yaml:"accountAddress,omitempty"`
	ContractAddress []byte `yaml:"contractAddress,omitempty"`
	CallData        []byte `yaml:"callData,omitempty"`
	TxHash          []byte `yaml:"txHash,omitempty"`
	FromBlock       int64  `yaml:"fromBlock,omitempty"`
	ToBlock         int64  `yaml:"toBlock,omitempty"`
	BlockNumber     int64  `yaml:"blockNumber,omitempty"`
}

// ChainWriteTargetConfig describes a single EVM write target.
type ChainWriteTargetConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	Receiver      []byte `yaml:"receiver"`
	GasLimit      uint64 `yaml:"gasLimit"`
}

// HTTPEndpointConfig describes an HTTP endpoint to call.
type HTTPEndpointConfig struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
	Body   string `yaml:"body,omitempty"`
}

// ConfHTTPEndpointConfig describes a confidential HTTP endpoint to call.
type ConfHTTPEndpointConfig struct {
	URL       string   `yaml:"url"`
	Method    string   `yaml:"method"`
	Body      string   `yaml:"body,omitempty"`
	SecretKey string   `yaml:"secretKey"`
}

// SecretConfig describes a secret to fetch.
type SecretConfig struct {
	ID        string `yaml:"id"`
	Namespace string `yaml:"namespace"`
}

// Config is the top-level workflow configuration.
// All arrays are sized to match CRE per-execution call limits.
type Config struct {
	// Trigger configuration
	CronSchedule       string                `yaml:"cronSchedule"`
	HTTPAuthorizedKeys []AuthorizedKeyConfig  `yaml:"httpAuthorizedKeys"`
	LogTriggers        []LogTriggerConfig     `yaml:"logTriggers"`

	// Chain read configuration (15 calls max)
	ChainReads []ChainReadCallConfig `yaml:"chainReads"`

	// Chain write configuration (10 targets max)
	ChainWrites []ChainWriteTargetConfig `yaml:"chainWrites"`

	// HTTP action configuration (5 calls max)
	HTTPEndpoints []HTTPEndpointConfig `yaml:"httpEndpoints"`

	// Confidential HTTP configuration (5 calls max)
	ConfHTTPEndpoints []ConfHTTPEndpointConfig `yaml:"confHttpEndpoints"`

	// Secrets configuration (5 calls max)
	Secrets []SecretConfig `yaml:"secrets"`

	// Consensus observation payload size target in bytes (up to 100 KB)
	ConsensusPayloadSize int `yaml:"consensusPayloadSize"`

	// Number of consensus rounds (up to 20)
	ConsensusRounds int `yaml:"consensusRounds"`

	// Number of log lines to emit (up to 999)
	LogEventCount int `yaml:"logEventCount"`

	// Report payload for chain writes (up to ~5 KB)
	ReportPayload []byte `yaml:"reportPayload"`

	// Metadata is a large key-value map used to approach the 1 MB config limit.
	Metadata map[string]string `yaml:"metadata"`
}
