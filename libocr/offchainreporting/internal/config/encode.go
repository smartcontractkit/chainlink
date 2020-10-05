package config

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

const EncodedConfigVersion = 1

type setConfigEncodedComponents struct {
	DeltaProgress           time.Duration
	DeltaResend             time.Duration
	DeltaRound              time.Duration
	DeltaGrace              time.Duration
	DeltaC                  time.Duration
	Alpha                   float64
	DeltaStage              time.Duration
	RMax                    uint8
	S                       []int
	OffchainPublicKeys      []types.OffchainPublicKey
	PeerIDs                 []string
	SharedSecretEncryptions SharedSecretEncryptions
}

type setConfigSerializationTypes struct {
	DeltaProgress           int64
	DeltaResend             int64
	DeltaRound              int64
	DeltaGrace              int64
	DeltaC                  int64
	Alpha                   uint64 	DeltaStage              int64
	RMax                    uint8
	S                       []uint8
	OffchainPublicKeys      []common.Hash 	PeerIDs                 string        	SharedSecretEncryptions sseSerializationTypes
}

type sseSerializationTypes struct {
	DiffieHellmanPoint common.Hash
	SharedSecretHash   common.Hash
	Encryptions        [][SharedSecretSize]byte
}

var encoding = getEncoding()

var configSizeBound = 20 * 1000

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
	var setConfig setConfigSerializationTypes
	if err := encoding.Unpack(&setConfig, b); err != nil {
		return o, errors.Wrapf(err, "could not deserialize setConfig binary blob")
	}
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
	if o.RMax < 0 {
		panic(fmt.Sprintf("rMax must be non-negative, got %d", o.RMax))
	}
	return setConfigSerializationTypes{
		int64(o.DeltaProgress),
		int64(o.DeltaResend),
		int64(o.DeltaRound),
		int64(o.DeltaGrace),
		int64(o.DeltaC),
		math.Float64bits(o.Alpha),
		int64(o.DeltaStage),
		o.RMax,
		transmitDelays,
		publicKeys,
		strings.Join(o.PeerIDs, ","),
		o.SharedSecretEncryptions.serializationRepresentation(),
	}
}

func (or setConfigSerializationTypes) golangRepresentation() setConfigEncodedComponents {
	alpha := math.Float64frombits(or.Alpha)
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
		alpha,
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

func init() { 	checkTupEntriesMatchStruct(encoding[0].Type, setConfigEncodedComponents{})
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

func init() { 	checkFieldNamesMatch(setConfigEncodedComponents{}, setConfigSerializationTypes{})
	checkFieldNamesMatch(SharedSecretEncryptions{}, sseSerializationTypes{})
}
