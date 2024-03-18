package presenters

import (
	"time"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

// ETHKeyResource represents a ETH key JSONAPI resource. It holds the hex
// representation of the address plus its ETH & LINK balances
type ETHKeyResource struct {
	JAID
	EVMChainID     big.Big            `json:"evmChainID"`
	Address        string             `json:"address"`
	EthBalance     *assets.Eth        `json:"ethBalance"`
	LinkBalance    *commonassets.Link `json:"linkBalance"`
	Disabled       bool               `json:"disabled"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
	MaxGasPriceWei *big.Big           `json:"maxGasPriceWei"`
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
		JAID:        NewPrefixedJAID(k.Address.Hex(), state.EVMChainID.String()),
		EVMChainID:  state.EVMChainID,
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

func SetETHKeyLinkBalance(linkBalance *commonassets.Link) NewETHKeyOption {
	return func(r *ETHKeyResource) {
		r.LinkBalance = linkBalance
	}
}

func SetETHKeyMaxGasPriceWei(maxGasPriceWei *big.Big) NewETHKeyOption {
	return func(r *ETHKeyResource) {
		r.MaxGasPriceWei = maxGasPriceWei
	}
}
