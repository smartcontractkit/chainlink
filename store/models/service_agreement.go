package models

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/utils"
)

// UnsignedServiceAgreement contains the information to sign a service agreement
type UnsignedServiceAgreement struct {
	Encumbrance             Encumbrance
	ID                      common.Hash
	RequestBody             string
	ServiceAgreementRequest ServiceAgreementRequest
}

// ServiceAgreement connects job specifications with on-chain encumbrances.
type ServiceAgreement struct {
	CreatedAt   Time        `json:"createdAt" storm:"index"`
	Encumbrance Encumbrance `json:"encumbrance" storm:"inline"`
	ID          string      `json:"id" storm:"id,unique"`
	JobSpecID   string      `json:"jobSpecID"`
	RequestBody string      `json:"requestBody"`
	Signature   Signature   `json:"signature"`
	JobSpec     JobSpec     // JobSpec is used during the initial SA creation.
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
	//sa.ID.SetString(value)
	sa.ID = value
	return nil
}

// Signer is used to produce a HMAC signature from an input digest
type Signer interface {
	Sign(input []byte) (Signature, error)
}

// BuildServiceAgreement builds a signed service agreement
func BuildServiceAgreement(us UnsignedServiceAgreement, signer Signer) (ServiceAgreement, error) {
	signature, err := signer.Sign(us.ID.Bytes())
	if err != nil {
		return ServiceAgreement{}, err
	}

	jobSpec := NewJob()
	jobSpec.Initiators = us.ServiceAgreementRequest.Initiators
	jobSpec.Tasks = us.ServiceAgreementRequest.Tasks
	jobSpec.EndAt = us.ServiceAgreementRequest.EndAt
	jobSpec.StartAt = us.ServiceAgreementRequest.StartAt

	return ServiceAgreement{
		ID:          us.ID.String(),
		CreatedAt:   Time{time.Now()},
		Encumbrance: us.Encumbrance,
		JobSpec:     jobSpec,
		RequestBody: us.RequestBody,
		Signature:   signature,
	}, nil
}

// NewUnsignedServiceAgreementFromRequest builds the information required to
// sign a service agreement
func NewUnsignedServiceAgreementFromRequest(reader io.Reader) (UnsignedServiceAgreement, error) {
	var sar ServiceAgreementRequest

	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}

	err = json.Unmarshal(input, &sar)
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}

	normalized, err := utils.NormalizedJSON(input)
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}

	requestDigest, err := utils.Keccak256([]byte(normalized))
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}

	encumbrance := Encumbrance{
		Payment:    sar.Payment,
		Expiration: sar.Expiration,
		Oracles:    sar.Oracles,
	}
	id, err := generateServiceAgreementID(encumbrance, common.BytesToHash(requestDigest))
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}

	us := UnsignedServiceAgreement{
		ID:                      id,
		Encumbrance:             encumbrance,
		RequestBody:             normalized,
		ServiceAgreementRequest: sar,
	}

	return us, err
}

func generateServiceAgreementID(e Encumbrance, digest common.Hash) (common.Hash, error) {
	buffer, err := serviceAgreementIDInputBuffer(e, digest)
	if err != nil {
		return common.Hash{}, nil
	}

	bytes, err := utils.Keccak256(buffer.Bytes())
	return common.BytesToHash(bytes), err
}

func serviceAgreementIDInputBuffer(encumbrance Encumbrance, digest common.Hash) (bytes.Buffer, error) {
	buffer := bytes.Buffer{}

	encumberanceBytes, err := encumbrance.ABI()
	if err != nil {
		return bytes.Buffer{}, err
	}
	_, err = buffer.Write(encumberanceBytes)
	if err != nil {
		return bytes.Buffer{}, err
	}

	_, err = buffer.Write(digest.Bytes())
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buffer, nil
}

// Encumbrance connects job specifications with on-chain encumbrances.
type Encumbrance struct {
	Payment    *assets.Link   `json:"payment"`
	Expiration uint64         `json:"expiration"`
	Oracles    []EIP55Address `json:"oracles"`
}

// ABI packs the encumberance as a byte array using the same technique as
// abi.encodePacked, meaning that addresses are padded with left 0s to match
// hashes in the oracle list
func (e Encumbrance) ABI() ([]byte, error) {
	buffer := bytes.Buffer{}
	var paymentHash common.Hash
	if e.Payment != nil {
		paymentHash = e.Payment.ToHash()
	}
	_, err := buffer.Write(paymentHash.Bytes())
	if err != nil {
		return []byte{}, err
	}
	expirationHash := common.BigToHash((*big.Int)(new(big.Int).SetUint64(e.Expiration)))
	_, err = buffer.Write(expirationHash.Bytes())
	if err != nil {
		return []byte{}, err
	}

	err = encodeOracles(&buffer, e.Oracles)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func encodeOracles(buffer *bytes.Buffer, oracles []EIP55Address) error {
	for _, o := range oracles {
		// XXX: Solidity packs addresses as hashes when doing abi.encodePacking, so mirror here
		oracleAddressHash := common.BytesToHash(o.Bytes())
		_, err := buffer.Write(oracleAddressHash.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

// ServiceAgreementRequest represents a service agreement as requested over the wire.
type ServiceAgreementRequest struct {
	Payment    *assets.Link   `json:"payment"`
	Expiration uint64         `json:"expiration"`
	Oracles    []EIP55Address `json:"oracles"`
	JobSpecRequest
}
