package vrfv2

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/devenv/products"
)

// Configurator implements the devenv Product interface for classic VRF Coordinator V2.
type Configurator struct {
	Config []*VRFv2 `toml:"vrfv2"`

	nodeEVMKeyAddr    string
	nodeEVMKeyEncJSON []byte
	nodeEVMKeyPass    string

	txKeyAddrs    []string
	txKeyEncJSONs [][]byte

	bhsKeyAddr    string
	bhsKeyEncJSON []byte
}

// VRFv2 holds per-instance configuration for vrfv2 product.
type VRFv2 struct {
	CLNodesFundingETH     float64              `toml:"cl_nodes_funding_eth"`
	CLNodeMaxGasPriceGWei int64                `toml:"cl_node_max_gas_price_gwei"`
	GasSettings           products.GasSettings `toml:"gas_settings"`

	NumTxKeys int `toml:"num_tx_keys"`

	MinimumConfirmations       uint16 `toml:"minimum_confirmations"`
	MaxGasLimitCoordinator     uint32 `toml:"max_gas_limit_coordinator"`
	StalenessSeconds           uint32 `toml:"staleness_seconds"`
	GasAfterPaymentCalculation uint32 `toml:"gas_after_payment_calculation"`
	FallbackWeiPerUnitLink     string `toml:"fallback_wei_per_unit_link"`

	FulfillmentFlatFeeLinkPPMTier1 uint32 `toml:"fulfillment_flat_fee_link_ppm_tier_1"`
	FulfillmentFlatFeeLinkPPMTier2 uint32 `toml:"fulfillment_flat_fee_link_ppm_tier_2"`
	FulfillmentFlatFeeLinkPPMTier3 uint32 `toml:"fulfillment_flat_fee_link_ppm_tier_3"`
	FulfillmentFlatFeeLinkPPMTier4 uint32 `toml:"fulfillment_flat_fee_link_ppm_tier_4"`
	FulfillmentFlatFeeLinkPPMTier5 uint32 `toml:"fulfillment_flat_fee_link_ppm_tier_5"`
	ReqsForTier2                   int64  `toml:"reqs_for_tier_2"`
	ReqsForTier3                   int64  `toml:"reqs_for_tier_3"`
	ReqsForTier4                   int64  `toml:"reqs_for_tier_4"`
	ReqsForTier5                   int64  `toml:"reqs_for_tier_5"`

	VRFJobForwardingAllowed             bool    `toml:"vrf_job_forwarding_allowed"`
	VRFJobEstimateGasMultiplier         float64 `toml:"vrf_job_estimate_gas_multiplier"`
	VRFJobBatchFulfillmentEnabled       bool    `toml:"vrf_job_batch_fulfillment_enabled"`
	VRFJobBatchFulfillmentGasMultiplier float64 `toml:"vrf_job_batch_fulfillment_gas_multiplier"`
	VRFJobPollPeriod                    string  `toml:"vrf_job_poll_period"`
	VRFJobRequestTimeout                string  `toml:"vrf_job_request_timeout"`
	VRFJobSimulationBlock               string  `toml:"vrf_job_simulation_block"`

	SubFundingAmountLink float64 `toml:"sub_funding_amount_link"`

	WrapperGasOverhead               uint32  `toml:"wrapper_gas_overhead"`
	CoordinatorGasOverhead           uint32  `toml:"coordinator_gas_overhead"`
	WrapperPremiumPercentage         uint8   `toml:"wrapper_premium_percentage"`
	WrapperMaxNumberOfWords          uint8   `toml:"wrapper_max_number_of_words"`
	WrapperConsumerFundingAmountLink float64 `toml:"wrapper_consumer_funding_amount_link"`

	NumberOfWords                             uint32 `toml:"number_of_words"`
	CallbackGasLimit                          uint32 `toml:"callback_gas_limit"`
	RandomnessRequestCountPerRequest          uint16 `toml:"randomness_request_count_per_request"`
	RandomnessRequestCountPerRequestDeviation uint16 `toml:"randomness_request_count_per_request_deviation"`
	RandomWordsFulfilledEventTimeout          string `toml:"random_words_fulfilled_event_timeout"`

	BatchCallbackGasLimit uint32 `toml:"batch_callback_gas_limit"`
	BatchTxGasBudget      uint32 `toml:"batch_tx_gas_budget"`

	EnableBHSJob         bool   `toml:"enable_bhs_job"`
	BHSJobWaitBlocks     int    `toml:"bhs_job_wait_blocks"`
	BHSJobLookbackBlocks int    `toml:"bhs_job_lookback_blocks"`
	BHSJobPollPeriod     string `toml:"bhs_job_poll_period"`
	BHSJobRunTimeout     string `toml:"bhs_job_run_timeout"`

	DeployedContracts DeployedContracts `toml:"deployed_contracts"`
	VRFKeyData        KeyOutput         `toml:"vrf_key_data"`
}

// DeployedContracts holds deployed contract addresses for VRF v2 smoke.
type DeployedContracts struct {
	LinkToken        string `toml:"link_token"`
	MockFeed         string `toml:"mock_feed"`
	BHS              string `toml:"bhs"`
	Coordinator      string `toml:"coordinator"`
	BatchCoordinator string `toml:"batch_coordinator"`
}

// KeyOutput is persisted to env-out for tests.
type KeyOutput struct {
	PubKeyCompressed   string   `toml:"pub_key_compressed"`
	PubKeyUncompressed string   `toml:"pub_key_uncompressed"`
	KeyHash            string   `toml:"key_hash"`
	VRFJobID           string   `toml:"vrf_job_id"`
	TxKeyAddresses     []string `toml:"tx_key_addresses"`
	BHSJobID           string   `toml:"bhs_job_id"`
}

// NewConfigurator returns a new vrfv2 configurator.
func NewConfigurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load vrfv2 product config: %w", err)
	}
	m.Config = cfg.Config
	return nil
}

func (m *Configurator) Store(path string, instanceIdx int) error {
	if err := products.Store(".", &Configurator{Config: []*VRFv2{m.Config[instanceIdx]}}); err != nil {
		return fmt.Errorf("failed to store vrfv2 product config: %w", err)
	}
	return nil
}
