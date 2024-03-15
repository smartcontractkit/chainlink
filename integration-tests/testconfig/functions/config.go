package functions

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/net"
)

const (
	ErrReadPerfConfig      = "failed to read TOML config for performance tests"
	ErrUnmarshalPerfConfig = "failed to unmarshal TOML config for performance tests"
)

type Config struct {
	Performance *Performance `toml:"Performance"`
	Common      *Common      `toml:"Common"`
}

type Common struct {
	NodeFunds                       *big.Float `toml:"node_funds"`
	SubFunds                        *big.Int   `toml:"sub_funds"`
	LINKTokenAddr                   *string    `toml:"link_token_addr"`
	Coordinator                     *string    `toml:"coordinator_addr"`
	Router                          *string    `toml:"router_addr"`
	LoadTestClient                  *string    `toml:"client_addr"`
	SubscriptionID                  *uint64    `toml:"subscription_id"`
	DONID                           *string    `toml:"don_id"`
	GatewayURL                      *string    `toml:"gateway_url"`
	Receiver                        *string    `toml:"receiver"`
	FunctionsCallPayloadHTTP        *string    `toml:"functions_call_payload_http"`
	FunctionsCallPayloadWithSecrets *string    `toml:"functions_call_payload_with_secrets"`
	FunctionsCallPayloadReal        *string    `toml:"functions_call_payload_real"`
	SecretsSlotID                   *uint8     `toml:"secrets_slot_id"`
	SecretsVersionID                *uint64    `toml:"secrets_version_id"`
	// Secrets these are for CI secrets
	Secrets *string `toml:"secrets"`
}

func (c *Common) Validate() error {
	if c.SubFunds == nil {
		return errors.New("sub_funds must be set")
	}
	if c.LINKTokenAddr == nil || *c.LINKTokenAddr == "" {
		return errors.New("link_token_addr must be set")
	}
	if !common.IsHexAddress(*c.LINKTokenAddr) {
		return errors.New("link_token_addr must be a valid address")
	}
	if c.Coordinator == nil || *c.Coordinator == "" {
		return errors.New("coordinator must be set")
	}
	if !common.IsHexAddress(*c.Coordinator) {
		return errors.New("coordinator must be a valid address")
	}
	if c.Router == nil || *c.Router == "" {
		return errors.New("router must be set")
	}
	if !common.IsHexAddress(*c.Router) {
		return errors.New("router must be a valid address")
	}
	if c.DONID == nil || *c.DONID == "" {
		return errors.New("don_id must be set")
	}
	if c.GatewayURL == nil || *c.GatewayURL == "" {
		return errors.New("gateway_url must be set")
	}
	if !net.IsValidURL(*c.GatewayURL) {
		return errors.New("gateway_url must be a valid URL")
	}
	if c.Receiver == nil || *c.Receiver == "" {
		return errors.New("receiver must be set")
	}
	if !common.IsHexAddress(*c.Receiver) {
		return errors.New("receiver must be a valid address")
	}
	return nil
}

type Performance struct {
	RPS             *int64                  `toml:"rps"`
	RequestsPerCall *uint32                 `toml:"requests_per_call"`
	Duration        *blockchain.StrDuration `toml:"duration"`
}

func (c *Performance) Validate() error {
	if c.RPS == nil || *c.RPS < 1 {
		return errors.New("rps must be greater than 0")
	}
	if c.RequestsPerCall != nil && *c.RequestsPerCall < 1 {
		return errors.New("requests_per_call must be greater than 0")
	}
	if c.Duration == nil || c.Duration.Duration < 1 {
		return errors.New("duration must be greater than 0")
	}
	return nil
}

func (c *Config) Validate() error {
	if c == nil {
		return nil
	}
	if c.Common != nil {
		if err := c.Common.Validate(); err != nil {
			return err
		}
	}
	if c.Performance != nil {
		if err := c.Performance.Validate(); err != nil {
			return err
		}
	}

	return nil
}
