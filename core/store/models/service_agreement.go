package models

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/utils"
)

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
	EncumbranceID uint        `json:"-"`
	RequestBody   string      `json:"requestBody"`
	Signature     Signature   `json:"signature" gorm:"type:varchar(255)"`
	JobSpec       JobSpec     `gorm:"foreignkey:JobSpecID"`
	JobSpecID     *ID         `json:"jobSpecId"`
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
	Sign(input []byte) (Signature, error)
}

// BuildServiceAgreement builds a signed service agreement
func BuildServiceAgreement(us UnsignedServiceAgreement, signer Signer) (ServiceAgreement, error) {
	signature, err := signer.Sign(us.ID.Bytes())
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

// NewUnsignedServiceAgreementFromRequest builds the information required to
// sign a service agreement
func NewUnsignedServiceAgreementFromRequest(reader io.Reader) (UnsignedServiceAgreement, error) {
	var jsr JobSpecRequest

	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return UnsignedServiceAgreement{}, errors.Wrap(err, "while reading service agreement JSON")
	}

	err = json.Unmarshal(input, &jsr)
	if err != nil {
		return UnsignedServiceAgreement{}, errors.Wrap(err, "while parsing job-spec JSON")
	}

	var encumbrance Encumbrance
	if err := json.Unmarshal(input, &encumbrance); err != nil {
		return UnsignedServiceAgreement{}, errors.Wrap(err, "while parsing service agreement JSON")
	}

	if encumbrance.Aggregator == "" {
		return UnsignedServiceAgreement{}, fmt.Errorf("must set aggregator contract address")
	}

	normalized, err := utils.NormalizedJSON(input)
	if err != nil {
		return UnsignedServiceAgreement{}, errors.Wrap(err, "while normalizing service agreement JSON")
	}

	requestDigest, err := utils.Keccak256([]byte(normalized))
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}

	id, err := generateServiceAgreementID(encumbrance, common.BytesToHash(requestDigest))
	if err != nil {
		return UnsignedServiceAgreement{}, err
	}

	us := UnsignedServiceAgreement{
		ID:             id,
		Encumbrance:    encumbrance,
		RequestBody:    normalized,
		JobSpecRequest: jsr,
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
	// Corresponds to requestDigest in solidity ServiceAgreement struct
	ID uint `json:"-" gorm:"primary_key;auto_increment"`
	// Price to request a report based on this agreement
	Payment    *assets.Link `json:"payment" gorm:"type:varchar(255)"`
	Expiration uint64       `json:"expiration"`
	// Agreement is valid until this time
	EndAt AnyTime `json:"endAt"`
	// Addresses of oracles committed to this agreement
	Oracles EIP55AddressCollection `json:"oracles" gorm:"type:text"`
	// Address of aggregator contract
	Aggregator EIP55Address `json:"aggregator" gorm:"not null"`
	// selector for initialization method on aggregator contract
	AggInitiateJobSelector FunctionSelector `json:"aggInitiateJobSelector" gorm:"not null"`
	// selector for fulfillment (oracle reporting) method on aggregator contract
	AggFulfillSelector FunctionSelector `json:"aggFulfillSelector" gorm:"not null"`
}

// ABI packs the encumberance as a byte array using the same technique as abi.encodePacked.
//
// Used only for constructing a stable hash which will be signed by all oracles,
// so it does not have to be easily parsed or unambiguous (e.g., re-ordering
// Oracles will result in output), just an injective function.
func (e Encumbrance) ABI() ([]byte, error) {
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
	endAt := e.EndAt.Time.Unix()
	if endAt > 0xffffffff { // Optimistically, this could be an issue in 2038...
		return nil, fmt.Errorf(
			"endat date %s is too late to fit in uint32",
			e.EndAt.Time)
	}
	endAtSerialised := make([]byte, 4)
	binary.BigEndian.PutUint32(endAtSerialised, uint32(endAt&math.MaxUint32))
	_, err = buffer.Write(endAtSerialised)
	if err != nil {
		return nil, errors.Wrap(err, "while writing endAt")
	}

	err = encodeOracles(&buffer, e.Oracles)
	if err != nil {
		return nil, errors.Wrap(err, "while writing oracles")
	}
	_, err = buffer.Write(address256Bits(e.Aggregator))
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

func encodeOracles(buffer *bytes.Buffer, oracles []EIP55Address) error {
	for _, o := range oracles {
		_, err := buffer.Write(address256Bits(o))
		if err != nil {
			return errors.Wrap(err, "while writing oracle address to buffer")
		}
	}
	return nil
}

// address256Bits Zero left-pads a to 32 bytes, as in Solidity's abi.encodePacking
func address256Bits(a EIP55Address) []byte {
	return a.Hash().Bytes()
}
