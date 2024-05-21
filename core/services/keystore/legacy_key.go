package keystore

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type rawLegacyKey []string
type rawLegacyKeys map[string]rawLegacyKey

type LegacyKeyStorage struct {
	legacyRawKeys rawLegacyKeys
}

func (rlk *rawLegacyKeys) len() (n int) {
	for _, v := range *rlk {
		n += len(v)
	}
	return n
}

func (rlk *rawLegacyKeys) has(name string) bool {
	for n := range *rlk {
		if n == name {
			return true
		}
	}
	return false
}

func (rlk *rawLegacyKeys) hasValueInField(fieldName, value string) bool {
	for _, v := range (*rlk)[fieldName] {
		if v == value {
			return true
		}
	}
	return false
}

// StoreUnsupported will store the raw keys that no longer have support in the node
// it will check if raw json contains keys that have not been added to the key ring
// and stores them internally
func (k *LegacyKeyStorage) StoreUnsupported(allRawKeysJson []byte, keyRing *keyRing) error {
	if keyRing == nil {
		return errors.New("keyring is nil")
	}
	supportedKeyRingJson, err := json.Marshal(keyRing.raw())
	if err != nil {
		return err
	}

	var (
		allKeys       = rawLegacyKeys{}
		supportedKeys = rawLegacyKeys{}
	)

	err = json.Unmarshal(allRawKeysJson, &allKeys)
	if err != nil {
		return err
	}
	err = json.Unmarshal(supportedKeyRingJson, &supportedKeys)
	if err != nil {
		return err
	}

	k.legacyRawKeys = rawLegacyKeys{}
	for fName, fValue := range allKeys {
		if !supportedKeys.has(fName) {
			k.legacyRawKeys[fName] = fValue
			continue
		}
		for _, v := range allKeys[fName] {
			if !supportedKeys.hasValueInField(fName, v) {
				k.legacyRawKeys[fName] = append(k.legacyRawKeys[fName], v)
			}
		}
	}

	return nil
}

// UnloadUnsupported will inject the unsupported keys into the raw key ring json
func (k *LegacyKeyStorage) UnloadUnsupported(supportedRawKeyRingJson []byte) ([]byte, error) {
	supportedKeys := rawLegacyKeys{}
	err := json.Unmarshal(supportedRawKeyRingJson, &supportedKeys)
	if err != nil {
		return nil, err
	}

	for fName, vals := range k.legacyRawKeys {
		if !supportedKeys.has(fName) {
			supportedKeys[fName] = vals
			continue
		}
		for _, v := range vals {
			if !supportedKeys.hasValueInField(fName, v) {
				supportedKeys[fName] = append(supportedKeys[fName], v)
			}
		}
	}

	allKeysJson, err := json.Marshal(supportedKeys)
	if err != nil {
		return nil, err
	}
	return allKeysJson, nil
}
