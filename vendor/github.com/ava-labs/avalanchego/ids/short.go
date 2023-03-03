// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/hashing"
)

// ShortEmpty is a useful all zero value
var ShortEmpty = ShortID{}

// ShortID wraps a 20 byte hash as an identifier
type ShortID [20]byte

// ToShortID attempt to convert a byte slice into an id
func ToShortID(bytes []byte) (ShortID, error) {
	return hashing.ToHash160(bytes)
}

// ShortFromString is the inverse of ShortID.String()
func ShortFromString(idStr string) (ShortID, error) {
	bytes, err := cb58.Decode(idStr)
	if err != nil {
		return ShortID{}, err
	}
	return ToShortID(bytes)
}

// ShortFromPrefixedString returns a ShortID assuming the cb58 format is
// prefixed
func ShortFromPrefixedString(idStr, prefix string) (ShortID, error) {
	if !strings.HasPrefix(idStr, prefix) {
		return ShortID{}, fmt.Errorf("ID: %s is missing the prefix: %s", idStr, prefix)
	}
	return ShortFromString(strings.TrimPrefix(idStr, prefix))
}

func (id ShortID) MarshalJSON() ([]byte, error) {
	str, err := cb58.Encode(id[:])
	if err != nil {
		return nil, err
	}
	return []byte("\"" + str + "\""), nil
}

func (id *ShortID) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == nullStr { // If "null", do nothing
		return nil
	} else if len(str) < 2 {
		return errMissingQuotes
	}

	lastIndex := len(str) - 1
	if str[0] != '"' || str[lastIndex] != '"' {
		return errMissingQuotes
	}

	// Parse CB58 formatted string to bytes
	bytes, err := cb58.Decode(str[1:lastIndex])
	if err != nil {
		return fmt.Errorf("couldn't decode ID to bytes: %w", err)
	}
	*id, err = ToShortID(bytes)
	return err
}

func (id *ShortID) UnmarshalText(text []byte) error {
	return id.UnmarshalJSON(text)
}

// Bytes returns the 20 byte hash as a slice. It is assumed this slice is not
// modified.
func (id ShortID) Bytes() []byte { return id[:] }

// Hex returns a hex encoded string of this id.
func (id ShortID) Hex() string { return hex.EncodeToString(id.Bytes()) }

func (id ShortID) String() string {
	// We assume that the maximum size of a byte slice that
	// can be stringified is at least the length of an ID
	str, _ := cb58.Encode(id.Bytes())
	return str
}

// PrefixedString returns the String representation with a prefix added
func (id ShortID) PrefixedString(prefix string) string {
	return prefix + id.String()
}

func (id ShortID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

type sortShortIDData []ShortID

func (ids sortShortIDData) Less(i, j int) bool {
	return bytes.Compare(
		ids[i].Bytes(),
		ids[j].Bytes()) == -1
}
func (ids sortShortIDData) Len() int      { return len(ids) }
func (ids sortShortIDData) Swap(i, j int) { ids[j], ids[i] = ids[i], ids[j] }

// SortShortIDs sorts the ids lexicographically
func SortShortIDs(ids []ShortID) { sort.Sort(sortShortIDData(ids)) }

// IsSortedAndUniqueShortIDs returns true if the ids are sorted and unique
func IsSortedAndUniqueShortIDs(ids []ShortID) bool {
	for i := 0; i < len(ids)-1; i++ {
		if bytes.Compare(ids[i].Bytes(), ids[i+1].Bytes()) != -1 {
			return false
		}
	}
	return true
}

// IsUniqueShortIDs returns true iff [ids] are unique
func IsUniqueShortIDs(ids []ShortID) bool {
	set := ShortSet{}
	set.Add(ids...)
	return set.Len() == len(ids)
}

// ShortIDsToStrings converts an array of shortIDs to an array of their string
// representations
func ShortIDsToStrings(ids []ShortID) []string {
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = id.String()
	}
	return idStrs
}
