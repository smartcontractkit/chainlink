package common

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/umbracle/fastrlp"
)

type Nonce [8]byte

var (
	nonceT = reflect.TypeOf(Nonce{})
)

func (n Nonce) String() string {
	return hexutil.Encode(n[:])
}

// MarshalText implements encoding.TextMarshaler
func (n Nonce) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

// UnmarshalJSON parses a nonce in hex syntax.
func (n *Nonce) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(nonceT, input, n[:])
}

type ExtraData string

func (e ExtraData) Decode() ([]byte, error) {
	return hexutil.Decode(string(e))
}

// Header represents a block header in the Ethereum blockchain.
type PolygonEdgeHeader struct {
	ParentHash   common.Hash    `json:"parentHash"       gencodec:"required"`
	Sha3Uncles   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Miner        common.Address `json:"miner"`
	StateRoot    common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxRoot       common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptsRoot common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	LogsBloom    types.Bloom    `json:"logsBloom"        gencodec:"required"`
	Difficulty   hexutil.Uint64 `json:"difficulty"       gencodec:"required"`
	Number       hexutil.Uint64 `json:"number"           gencodec:"required"`
	GasLimit     hexutil.Uint64 `json:"gasLimit"         gencodec:"required"`
	GasUsed      hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
	Timestamp    hexutil.Uint64 `json:"timestamp"        gencodec:"required"`
	ExtraData    ExtraData      `json:"extraData"        gencodec:"required"`
	MixHash      common.Hash    `json:"mixHash"`
	Nonce        Nonce          `json:"nonce"`
	Hash         common.Hash    `json:"hash"`

	// baseFeePerGas is the response format from go-ethereum. Polygon-Edge
	// seems to have fixed this in this commit:
	// https://github.com/0xPolygon/polygon-edge/commit/e859acf7e7f0286ceeecce022b978c8fdb57d71b
	// But node operators dont seem to have updated their polygon-edge client
	// version and still send baseFee instead of baseFeePerGas.
	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee    hexutil.Uint64 `json:"baseFeePerGas"`
	BaseFeeAlt hexutil.Uint64 `json:"baseFee,omitempty"`
}

func GetPolygonEdgeRLPHeader(jsonRPCClient *rpc.Client, blockNum *big.Int) (rlpHeader []byte, hash string, err error) {
	var h PolygonEdgeHeader
	err = jsonRPCClient.Call(&h, "eth_getBlockByNumber", "0x"+blockNum.Text(16), true)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get poloygon-edge header: %+v", err)
	}

	ar := &fastrlp.Arena{}
	val, err := MarshalRLPWith(ar, &h)
	if err != nil {
		return nil, "", err
	}

	dst := make([]byte, 0)
	dst = val.MarshalTo(dst)

	return dst, h.Hash.String(), err
}

// MarshalRLPWith marshals the header to RLP with a specific fastrlp.Arena
// Adding polygon-edge as a dependency caused a lot of issues with conflicting
// dependency version with other libraries in this repo and some methods being
// referenced from older versions
// Reference: https://github.com/0xPolygon/polygon-edge/blob/develop/types/rlp_marshal.go#L73C50-L73C53
func MarshalRLPWith(arena *fastrlp.Arena, h *PolygonEdgeHeader) (*fastrlp.Value, error) {
	vv := arena.NewArray()

	vv.Set(arena.NewCopyBytes(h.ParentHash.Bytes()))
	vv.Set(arena.NewCopyBytes(h.Sha3Uncles.Bytes()))
	vv.Set(arena.NewCopyBytes(h.Miner[:]))
	vv.Set(arena.NewCopyBytes(h.StateRoot.Bytes()))
	vv.Set(arena.NewCopyBytes(h.TxRoot.Bytes()))
	vv.Set(arena.NewCopyBytes(h.ReceiptsRoot.Bytes()))
	vv.Set(arena.NewCopyBytes(h.LogsBloom[:]))

	vv.Set(arena.NewUint(uint64(h.Difficulty)))
	vv.Set(arena.NewUint(uint64(h.Number)))
	vv.Set(arena.NewUint(uint64(h.GasLimit)))
	vv.Set(arena.NewUint(uint64(h.GasUsed)))
	vv.Set(arena.NewUint(uint64(h.Timestamp)))

	extraDataBytes, err := h.ExtraData.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to hex decode polygon-edge ExtraData: %+v", err)
	}
	extraDataBytes, err = GetIbftExtraClean(extraDataBytes)
	if err != nil {
		return nil, fmt.Errorf("GetIbftExtraClean error : %+v", err)
	}
	vv.Set(arena.NewCopyBytes(extraDataBytes))
	vv.Set(arena.NewCopyBytes(h.MixHash.Bytes()))

	nonceHexString := h.Nonce.String()
	nonceBytes, err := hexutil.Decode(nonceHexString)
	if err != nil {
		return nil, fmt.Errorf("failed to hex decode polygon-edge ExtraData: %+v", err)
	}
	vv.Set(arena.NewCopyBytes(nonceBytes))

	baseFee := h.BaseFee
	if h.BaseFeeAlt > 0 {
		baseFee = h.BaseFeeAlt
	}
	vv.Set(arena.NewUint(uint64(baseFee)))

	return vv, nil
}

// Remove blockHeader.ExtraData.Committed without unpacking ExtraData into
// its full fledged type, which needs the full import of the package
// github.com/0xPolygon/polygon-edge. polygon-edge is a node implementation,
// and not a client. Adding polygon-edge as a dependency caused a lot of
// issues with conflicting dependency version with other libraries in this
// repo and some methods being referenced from older versions.
func GetIbftExtraClean(extra []byte) (cleanedExtra []byte, err error) {
	// Capture prefix 0's sent by nexon supernet
	hexExtra := hex.EncodeToString(extra)
	prefix := ""
	for _, s := range hexExtra {
		if s != '0' {
			break
		}
		prefix = prefix + "0"
	}

	hexExtra = strings.TrimLeft(hexExtra, "0")
	extra, err = hex.DecodeString(hexExtra)
	if err != nil {
		return nil, fmt.Errorf("invalid extra data in polygon-edge chain: %+v", err)
	}

	var extraData []interface{}
	err = rlp.DecodeBytes(extra, &extraData)
	if err != nil {
		return nil, err
	}

	// Remove Committed from blockHeader.ExtraData, because it holds signatures for
	// the current block which is not finalized until the next block. So this gets
	// ignored when calculating the hash
	// Reference: https://github.com/0xPolygon/polygon-edge/blob/develop/consensus/polybft/hash.go#L20-L27.
	if len(extraData) > 3 {
		extraData[2] = []interface{}{[]byte{}, []byte{}}
	}

	cleanedExtra, err = rlp.EncodeToBytes(extraData)
	if err != nil {
		return nil, err
	}

	// Add prefix 0's sent by nexon supernet before sending output
	hexExtra = prefix + hex.EncodeToString(cleanedExtra)
	cleanedExtra, err = hex.DecodeString(hexExtra)
	return cleanedExtra, err
}
