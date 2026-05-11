package main

import (
	"fmt"
	"hash/fnv"
)

func shardForPackage(pkg string, shardCount int) int {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(pkg))
	return int(int64(hasher.Sum32()) % int64(shardCount))
}

func main() {
	variations := []string{
		"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana",
		"github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana",
	}
	for _, v := range variations {
		fmt.Printf("Pkg: %s, Shard: %d/7\n", v, shardForPackage(v, 7))
	}
}
