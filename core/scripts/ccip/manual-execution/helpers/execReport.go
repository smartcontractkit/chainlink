package helpers

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

// Hash contains all supported hash formats.
// Add additional hash types e.g. [20]byte as needed here.
type Hash interface {
	[32]byte
}

type Ctx[H Hash] interface {
	Hash(l []byte) H
	HashInternal(a, b H) H
	ZeroHash() H
}

type keccakCtx struct {
	InternalDomainSeparator [32]byte
}

func NewKeccakCtx() Ctx[[32]byte] {
	return keccakCtx{
		InternalDomainSeparator: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}

// Hash hashes a byte array with Keccak256
func (k keccakCtx) Hash(l []byte) [32]byte {
	// Note this Keccak256 cannot error https://github.com/golang/crypto/blob/master/sha3/sha3.go#L126
	// if we start supporting hashing algos which do, we can change this API to include an error.
	return Keccak256Fixed(l)
}

// HashInternal orders two [32]byte values and prepends them with
// a separator before hashing them.
func (k keccakCtx) HashInternal(a, b [32]byte) [32]byte {
	if bytes.Compare(a[:], b[:]) < 0 {
		return k.Hash(append(k.InternalDomainSeparator[:], append(a[:], b[:]...)...))
	}
	return k.Hash(append(k.InternalDomainSeparator[:], append(b[:], a[:]...)...))
}

// ZeroHash returns the zero hash: 0xFF..FF
// We use bytes32 0xFF..FF for zeroHash in the CCIP research spec, this needs to match.
// This value is chosen since it is unlikely to be the result of a hash, and cannot match any internal node preimage.
func (k keccakCtx) ZeroHash() [32]byte {
	var zeroes [32]byte
	for i := 0; i < 32; i++ {
		zeroes[i] = 0xFF
	}
	return zeroes
}

var (
	LeafDomainSeparator = [1]byte{0x00}
)

type LeafHasher struct {
	geABI        abi.ABI
	metaDataHash [32]byte
	ctx          Ctx[[32]byte]
}

func (t *LeafHasher) HashLeaf(log types.Log) ([32]byte, error) {
	event, err := t.ParseEVM2EVMLog(log)
	if err != nil {
		return [32]byte{}, err
	}
	encodedTokens, err := ABIEncode(`[{"components": [{"name": "token","type": "address"}, {"name": "amount", "type": "uint256"}],"type": "tuple[]"}]`, event.Message.TokenAmounts)
	if err != nil {
		return [32]byte{}, err
	}

	bytesArray, err := abi.NewType("bytes[]", "bytes[]", nil)
	if err != nil {
		return [32]byte{}, err
	}

	encodedSourceTokenData, err := abi.Arguments{abi.Argument{Type: bytesArray}}.PackValues([]interface{}{event.Message.SourceTokenData})
	if err != nil {
		return [32]byte{}, err
	}

	packedFixedSizeValues, err := ABIEncode(
		`[
{"name": "sender", "type":"address"},
{"name": "receiver", "type":"address"},
{"name": "sequenceNumber", "type":"uint64"},
{"name": "gasLimit", "type":"uint256"},
{"name": "strict", "type":"bool"},
{"name": "nonce", "type":"uint64"},
{"name": "feeToken","type": "address"},
{"name": "feeTokenAmount","type": "uint256"}
]`,
		event.Message.Sender,
		event.Message.Receiver,
		event.Message.SequenceNumber,
		event.Message.GasLimit,
		event.Message.Strict,
		event.Message.Nonce,
		event.Message.FeeToken,
		event.Message.FeeTokenAmount,
	)
	if err != nil {
		return [32]byte{}, err
	}
	fixedSizeValuesHash := t.ctx.Hash(packedFixedSizeValues)

	packedValues, err := ABIEncode(
		`[
{"name": "leafDomainSeparator","type":"bytes1"},
{"name": "metadataHash", "type":"bytes32"},
{"name": "fixedSizeValuesHash", "type":"bytes32"},
{"name": "dataHash", "type":"bytes32"},
{"name": "tokenAmountsHash", "type":"bytes32"},
{"name": "sourceTokenDataHash", "type":"bytes32"}
]`,
		LeafDomainSeparator,
		t.metaDataHash,
		fixedSizeValuesHash,
		t.ctx.Hash(event.Message.Data),
		t.ctx.Hash(encodedTokens),
		t.ctx.Hash(encodedSourceTokenData),
	)
	if err != nil {
		return [32]byte{}, err
	}
	return t.ctx.Hash(packedValues), nil
}

func (t *LeafHasher) ParseEVM2EVMLog(log types.Log) (*SendRequestedEvent, error) {
	event := new(SendRequestedEvent)
	err := bind.NewBoundContract(common.Address{}, t.geABI, nil, nil, nil).UnpackLog(event, "CCIPSendRequested", log)
	return event, err
}

func NewLeafHasher(sourceChainID uint64, destChainId uint64, onRampId common.Address, ctx Ctx[[32]byte]) *LeafHasher {
	geABI, _ := abi.JSON(strings.NewReader(OnRampABI))
	return &LeafHasher{
		geABI:        geABI,
		metaDataHash: getMetaDataHash(ctx, ctx.Hash([]byte("EVM2EVMMessageHashV2")), sourceChainID, onRampId, destChainId),
		ctx:          ctx,
	}
}

func Keccak256Fixed(in []byte) [32]byte {
	hash := sha3.NewLegacyKeccak256()
	// Note this Keccak256 cannot error https://github.com/golang/crypto/blob/master/sha3/sha3.go#L126
	// if we start supporting hashing algos which do, we can change this API to include an error.
	hash.Write(in)
	var h [32]byte
	copy(h[:], hash.Sum(nil))
	return h
}

func getMetaDataHash[H Hash](ctx Ctx[H], prefix [32]byte, sourceChainId uint64, onRampId common.Address, destChainId uint64) H {
	paddedOnRamp := onRampId.Hash()
	return ctx.Hash(ConcatBytes(prefix[:],
		math.U256Bytes(big.NewInt(0).SetUint64(sourceChainId)),
		math.U256Bytes(big.NewInt(0).SetUint64(destChainId)), paddedOnRamp[:]))
}

// ConcatBytes appends a bunch of byte arrays into a single byte array
func ConcatBytes(bufs ...[]byte) []byte {
	return bytes.Join(bufs, []byte{})
}

// ABIEncode is the equivalent of abi.encode.
// See a full set of examples https://github.com/ethereum/go-ethereum/blob/420b78659bef661a83c5c442121b13f13288c09f/accounts/abi/packing_test.go#L31
func ABIEncode(abiStr string, values ...interface{}) ([]byte, error) {
	// Create a dummy method with arguments
	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "inputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}
	res, err := inAbi.Pack("method", values...)
	if err != nil {
		return nil, err
	}
	return res[4:], nil
}

const (
	SourceFromHashes = true
	SourceFromProof  = false
)

type Proof[H Hash] struct {
	Hashes      []H    `json:"hashes"`
	SourceFlags []bool `json:"source_flags"`
}

type singleLayerProof[H Hash] struct {
	nextIndices []int
	subProof    []H
	sourceFlags []bool
}

type Tree[H Hash] struct {
	layers [][]H
}

func NewTree[H Hash](ctx Ctx[H], leafHashes []H) (*Tree[H], error) {
	if len(leafHashes) == 0 {
		return nil, errors.New("Cannot construct a tree without leaves")
	}
	var layer = make([]H, len(leafHashes))
	copy(layer, leafHashes)
	var layers = [][]H{layer}
	var curr int
	for len(layer) > 1 {
		paddedLayer, nextLayer := computeNextLayer(ctx, layer)
		layers[curr] = paddedLayer
		curr++
		layers = append(layers, nextLayer)
		layer = nextLayer
	}
	return &Tree[H]{
		layers: layers,
	}, nil
}

// Revive appears confused with the generics "receiver name t should be consistent with previous receiver name p for invalid-type"
//
//revive:disable:receiver-naming
func (t *Tree[H]) String() string {
	b := strings.Builder{}
	for _, layer := range t.layers {
		b.WriteString(fmt.Sprintf("%v", layer))
	}
	return b.String()
}

func (t *Tree[H]) Root() H {
	return t.layers[len(t.layers)-1][0]
}

func (t *Tree[H]) Prove(indices []int) Proof[H] {
	var proof Proof[H]
	for _, layer := range t.layers[:len(t.layers)-1] {
		res := proveSingleLayer(layer, indices)
		indices = res.nextIndices
		proof.Hashes = append(proof.Hashes, res.subProof...)
		proof.SourceFlags = append(proof.SourceFlags, res.sourceFlags...)
	}
	return proof
}

func computeNextLayer[H Hash](ctx Ctx[H], layer []H) ([]H, []H) {
	if len(layer) == 1 {
		return layer, layer
	}
	if len(layer)%2 != 0 {
		layer = append(layer, ctx.ZeroHash())
	}
	var nextLayer []H
	for i := 0; i < len(layer); i += 2 {
		nextLayer = append(nextLayer, ctx.HashInternal(layer[i], layer[i+1]))
	}
	return layer, nextLayer
}

func parentIndex(idx int) int {
	return idx / 2
}

func siblingIndex(idx int) int {
	return idx ^ 1
}

func proveSingleLayer[H Hash](layer []H, indices []int) singleLayerProof[H] {
	var (
		authIndices []int
		nextIndices []int
		sourceFlags []bool
	)
	j := 0
	for j < len(indices) {
		x := indices[j]
		nextIndices = append(nextIndices, parentIndex(x))
		if j+1 < len(indices) && indices[j+1] == siblingIndex(x) {
			j++
			sourceFlags = append(sourceFlags, SourceFromHashes)
		} else {
			authIndices = append(authIndices, siblingIndex(x))
			sourceFlags = append(sourceFlags, SourceFromProof)
		}
		j++
	}
	subProof := make([]H, 0, len(authIndices))
	for _, i := range authIndices {
		subProof = append(subProof, layer[i])
	}
	return singleLayerProof[H]{
		nextIndices: nextIndices,
		subProof:    subProof,
		sourceFlags: sourceFlags,
	}
}

// ProofFlagsToBits transforms a list of boolean proof flags to a *big.Int
// encoded number.
func ProofFlagsToBits(proofFlags []bool) *big.Int {
	encodedFlags := big.NewInt(0)
	for i := 0; i < len(proofFlags); i++ {
		if proofFlags[i] {
			encodedFlags.SetBit(encodedFlags, i, 1)
		}
	}
	return encodedFlags
}
