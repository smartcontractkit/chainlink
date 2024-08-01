package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = &Params{}

// MaxNumOracles denotes number of oracles the offchain reporting protocol is designed for
var MaxNumOracles = 31

// Parameter keys
var (
	KeyLinkDenom      = []byte("LinkDenom")
	KeyPayoutInterval = []byte("PayoutInterval")
	KeyModuleAdmin    = []byte("ModuleAdmin")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	linkDenom string,
	payoutInterval uint64,
	moduleAdmin sdk.AccAddress,
) Params {
	return Params{
		LinkDenom:           linkDenom,
		PayoutBlockInterval: payoutInterval,
		ModuleAdmin:         moduleAdmin.String(),
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLinkDenom, &p.LinkDenom, validateLinkDenom),
		paramtypes.NewParamSetPair(KeyPayoutInterval, &p.PayoutBlockInterval, validatePayoutInterval),
		paramtypes.NewParamSetPair(KeyModuleAdmin, &p.ModuleAdmin, validateModuleAdmin),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		LinkDenom:           "peggy0x514910771AF9Ca656af840dff83E8264EcF986CA",
		PayoutBlockInterval: 100000,
		ModuleAdmin:         "",
	}
}

// Validate performs basic validation on insurance parameters.
func (p Params) Validate() error {
	return validateLinkDenom(p.LinkDenom)
}

func validateLinkDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("linkDenom param cannot be empty: %v", v)
	}

	return nil
}

func validatePayoutInterval(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("PayoutInterval must be positive: %d", v)
	}

	return nil
}

func validateModuleAdmin(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return nil
	}

	if _, err := sdk.AccAddressFromBech32(v); err != nil {
		return err
	}

	return nil
}
