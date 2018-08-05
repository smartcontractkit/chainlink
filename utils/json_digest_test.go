package utils

import (
	"fmt"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNormalizedJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		object   interface{}
		output   string
		didError bool
	}{
		{"empty object", struct{}{}, "{}", false},
		{"empty array", []string{}, "[]", false},
		{"null", nil, "null", false},
		{"float", 1510599740287532257480015872.0, "1.510600e+27", false},
		{"bool", true, "true", false},
		{"string", "string", "\"string\"", false},
		{"array with one item", []string{"item"}, "[\"item\"]", false},
		{"map with one item", map[string]string{"item": "value"}, "{\"item\":\"value\"}", false},
		// See https://en.wikipedia.org/wiki/Precomposed_character
		{"string with decomposed characters",
			"\u0041\u030a\u0073\u0074\u0072\u006f\u0308\u006d",
			"\"\u00c5\u0073\u0074\u0072\u00f6\u006d\"",
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str, err := NormalizedJSON(test.object)
			if test.didError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.output, str)
		})
	}
}

func TestObjectDigest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		object   interface{}
		hash     string
		didError bool
	}{
		{"empty object", struct{}{}, "0xb48d38f93eaa084033fc5970bf96e559c33c4cdc07d889ab00b4d63f9590739d", false},
		{"empty array", []string{}, "0x518674ab2b227e5f11e9084f615d57663cde47bce1ba168b4c19c7ee22a73d70", false},
		{"null", nil, "0xefbde2c3aee204a69b7696d4b10ff31137fe78e3946306284f806e2dfc68b805", false},
		{"float", 0.02, "0xec649254e62c28ebf60f286722c5899dac36060f1979ad7600e93f8d6c2086cd", false},
		{"bool", true, "0x6273151f959616268004b58dbb21e5c851b7b8d04498b4aabee12291d22fc034", false},
		{"string", "string", "0x906a2df8e83ef1c8ca8d9bb50fe4f92f30d85286e0cccd65b363f5d3a1520c91", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			digest, err := ObjectDigest(test.object)
			if test.didError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.hash, common.ToHex(digest))
		})
	}
}

func TestObjectDigest_ProducesDeterministicResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		objects []interface{}
	}{
		{"float representations",
			[]interface{}{
				[]float64{1.5105997402875323e+27},
				[]float64{math.Pow(1147404288.0, 3.0)},
				[]float64{1510599740287532257480015872.0},
			},
		},
		// FIXME: this is not a very good test because we have no guarantee that
		// the keys will iterate over in the order they appear below
		{"object ordering",
			[]interface{}{
				map[string]interface{}{
					"a": nil, "b": nil, "c": nil,
				},
				map[string]interface{}{
					"c": nil, "b": nil, "a": nil,
				},
				map[string]interface{}{
					"b": nil, "c": nil, "a": nil,
				},
			},
		},
		// See https://en.wikipedia.org/wiki/Precomposed_character
		{"utf-8 precomposed vs decomposed",
			[]interface{}{
				[]string{"\u00c5\u0073\u0074\u0072\u00f6\u006d"},
				[]string{"\u0041\u030a\u0073\u0074\u0072\u006f\u0308\u006d"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var firstHash string
			for _, object := range test.objects {
				digest, err := ObjectDigest(object)
				assert.NoError(t, err)

				hash := common.ToHex(digest)
				if firstHash == "" {
					firstHash = hash
				}
				assert.Equal(t, firstHash, hash, fmt.Sprintf("When creating digest for %+v", object))
			}
			assert.NotEqual(t, firstHash, "")
		})
	}
}

func TestObjectDigest_DifferentTypesProduceDifferentHashes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		firstObject  interface{}
		secondObject interface{}
	}{
		{"map vs array", []string{"a", "b", "c", "d"}, map[string]string{"a": "b", "c": "d"}},
		{"float vs string", 2.400000e+00, "2.400000e+00"},
		{"null vs string", nil, "null"},
		{"delimiters in array", []string{"a", "b"}, []string{"ab"}},
		{"delimiters in map", map[string]string{"a": "b", "c": "d"}, map[string]string{"ab": "cd"}},
		{"separators in map", map[string]string{"abc": "d"}, map[string]string{"a": "bcd"}},
		{"string escaping", []string{"a", "b"}, []string{"a\",\"b"}},
		{"bool vs string", true, "true"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			firstDigest, err := ObjectDigest(test.firstObject)
			assert.NoError(t, err)
			firstHash := common.ToHex(firstDigest)

			secondDigest, err := ObjectDigest(test.secondObject)
			assert.NoError(t, err)
			secondHash := common.ToHex(secondDigest)

			assert.NotEqual(t, firstHash, secondHash, fmt.Sprintf("%+v should not produce the same hash as %+v", test.firstObject, test.secondObject))
		})
	}
}
