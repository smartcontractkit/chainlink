package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// ObjectDigest produces a solidity compatible Keccak256 signature from a
// consistent JSON encoding of any go object.
func ObjectDigest(object interface{}) ([]byte, error) {
	str, err := NormalizedJSON(object)
	if err != nil {
		return nil, err
	}
	return Keccak256([]byte(str))
}

// NormalizedJSON returns a JSON representation of an object that has been
// normalized to produce a consistent output for hashing.
//
// NOTE: If this string is unmarshalled again, there is no guarantee that the
// final representation will be consistent with the string produced by this
// function due to differences in JSON implementations and information loss.
// e.g:
// 	JSON does not have a requirement to respect object key ordering.
func NormalizedJSON(object interface{}) (string, error) {
	// First marshal into a JSON string
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	return NormalizedJSONString(jsonBytes)
}

// NormalizedJSONString takes stirng of JSON representation and produces a
// normalized version for consistent output for hashing.
func NormalizedJSONString(val []byte) (string, error) {
	// Unmarshal into a generic interface{}
	var data interface{}
	if err := json.Unmarshal(val, &data); err != nil {
		return "", err
	}

	buffer := &strings.Builder{}
	writer := bufio.NewWriter(buffer)

	// Wrap the buffer in a normalization writer
	wc := norm.NFC.Writer(writer)
	defer wc.Close()

	// Now marshal the generic interface
	if err := marshal(wc, data); err != nil {
		return "", err
	}
	wc.Close()
	writer.Flush()
	return buffer.String(), nil
}

// recursively write elements of the JSON to the hash, making sure to sort
// objects and to represent floats in exponent form
func marshal(writer io.Writer, data interface{}) error {
	switch element := data.(type) {
	case map[string]interface{}:
		return marshalObject(writer, element)
	case []interface{}:
		return marshalArray(writer, element)
	case float64:
		return marshalFloat(writer, element)
	case string:
		return marshalPrimitive(writer, element)
	case bool:
		return marshalPrimitive(writer, element)
	case nil:
		return marshalPrimitive(writer, element)
	default:
		panic(fmt.Sprintf("type '%T' in JSON input not handled", data))
	}
}

func marshalObject(writer io.Writer, data map[string]interface{}) error {
	_, err := fmt.Fprintf(writer, "{")
	if err != nil {
		return err
	}

	err = marshalMapOrderedKeys(writer, orderedKeys(data), data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(writer, "}")
	return err
}

func orderedKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func marshalMapOrderedKeys(writer io.Writer, orderedKeys []string, data map[string]interface{}) error {
	for index, key := range orderedKeys {
		err := marshal(writer, key)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(writer, ":")
		if err != nil {
			return err
		}

		value := data[key]
		err = marshal(writer, value)
		if err != nil {
			return err
		}

		if index == len(orderedKeys)-1 {
			break
		}

		_, err = fmt.Fprintf(writer, ",")
		if err != nil {
			return err
		}
	}
	return nil
}

func marshalArray(writer io.Writer, data []interface{}) error {
	_, err := fmt.Fprintf(writer, "[")
	if err != nil {
		return err
	}

	for index, item := range data {
		err := marshal(writer, item)
		if err != nil {
			return err
		}

		if index == len(data)-1 {
			break
		}

		_, err = fmt.Fprintf(writer, ",")
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(writer, "]")
	return err
}

func marshalPrimitive(writer io.Writer, data interface{}) error {
	output, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = writer.Write(output)
	return err
}

func marshalFloat(writer io.Writer, data float64) error {
	_, err := fmt.Fprintf(writer, "%e", data)
	return err
}
