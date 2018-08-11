package models

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/utils"
)

// ServiceAgreement connects job specifications with on-chain encumbrances.
type ServiceAgreement struct {
	Encumbrance Encumbrance `json:"encumbrance" storm:"inline"`
	ID          string      `json:"id" storm:"id,unique"`
	JobSpecID   string      `json:"jobSpecID"`
	RequestBody string      `json:"requestBody"`
	jobSpec     JobSpec     // jobSpec is used during the initial SA creation.
	// If needed later, it can be retrieved from the database with JobSpecID.
}

// GetID returns the ID of this structure for jsonapi serialization.
func (sa ServiceAgreement) GetID() string {
	return sa.ID
}

// NewServiceAgreementFromRequest builds a new ServiceAgreement.
func NewServiceAgreementFromRequest(sar ServiceAgreementRequest) (ServiceAgreement, error) {
	id, err := generateServiceAgreementID(sar.Encumbrance, sar.Digest)

	return ServiceAgreement{
		Encumbrance: sar.Encumbrance,
		RequestBody: sar.NormalizedBody,
		ID:          id,
		jobSpec:     sar.JobSpec,
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
	Payment    *big.Int `json:"payment"`
	Expiration uint64   `json:"expiration"`
}

// ABI returns the encumbrance ABI encoded as a hex string.
func (e Encumbrance) ABI() string {
	if e.Payment == nil {
		e.Payment = big.NewInt(0)
	}
	return fmt.Sprintf("%064x%064x", e.Payment, e.Expiration)
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
	js, err := jobSpecFromSARequest(input)
	if err != nil {
		return err
	}

	var en Encumbrance
	if err := json.Unmarshal(input, &en); err != nil {
		return err
	}

	normalized, err := utils.NormalizedJSONString(input)
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

func jobSpecFromSARequest(input []byte) (JobSpec, error) {
	var jsr JobSpecRequest
	if err := json.Unmarshal(input, &jsr); err != nil {
		return JobSpec{}, err
	}

	return NewJobFromRequest(jsr)
}
