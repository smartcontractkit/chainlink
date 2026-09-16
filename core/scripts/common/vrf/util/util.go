package util

import (
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/model"
)

func MapToSendingKeyArr(nodeSendingKeys []string) []model.SendingKey {
	sendingKeys := make([]model.SendingKey, 0, len(nodeSendingKeys))

	for _, key := range nodeSendingKeys {
		sendingKeys = append(sendingKeys, model.SendingKey{Address: key})
	}
	return sendingKeys
}

func MapToAddressArr(sendingKeys []model.SendingKey) []string {
	sendingKeysString := make([]string, 0, len(sendingKeys))
	for _, sendingKey := range sendingKeys {
		sendingKeysString = append(sendingKeysString, sendingKey.Address)
	}
	return sendingKeysString
}
