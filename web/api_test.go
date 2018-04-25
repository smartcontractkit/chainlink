package web

import (
	"net/url"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestApi_ParsePaginatedRequest(t *testing.T) {
	tests := []struct {
		name      string
		sizeParam string
		pageParam string
		err       bool
		size      int
		page      int
		offset    int
	}{
		{"blank values", "", "", false, 25, 1, 0},
		{"valid sizeParam", "10", "", false, 10, 1, 0},
		{"valid pageParam", "", "3", false, 25, 3, 50},
		{"invalid sizeParam", "xhje", "", true, 0, 0, 0},
		{"invalid pageParam", "", "ewjh", true, 0, 0, 0},
		{"small sizeParam", "0", "", true, 0, 0, 0},
		{"negative pageParam", "", "-1", true, 0, 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			size, page, offset, err := ParsePaginatedRequest(test.sizeParam, test.pageParam)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.size, size)
			assert.Equal(t, test.page, page)
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
		page     int
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
			"/v2/index", 5, 1, 7, []TestResource{TestResource{Title: "Item 1"}},
			false, `{"links":{"next":"/v2/index?page=2\u0026size=5"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}}]}`,
		},
		{
			"middle page of collection results",
			"/v2/index", 5, 2, 13, []TestResource{TestResource{Title: "Item 2"}},
			false, `{"links":{"next":"/v2/index?page=3\u0026size=5","prev":"/v2/index?page=1\u0026size=5"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 2"}}]}`,
		},
		{
			"end page of collection results",
			"/v2/index", 5, 3, 13, []TestResource{TestResource{Title: "Item 3"}},
			false, `{"links":{"prev":"/v2/index?page=2\u0026size=5"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 3"}}]}`,
		},
		{
			"path with existing query",
			"/v2/index?authToken=3123", 1, 0, 2, []TestResource{TestResource{Title: "Item 1"}},
			false, `{"links":{"next":"/v2/index?authToken=3123\u0026page=1\u0026size=1"},"data":[{"type":"testResources","id":"1","attributes":{"Title":"Item 1"}}]}`,
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
			buffer, err := NewPaginatedResponse(*url, test.size, test.page, test.count, test.resource)
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
