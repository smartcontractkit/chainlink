package secrets

import "strconv"

var chainIdToRPC = map[uint64]string{}

func GetRPC(chainID uint64) string {
	if rpc, ok := chainIdToRPC[chainID]; ok {
		return rpc
	}
	panic("RPC not found. Please check secrets.go for chainID " + strconv.FormatUint(chainID, 10))
}
