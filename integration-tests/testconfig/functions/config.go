package functions

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

const (
	ErrReadPerfConfig      = "failed to read TOML config for performance tests"
	ErrUnmarshalPerfConfig = "failed to unmarshal TOML config for performance tests"
)

type Config struct {
	Soak            *Soak            `toml:"Soak"`
	SecretsSoak     *SecretsSoak     `toml:"SecretsSoak"`
	RealSoak        *RealSoak        `toml:"RealSoak"`
	Stress          *Stress          `toml:"Stress"`
	SecretsStress   *SecretsStress   `toml:"SecretsStress"`
	RealStress      *RealStress      `toml:"RealStress"`
	GatewayListSoak *GatewayListSoak `toml:"GatewayListSoak"`
	GatewaySetSoak  *GatewaySetSoak  `toml:"GatewaySetSoak"`
	Common          *Common          `toml:"Common"`
}

type Common struct {
	Funding
	LINKTokenAddr                   string `toml:"link_token_addr"`
	Coordinator                     string `toml:"coordinator_addr"`
	Router                          string `toml:"router_addr"`
	LoadTestClient                  string `toml:"client_addr"`
	SubscriptionID                  uint64 `toml:"subscription_id"`
	DONID                           string `toml:"don_id"`
	GatewayURL                      string `toml:"gateway_url"`
	Receiver                        string `toml:"receiver"`
	FunctionsCallPayloadHTTP        string `toml:"functions_call_payload_http"`
	FunctionsCallPayloadWithSecrets string `toml:"functions_call_payload_with_secrets"`
	FunctionsCallPayloadReal        string `toml:"functions_call_payload_real"`
	SecretsSlotID                   uint8  `toml:"secrets_slot_id"`
	SecretsVersionID                uint64 `toml:"secrets_version_id"`
	// Secrets these are for CI secrets
	Secrets string `toml:"secrets"`
}

type Funding struct {
	NodeFunds *big.Float `toml:"node_funds"`
	SubFunds  *big.Int   `toml:"sub_funds"`
}

// TODO remove all these types, move these fields to Common and just make sure
// that we read correct TOML file for each test type, although that's a bit problematic
// because we kind of have test types here that do not exist for any other product and adding
// general support for something so specifici is not to my liking
type Soak struct {
	RPS             int64            `toml:"rps"`
	RequestsPerCall uint32           `toml:"requests_per_call"`
	Duration        *models.Duration `toml:"duration"`
}

type SecretsSoak struct {
	RPS             int64            `toml:"rps"`
	RequestsPerCall uint32           `toml:"requests_per_call"`
	Duration        *models.Duration `toml:"duration"`
}

type RealSoak struct {
	RPS             int64            `toml:"rps"`
	RequestsPerCall uint32           `toml:"requests_per_call"`
	Duration        *models.Duration `toml:"duration"`
}

type Stress struct {
	RPS             int64            `toml:"rps"`
	RequestsPerCall uint32           `toml:"requests_per_call"`
	Duration        *models.Duration `toml:"duration"`
}

type SecretsStress struct {
	RPS             int64            `toml:"rps"`
	RequestsPerCall uint32           `toml:"requests_per_call"`
	Duration        *models.Duration `toml:"duration"`
}

type RealStress struct {
	RPS             int64            `toml:"rps"`
	RequestsPerCall uint32           `toml:"requests_per_call"`
	Duration        *models.Duration `toml:"duration"`
}

type GatewayListSoak struct {
	RPS      int64            `toml:"rps"`
	Duration *models.Duration `toml:"duration"`
}

type GatewaySetSoak struct {
	RPS      int64            `toml:"rps"`
	Duration *models.Duration `toml:"duration"`
}

func (c *Config) ApplyOverrides(from interface{}) error {
	//TODO implement me
	return nil
}
