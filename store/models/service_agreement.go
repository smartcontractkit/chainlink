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

	encumbrance := Encumbrance{
		Payment:    sar.Payment,
		Expiration: sar.Expiration,
		Oracles:    sar.Oracles,
	}

	id, err := generateServiceAgreementID(encumbrance, common.BytesToHash(requestDigest))
	if err != nil {
		return ServiceAgreement{}, err
	}

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
		ID:          id.String(),
		CreatedAt:   Time{time.Now()},
		Encumbrance: encumbrance,
		JobSpec:     jobSpec,
		RequestBody: normalized,
		Signature:   signature,
	}, nil
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
