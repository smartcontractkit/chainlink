package vrfv2plus

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/devenv/products"
)

// Configurator implements the devenv Product interface for vrfv2_plus.
type Configurator struct {
	Config []*VRFv2Plus `toml:"vrfv2_plus"`

	// nodeEVMKey: always pre-generated; index 0 of fromAddresses.
	nodeEVMKeyAddr    string
	nodeEVMKeyEncJSON []byte
	nodeEVMKeyPass    string

	// Extra TX keys (len = cfg.NumTxKeys); appended after nodeEVMKey in fromAddresses.
	txKeyAddrs    []string
	txKeyEncJSONs [][]byte

	// BHS node TX key (only when EnableBHSJob)
	bhsKeyAddr    string
	bhsKeyEncJSON []byte

	// BHF node TX key (only when EnableBHFJob)
	bhfKeyAddr    string
	bhfKeyEncJSON []byte
}

// VRFv2Plus holds the per-instance configuration for the vrfv2_plus product.
type VRFv2Plus struct {
	CLNodesFundingETH     float64              `toml:"cl_nodes_funding_eth"`
	CLNodeMaxGasPriceGWei int64                `toml:"cl_node_max_gas_price_gwei"`
	GasSettings           products.GasSettings `toml:"gas_settings"`

	// TX keys: N extra keys generated in addition to the node's own EVM key (default 0)
	NumTxKeys int `toml:"num_tx_keys"`

	// Coordinator config
	MinimumConfirmations    uint16 `toml:"minimum_confirmations"`
	MaxGasLimitCoordinator  uint32 `toml:"max_gas_limit_coordinator"`
	FlatFeeNativePPM        uint32 `toml:"flat_fee_native_ppm"`
	FlatFeeLinkDiscountPPM  uint32 `toml:"flat_fee_link_discount_ppm"`
	NativePremiumPercentage uint8  `toml:"native_premium_percentage"`
	LinkPremiumPercentage   uint8  `toml:"link_premium_percentage"`

	// Job config
	VRFJobPollPeriod     string `toml:"vrf_job_poll_period"`
	VRFJobRequestTimeout string `toml:"vrf_job_request_timeout"`

	// Subscription funding defaults (used by tests)
	SubFundingAmountLink   float64 `toml:"sub_funding_amount_link"`
	SubFundingAmountNative float64 `toml:"sub_funding_amount_native"`

	// Wrapper config
	WrapperGasOverhead            uint32 `toml:"wrapper_gas_overhead"`
	CoordinatorGasOverheadNative  uint32 `toml:"coordinator_gas_overhead_native"`
	CoordinatorGasOverheadLink    uint32 `toml:"coordinator_gas_overhead_link"`
	CoordinatorNativePremiumPct   uint8  `toml:"coordinator_native_premium_pct"`
	CoordinatorLinkPremiumPct     uint8  `toml:"coordinator_link_premium_pct"`
	CoordinatorGasOverheadPerWord uint16 `toml:"coordinator_gas_overhead_per_word"`

	// Wrapper consumer funding (hardcoded; not in TOML)
	WrapperConsumerFundLinkJuels *big.Int `toml:"-"`
	WrapperConsumerFundNativeWei *big.Int `toml:"-"`

	// Batch fulfillment
	BatchFulfillmentEnabled       bool    `toml:"batch_fulfillment_enabled"`
	BatchFulfillmentGasMultiplier float64 `toml:"batch_fulfillment_gas_multiplier"`
	BatchCallbackGasLimit         uint32  `toml:"batch_callback_gas_limit"`
	BatchTxGasBudget              uint32  `toml:"batch_tx_gas_budget"`

	// BHS job (only when EnableBHSJob = true)
	EnableBHSJob         bool   `toml:"enable_bhs_job"`
	BHSJobWaitBlocks     int    `toml:"bhs_job_wait_blocks"`
	BHSJobLookbackBlocks int    `toml:"bhs_job_lookback_blocks"`
	BHSJobPollPeriod     string `toml:"bhs_job_poll_period"`
	BHSJobRunTimeout     string `toml:"bhs_job_run_timeout"`

	// BHF job (only when EnableBHFJob = true)
	EnableBHFJob         bool   `toml:"enable_bhf_job"`
	BHFJobWaitBlocks     int    `toml:"bhf_job_wait_blocks"`
	BHFJobLookbackBlocks int    `toml:"bhf_job_lookback_blocks"`
	BHFJobPollPeriod     string `toml:"bhf_job_poll_period"`
	BHFJobRunTimeout     string `toml:"bhf_job_run_timeout"`

	// VRF job pipeline simulation block ("pending" or "")
	VRFJobSimulationBlock string `toml:"vrf_job_simulation_block"`

	DeployedContracts VRFDeployedContracts `toml:"deployed_contracts"`
	VRFKeyData        VRFKeyOutput         `toml:"vrf_key_data"`
}

// VRFDeployedContracts holds addresses of all deployed VRF-related contracts.
type VRFDeployedContracts struct {
	LinkToken        string `toml:"link_token"`
	MockFeed         string `toml:"mock_feed"`
	BHS              string `toml:"bhs"`
	BatchBHS         string `toml:"batch_bhs"`
	Coordinator      string `toml:"coordinator"`
	BatchCoordinator string `toml:"batch_coordinator"`
	Wrapper          string `toml:"wrapper"`
	WrapperConsumer  string `toml:"wrapper_consumer"`
	WrapperSubID     string `toml:"wrapper_sub_id"`
}

// VRFKeyOutput holds VRF key data and the job ID, written to the output TOML.
type VRFKeyOutput struct {
	PubKeyCompressed   string   `toml:"pub_key_compressed"`
	PubKeyUncompressed string   `toml:"pub_key_uncompressed"` // full uncompressed key; needed for RegisterProvingKey
	KeyHash            string   `toml:"key_hash"`
	VRFJobID           string   `toml:"vrf_job_id"`
	TxKeyAddresses     []string `toml:"tx_key_addresses"` // [nodeEVMKey, ...extraKeys]; len = 1 + num_tx_keys
	BHSJobID           string   `toml:"bhs_job_id"`       // populated if EnableBHSJob
	BHFJobID           string   `toml:"bhf_job_id"`       // populated if EnableBHFJob
}

func NewConfigurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load vrfv2plus product config: %w", err)
	}
	m.Config = cfg.Config
	return nil
}

func (m *Configurator) Store(path string, instanceIdx int) error {
	if err := products.Store(".", &Configurator{Config: []*VRFv2Plus{m.Config[instanceIdx]}}); err != nil {
		return fmt.Errorf("failed to store vrfv2plus product config: %w", err)
	}
	return nil
}
