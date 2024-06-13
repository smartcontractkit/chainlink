package utils

import (
	"math/big"
	"net"
	"os"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
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

// GetenvOrDefault returns the value of an environment variable or a default value if not set.
func GetenvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
