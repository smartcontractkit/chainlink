package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

const EncodedConfigVersion = 1

// setConfigEncodedComponents contains the contents of the oracle Config objects
// which need to be serialized
type setConfigEncodedComponents struct {
	DeltaProgress           time.Duration
	DeltaResend             time.Duration
	DeltaRound              time.Duration
	DeltaGrace              time.Duration
	DeltaC                  time.Duration
	AlphaPPB                uint64
	DeltaStage              time.Duration
	RMax                    uint8
	S                       []int
	OffchainPublicKeys      []types.OffchainPublicKey
	PeerIDs                 []string
	SharedSecretEncryptions SharedSecretEncryptions
}

// setConfigSerializationTypes gives the types used to represent a
// setConfigEncodedComponents to abiencode. The field names must match those of
// setConfigEncodedComponents.
type setConfigSerializationTypes struct {
	DeltaProgress           int64
	DeltaResend             int64
	DeltaRound              int64
	DeltaGrace              int64
	DeltaC                  int64
	AlphaPPB                uint64
	DeltaStage              int64
	RMax                    uint8
	S                       []uint8
	OffchainPublicKeys      []common.Hash // Each key is a bytes32
	PeerIDs                 string        // comma-separated
	SharedSecretEncryptions sseSerializationTypes
}

// sseSerializationTypes gives the types used to represent an
// SharedSecretEncryptions to abiencode. The field names must match those of
// SharedSecretEncryptions.
type sseSerializationTypes struct {
	DiffieHellmanPoint common.Hash
	SharedSecretHash   common.Hash
	Encryptions        [][SharedSecretSize]byte
}

// encoding is the ABI schema used to encode a setConfigEncodedComponents, taken
// from setConfigEncodedComponentsABI in ./abiencode.go (in this package directory.)
var encoding = getEncoding()

// Serialized configs must be no larger than this (arbitrary bound, to prevent
// resource exhaustion attacks)
var configSizeBound = 20 * 1000

// Encode returns a binary serialization of o
func (o setConfigEncodedComponents) encode() []byte {
	rv, err := encoding.Pack(o.serializationRepresentation())
	if err != nil {
		panic(err)
	}
	if len(rv) > configSizeBound {
		panic("config serialization too large")
	}
	return rv
}

func decodeContractSetConfigEncodedComponents(
	b []byte,
) (o setConfigEncodedComponents, err error) {
	if len(b) > configSizeBound {
		return o, errors.Errorf(
			"attempt to deserialize a too-long config (%d bytes)", len(b),
		)
	}
	var vals []interface{}
	if vals, err = encoding.Unpack(b); err != nil {
		return o, errors.Wrapf(err, "could not deserialize setConfig binary blob")
	}
	setConfig := abi.ConvertType(vals[0], &setConfigSerializationTypes{}).(*setConfigSerializationTypes)
	return setConfig.golangRepresentation(), nil
}

func (o setConfigEncodedComponents) serializationRepresentation() setConfigSerializationTypes {
	transmitDelays := make([]uint8, len(o.S))
	for i, d := range o.S {
		transmitDelays[i] = uint8(d)
	}
	publicKeys := make([]common.Hash, len(o.OffchainPublicKeys))
	for i, k := range o.OffchainPublicKeys {
		publicKeys[i] = common.BytesToHash(k)
	}
	return setConfigSerializationTypes{
		int64(o.DeltaProgress),
		int64(o.DeltaResend),
		int64(o.DeltaRound),
		int64(o.DeltaGrace),
		int64(o.DeltaC),
		o.AlphaPPB,
		int64(o.DeltaStage),
		o.RMax,
		transmitDelays,
		publicKeys,
		strings.Join(o.PeerIDs, ","),
		o.SharedSecretEncryptions.serializationRepresentation(),
	}
}

func (or setConfigSerializationTypes) golangRepresentation() setConfigEncodedComponents {
	transmitDelays := make([]int, len(or.S))
	for i, d := range or.S {
		transmitDelays[i] = int(d)
	}
	keys := make([]types.OffchainPublicKey, len(or.OffchainPublicKeys))
	for i, k := range or.OffchainPublicKeys {
		keys[i] = types.OffchainPublicKey(k.Bytes())
	}
	var peerIDs []string
	if len(or.PeerIDs) > 0 {
		peerIDs = strings.Split(or.PeerIDs, ",")
	}
	return setConfigEncodedComponents{
		time.Duration(or.DeltaProgress),
		time.Duration(or.DeltaResend),
		time.Duration(or.DeltaRound),
		time.Duration(or.DeltaGrace),
		time.Duration(or.DeltaC),
		or.AlphaPPB,
		time.Duration(or.DeltaStage),
		or.RMax,
		transmitDelays,
		keys,
		peerIDs,
		or.SharedSecretEncryptions.golangRepresentation(),
	}
}

func (e SharedSecretEncryptions) serializationRepresentation() sseSerializationTypes {
	encs := make([][SharedSecretSize]byte, len(e.Encryptions))
	for i, enc := range e.Encryptions {
		encs[i] = enc
	}
	return sseSerializationTypes{
		common.Hash(e.DiffieHellmanPoint),
		e.SharedSecretHash,
		encs,
	}
}

func (er sseSerializationTypes) golangRepresentation() SharedSecretEncryptions {
	encs := make([]encryptedSharedSecret, len(er.Encryptions))
	for i, enc := range er.Encryptions {
		encs[i] = encryptedSharedSecret(enc)
	}
	return SharedSecretEncryptions{
		[32]byte(er.DiffieHellmanPoint),
		er.SharedSecretHash,
		encs,
	}
}

func getEncoding() abi.Arguments {
	// Trick used in abi's TestPack, to parse a list of arguments: make a JSON
	// representation of a method which has the target list as the inputs, then
	// pull the parsed argument list out of that method.
	aBI, err := abi.JSON(strings.NewReader(fmt.Sprintf(
		`[{ "name" : "method", "type": "function", "inputs": %s}]`,
		setConfigEncodedComponentsABI)))
	if err != nil {
		panic(err)
	}
	return aBI.Methods["method"].Inputs
}

func checkFieldNamesAgainstStruct(fields map[string]bool, i interface{}) {
	s := reflect.ValueOf(i).Type()
	for i := 0; i < s.NumField(); i++ {
		fieldName := s.Field(i).Name
		if !fields[fieldName] {
			panic("no encoding found for " + fieldName)
		}
		fields[fieldName] = false
	}
	for name, unseen := range fields {
		if unseen {
			panic("extra field found in abiencode schema, " + name)
		}
	}
}

func checkTupEntriesMatchStruct(t abi.Type, i interface{}) {
	if t.T != abi.TupleTy {
		panic("tuple required")
	}
	fields := make(map[string]bool)
	for _, fieldName := range t.TupleRawNames {
		capitalizedName := strings.ToUpper(fieldName[:1]) + fieldName[1:]
		fields[capitalizedName] = true
	}
	checkFieldNamesAgainstStruct(fields, i)
}

func init() { // check that abiencode fields match those of config structs
	checkTupEntriesMatchStruct(encoding[0].Type, setConfigEncodedComponents{})
	components := encoding[0].Type.TupleElems
	essName := encoding[0].Type.TupleRawNames[len(components)-1]
	if essName != "sharedSecretEncryptions" {
		panic("expecting sharedSecretEncryptions in last position, got " + essName)
	}
	ess := components[len(components)-1]
	checkTupEntriesMatchStruct(*ess, SharedSecretEncryptions{})
}

func checkFieldNamesMatch(s, t interface{}) {
	st, tt := reflect.ValueOf(s).Type(), reflect.ValueOf(t).Type()
	if st.NumField() != tt.NumField() {
		panic(fmt.Sprintf("number of fields differ: %T has %d, %T has %d",
			s, st.NumField(),
			t, tt.NumField()))
	}
	for i := 0; i < st.NumField(); i++ {
		if st.Field(i).Name != tt.Field(i).Name {
			panic(fmt.Sprintf("field name mismatch on %T vs %T: %s vs %s",
				s, t, st.Field(i).Name, tt.Field(i).Name))
		}
	}
}

func init() { // Check that serialization fields match those of target structs
	checkFieldNamesMatch(setConfigEncodedComponents{}, setConfigSerializationTypes{})
	checkFieldNamesMatch(SharedSecretEncryptions{}, sseSerializationTypes{})
}
