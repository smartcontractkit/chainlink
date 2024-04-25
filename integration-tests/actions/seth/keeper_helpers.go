package actions_seth

import (
	"github.com/ethereum/go-ethereum/common"
)

var ZeroAddress = common.Address{}

var INSUFFICIENT_EPHEMERAL_KEYS = `
Error: Insufficient Ephemeral Addresses for Simulated Network

To operate on a simulated network, you must configure at least one ephemeral address. Currently, %d ephemeral address(es) are set. Please update your TOML configuration file as follows to meet this requirement:
[Seth] ephemeral_addresses_number = 1

This adjustment ensures that your setup is minimaly viable. Although it is highly recommended to use at least 20 ephemeral addresses.
`

var INSUFFICIENT_STATIC_KEYS = `
Error: Insufficient Private Keys for Live Network

To run this test on a live network, you must either:
1. Set at least two private keys in the '[Network.WalletKeys]' section of your TOML configuration file. Example format:
   [Network.WalletKeys]
   NETWORK_NAME=["PRIVATE_KEY_1", "PRIVATE_KEY_2"]
2. Set at least two private keys in the '[Network.EVMNetworks.NETWORK_NAME] section of your TOML configuration file. Example format:
   evm_keys=["PRIVATE_KEY_1", "PRIVATE_KEY_2"]

Currently, only %d private key/s is/are set.

Recommended Action:
Distribute your funds across multiple private keys and update your configuration accordingly. Even though 1 private key is sufficient for testing, it is highly recommended to use at least 10 private keys.
`
