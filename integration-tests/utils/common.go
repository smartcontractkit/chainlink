package utils

import (
	"net"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func Ptr[T any](t T) *T { return &t }

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
