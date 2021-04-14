package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

// ETHKeyResource represents a ETH key JSONAPI resource. It holds the hex
// representation of the address plus it's ETH & LINK balances
type ETHKeyResource struct {
	JAID
	Address     string         `json:"address"`
	EthBalance  *assets.Eth    `json:"ethBalance"`
	LinkBalance *assets.Link   `json:"linkBalance"`
	NextNonce   int64          `json:"nextNonce"`
	LastUsed    *time.Time     `json:"lastUsed"`
	IsFunding   bool           `json:"isFunding"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
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
func NewETHKeyResource(k models.Key, opts ...NewETHKeyOption) (*ETHKeyResource, error) {
	r := &ETHKeyResource{
		JAID:        NewJAID(k.Address.Hex()),
		Address:     k.Address.Hex(),
		EthBalance:  nil,
		LinkBalance: nil,
		NextNonce:   k.NextNonce,
		LastUsed:    k.LastUsed,
		IsFunding:   k.IsFunding,
		CreatedAt:   k.CreatedAt,
		UpdatedAt:   k.UpdatedAt,
		DeletedAt:   k.DeletedAt,
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
