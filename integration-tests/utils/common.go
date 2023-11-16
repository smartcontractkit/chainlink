package utils

import (
	"math/big"
	"net"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func MustURL(s string) *models.URL {
	var u models.URL
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
		// Return background context if testing.T not set
		if t == nil {
			return ctx
		}
