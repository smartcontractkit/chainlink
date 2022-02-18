package pipeline

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpAllowedIPS_isRestrictedIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ip           net.IP
		isRestricted bool
	}{
		{net.ParseIP("1.1.1.1"), false},
		{net.ParseIP("216.239.32.10"), false},
		{net.ParseIP("2001:4860:4860::8888"), false},
		{net.ParseIP("127.0.0.1"), true},
		{net.ParseIP("255.255.255.255"), true},
		{net.ParseIP("224.0.0.1"), true},
		{net.ParseIP("224.0.0.2"), true},
		{net.ParseIP("224.1.1.1"), true},
		{net.ParseIP("0.0.0.0"), true},
		{net.ParseIP("192.168.0.1"), true},
		{net.ParseIP("192.168.1.255"), true},
		{net.ParseIP("255.255.255.255"), true},
		{net.ParseIP("10.0.0.1"), true},
		{net.ParseIP("::1"), true},
		{net.ParseIP("fd57:03f9:9ef5:8a81::1"), true},
		{net.ParseIP("FD00::1"), true},
		{net.ParseIP("FF02::1"), true},
		{net.ParseIP("FE80:0000:0000:0000:abcd:abcd:abcd:abcd"), true},
		{net.IP{0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}, true},
		{net.IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}, true},
		{net.IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02}, true},
	}

	for _, test := range tests {
		t.Run(test.ip.String(), func(t *testing.T) {
			assert.Equal(t, test.isRestricted, isRestrictedIP(test.ip))
		})
	}
}
