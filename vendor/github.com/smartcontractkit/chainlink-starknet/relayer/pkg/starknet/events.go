package starknet

import (
	caigogw "github.com/dontpanicdao/caigo/gateway"
	caigotypes "github.com/dontpanicdao/caigo/types"
)

func IsEventFromContract(event *caigogw.Event, address caigotypes.Hash, eventName string) bool {
	eventKey := caigotypes.GetSelectorFromName(eventName)
	// encoded event name guaranteed to be at index 0
	return CompareAddress(event.FromAddress, address.String()) && event.Keys[0].Cmp(eventKey) == 0
}
