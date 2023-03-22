package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ETHKeyResource represents a ETH key JSONAPI resource. It holds the hex
// representation of the address plus its ETH & LINK balances
type ETHKeyResource struct {
	JAID
	EVMChainID     utils.Big    `json:"evmChainID"`
	Address        string       `json:"address"`
	NextNonce      int64        `json:"nextNonce"`
	EthBalance     *assets.Eth  `json:"ethBalance"`
	LinkBalance    *assets.Link `json:"linkBalance"`
	Disabled       bool         `json:"disabled"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	MaxGasPriceWei *utils.Big   `json:"maxGasPriceWei"`
}

// GetName implements the api2go EntityNamer interface
//
// This is named as such for backwards compatibility with the operator ui
// TODO - Standardise this to ethKeys
func (r ETHKeyResource) GetName() string {
	return "eTHKeys"
}

// NewETHKeyOption defines a functional option which allows customisation of the
// EthKeyResource
type NewETHKeyOption func(*ETHKeyResource)

// NewETHKeyResource constructs a new ETHKeyResource from a Key.
//
// Use the functional options to inject the ETH and LINK balances
func NewETHKeyResource(k ethkey.KeyV2, state ethkey.State, opts ...NewETHKeyOption) *ETHKeyResource {
	r := &ETHKeyResource{
		JAID:        NewJAID(k.Address.Hex()),
		EVMChainID:  state.EVMChainID,
		NextNonce:   state.NextNonce,
		Address:     k.Address.Hex(),
		EthBalance:  nil,
		LinkBalance: nil,
		Disabled:    state.Disabled,
		CreatedAt:   state.CreatedAt,
		UpdatedAt:   state.UpdatedAt,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func SetETHKeyEthBalance(ethBalance *assets.Eth) NewETHKeyOption {
	return func(r *ETHKeyResource) {
		r.EthBalance = ethBalance
	}
}

func SetETHKeyLinkBalance(linkBalance *assets.Link) NewETHKeyOption {
	return func(r *ETHKeyResource) {
		r.LinkBalance = linkBalance
	}
}

func SetETHKeyMaxGasPriceWei(maxGasPriceWei *utils.Big) NewETHKeyOption {
	return func(r *ETHKeyResource) {
		r.MaxGasPriceWei = maxGasPriceWei
	}
}
