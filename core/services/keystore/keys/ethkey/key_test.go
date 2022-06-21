package ethkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEthKey_Type(t *testing.T) {
	k := Key{
		IsFunding: true,
	}
	k2 := Key{
		IsFunding: false,
	}

	assert.Equal(t, k.Type(), "funding")
	assert.Equal(t, k2.Type(), "sending")
}
