package evm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	ol := OffchainLookup{
		url:              "",
		extraData:        nil,
		fields:           []string{"gender", "name", "count", "probability"},
		callbackFunction: [4]byte{},
	}
	body := []byte(`{"count":256938,"gender":"male","name":"chris","probability":0.92}`)

	fmt.Println(string(body))
	values, err := ol.parseJson(body)

	assert.Equal(t, nil, err)
	assert.Equal(t, "male", values[0], "gender")
	assert.Equal(t, "chris", values[1], "name")
	assert.Equal(t, "256938", values[2], "count")
	assert.Equal(t, "0.92", values[3], "probability")
	fmt.Printf("%+v\n", values)
}
