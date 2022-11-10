package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"testing"
)

func GetRandomPort() string {
	r, err := rand.Int(rand.Reader, big.NewInt(65535-1023))
	if err != nil {
		panic(fmt.Errorf("unexpected error generating random port: %w", err))
	}

	return strconv.Itoa(int(r.Int64() + 1024))
}

func IsPortOpen(t *testing.T, port string) bool {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		t.Log("error in checking port: ", err.Error())
		return false
	}
	defer l.Close()
	return true
}

func MustRandomPort(t *testing.T) string {
	for i := 0; i < 5; i++ {
		port := GetRandomPort()

		// check port if port is open
		if IsPortOpen(t, port) {
			t.Log("found open port: " + port)
			return port
		}
	}

	panic("unable to find open port")
}
