package models

import (
	"crypto/subtle"
	"fmt"

	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// BridgeTypeRequest is the incoming record used to create a BridgeType
type BridgeTypeRequest struct {
	Name                   TaskType    `json:"name"`
	URL                    WebURL      `json:"url"`
	Confirmations          uint64      `json:"confirmations"`
	OutgoingToken          string      `json:"outgoingToken"`
	MinimumContractPayment assets.Link `json:"minimumContractPayment"`
}

// BridgeType is used for external adapters and has fields for
// the name of the adapter and its URL.
type BridgeType struct {
	Name                   TaskType `json:"name" gorm:"primary_key"`
	URL                    WebURL   `json:"url"`
	Confirmations          uint64   `json:"confirmations"`
	IncomingTokenHash      string
	Salt                   string
	OutgoingToken          string
	MinimumContractPayment assets.Link `json:"minimumContractPayment" gorm:"type:varchar(255)"`
}

// BridgeAuthentication is the record returned in response to a request to create a BridgeType
type BridgeTypeAuthentication struct {
	Name                   TaskType `json:"name"`
	URL                    WebURL   `json:"url"`
	Confirmations          uint64   `json:"confirmations"`
	IncomingToken          string
	MinimumContractPayment assets.Link `json:"minimumContractPayment"`
}

// NewBridgeType returns a bridge bridge type authentication (with plaintext
// password) and a bridge type (with hashed password, for persisting)
func NewBridgeType(btr *BridgeTypeRequest) (*BridgeTypeAuthentication,
	*BridgeType, error) {
	//if err := services.ValidateAdapter(bt, store); err != nil {
	//return nil, nil, NewValidationError(err.Error())
	//}

	incomingToken := utils.NewBytes32ID()
	outgoingToken := utils.NewBytes32ID()
	hash, err := utils.Sha256(incomingToken)
	if err != nil {
		return nil, nil, err
	}

	return &BridgeTypeAuthentication{
			Name:                   btr.Name,
			URL:                    btr.URL,
			Confirmations:          btr.Confirmations,
			IncomingToken:          incomingToken,
			MinimumContractPayment: btr.MinimumContractPayment,
		}, &BridgeType{
			Name:                   btr.Name,
			URL:                    btr.URL,
			Confirmations:          btr.Confirmations,
			IncomingTokenHash:      hash,
			OutgoingToken:          outgoingToken,
			MinimumContractPayment: btr.MinimumContractPayment,
		}, nil
}

// GetID returns the ID of this structure for jsonapi serialization.
func (bt BridgeType) GetID() string {
	return bt.Name.String()
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (bt BridgeType) GetName() string {
	return "bridges"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (bt *BridgeType) SetID(value string) error {
	name, err := NewTaskType(value)
	bt.Name = name
	return err
}

// AuthenticateBridgeType returns true if the passed token matches its
// IncomingToken, or returns false with an error.
func AuthenticateBridgeType(bt *BridgeType, token string) (bool, error) {
	input := fmt.Sprintf("%s-%s", token, bt.Salt)
	hash, err := utils.Sha256(input)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare([]byte(hash), []byte(bt.IncomingTokenHash)) == 1, nil
}
