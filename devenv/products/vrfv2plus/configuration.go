package vrfv2plus

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/devenv/products"
)

// Configurator implements the devenv Product interface for vrfv2_plus.
type Configurator struct {
	Config []*VRFv2Plus `toml:"vrfv2_plus"`

	// Pre-generated EVM key fields; populated by ModifyNodesets(), used by
	// ConfigureJobsAndContracts().
	txKeyAddr    string
	txKeyEncJSON []byte
	txKeyPass    string
}

// VRFv2Plus holds the per-instance configuration for the vrfv2_plus product.
type VRFv2Plus struct {
	CLNodesFundingETH     float64              `toml:"cl_nodes_funding_eth"`
	CLNodeMaxGasPriceGWei int64                `toml:"cl_node_max_gas_price_gwei"`
	GasSettings           products.GasSettings `toml:"gas_settings"`

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
	PubKeyCompressed string `toml:"pub_key_compressed"`
	KeyHash          string `toml:"key_hash"`
	VRFJobID         string `toml:"vrf_job_id"`
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
