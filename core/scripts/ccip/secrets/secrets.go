package secrets

import (
	"fmt"
	"os"
	"strconv"
)

func GetRPC(chainID uint64) string {
	envVariable := "RPC_" + strconv.FormatUint(chainID, 10)
	rpc := os.Getenv(envVariable)
	if rpc != "" {
		return rpc
	}
	panic(fmt.Errorf("RPC not found. Please set the environment variable for chain %d e.g. RPC_420=https://rpc.420.com", chainID))
}
