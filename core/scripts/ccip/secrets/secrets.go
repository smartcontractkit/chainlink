package secrets

import (
	"os"
	"strconv"
)

func GetRPC(chainID uint64) string {
	envVariable := "RPC_" + strconv.FormatUint(chainID, 10)
	rpc := os.Getenv(envVariable)
	if rpc != "" {
		return rpc
	}
	panic("RPC not found. Please check secrets.go for chainID " + strconv.FormatUint(chainID, 10))
}
