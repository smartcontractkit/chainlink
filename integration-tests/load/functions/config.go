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
	SecretsSoak     *SecretsSoak     `toml:"SecretsSoak"`
	RealSoak        *RealSoak        `toml:"RealSoak"`
	Stress          *Stress          `toml:"Stress"`
	SecretsStress   *SecretsStress   `toml:"SecretsStress"`
	RealStress      *RealStress      `toml:"RealStress"`
	GatewayListSoak *GatewayListSoak `toml:"GatewayListSoak"`
	GatewaySetSoak  *GatewaySetSoak  `toml:"GatewaySetSoak"`
	Common          *Common          `toml:"Common"`
	Networks        *Networks        `toml:"Networks"`
	SelectedNetwork *Network         `toml:"SelectedNetwork"`
	PrivateKey      string
}

type Networks struct {
	MumbaiStaging *Network `toml:"MumbaiStaging"`
	Mumbai        *Network `toml:"Mumbai"`
	Fuji          *Network `toml:"Fuji"`
	Sepolia       *Network `toml:"Sepolia"`
}

type Network struct {
	DONID          string `toml:"don_id"`
	LINKTokenAddr  string `toml:"link_token_addr"`
	Coordinator    string `toml:"coordinator_addr"`
	Router         string `toml:"router_addr"`
	LoadTestClient string `toml:"client_addr"`
	SubscriptionID uint64 `toml:"subscription_id"`
}

type Common struct {
	Funding
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
	var keys, urls string
	snet := os.Getenv("SELECTED_NETWORKS")
	log.Warn().Str("NET", snet).Send()
	switch snet {
	case "MUMBAI":
		keys = os.Getenv("MUMBAI_KEYS")
		urls = os.Getenv("MUMBAI_URLS")
		_, isSet := os.LookupEnv("FUNCTIONS_STAGING")
		if isSet {
			cfg.SelectedNetwork = cfg.Networks.MumbaiStaging
		} else {
			cfg.SelectedNetwork = cfg.Networks.Mumbai
		}
	case "AVALANCHE_FUJI":
		keys = os.Getenv("AVALANCHE_FUJI_KEYS")
		urls = os.Getenv("AVALANCHE_FUJI_URLS")
		cfg.SelectedNetwork = cfg.Networks.Fuji
		//networks.AvalancheFuji = blockchain.EVMNetwork{
		//	Name:                      "Avalanche Fuji",
		//	SupportsEIP1559:           true,
		//	ClientImplementation:      blockchain.EthereumClientImplementation,
		//	ChainID:                   43113,
		//	Simulated:                 false,
		//	ChainlinkTransactionLimit: 100000000,
		//	Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		//	MinimumConfirmations:      1,
		//	GasEstimationBuffer:       100000000,
		//	FinalityDepth:             35,
		//	DefaultGasLimit:           900000000000,
		//}
	case "SEPOLIA":
		keys = os.Getenv("SEPOLIA_KEYS")
		urls = os.Getenv("SEPOLIA_URLS")
		cfg.SelectedNetwork = cfg.Networks.Sepolia
	}
	if keys == "" || urls == "" || snet == "" {
		return nil, errors.New(
			"ensure variables are set:\nMUMBAI_KEYS variable, private keys, comma separated\nSELECTED_NETWORKS=MUMBAI\nMUMBAI_URLS variable, websocket urls, comma separated",
		)
	} else {
		cfg.PrivateKey = keys
	}
	log.Debug().Interface("PerformanceConfig", cfg).Msg("Parsed performance config")
	return cfg, nil
}
