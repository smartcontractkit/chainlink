package util

import (
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/model"
)

func MapToSendingKeyArr(nodeSendingKeys []string) []model.SendingKey {
	var sendingKeys []model.SendingKey

	for _, key := range nodeSendingKeys {
		sendingKeys = append(sendingKeys, model.SendingKey{Address: key})
	}
	return sendingKeys
}

func MapToAddressArr(sendingKeys []model.SendingKey) []string {
	var sendingKeysString []string
	for _, sendingKey := range sendingKeys {
		sendingKeysString = append(sendingKeysString, sendingKey.Address)
	}
	return sendingKeysString
}
