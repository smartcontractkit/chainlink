package forms

import (
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// NewUpdateBridgeType initializes the form attributes with the existing
// attributes from the bridge
func NewUpdateBridgeType(store *store.Store, bridgeName string) (UpdateBridgeType, error) {
	bt, err := store.FindBridge(bridgeName)
	if err != nil {
		return UpdateBridgeType{}, err
	}

	form := UpdateBridgeType{
		store:                  store,
		bridgeName:             bridgeName,
		URL:                    bt.URL,
		Confirmations:          bt.Confirmations,
		MinimumContractPayment: bt.MinimumContractPayment,
	}
	return form, nil
}

// UpdateBridgeType whitelists attributes that can be updated on the bridge
type UpdateBridgeType struct {
	store                  *store.Store
	bridgeName             string
	URL                    models.WebURL `json:"url"`
	Confirmations          uint64        `json:"confirmations"`
	MinimumContractPayment assets.Link   `json:"minimumContractPayment"`
}

// Save updates the whitelisted attributes on the bridge
func (ubt UpdateBridgeType) Save() error {
	bt, err := ubt.findBridge()
	if err != nil {
		return err
	}

	bt.URL = ubt.URL
	bt.Confirmations = ubt.Confirmations
	bt.MinimumContractPayment = ubt.MinimumContractPayment
	return ubt.store.SaveBridgeType(&bt)
}

// Marshal encodes the bridge with the JSON-API presenter
func (ubt UpdateBridgeType) Marshal() ([]byte, error) {
	bt, err := ubt.findBridge()
	if err != nil {
		return []byte{}, err
	}

	return jsonapi.Marshal(
		presenters.BridgeType{BridgeType: bt},
	)
}

func (ubt UpdateBridgeType) findBridge() (models.BridgeType, error) {
	return ubt.store.FindBridge(ubt.bridgeName)
}
