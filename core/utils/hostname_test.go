package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hostnameport_validator(t *testing.T) {
	type testInput struct {
		data     string
		expected bool
	}
	testData := []testInput{
		{"bad..domain.name:234", false},
		{"extra.dot.com.", false},
		{"localhost:1234", true},
		{"192.168.1.1:1234", true},
		{":1234", true},
		{"domain.com:1334", true},
		{"this.domain.com:234", true},
		{"domain:75000", false},
		{"missing.port", false},
	}
	for _, td := range testData {
		valid := IsHostnamePort(td.data)
		assert.Equal(t, td.expected, valid)
	}
}
