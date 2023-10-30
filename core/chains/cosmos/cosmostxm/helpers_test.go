package cosmostxm

import "golang.org/x/exp/maps"

func (ka *keystoreAdapter) Accounts() ([]string, error) {
	ka.mutex.Lock()
	defer ka.mutex.Unlock()
	err := ka.updateMappingLocked()
	if err != nil {
		return nil, err
	}
	addresses := maps.Keys(ka.addressToPubKey)

	return addresses, nil
}
