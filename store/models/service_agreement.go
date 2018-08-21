package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/utils"
)

// ServiceAgreement connects job specifications with on-chain encumbrances.
type ServiceAgreement struct {
	CreatedAt   Time        `json:"createdAt" storm:"index"`
	Encumbrance Encumbrance `json:"encumbrance" storm:"inline"`
	ID          string      `json:"id" storm:"id,unique"`
	JobSpecID   string      `json:"jobSpecID"`
	RequestBody string      `json:"requestBody"`
	Signature   string      `json:"signature"`
	jobSpec     JobSpec     // jobSpec is used during the initial SA creation.
	// If needed later, it can be retrieved from the database with JobSpecID.
}

// GetID returns the ID of this structure for jsonapi serialization.
func (sa ServiceAgreement) GetID() string {
	return sa.ID
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (sa ServiceAgreement) GetName() string {
	return "service_agreements"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (sa *ServiceAgreement) SetID(value string) error {
	sa.ID = value
	return nil
}

// Signer is used to produce a HMAC signature from an input digest
type Signer interface {
	Sign(input []byte) (string, error)
}

// NewServiceAgreementFromRequest builds a new ServiceAgreement.
func NewServiceAgreementFromRequest(sar ServiceAgreementRequest, signer Signer) (ServiceAgreement, error) {
	id, err := generateServiceAgreementID(sar.Encumbrance, sar.Digest)
	if err != nil {
		return ServiceAgreement{}, err
	}

	signature, err := signer.Sign([]byte(id))
	if err != nil {
		return ServiceAgreement{}, err
	}

	return ServiceAgreement{
		CreatedAt:   Time{time.Now()},
		Encumbrance: sar.Encumbrance,
		ID:          id,
		jobSpec:     sar.JobSpec,
		RequestBody: sar.NormalizedBody,
		Signature:   signature,
	}, err
}

func generateServiceAgreementID(e Encumbrance, digest string) (string, error) {
	b, err := utils.HexToBytes(e.ABI(), digest)
	if err != nil {
		return "", err
	}
	bytesID, err := utils.Keccak256(b)
	return common.ToHex(bytesID), err
}

// Encumbrance connects job specifications with on-chain encumbrances.
type Encumbrance struct {
	Payment    *assets.Link `json:"payment"`
	Expiration uint64       `json:"expiration"`
}

// ABI returns the encumbrance ABI encoded as a hex string.
func (e Encumbrance) ABI() string {
	payment := e.Payment
	if payment == nil {
		payment = assets.NewLink(0)
	}
	return fmt.Sprintf("%064s%064x", payment.Text(16), e.Expiration)
}

// ServiceAgreementRequest represents a service agreement as requested over the wire.
type ServiceAgreementRequest struct {
	JobSpec        JobSpec
	Encumbrance    Encumbrance
	NormalizedBody string
	Digest         string
}

// UnmarshalJSON fulfills Go's built in JSON unmarshaling interface.
func (sar *ServiceAgreementRequest) UnmarshalJSON(input []byte) error {
	js := NewJob()
	if err := json.Unmarshal(input, &js); err != nil {
		return err
	}

	var en Encumbrance
	if err := json.Unmarshal(input, &en); err != nil {
		return err
	}

	normalized, err := utils.NormalizedJSON(input)
	if err != nil {
		return err
	}
	requestDigest, err := utils.Keccak256([]byte(normalized))

	sar.JobSpec = js
	sar.Encumbrance = en
	sar.NormalizedBody = normalized
	sar.Digest = common.ToHex(requestDigest)
	return err
}
