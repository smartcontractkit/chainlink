package abi

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/umbracle/ethgo"
)

// ParseLog parses an event log
func ParseLog(args *Type, log *ethgo.Log) (map[string]interface{}, error) {
	var indexed, nonIndexed []*TupleElem

	for _, arg := range args.TupleElems() {
		if arg.Indexed {
			indexed = append(indexed, arg)
		} else {
			nonIndexed = append(nonIndexed, arg)
		}
	}

	// decode indexed fields
	indexedObjs, err := ParseTopics(&Type{kind: KindTuple, tuple: indexed}, log.Topics[1:])
	if err != nil {
		return nil, err
	}

	var nonIndexedObjs map[string]interface{}
	if len(nonIndexed) > 0 {
		nonIndexedRaw, err := Decode(&Type{kind: KindTuple, tuple: nonIndexed}, log.Data)
		if err != nil {
			return nil, err
		}
		raw, ok := nonIndexedRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("bad decoding")
		}
		nonIndexedObjs = raw
	}

	res := map[string]interface{}{}
	for _, arg := range args.TupleElems() {
		if arg.Indexed {
			res[arg.Name] = indexedObjs[0]
			indexedObjs = indexedObjs[1:]
		} else {
			res[arg.Name] = nonIndexedObjs[arg.Name]
		}
	}

	return res, nil
}

// ParseTopics parses topics from a log event
func ParseTopics(args *Type, topics []ethgo.Hash) ([]interface{}, error) {
	if args.kind != KindTuple {
		return nil, fmt.Errorf("expected a tuple type")
	}
	if len(args.TupleElems()) != len(topics) {
		return nil, fmt.Errorf("bad length")
	}

	elems := []interface{}{}
	for indx, arg := range args.TupleElems() {
		elem, err := ParseTopic(arg.Elem, topics[indx])
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}

	return elems, nil
}

// ParseTopic parses an individual topic
func ParseTopic(t *Type, topic ethgo.Hash) (interface{}, error) {
	switch t.kind {
	case KindBool:
		if bytes.Equal(topic[:], topicTrue[:]) {
			return true, nil
		} else if bytes.Equal(topic[:], topicFalse[:]) {
			return false, nil
		}
		return true, fmt.Errorf("is not a boolean")

	case KindInt, KindUInt:
		return readInteger(t, topic[:]), nil

	case KindAddress:
		return readAddr(topic[:])

	case KindFixedBytes:
		return readFixedBytes(t, topic[:])

	default:
		return nil, fmt.Errorf("topic parsing for type %s not supported", t.String())
	}
}

// EncodeTopic encodes a topic
func EncodeTopic(t *Type, val interface{}) (ethgo.Hash, error) {
	return encodeTopic(t, reflect.ValueOf(val))
}

func encodeTopic(t *Type, val reflect.Value) (ethgo.Hash, error) {
	switch t.kind {
	case KindBool:
		return encodeTopicBool(val)

	case KindUInt, KindInt:
		return encodeTopicNum(t, val)

	case KindAddress:
		return encodeTopicAddress(val)

	}
	return ethgo.Hash{}, fmt.Errorf("not found")
}

var topicTrue, topicFalse ethgo.Hash

func init() {
	topicTrue[31] = 1
}

func encodeTopicAddress(val reflect.Value) (res ethgo.Hash, err error) {
	var b []byte
	b, err = encodeAddress(val)
	if err != nil {
		return
	}
	copy(res[:], b[:])
	return
}

func encodeTopicNum(t *Type, val reflect.Value) (res ethgo.Hash, err error) {
	var b []byte
	b, err = encodeNum(val)
	if err != nil {
		return
	}
	copy(res[:], b[:])
	return
}

func encodeTopicBool(v reflect.Value) (res ethgo.Hash, err error) {
	if v.Kind() != reflect.Bool {
		return ethgo.Hash{}, encodeErr(v, "bool")
	}
	if v.Bool() {
		return topicTrue, nil
	}
	return topicFalse, nil
}
