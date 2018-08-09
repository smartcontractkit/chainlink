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
	Normalized  string      `json:"normalizedRequest"`
	jobSpec     JobSpec     // jobSpec is used during the initial SA creation.
	// If needed later, it can be retrieved from the database with JobSpecID.
}

// GetID returns the ID of this structure for jsonapi serialization.
func (sa ServiceAgreement) GetID() string {
	return sa.ID
}

// NewServiceAgreementFromRequest builds a new ServiceAgreement.
func NewServiceAgreementFromRequest(sar ServiceAgreementRequest) (ServiceAgreement, error) {
	sa := ServiceAgreement{}

	sa.Encumbrance = sar.Encumbrance
	sa.jobSpec = sar.JobSpec

	b, err := utils.HexToBytes(sa.Encumbrance.ABI(), sa.jobSpec.Digest)
	if err != nil {
		return sa, err
	}
	digest, err := utils.Keccak256(b)
	if err != nil {
		return sa, err
	}
	sa.ID = common.ToHex(digest)
	sa.Normalized = sar.Normalized

	return sa, nil
}

// Encumbrance connects job specifications with on-chain encumbrances.
type Encumbrance struct {
	Payment    *big.Int `json:"payment"`
	Expiration *big.Int `json:"expiration"`
}

// ABI returns the encumbrance ABI encoded as a hex string.
func (e Encumbrance) ABI() string {
	if e.Payment == nil {
		e.Payment = big.NewInt(0)
	}
	if e.Expiration == nil {
		e.Expiration = big.NewInt(0)
	}
	return fmt.Sprintf("%064x%064x", e.Payment, e.Expiration)
}

// ServiceAgreementRequest represents a service agreement as requested over the wire.
type ServiceAgreementRequest struct {
	JobSpec     JobSpec
	Encumbrance Encumbrance
	Normalized  string
}

// UnmarshalJSON fulfills Go's built in JSON unmarshaling interface.
func (sar *ServiceAgreementRequest) UnmarshalJSON(input []byte) error {
	var jsr JobSpecRequest
	if err := json.Unmarshal(input, &jsr); err != nil {
		return err
	}

	js, err := NewJobFromRequest(jsr)
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

	sar.JobSpec = js
	sar.Encumbrance = en
	sar.Normalized = normalized

	return nil
}
