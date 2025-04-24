package secrets

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func GenerateSecrets(input *cretypes.GenerateSecretsInput) (cretypes.NodeIndexToSecretsOverride, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	overrides := make(cretypes.NodeIndexToSecretsOverride)

	for i := range input.DonMetadata.NodesMetadata {
		nodeSecret := types.NodeSecret{}
		if input.EVMKeys != nil {
			nodeSecret.EthKeys = types.NodeEthKeyWrapper{}
			for chainID, evmKeys := range input.EVMKeys {
				nodeSecret.EthKeys.EthKeys = append(nodeSecret.EthKeys.EthKeys, types.NodeEthKey{
					JSON:     string(evmKeys.EncryptedJSONs[i]),
					Password: evmKeys.Password,
					ChainID:  chainID,
				})
			}
		}

		if input.P2PKeys != nil {
			nodeSecret.P2PKey = types.NodeP2PKey{
				JSON:     string(input.P2PKeys.EncryptedJSONs[i]),
				Password: input.P2PKeys.Password,
			}
		}

		nodeSecretString, err := toml.Marshal(nodeSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal node secrets")
		}

		overrides[i] = string(nodeSecretString)
	}

	return overrides, nil
}

func AddKeysToTopology(topology *cretypes.Topology, keys *cretypes.GenerateKeysOutput) (*cretypes.Topology, error) {
	if topology == nil {
		return nil, errors.New("topology is nil")
	}

	if keys == nil {
		return nil, errors.New("keys is nil")
	}

	if len(keys.P2PKeys) != len(topology.DonsMetadata) {
		return nil, fmt.Errorf("number of P2P keys does not match the number of DONs. Expected %d, got %d", len(topology.DonsMetadata), len(keys.P2PKeys))
	}

	if len(keys.EVMKeys) != len(topology.DonsMetadata) {
		return nil, fmt.Errorf("number of EVM keys does not match the number of DONs. Expected %d, got %d", len(topology.DonsMetadata), len(keys.EVMKeys))
	}

	for _, donMetadata := range topology.DonsMetadata {
		if _, ok := keys.P2PKeys[donMetadata.ID]; !ok {
			return nil, fmt.Errorf("no P2P keys found for DON %d", donMetadata.ID)
		}

		p2pKeys := keys.P2PKeys[donMetadata.ID]
		if len(p2pKeys.PeerIDs) != len(donMetadata.NodesMetadata) {
			return nil, fmt.Errorf("number of P2P keys for DON %d does not match the number of nodes. Expected %d, got %d", donMetadata.ID, len(donMetadata.NodesMetadata), len(p2pKeys.PeerIDs))
		}
		for idx, nodeMetadata := range donMetadata.NodesMetadata {
			nodeMetadata.Labels = append(nodeMetadata.Labels, &cretypes.Label{
				Key:   node.NodeP2PIDKey,
				Value: p2pKeys.PeerIDs[idx],
			})
		}

		if _, ok := keys.EVMKeys[donMetadata.ID]; !ok {
			return nil, fmt.Errorf("no EVM keys found for DON %d", donMetadata.ID)
		}

		chainIDsToEVMKeys := keys.EVMKeys[donMetadata.ID]

		// Now add the EVM addresses to the node metadata
		for chainID, evmKeys := range chainIDsToEVMKeys {
			chainSelector, selectorErr := chainselectors.SelectorFromChainId(libc.MustSafeUint64(int64(chainID)))
			if selectorErr != nil {
				return nil, errors.Wrapf(selectorErr, "failed to get chain selector for chain ID %d", chainID)
			}
			if len(evmKeys.PublicAddresses) != len(donMetadata.NodesMetadata) {
				return nil, fmt.Errorf("number of EVM keys for DON %d and chain ID %d does not match the number of nodes. Expected %d, got %d", donMetadata.ID, chainID, len(donMetadata.NodesMetadata), len(evmKeys.PublicAddresses))
			}
			for idx, nodeMetadata := range donMetadata.NodesMetadata {
				nodeMetadata.Labels = append(nodeMetadata.Labels, &cretypes.Label{
					Key:   node.AddressKeyFromSelector(chainSelector),
					Value: evmKeys.PublicAddresses[idx].Hex(),
				})
			}
		}
	}

	return topology, nil
}

// secrets struct mirrors `Secrets` struct in "github.com/smartcontractkit/chainlink/v2/core/config/toml"
// we use a copy to avoid depending on the core config package, we consider it safe, because that struct changes very rarely
type secrets struct {
	EVM    ethKeys `toml:",omitempty"` // choose EVM as the TOML field name to align with relayer config convention
	P2PKey p2PKey  `toml:",omitempty"`
}

type p2PKey struct {
	JSON     *string
	Password *string
}

type ethKeys struct {
	Keys []*ethKey
}

type ethKey struct {
	JSON     *string
	ID       *int
	Password *string
}

// struct required for reading "address" from this bit of encrypted JSON:
// JSON = '{"address":"e753ac0b6e175ce3a939c55433a0109c5a6f8777"}'
type evmJSON struct {
	Address string `json:"address"`
}

var publicEVMAddressFromEncryptedJSON = func(jsonString string) (string, error) {
	var eJSON evmJSON
	err := json.Unmarshal([]byte(jsonString), &eJSON)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal evm json")
	}
	return eJSON.Address, nil
}

// struct required for reading "peerID" from this bit of encrypted JSON:
// JSON = '{"keyType":"P2P","publicKey":"f3c458c9064bdde449a3904ba8d3f8f5ebf79623077430325252c3368f920199","peerID":"p2p_12D3KooWSDvtYVF3FoyGeMrmDxYeJZMzbEyMHRwmf5GUSqgJhST2"}'
type p2pJSON struct {
	PeerID string `json:"peerID"`
}

var publicP2PAddressFromEncryptedJSON = func(jsonString string) (string, error) {
	var pJSON p2pJSON
	err := json.Unmarshal([]byte(jsonString), &pJSON)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal p2p json")
	}
	return pJSON.PeerID, nil
}

func KeysOutputFromConfig(nodeSets []*cretypes.CapabilitiesAwareNodeSet) (*cretypes.GenerateKeysOutput, error) {
	output := &cretypes.GenerateKeysOutput{
		EVMKeys: make(cretypes.DonsToEVMKeys),
		P2PKeys: make(cretypes.DonsToP2PKeys),
	}
	p2pKeysFoundPerDon := make(map[uint32]int)
	evmKeysFoundPerDon := make(map[uint32]int)
	for donIdx, nodeSet := range nodeSets {
		donIdxUint32 := uint32(donIdx) // #nosec G115: ignore as this will NEVER happen, we don't have zillions of DONs
		p2pKeys := types.P2PKeys{}
		evmKeysPerChainID := make(cretypes.ChainIDToEVMKeys)
		for nodeIdx, nodeSpec := range nodeSet.NodeSpecs {
			if nodeSpec.Node.TestSecretsOverrides != "" {
				var sSecrets secrets
				unmarshallErr := toml.Unmarshal([]byte(nodeSpec.Node.TestSecretsOverrides), &sSecrets)
				if unmarshallErr != nil {
					return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal secrets for node %d in DON %d", nodeIdx, donIdx)
				}

				// For simplicity we will allow importing only both P2P keys and EVM keys, not just one of them
				if sSecrets.P2PKey.JSON == nil || sSecrets.P2PKey.Password == nil {
					return nil, fmt.Errorf("P2P key or password is nil for node %d in DON %d", nodeIdx, donIdx)
				}
				p2pKeys.EncryptedJSONs = append(p2pKeys.EncryptedJSONs, []byte(*sSecrets.P2PKey.JSON))
				p2pKeys.Password = *sSecrets.P2PKey.Password
				peerID, peerIDErr := publicP2PAddressFromEncryptedJSON(*sSecrets.P2PKey.JSON)
				if peerIDErr != nil {
					return nil, errors.Wrapf(peerIDErr, "failed to get public p2p address for node %d in DON %d from encrypted JSON", nodeIdx, donIdx)
				}
				p2pKeys.PeerIDs = append(p2pKeys.PeerIDs, peerID)
				p2pKeysFoundPerDon[donIdxUint32]++
				if len(sSecrets.EVM.Keys) == 0 {
					return nil, fmt.Errorf("EVM keys is nil for node %d in DON %d", nodeIdx, donIdx)
				}

				for _, evmKey := range sSecrets.EVM.Keys {
					if evmKey.JSON == nil || evmKey.Password == nil || evmKey.ID == nil {
						return nil, fmt.Errorf("EVM key or password or ID is nil for node %d in DON %d", nodeIdx, donIdx)
					}

					publicEVMAddress, publicEVMAddressErr := publicEVMAddressFromEncryptedJSON(*evmKey.JSON)
					if publicEVMAddressErr != nil {
						return nil, errors.Wrapf(publicEVMAddressErr, "failed to get public evm address for node %d in DON %d from encrypted JSON", nodeIdx, donIdx)
					}

					if _, ok := evmKeysPerChainID[*evmKey.ID]; !ok {
						evmKeysPerChainID[*evmKey.ID] = &types.EVMKeys{}
					}

					evmKeysPerChainID[*evmKey.ID].EncryptedJSONs = append(evmKeysPerChainID[*evmKey.ID].EncryptedJSONs, []byte(*evmKey.JSON))
					evmKeysPerChainID[*evmKey.ID].PublicAddresses = append(evmKeysPerChainID[*evmKey.ID].PublicAddresses, common.HexToAddress(publicEVMAddress))
					evmKeysPerChainID[*evmKey.ID].Password = *evmKey.Password
				}
				evmKeysFoundPerDon[donIdxUint32]++
			}
		}
		// +1 because we use 1-based indexing in the CRE
		donIndexToUse := uint32(donIdx + 1) // #nosec G115
		output.P2PKeys[donIndexToUse] = &p2pKeys
		output.EVMKeys[donIndexToUse] = evmKeysPerChainID
	}

	anyFound := false
	// Validate that we found keys for all nodes in all DONs
	for donIdx, nodeSet := range nodeSets {
		donIdxUint32 := uint32(donIdx) // #nosec G115
		if p2pKeysFoundPerDon[donIdxUint32] != 0 && len(nodeSet.NodeSpecs) != p2pKeysFoundPerDon[donIdxUint32] {
			return nil, fmt.Errorf("number of P2P keys found for DON %d does not match the number of nodes. Expected %d, got %d", donIdx, len(nodeSet.NodeSpecs), p2pKeysFoundPerDon[donIdxUint32])
		}
		if evmKeysFoundPerDon[donIdxUint32] != 0 && len(nodeSet.NodeSpecs) != evmKeysFoundPerDon[donIdxUint32] {
			return nil, fmt.Errorf("number of EVM keys found for DON %d does not match the number of nodes. Expected %d, got %d", donIdx, len(nodeSet.NodeSpecs), evmKeysFoundPerDon[donIdxUint32])
		}
		if p2pKeysFoundPerDon[donIdxUint32] != 0 && evmKeysFoundPerDon[donIdxUint32] != 0 {
			anyFound = true
		}
	}

	if !anyFound {
		// If no keys were found for any DON, we can return empty output
		return nil, nil
	}

	return output, nil
}

func GenereteKeys(input *cretypes.GenerateKeysInput) (*cretypes.GenerateKeysOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	if input.Out != nil {
		return input.Out, nil
	}

	output := &cretypes.GenerateKeysOutput{
		EVMKeys: make(cretypes.DonsToEVMKeys),
		P2PKeys: make(cretypes.DonsToP2PKeys),
	}

	for _, donMetadata := range input.Topology.DonsMetadata {
		if input.GenerateP2PKeys {
			p2pKeys, err := crypto.GenerateP2PKeys(input.Password, len(donMetadata.NodesMetadata))
			if err != nil {
				return nil, errors.Wrap(err, "failed to generate P2P keys")
			}
			output.P2PKeys[donMetadata.ID] = p2pKeys
		}

		if len(input.GenerateEVMKeysForChainIDs) > 0 {
			for _, chainID := range input.GenerateEVMKeysForChainIDs {
				evmKeys, err := crypto.GenerateEVMKeys(input.Password, len(donMetadata.NodesMetadata))
				if err != nil {
					return nil, errors.Wrap(err, "failed to generate EVM keys")
				}
				if _, ok := output.EVMKeys[donMetadata.ID]; !ok {
					output.EVMKeys[donMetadata.ID] = make(cretypes.ChainIDToEVMKeys)
				}
				output.EVMKeys[donMetadata.ID][chainID] = evmKeys
			}
		}
	}

	return output, nil
}
