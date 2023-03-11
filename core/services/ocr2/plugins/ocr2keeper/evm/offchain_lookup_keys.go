package evm

import (
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"
)

type OffchainAPIKeys struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	// header,param,hmac
	Type string `json:"type"`
	// for header type this will be put as the Header.Key
	// for params type this will be used as url substitution key ie {name}
	Name string `json:"name"`
	// encrypted(cbor([upkeepID,value]))
	// for header the inner value will be used as Header.Value
	// for param the inner value will be what is replaced in the url
	Value []byte `json:"value"`
	// this is the value after peeling away the encryption the cbor and checking the upkeepID
	DecryptVal string
}

type DecryptedValue struct {
	UpkeepID *big.Int
	Value    string
}

// getAPIKeys takes the offchain config and attempts to get any api keys that it may have
func getAPIKeys(upkeepID *big.Int, offchainConfig []byte) (OffchainAPIKeys, error) {
	offchainAPIKeys := OffchainAPIKeys{}
	if len(offchainConfig) == 0 {
		return OffchainAPIKeys{}, nil
	}
	err := cbor.Unmarshal(offchainConfig, &offchainAPIKeys)
	if err != nil {
		return OffchainAPIKeys{}, err
	}
	for i, key := range offchainAPIKeys.Keys {
		decodedKey, decryptErr := decodeValue("TODO_KEY", key.Value)
		if decryptErr != nil {
			continue
		}
		if decodedKey.UpkeepID.String() != upkeepID.String() {
			// should we error here...maybe
			continue // we won't use the key if it's not specifically for this upkeepID
		}
		// replace with decrypted value
		offchainAPIKeys.Keys[i].DecryptVal = decodedKey.Value
	}

	return offchainAPIKeys, nil
}

// decodeValue needs the node key to do the decryption of the input string.
// The underlying format of the string is expected to be cbor encoded and
// in the format [upkeepID []byte, apiKey string]
func decodeValue(nodeKey string, input []byte) (DecryptedValue, error) {
	// TODO decrypt value with key

	var plaintext [2]interface{}
	err := cbor.Unmarshal(input, &plaintext)
	if err != nil {
		return DecryptedValue{}, err
	}

	id := big.NewInt(0).SetBytes(plaintext[0].([]byte))
	if id == nil {
		return DecryptedValue{}, err
	}

	value := fmt.Sprint(plaintext[1])

	return DecryptedValue{id, value}, nil
}
