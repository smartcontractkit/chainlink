package models

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	null "gopkg.in/guregu/null.v4"
)

// Encumbrance connects job specifications with on-chain encumbrances.
type Encumbrance struct {
	// Corresponds to requestDigest in solidity ServiceAgreement struct
	ID int64 `json:"-" gorm:"primary_key;auto_increment"`
	// Price to request a report based on this agreement
	Payment *assets.Link `json:"payment,omitempty"`
	// Expiration is the amount of time an oracle has to answer a request
	Expiration uint64 `json:"expiration"`
	// Agreement is valid until this time
	EndAt AnyTime `json:"endAt"`
	// Addresses of oracles committed to this agreement
	Oracles ethkey.EIP55AddressCollection `json:"oracles" gorm:"type:text"`
	// Address of aggregator contract
	Aggregator ethkey.EIP55Address `json:"aggregator" gorm:"not null"`
	// selector for initialization method on aggregator contract
	AggInitiateJobSelector FunctionSelector `json:"aggInitiateJobSelector" gorm:"not null"`
	// selector for fulfillment (oracle reporting) method on aggregator contract
	AggFulfillSelector FunctionSelector `json:"aggFulfillSelector" gorm:"not null"`
	CreatedAt          time.Time        `json:"-"`
	UpdatedAt          time.Time        `json:"-"`
}

// UnsignedServiceAgreement contains the information to sign a service agreement
type UnsignedServiceAgreement struct {
	Encumbrance    Encumbrance
	ID             common.Hash
	RequestBody    string
	JobSpecRequest JobSpecRequest
}

// ServiceAgreement connects job specifications with on-chain encumbrances.
type ServiceAgreement struct {
	ID            string      `json:"id" gorm:"primary_key"`
	CreatedAt     time.Time   `json:"createdAt" gorm:"index"`
	Encumbrance   Encumbrance `json:"encumbrance"`
	EncumbranceID int64       `json:"-"`
	RequestBody   string      `json:"requestBody"`
	Signature     Signature   `json:"signature" gorm:"type:varchar(255)"`
	JobSpec       JobSpec     `gorm:"foreignkey:JobSpecID"`
	JobSpecID     JobID       `json:"jobSpecId"`
	UpdatedAt     time.Time   `json:"-"`
}

// ServiceAgreementRequest encodes external ServiceAgreement json representation.
type ServiceAgreementRequest struct {
	Initiators             []InitiatorRequest            `json:"initiators"`
	Tasks                  []TaskSpecRequest             `json:"tasks"`
	Payment                *assets.Link                  `json:"payment,omitempty"`
	Expiration             uint64                        `json:"expiration"`
	EndAt                  AnyTime                       `json:"endAt"`
	Oracles                ethkey.EIP55AddressCollection `json:"oracles"`
	Aggregator             ethkey.EIP55Address           `json:"aggregator"`
	AggInitiateJobSelector FunctionSelector              `json:"aggInitiateJobSelector"`
	AggFulfillSelector     FunctionSelector              `json:"aggFulfillSelector"`
	StartAt                AnyTime                       `json:"startAt"`
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
	SignHash(hash common.Hash) (Signature, error)
}

type NullSigner struct{}

func (NullSigner) SignHash(common.Hash) (Signature, error) {
	return Signature{}, nil
}

// BuildServiceAgreement builds a signed service agreement
func BuildServiceAgreement(us UnsignedServiceAgreement, signer Signer) (ServiceAgreement, error) {
	signature, err := signer.SignHash(us.ID)
	if err != nil {
		return ServiceAgreement{}, err
	}
	return ServiceAgreement{
		ID:          us.ID.String(),
		Encumbrance: us.Encumbrance,
		JobSpec:     NewJobFromRequest(us.JobSpecRequest),
		RequestBody: us.RequestBody,
		Signature:   signature,
	}, nil
}

func parseServiceAgreementJSON(reader io.Reader) (
	string, ServiceAgreementRequest, error) {
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", ServiceAgreementRequest{}, errors.Wrap(err,
			"while reading service agreement JSON")
	}
	var sar ServiceAgreementRequest
	err = json.Unmarshal(input, &sar)
	if err != nil {
		return "", ServiceAgreementRequest{}, errors.Wrap(err,
			"while parsing service agreement JSON")
	}
	normalized, err := utils.NormalizedJSON(input)
	if err != nil {
		return "", ServiceAgreementRequest{}, errors.Wrap(err,
			"while normalizing service agreement JSON")
	}
	return normalized, sar, nil
}

// NewUnsignedServiceAgreementFromRequest builds the information required to
// sign a service agreement
func NewUnsignedServiceAgreementFromRequest(reader io.Reader) (UnsignedServiceAgreement, error) {
	normalized, sar, err := parseServiceAgreementJSON(reader)
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}
	us := UnsignedServiceAgreement{
		Encumbrance: Encumbrance{
			Payment:                sar.Payment,
			Expiration:             sar.Expiration,
			EndAt:                  sar.EndAt,
			Oracles:                sar.Oracles,
			Aggregator:             sar.Aggregator,
			AggInitiateJobSelector: sar.AggInitiateJobSelector,
			AggFulfillSelector:     sar.AggFulfillSelector,
		},
		RequestBody: normalized,
		JobSpecRequest: JobSpecRequest{
			Initiators: sar.Initiators,
			Tasks:      sar.Tasks,
			StartAt:    null.NewTime(sar.StartAt.Time, sar.StartAt.Valid),
			EndAt:      null.NewTime(sar.EndAt.Time, sar.EndAt.Valid),
			MinPayment: sar.Payment,
		},
	}

	requestDigest, err := utils.Keccak256([]byte(normalized))
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}

	us.ID, err = generateSAID(us.Encumbrance,
		common.BytesToHash(requestDigest))
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}
	return us, nil
}

func generateSAID(e Encumbrance, digest common.Hash) (common.Hash, error) {
	saBytes, err := e.ABI(digest)
	if err != nil {
		return common.Hash{}, nil
	}
	bytes, err := utils.Keccak256(saBytes)
	return common.BytesToHash(bytes), err
}

// ABI packs the encumberance as a byte array using the same rules as the
// abi.encodePacked in Coordinator#getId.
//
// Used only for constructing a stable hash which will be signed by all oracles,
// so it does not have to be easily parsed or unambiguous (e.g., re-ordering
// Oracles will result in different output.) It just has to be an injective
// function.
func (e Encumbrance) ABI(digest common.Hash) ([]byte, error) {
	buffer := bytes.Buffer{}
	var paymentHash common.Hash
	if e.Payment != nil {
		paymentHash = e.Payment.ToHash()
	}
	_, err := buffer.Write(paymentHash.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "while writing payment")
	}
	expirationHash := common.BigToHash(new(big.Int).SetUint64(e.Expiration))
	_, err = buffer.Write(expirationHash.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "while writing expiration")
	}

	// Absolute end date as a big-endian uint32 (unix seconds)
	var endAt uint64
	if e.EndAt.Valid {
		endAt = uint64(e.EndAt.Time.Unix())
	}
	endAtBytes := common.BigToHash(new(big.Int).SetUint64(endAt))
	_, err = buffer.Write(endAtBytes.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "while writing endAt")
	}

	err = encodeOracles(&buffer, e.Oracles)
	if err != nil {
		return nil, errors.Wrap(err, "while writing oracles")
	}
	_, err = buffer.Write(digest.Bytes())
	if err != nil {
		return nil, err
	}
	_, err = buffer.Write(e.Aggregator.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "while writing aggregator address")
	}
	_, err = buffer.Write(e.AggInitiateJobSelector[:])
	if err != nil {
		return nil, errors.Wrap(err, "while writing aggregator initiation method selector")
	}
	_, err = buffer.Write(e.AggFulfillSelector[:])
	if err != nil {
		return nil, errors.Wrap(err, "while writing aggregator fulfill method selector")
	}
	return buffer.Bytes(), nil
}

func encodeOracles(buffer *bytes.Buffer, oracles []ethkey.EIP55Address) error {
	for _, o := range oracles {
		_, err := buffer.Write(address256Bits(o))
		if err != nil {
			return errors.Wrap(err, "while writing oracle address to buffer")
		}
	}
	return nil
}

// address256Bits Zero left-pads a to 32 bytes, as in Solidity's abi.encodePacking
func address256Bits(a ethkey.EIP55Address) []byte {
	return a.Hash().Bytes()
}
