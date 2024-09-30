package ocr

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
)

type Config struct {
	Soak      *SoakConfig `toml:"Soak"`
	Load      *Load       `toml:"Load"`
	Volume    *Volume     `toml:"Volume"`
	Common    *Common     `toml:"Common"`
	Contracts *Contracts  `toml:"Contracts"`
}

func (o *Config) Validate() error {
	if o.Common != nil {
		if err := o.Common.Validate(); err != nil {
			return err
		}
	}
	if o.Soak != nil {
		if err := o.Soak.Validate(); err != nil {
			return err
		}
	}
	if o.Volume != nil {
		if err := o.Volume.Validate(); err != nil {
			return err
		}
	}
	if o.Contracts != nil {
		if err := o.Contracts.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type Common struct {
	NumberOfContracts *int                    `toml:"number_of_contracts"`
	ETHFunds          *int                    `toml:"eth_funds"`
	TestDuration      *blockchain.StrDuration `toml:"test_duration"`
}

func (o *Common) Validate() error {
	if o.NumberOfContracts != nil && *o.NumberOfContracts < 1 {
		return errors.New("when number_of_contracts is set, it must be greater than 0")
	}
	if o.ETHFunds != nil && *o.ETHFunds < 0 {
		return errors.New("eth_funds must be set and cannot be negative")
	}
	return nil
}

type Load struct {
	Rate                  *int64                  `toml:"rate"`
	RequestsPerUnit       *int                    `toml:"requests_per_unit"`
	RateLimitUnitDuration *blockchain.StrDuration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *blockchain.StrDuration `toml:"verification_interval"`
	VerificationTimeout   *blockchain.StrDuration `toml:"verification_timeout"`
	EAChangeInterval      *blockchain.StrDuration `toml:"ea_change_interval"`
	TestDuration          *blockchain.StrDuration `toml:"test_duration"`
}

func (o *Load) Validate() error {
	if o.TestDuration == nil {
		return errors.New("load test duration must be set")
	}
	if o.Rate == nil || *o.Rate <= 0 {
		return errors.New("rate must be set and be a positive integer")
	}
	if o.RequestsPerUnit == nil || *o.RequestsPerUnit <= 0 {
		return errors.New("vu_requests_per_unit must be set and be a positive integer")
	}
	if o.RateLimitUnitDuration == nil || o.RateLimitUnitDuration.Duration == 0 {
		return errors.New("rate_limit_unit_duration must be set and be a positive integer")
	}
	if o.VerificationInterval == nil || o.VerificationInterval.Duration == 0 {
		return errors.New("verification_interval must be set and be a positive integer")
	}
	if o.VerificationTimeout == nil || o.VerificationTimeout.Duration == 0 {
		return errors.New("verification_timeout must be set and be a positive integer")
	}
	if o.EAChangeInterval == nil || o.EAChangeInterval.Duration == 0 {
		return errors.New("ea_change_interval must be set and be a positive integer")
	}

	return nil
}

type Volume struct {
	Rate                  *int64                  `toml:"rate"`
	VURequestsPerUnit     *int                    `toml:"vu_requests_per_unit"`
	RateLimitUnitDuration *blockchain.StrDuration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *blockchain.StrDuration `toml:"verification_interval"`
	VerificationTimeout   *blockchain.StrDuration `toml:"verification_timeout"`
	EAChangeInterval      *blockchain.StrDuration `toml:"ea_change_interval"`
	TestDuration          *blockchain.StrDuration `toml:"test_duration"`
}

func (o *Volume) Validate() error {
	if o.TestDuration == nil {
		return errors.New("volume test duration must be set")
	}
	if o.Rate == nil || *o.Rate <= 0 {
		return errors.New("rate must be set and be a positive integer")
	}
	if o.VURequestsPerUnit == nil || *o.VURequestsPerUnit <= 0 {
		return errors.New("vu_requests_per_unit must be set and be a positive integer")
	}
	if o.RateLimitUnitDuration == nil || o.RateLimitUnitDuration.Duration == 0 {
		return errors.New("rate_limit_unit_duration must be set and be a positive integer")
	}
	if o.VerificationInterval == nil || o.VerificationInterval.Duration == 0 {
		return errors.New("verification_interval must be set and be a positive integer")
	}
	if o.VerificationTimeout == nil || o.VerificationTimeout.Duration == 0 {
		return errors.New("verification_timeout must be set and be a positive integer")
	}
	if o.EAChangeInterval == nil || o.EAChangeInterval.Duration == 0 {
		return errors.New("ea_change_interval must be set and be a positive integer")
	}

	return nil
}

type SoakConfig struct {
	TimeBetweenRounds *blockchain.StrDuration `toml:"time_between_rounds"`
}

func (o *SoakConfig) Validate() error {
	if o.TimeBetweenRounds == nil || o.TimeBetweenRounds.Duration == 0 {
		return errors.New("time_between_rounds must be set and be a positive integer")
	}
	return nil
}

// For more information on the configuration of contracts, see https://smartcontract-it.atlassian.net/wiki/spaces/TT/pages/828407894/Contracts+addresses+in+TOML+convention
type Contracts struct {
	ShouldBeUsed                *bool                      `toml:"use"`
	LinkTokenAddress            *string                    `toml:"link_token"`
	OffChainAggregatorAddresses []string                   `toml:"offchain_aggregators"`
	Settings                    map[string]ContractSetting `toml:"Settings"`
}

func (o *Contracts) Validate() error {
	if o.LinkTokenAddress != nil && !common.IsHexAddress(*o.LinkTokenAddress) {
		return errors.New("link_token must be a valid ethereum address")
	}
	if o.OffChainAggregatorAddresses != nil {
		allEnabled := make(map[bool]int)
		allConfigure := make(map[bool]int)
		for _, address := range o.OffChainAggregatorAddresses {
			if !common.IsHexAddress(address) {
				return fmt.Errorf("offchain_aggregators must be valid ethereum addresses, but %s is not", address)
			}

			if v, ok := o.Settings[address]; ok {
				if v.ShouldBeUsed != nil {
					allEnabled[*v.ShouldBeUsed]++
				} else {
					allEnabled[true]++
				}
				if v.Configure != nil {
					allConfigure[*v.Configure]++
				} else {
					allConfigure[true]++
				}
			}
		}

		if allEnabled[true] > 0 && allEnabled[false] > 0 {
			return errors.New("either all or none offchain_aggregators must be used")
		}

		if allConfigure[true] > 0 && allConfigure[false] > 0 {
			return errors.New("either all or none offchain_aggregators must be configured")
		}
	}

	return nil
}

func (o *Config) UseExistingContracts() bool {
	if o.Contracts == nil {
		return false
	}

	if o.Contracts.ShouldBeUsed != nil {
		return *o.Contracts.ShouldBeUsed
	}

	return false
}

func (o *Config) LinkTokenContractAddress() (common.Address, error) {
	if o.Contracts != nil && o.Contracts.LinkTokenAddress != nil {
		return common.HexToAddress(*o.Contracts.LinkTokenAddress), nil
	}

	return common.Address{}, errors.New("link token address must be set")
}

func (o *Config) UseExistingLinkTokenContract() bool {
	if !o.UseExistingContracts() {
		return false
	}

	if o.Contracts.LinkTokenAddress == nil {
		return false
	}

	if len(o.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := o.Contracts.Settings[*o.Contracts.LinkTokenAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

type ContractSetting struct {
	ShouldBeUsed *bool `toml:"use"`
	Configure    *bool `toml:"configure"`
}

type OffChainAggregatorsConfig interface {
	OffChainAggregatorsContractsAddresses() []common.Address
	UseExistingOffChainAggregatorsContracts() bool
	ConfigureExistingOffChainAggregatorsContracts() bool
	NumberOfContractsToDeploy() int
}

func (o *Config) UseExistingOffChainAggregatorsContracts() bool {
	if !o.UseExistingContracts() {
		return false
	}

	if len(o.Contracts.OffChainAggregatorAddresses) == 0 {
		return false
	}

	if len(o.Contracts.Settings) == 0 {
		return true
	}

	for _, address := range o.Contracts.OffChainAggregatorAddresses {
		if v, ok := o.Contracts.Settings[address]; ok {
			return v.ShouldBeUsed != nil && *v.ShouldBeUsed
		}
	}

	return true
}

func (o *Config) OffChainAggregatorsContractsAddresses() []common.Address {
	var ocrInstanceAddresses []common.Address
	if !o.UseExistingOffChainAggregatorsContracts() {
		return ocrInstanceAddresses
	}

	for _, address := range o.Contracts.OffChainAggregatorAddresses {
		ocrInstanceAddresses = append(ocrInstanceAddresses, common.HexToAddress(address))
	}

	return ocrInstanceAddresses
}

func (o *Config) ConfigureExistingOffChainAggregatorsContracts() bool {
	if !o.UseExistingOffChainAggregatorsContracts() {
		return true
	}

	for _, address := range o.Contracts.OffChainAggregatorAddresses {
		for maybeOcrAddress, setting := range o.Contracts.Settings {
			if maybeOcrAddress == address {
				return setting.Configure != nil && *setting.Configure
			}
		}
	}

	return true
}

func (o *Config) NumberOfContractsToDeploy() int {
	if o.Common != nil && o.Common.NumberOfContracts != nil {
		return *o.Common.NumberOfContracts
	}

	return 0
}
