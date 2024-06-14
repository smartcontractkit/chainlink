package loadfunctions

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"math/big"
	"os"
)

const (
	DefaultConfigFilename = "config.toml"

	ErrReadPerfConfig      = "failed to read TOML config for performance tests"
	ErrUnmarshalPerfConfig = "failed to unmarshal TOML config for performance tests"
)

type PerformanceConfig struct {
	Soak             *Soak            `toml:"Soak"`
	SecretsSoak      *SecretsSoak     `toml:"SecretsSoak"`
	RealSoak         *RealSoak        `toml:"RealSoak"`
	Stress           *Stress          `toml:"Stress"`
	SecretsStress    *SecretsStress   `toml:"SecretsStress"`
	RealStress       *RealStress      `toml:"RealStress"`
	GatewayListSoak  *GatewayListSoak `toml:"GatewayListSoak"`
	GatewaySetSoak   *GatewaySetSoak  `toml:"GatewaySetSoak"`
	Common           *Common          `toml:"Common"`
	MumbaiPrivateKey string
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

func ReadConfig() (*PerformanceConfig, error) {
	var cfg *PerformanceConfig
	d, err := os.ReadFile(DefaultConfigFilename)
	if err != nil {
		return nil, errors.Wrap(err, ErrReadPerfConfig)
	}
	err = toml.Unmarshal(d, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalPerfConfig)
	}
	log.Debug().Interface("PerformanceConfig", cfg).Msg("Parsed performance config")
	mpk := os.Getenv("MUMBAI_KEYS")
	murls := os.Getenv("MUMBAI_URLS")
	snet := os.Getenv("SELECTED_NETWORKS")
	if mpk == "" || murls == "" || snet == "" {
		return nil, errors.New(
			"ensure variables are set:\nMUMBAI_KEYS variable, private keys, comma separated\nSELECTED_NETWORKS=MUMBAI\nMUMBAI_URLS variable, websocket urls, comma separated",
		)
	} else {
		cfg.MumbaiPrivateKey = mpk
	}
	return cfg, nil
}
