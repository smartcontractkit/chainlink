package utils

import "github.com/aptos-labs/aptos-go-sdk"

func IsModuleDeployed(client aptos.AptosRpcClient, address aptos.AccountAddress, moduleName string) (bool, error) {
	module, err := client.AccountModule(address, moduleName)
	if err != nil {
		return false, err
	}
	if module == nil {
		return false, nil
	}
	return true, nil
}
