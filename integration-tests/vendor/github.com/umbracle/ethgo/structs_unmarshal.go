package ethgo

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/valyala/fastjson"
)

var defaultPool fastjson.ParserPool

// UnmarshalJSON implements the unmarshal interface
func (b *Block) UnmarshalJSON(buf []byte) error {
	p := defaultPool.Get()
	defer defaultPool.Put(p)

	v, err := p.Parse(string(buf))
	if err != nil {
		return err
	}

	if err := decodeHash(&b.Hash, v, "hash"); err != nil {
		return err
	}
	if err := decodeHash(&b.ParentHash, v, "parentHash"); err != nil {
		return err
	}
	if err := decodeHash(&b.Sha3Uncles, v, "sha3Uncles"); err != nil {
		return err
	}
	if err := decodeHash(&b.TransactionsRoot, v, "transactionsRoot"); err != nil {
		return err
	}
	if err := decodeHash(&b.StateRoot, v, "stateRoot"); err != nil {
		return err
	}
	if err := decodeHash(&b.ReceiptsRoot, v, "receiptsRoot"); err != nil {
		return err
	}
	if err := decodeAddr(&b.Miner, v, "miner"); err != nil {
		return err
	}
	if b.Number, err = decodeUint(v, "number"); err != nil {
		return err
	}
	if b.GasLimit, err = decodeUint(v, "gasLimit"); err != nil {
		return err
	}
	if b.GasUsed, err = decodeUint(v, "gasUsed"); err != nil {
		return err
	}
	if b.Timestamp, err = decodeUint(v, "timestamp"); err != nil {
		return err
	}
	if b.Difficulty, err = decodeBigInt(b.Difficulty, v, "difficulty"); err != nil {
		return err
	}
	if b.ExtraData, err = decodeBytes(b.ExtraData[:0], v, "extraData"); err != nil {
		return err
	}

	b.TransactionsHashes = b.TransactionsHashes[:0]
	b.Transactions = b.Transactions[:0]

	elems := v.GetArray("transactions")
	if len(elems) != 0 {
		if elems[0].Type() == fastjson.TypeString {
			// hashes (non full block)
			for _, elem := range elems {
				var h Hash
				if err := h.UnmarshalText(elem.GetStringBytes()); err != nil {
					return err
				}
				b.TransactionsHashes = append(b.TransactionsHashes, h)
			}
		} else {
			// structs (full block)
			for _, elem := range elems {
				txn := new(Transaction)
				if err := txn.unmarshalJSON(elem); err != nil {
					panic(err)
				}
				b.Transactions = append(b.Transactions, txn)
			}
		}
	}

	// uncles
	b.Uncles = b.Uncles[:0]
	for _, elem := range v.GetArray("uncles") {
		var h Hash
		if err := h.UnmarshalText(elem.GetStringBytes()); err != nil {
			return err
		}
		b.Uncles = append(b.Uncles, h)
	}

	return nil
}

// UnmarshalJSON implements the unmarshal interface
func (t *Transaction) UnmarshalJSON(buf []byte) error {
	p := defaultPool.Get()
	defer defaultPool.Put(p)

	v, err := p.Parse(string(buf))
	if err != nil {
		return err
	}
	return t.unmarshalJSON(v)
}

// isKeySet is a helper function for checking if a key has any value != nil,
// or if it's been set at all
func isKeySet(v *fastjson.Value, key string) bool {
	value := v.Get(key)
	return value != nil && value.Type() != fastjson.TypeNull
}

func (t *Transaction) unmarshalJSON(v *fastjson.Value) error {
	exists := func(names ...string) error {
		for _, name := range names {
			if !v.Exists(name) {
				return fmt.Errorf("'%s' not found", name)
			}
		}
		return nil
	}

	// detect transaction type
	var typ TransactionType
	if isKeySet(v, "chainId") {
		if isKeySet(v, "maxFeePerGas") {
			typ = TransactionDynamicFee
		} else {
			typ = TransactionAccessList
		}
	} else {
		typ = TransactionLegacy
	}
	t.Type = typ

	var err error
	if err := decodeHash(&t.Hash, v, "hash"); err != nil {
		return err
	}
	if err = decodeAddr(&t.From, v, "from"); err != nil {
		return err
	}
	if t.GasPrice, err = decodeUint(v, "gasPrice"); err != nil {
		return err
	}
	if t.Input, err = decodeBytes(t.Input[:0], v, "input"); err != nil {
		return err
	}
	if t.Value, err = decodeBigInt(t.Value, v, "value"); err != nil {
		return err
	}
	if t.Nonce, err = decodeUint(v, "nonce"); err != nil {
		return err
	}

	{
		// Do not decode 'to' if it doesn't exist.
		if err := exists("to"); err == nil {
			if v.Get("to").String() != "null" {
				var to Address
				if err = decodeAddr(&to, v, "to"); err != nil {
					return err
				}
				t.To = &to
			}
		}
	}

	if t.V, err = decodeBytes(t.V[:0], v, "v"); err != nil {
		return err
	}
	if t.R, err = decodeBytes(t.R[:0], v, "r"); err != nil {
		return err
	}
	if t.S, err = decodeBytes(t.S[:0], v, "s"); err != nil {
		return err
	}

	if typ != TransactionLegacy {
		if t.ChainID, err = decodeBigInt(t.ChainID, v, "chainId"); err != nil {
			return err
		}
		if isKeySet(v, "accessList") {
			if err := t.AccessList.unmarshalJSON(v.Get("accessList")); err != nil {
				return err
			}
		}
	}

	if t.Gas, err = decodeUint(v, "gas"); err != nil {
		return err
	}

	if typ == TransactionDynamicFee {
		if t.MaxPriorityFeePerGas, err = decodeBigInt(t.MaxPriorityFeePerGas, v, "maxPriorityFeePerGas"); err != nil {
			return err
		}
		if t.MaxFeePerGas, err = decodeBigInt(t.MaxFeePerGas, v, "maxFeePerGas"); err != nil {
			return err
		}
	}

	// Check if the block hash field is set
	// If it's not -> the transaction is a pending txn, so these fields should be omitted
	// If it is -> the transaction is a sealed txn, so these fields should be included
	if isKeySet(v, "blockHash") {
		// The transaction is not a pending transaction, read data

		// Grab the block hash
		if err = decodeHash(&t.BlockHash, v, "blockHash"); err != nil {
			return err
		}

		// Grab the block number
		if t.BlockNumber, err = decodeUint(v, "blockNumber"); err != nil {
			return err
		}

		// Grab the transaction index
		if t.TxnIndex, err = decodeUint(v, "transactionIndex"); err != nil {
			return err
		}
	}

	return nil
}

func (t *AccessList) unmarshalJSON(v *fastjson.Value) error {
	elems, err := v.Array()
	if err != nil {
		return err
	}
	for _, elem := range elems {
		entry := AccessEntry{}
		if err = decodeAddr(&entry.Address, elem, "address"); err != nil {
			return err
		}
		storage, err := elem.Get("storageKeys").Array()
		if err != nil {
			return err
		}

		entry.Storage = make([]Hash, len(storage))
		for indx, stg := range storage {
			b, err := stg.StringBytes()
			if err != nil {
				return err
			}
			if err := entry.Storage[indx].UnmarshalText(b); err != nil {
				return err
			}
		}
		*t = append(*t, entry)
	}
	return nil
}

// UnmarshalJSON implements the unmarshal interface
func (r *Receipt) UnmarshalJSON(buf []byte) error {
	p := defaultPool.Get()
	defer defaultPool.Put(p)

	v, err := p.Parse(string(buf))
	if err != nil {
		return nil
	}

	if err := decodeAddr(&r.From, v, "from"); err != nil {
		return err
	}
	if fieldNotFull(v, "contractAddress") {
		if err := decodeAddr(&r.ContractAddress, v, "contractAddress"); err != nil {
			return err
		}
	}
	if err := decodeHash(&r.TransactionHash, v, "transactionHash"); err != nil {
		return err
	}
	if err := decodeHash(&r.BlockHash, v, "blockHash"); err != nil {
		return err
	}
	if r.TransactionIndex, err = decodeUint(v, "transactionIndex"); err != nil {
		return err
	}
	if r.BlockNumber, err = decodeUint(v, "blockNumber"); err != nil {
		return err
	}
	if r.GasUsed, err = decodeUint(v, "gasUsed"); err != nil {
		return err
	}
	if r.CumulativeGasUsed, err = decodeUint(v, "cumulativeGasUsed"); err != nil {
		return err
	}
	if r.LogsBloom, err = decodeBytes(r.LogsBloom[:0], v, "logsBloom", 256); err != nil {
		return err
	}
	if r.Status, err = decodeUint(v, "status"); err != nil {
		return err
	}

	// logs
	r.Logs = r.Logs[:0]
	for _, elem := range v.GetArray("logs") {
		log := new(Log)
		if err := log.unmarshalJSON(elem); err != nil {
			return err
		}
		r.Logs = append(r.Logs, log)
	}

	return nil
}

func (lf *LogFilter) UnmarshalJSON(buf []byte) error {
	p := defaultPool.Get()
	defer defaultPool.Put(p)

	v, err := p.Parse(string(buf))
	if err != nil {
		return fmt.Errorf("unable to parse input, %w", err)
	}

	// Unmarshal the address field
	lf.Address = lf.Address[:0]

	appendAddress := func(addressRaw []byte) error {
		address := new(Address)
		if err := address.UnmarshalText(addressRaw); err != nil {
			return err
		}

		lf.Address = append(lf.Address, *address)

		return nil
	}

	for _, addressValue := range v.GetArray("address") {
		addressRaw, err := addressValue.StringBytes()
		if err != nil {
			return err
		}

		if err := appendAddress(addressRaw); err != nil {
			return err
		}
	}

	// The address field can also be a single value
	if addressRaw := v.GetStringBytes("address"); addressRaw != nil {
		if err := appendAddress(addressRaw); err != nil {
			return err
		}
	}

	// Unmarshal the block hash
	lf.BlockHash = nil
	if v.Exists("blockHash") {
		extractedHash := &Hash{}
		if err := decodeHash(extractedHash, v, "blockHash"); err != nil {
			return err
		}

		lf.BlockHash = extractedHash
	}

	// decodeBlockNum is a helper method for extracting a BlockNumber
	decodeBlockNum := func(key string) (*BlockNumber, error) {
		numRaw, err := decodeInt64(v, key)
		if err != nil {
			return nil, err
		}

		blockNum := BlockNumber(numRaw)

		return &blockNum, nil
	}

	// Unmarshal the from field
	lf.From = nil
	if v.Exists("fromBlock") {
		if lf.From, err = decodeBlockNum("fromBlock"); err != nil {
			return err
		}
	}

	// Unmarshal the to field
	lf.To = nil
	if v.Exists("toBlock") {
		if lf.To, err = decodeBlockNum("toBlock"); err != nil {
			return err
		}
	}

	// Unmarshal the topics
	lf.Topics = lf.Topics[:0]
	for _, topicsValue := range v.GetArray("topics") {
		// Check if the index is set
		if topicsValue == nil || topicsValue.String() == "null" {
			lf.Topics = append(lf.Topics, nil)

			continue
		}

		innerTopics, err := topicsValue.Array()
		if err != nil {
			return err
		}

		resTopics := make([]*Hash, 0)
		for _, innerTopic := range innerTopics {
			hashValRaw, err := innerTopic.StringBytes()
			if err != nil {
				return err
			}

			hashVal := &Hash{}
			if err := hashVal.UnmarshalText(hashValRaw); err != nil {
				return err
			}

			resTopics = append(resTopics, hashVal)
		}

		lf.Topics = append(lf.Topics, resTopics)
	}

	return nil
}

// UnmarshalJSON implements the unmarshal interface
func (r *Log) UnmarshalJSON(buf []byte) error {
	p := defaultPool.Get()
	defer defaultPool.Put(p)

	v, err := p.Parse(string(buf))
	if err != nil {
		return nil
	}
	return r.unmarshalJSON(v)
}

func (r *Log) unmarshalJSON(v *fastjson.Value) error {
	var err error
	if v.Exists("removed") {
		// it is empty in etherscan API endpoint
		if r.Removed, err = decodeBool(v, "removed"); err != nil {
			return err
		}
	}
	if r.LogIndex, err = decodeUint(v, "logIndex"); err != nil {
		return err
	}
	if r.BlockNumber, err = decodeUint(v, "blockNumber"); err != nil {
		return err
	}
	if r.TransactionIndex, err = decodeUint(v, "transactionIndex"); err != nil {
		return err
	}
	if err := decodeHash(&r.TransactionHash, v, "transactionHash"); err != nil {
		return err
	}
	if v.Exists("blockHash") {
		// it is empty in etherscan API endpoint
		if err := decodeHash(&r.BlockHash, v, "blockHash"); err != nil {
			return err
		}
	}
	if err := decodeAddr(&r.Address, v, "address"); err != nil {
		return err
	}
	if r.Data, err = decodeBytes(r.Data[:0], v, "data"); err != nil {
		return err
	}

	r.Topics = r.Topics[:0]
	for _, topic := range v.GetArray("topics") {
		var t Hash
		b, err := topic.StringBytes()
		if err != nil {
			return err
		}
		if err := t.UnmarshalText(b); err != nil {
			return err
		}
		r.Topics = append(r.Topics, t)
	}
	return nil
}

func fieldNotFull(v *fastjson.Value, key string) bool {
	vv := v.Get(key)
	if vv == nil {
		return false
	}
	if vv.String() == "null" {
		return false
	}
	return true
}

func decodeBigInt(b *big.Int, v *fastjson.Value, key string) (*big.Int, error) {
	vv := v.Get(key)
	if vv == nil {
		return nil, fmt.Errorf("field '%s' not found", key)
	}
	str := vv.String()
	str = strings.Trim(str, "\"")

	if !strings.HasPrefix(str, "0x") {
		return nil, fmt.Errorf("field '%s' does not have 0x prefix: '%s'", key, str)
	}
	if b == nil {
		b = new(big.Int)
	}

	var ok bool
	b, ok = b.SetString(str[2:], 16)
	if !ok {
		return nil, fmt.Errorf("field '%s' failed to decode big int: '%s'", key, str)
	}
	return b, nil
}

func decodeBytes(dst []byte, v *fastjson.Value, key string, bits ...int) ([]byte, error) {
	vv := v.Get(key)
	if vv == nil {
		return nil, fmt.Errorf("field '%s' not found", key)
	}
	str := vv.String()
	str = strings.Trim(str, "\"")

	if !strings.HasPrefix(str, "0x") {
		return nil, fmt.Errorf("field '%s' does not have 0x prefix: '%s'", key, str)
	}
	str = str[2:]
	if len(str)%2 != 0 {
		str = "0" + str
	}
	buf, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	if len(bits) > 0 && bits[0] != len(buf) {
		return nil, fmt.Errorf("field '%s' invalid length, expected %d but found %d: %s", key, bits[0], len(buf), str)
	}
	dst = append(dst, buf...)
	return dst, nil
}

func decodeUint(v *fastjson.Value, key string) (uint64, error) {
	vv := v.Get(key)
	if vv == nil {
		return 0, fmt.Errorf("field '%s' not found", key)
	}
	str := vv.String()
	str = strings.Trim(str, "\"")

	if !strings.HasPrefix(str, "0x") {
		return 0, fmt.Errorf("field '%s' does not have 0x prefix: '%s'", key, str)
	}
	str = str[2:]
	if str == "" {
		str = "0"
	}

	num, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("field '%s' failed to decode uint: %s", key, str)
	}
	return num, nil
}

func decodeInt64(v *fastjson.Value, key string) (int64, error) {
	vv := v.Get(key)
	if vv == nil {
		return 0, fmt.Errorf("field '%s' not found", key)
	}
	str := vv.String()
	str = strings.Trim(str, "\"")

	if !strings.HasPrefix(str, "0x") {
		return 0, fmt.Errorf("field '%s' does not have 0x prefix: '%s'", key, str)
	}
	str = str[2:]
	if str == "" {
		str = "0"
	}

	num, err := strconv.ParseInt(str, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("field '%s' failed to decode int64: %s", key, str)
	}
	return num, nil
}

func decodeHash(h *Hash, v *fastjson.Value, key string) error {
	b := v.GetStringBytes(key)
	if len(b) == 0 {
		return fmt.Errorf("field '%s' not found", key)
	}

	// Make sure the memory location is initialized
	if h == nil {
		h = &Hash{}
	}

	h.UnmarshalText(b)
	return nil
}

func decodeAddr(a *Address, v *fastjson.Value, key string) error {
	b := v.GetStringBytes(key)
	if len(b) == 0 {
		return fmt.Errorf("field '%s' not found", key)
	}
	a.UnmarshalText(b)
	return nil
}

func decodeBool(v *fastjson.Value, key string) (bool, error) {
	vv := v.Get(key)
	if vv == nil {
		return false, fmt.Errorf("field '%s' not found", key)
	}
	str := vv.String()
	if str == "false" {
		return false, nil
	} else if str == "true" {
		return true, nil
	}
	return false, fmt.Errorf("field '%s' with content '%s' cannot be decoded as bool", key, str)
}

func unmarshalTextByte(dst, src []byte, size int) error {
	str := string(src)

	str = strings.Trim(str, "\"")
	if !strings.HasPrefix(str, "0x") {
		return fmt.Errorf("0x prefix not found")
	}
	str = str[2:]
	b, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	if len(b) != size {
		return fmt.Errorf("length %d is not correct, expected %d", len(b), size)
	}
	copy(dst, b)
	return nil
}
