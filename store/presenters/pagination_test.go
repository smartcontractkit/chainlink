package presenters

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination_UnmarshalJSON(t *testing.T) {
	doc := JobSpecsDocument{}
	err := json.Unmarshal([]byte(`{"data":[{"type":"specs","id":"1","attributes":{}}]}`), &doc)
	assert.NoError(t, err)

	// Typo in "type"
	err = json.Unmarshal([]byte(`{"data":[{"type":"spcs","id":"1","attributes":{}}]}`), &doc)
	assert.Error(t, err)

	// Typo in "links"
	err = json.Unmarshal([]byte(`{"links":[],"data":[{"type":"specs","id":"1","attributes":{}}]}`), &doc)
	assert.Error(t, err)
}
