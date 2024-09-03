package utils

import (
	"math/big"
	"net"
	"path/filepath"
	"runtime"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// ProjectRoot Root folder of this project
	ProjectRoot = filepath.Join(filepath.Dir(b), "/../..")
	// ContractsDir path to our contracts
	ContractsDir = filepath.Join(ProjectRoot, "contracts", "target", "deploy")
)

func MustURL(s string) *commonconfig.URL {
	var u commonconfig.URL
	if err := u.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return &u
}

func MustIP(s string) *net.IP {
	var ip net.IP
	if err := ip.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return &ip
}

func BigIntSliceContains(slice []*big.Int, b *big.Int) bool {
	for _, a := range slice {
		if b.Cmp(a) == 0 {
			return true
		}
	}
	return false
}
