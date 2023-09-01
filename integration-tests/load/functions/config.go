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
	Soak            *Soak            `toml:"Soak"`
	Load            *Load            `toml:"Load"`
	GatewayListSoak *GatewayListSoak `toml:"GatewayListSoak"`
	GatewaySetSoak  *GatewaySetSoak  `toml:"GatewaySetSoak"`
	Common          *Common          `toml:"Common"`
}

type Common struct {
	Funding
	LINKTokenAddr   string `toml:"link_token_addr"`
	Coordinator     string `toml:"coordinator_addr"`
	Router          string `toml:"router_addr"`
	LoadTestClient  string `toml:"client_addr"`
	SubscriptionID  uint64 `toml:"subscription_id"`
	DONID           string `toml:"don_id"`
	GatewayURL      string `toml:"gateway_url"`
	RequestsPerCall uint32 `toml:"requests_per_call"`
	Receiver        string `toml:"receiver"`
}

type Funding struct {
	NodeFunds *big.Float `toml:"node_funds"`
	SubFunds  *big.Int   `toml:"sub_funds"`
}

type Soak struct {
	RPS      int64            `toml:"rps"`
	Duration *models.Duration `toml:"duration"`
}

type Load struct {
	RPSFrom     int64            `toml:"rps_from"`
	RPSIncrease int64            `toml:"rps_increase"`
	RPSSteps    int              `toml:"rps_steps"`
	Duration    *models.Duration `toml:"duration"`
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
	return cfg, nil
}
