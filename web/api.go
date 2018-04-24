package web

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/manyminds/api2go/jsonapi"
)

const (
	// PaginationDefault is the number of records to supply from a paginated
	// request when no size param is supplied.
	PaginationDefault = 25

	// MediaType is the response header for JSONAPI documents.
	MediaType = "application/vnd.api+json"

	KeyNextLink     = "next"
	KeyPreviousLink = "prev"
)

// ParsePaginatedRequest parses the parameters that control pagination for a
// collection request, returning the size and offset if specified, or a
// sensible default.
func ParsePaginatedRequest(sizeParam, offsetParam string) (int, int, error) {
	var err error
	var offset int
	size := PaginationDefault

	if sizeParam != "" {
		if size, err = strconv.Atoi(sizeParam); err != nil || size <= 0 {
			return 0, 0, fmt.Errorf("invalid size param, error: %+v", err)
		}
	}

	if offsetParam != "" {
		if offset, err = strconv.Atoi(offsetParam); err != nil || offset < 0 {
			return 0, 0, fmt.Errorf("invalid offset param, error: %+v", err)
		}
	}

	return size, offset, nil
}

// NewPaginatedResponse returns a HALResource with links to next and previous collection pages
func NewPaginatedResponse(path string, size, offset, count int, resource interface{}) ([]byte, error) {
	document, err := jsonapi.MarshalToStruct(resource, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal jobs to struct: %+v", err)
	}

	document.Links = make(jsonapi.Links)
	if count > size {
		if offset+size < count {
			nextURI := fmt.Sprintf("%s?size=%d&offset=%d", path, size, offset+size)
			document.Links["next"] = jsonapi.Link{Href: nextURI}
		}
		if offset > 0 {
			prevURI := fmt.Sprintf("%s?size=%d&offset=%d", path, size, offset-size)
			document.Links["prev"] = jsonapi.Link{Href: prevURI}
		}
	}

	return json.Marshal(document)
}
