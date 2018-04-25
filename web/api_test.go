package web

import (
	"net/url"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestApi_ParsePaginatedRequest(t *testing.T) {
	tests := []struct {
		name        string
		sizeParam   string
		offsetParam string
		err         bool
		size        int
		offset      int
	}{
		{"blank values", "", "", false, 25, 0},
		{"valid sizeParam", "10", "", false, 10, 0},
		{"valid offsetParam", "", "10", false, 25, 10},
		{"invalid sizeParam", "xhje", "", true, 0, 0},
		{"invalid offsetParam", "", "ewjh", true, 0, 0},
		{"small sizeParam", "0", "", true, 0, 0},
		{"negative offsetParam", "", "-1", true, 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			size, offset, err := ParsePaginatedRequest(test.sizeParam, test.offsetParam)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.size, size)
			assert.Equal(t, test.offset, offset)
		})
	}
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
	tests := []struct {
		name     string
		path     string
		size     int
		offset   int
		count    int
		resource interface{}
		err      bool
		output   string
	}{
		{
			"a single resource",
			"/v2/index", 1, 0, 0, TestResource{Title: "Item"},
			false, `{"data":{"type":"testResources","id":"1","attributes":{"Title":"Item"}}}`,
		},
		{
			"a resource collection",
			"/v2/index", 1, 0, 0, []TestResource{TestResource{Title: "Item 1"}, TestResource{Title: "Item 2"}},
			false, `{"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}},{"type":"testResources","id":"1","attributes":{"Title":"Item 2"}}]}`,
		},
		{
			"first page of collection results",
			"/v2/index", 1, 0, 3, []TestResource{TestResource{Title: "Item 1"}},
			false, `{"links":{"next":"/v2/index?offset=1\u0026size=1"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}}]}`,
		},
		{
			"middle page of collection results",
			"/v2/index", 1, 1, 3, []TestResource{TestResource{Title: "Item 2"}},
			false, `{"links":{"next":"/v2/index?offset=2\u0026size=1","prev":"/v2/index?offset=0\u0026size=1"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 2"}}]}`,
		},
		{
			"end page of collection results",
			"/v2/index", 1, 2, 3, []TestResource{TestResource{Title: "Item 3"}},
			false, `{"links":{"prev":"/v2/index?offset=1\u0026size=1"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 3"}}]}`,
		},
		{
			"path with existing query",
			"/v2/index?authToken=3123", 1, 0, 2, []TestResource{TestResource{Title: "Item 1"}},
			false, `{"links":{"next":"/v2/index?authToken=3123\u0026offset=1\u0026size=1"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}}]}`,
		},
		{
			"json marshalling failure",
			"/v2/index", 1, 0, 0, "",
			true, ``,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url, err := url.Parse(test.path)
			assert.NoError(t, err)
			buffer, err := NewPaginatedResponse(*url, test.size, test.offset, test.count, test.resource)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.output, string(buffer))
		})
	}
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
