package web

import (
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestApi_ParsePaginatedRequest(t *testing.T) {
	var size int
	var offset int
	var err error

	size, offset, err = ParsePaginatedRequest("", "")
	assert.NoError(t, err)
	assert.Equal(t, size, 25)
	assert.Equal(t, offset, 0)

	size, offset, err = ParsePaginatedRequest("10", "")
	assert.NoError(t, err)
	assert.Equal(t, size, 10)
	assert.Equal(t, offset, 0)

	size, offset, err = ParsePaginatedRequest("", "10")
	assert.NoError(t, err)
	assert.Equal(t, size, 25)
	assert.Equal(t, offset, 10)

	size, offset, err = ParsePaginatedRequest("x!", "")
	assert.Error(t, err)

	size, offset, err = ParsePaginatedRequest("0", "")
	assert.Error(t, err)

	size, offset, err = ParsePaginatedRequest("", "hh")
	assert.Error(t, err)

	size, offset, err = ParsePaginatedRequest("", "-1")
	assert.Error(t, err)
}

type TestResource struct {
	Title string
}

func (r TestResource) GetID() string {
	return "1"
}

func (r *TestResource) SetID(value string) error {
	return nil
}

func TestApi_NewPaginatedResponse(t *testing.T) {
	var buffer []byte
	var err error

	resource := TestResource{Title: "Item"}

	buffer, err = NewPaginatedResponse("/v2/index", 1, 0, 0, resource)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":{"type":"testResources","id":"1","attributes":{"Title":"Item"}}}`, string(buffer))

	resources := []TestResource{TestResource{Title: "Item 1"}, TestResource{Title: "Item 2"}}

	buffer, err = NewPaginatedResponse("/v2/index", 25, 0, 2, resources)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}},{"type":"testResources","id":"1","attributes":{"Title":"Item 2"}}]}`, string(buffer))

	resources = []TestResource{TestResource{Title: "Item 1"}}

	buffer, err = NewPaginatedResponse("/v2/index", 1, 0, 3, resources)
	assert.NoError(t, err)
	assert.Equal(t, `{"links":{"next":"/v2/index?size=1\u0026offset=1"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}}]}`, string(buffer))

	buffer, err = NewPaginatedResponse("/v2/index", 1, 1, 3, resources)
	assert.NoError(t, err)
	assert.Equal(t, `{"links":{"next":"/v2/index?size=1\u0026offset=2","prev":"/v2/index?size=1\u0026offset=0"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}}]}`, string(buffer))

	buffer, err = NewPaginatedResponse("/v2/index", 1, 2, 3, resources)
	assert.NoError(t, err)
	assert.Equal(t, `{"links":{"prev":"/v2/index?size=1\u0026offset=1"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}}]}`, string(buffer))
}

func TestPagination_ParsePaginatedResponse(t *testing.T) {
	var docs []TestResource
	var links jsonapi.Links

	err := ParsePaginatedResponse([]byte(`{"data":[{"type":"testResources","id":"1","attributes":{"Title":"album 1"}}]}`), &docs, &links)
	assert.NoError(t, err)
	assert.Equal(t, "album 1", docs[0].Title)

	// Typo in "type"
	err = ParsePaginatedResponse([]byte(`{"data":[{"type":"testNotResources","id":"1","attributes":{}}]}`), &docs, &links)
	assert.Error(t, err)

	// Typo in "links"
	err = ParsePaginatedResponse([]byte(`{"links":[],"data":[{"type":"testResources","id":"1","attributes":{}}]}`), &docs, &links)
	assert.Error(t, err)
}
