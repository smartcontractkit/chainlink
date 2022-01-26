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
	EthBalance     *assets.Eth  `json:"ethBalance"`
	LinkBalance    *assets.Link `json:"linkBalance"`
	IsFunding      bool         `json:"isFunding"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	MaxGasPriceWei utils.Big    `json:"maxGasPriceWei"`
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
type NewETHKeyOption func(*ETHKeyResource) error

// NewETHKeyResource constructs a new ETHKeyResource from a Key.
//
// Use the functional options to inject the ETH and LINK balances
func NewETHKeyResource(k ethkey.KeyV2, state ethkey.State, opts ...NewETHKeyOption) (*ETHKeyResource, error) {
	r := &ETHKeyResource{
		JAID:        NewJAID(k.Address.Hex()),
		EVMChainID:  state.EVMChainID,
		Address:     k.Address.Hex(),
		EthBalance:  nil,
		LinkBalance: nil,
		IsFunding:   state.IsFunding,
		CreatedAt:   state.CreatedAt,
		UpdatedAt:   state.UpdatedAt,
	}

	for _, opt := range opts {
		err := opt(r)

		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func SetETHKeyEthBalance(ethBalance *assets.Eth) NewETHKeyOption {
	return func(r *ETHKeyResource) error {
		r.EthBalance = ethBalance

		return nil
	}
}

func SetETHKeyLinkBalance(linkBalance *assets.Link) NewETHKeyOption {
	return func(r *ETHKeyResource) error {
		r.LinkBalance = linkBalance

		return nil
	}
}

func SetETHKeyMaxGasPriceWei(maxGasPriceWei utils.Big) NewETHKeyOption {
	return func(r *ETHKeyResource) error {
		r.MaxGasPriceWei = maxGasPriceWei

		return nil
	}
}
