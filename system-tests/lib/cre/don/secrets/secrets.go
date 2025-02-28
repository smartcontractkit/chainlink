package secrets

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
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
			for _, chainID := range input.EVMKeys.ChainIDs {
				nodeSecret.EthKeys.EthKeys = append(nodeSecret.EthKeys.EthKeys, types.NodeEthKey{
					JSON:     string(input.EVMKeys.EncryptedJSONs[i]),
					Password: input.EVMKeys.Password,
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
