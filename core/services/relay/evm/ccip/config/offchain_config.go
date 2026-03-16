package config

import (
	"encoding/json"
)

type OffchainConfig interface {
	Validate() error
}

func DecodeOffchainConfig[T OffchainConfig](encodedConfig []byte) (T, error) {
	var result T
	err := json.Unmarshal(encodedConfig, &result)
	if err != nil {
		return result, err
	}
	err = result.Validate()
	if err != nil {
		return result, err
	}
	return result, nil
}

func EncodeOffchainConfig[T OffchainConfig](occ T) ([]byte, error) {
	return json.Marshal(occ)
}
