package models

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/utils"
)

// ServiceAgreement connects job specifications with on-chain encumbrances.
type ServiceAgreement struct {
	CreatedAt   Time        `json:"createdAt" storm:"index"`
	Encumbrance Encumbrance `json:"encumbrance" storm:"inline"`
	ID          common.Hash `json:"id" storm:"id,unique"`
	JobSpecID   string      `json:"jobSpecID"`
	RequestBody string      `json:"requestBody"`
	Signature   common.Hash `json:"signature"`
	jobSpec     JobSpec     // jobSpec is used during the initial SA creation.
	// If needed later, it can be retrieved from the database with JobSpecID.
}

// GetID returns the ID of this structure for jsonapi serialization.
func (sa ServiceAgreement) GetID() string {
	return sa.ID.String()
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (sa ServiceAgreement) GetName() string {
	return "service_agreements"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (sa *ServiceAgreement) SetID(value string) error {
	sa.ID.SetString(value)
	return nil
}

// Signer is used to produce a HMAC signature from an input digest
type Signer interface {
	Sign(input []byte) ([]byte, error)
}

// NewServiceAgreementFromRequest builds a new ServiceAgreement.
func NewServiceAgreementFromRequest(reader io.Reader, signer Signer) (ServiceAgreement, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return ServiceAgreement{}, err
	}

	var sar ServiceAgreementRequest
	err = json.Unmarshal(input, &sar)
	if err != nil {
		return ServiceAgreement{}, err
	}

	normalized, err := utils.NormalizedJSON(input)
	if err != nil {
		return ServiceAgreement{}, err
	}

	requestDigest, err := utils.Keccak256([]byte(normalized))
	//logger.Debug("SA requestDigest", requestDigest)
	fmt.Println("SA requestDigest", requestDigest)
	digest := common.ToHex(requestDigest)

	encumbrance := Encumbrance{
		Payment:    sar.Payment,
		Expiration: sar.Expiration,
		Oracles:    sar.Oracles,
	}

	id, err := generateServiceAgreementID(encumbrance, digest)
	if err != nil {
		return ServiceAgreement{}, err
	}

	//logger.Debug("SAID", id)
	fmt.Println("SAID", id)

	signature, err := signer.Sign(id.Bytes())
	if err != nil {
		return ServiceAgreement{}, err
	}

	jobSpec := NewJob()
	jobSpec.Initiators = sar.Initiators
	jobSpec.Tasks = sar.Tasks
	jobSpec.EndAt = sar.EndAt
	jobSpec.StartAt = sar.StartAt

	return ServiceAgreement{
		ID:          id,
		CreatedAt:   Time{time.Now()},
		Encumbrance: encumbrance,
		jobSpec:     jobSpec,
		RequestBody: normalized,
		Signature:   common.BytesToHash(signature),
	}, nil
}

func generateServiceAgreementID(e Encumbrance, digest string) (common.Hash, error) {
	b, err := utils.HexToBytes(e.ABI(), digest)
	if err != nil {
		return common.Hash{}, err
	}
	bytes, err := utils.Keccak256(b)
	return common.BytesToHash(bytes), err
}

// Encumbrance connects job specifications with on-chain encumbrances.
type Encumbrance struct {
	Payment    *assets.Link   `json:"payment"`
	Expiration uint64         `json:"expiration"`
	Oracles    []EIP55Address `json:"oracles"`
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
	Payment    *assets.Link   `json:"payment"`
	Expiration uint64         `json:"expiration"`
	Oracles    []EIP55Address `json:"oracles"`
	JobSpecRequest
}
