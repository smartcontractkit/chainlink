package web

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/manyminds/api2go/jsonapi"
)

const (
	// PaginationDefault is the number of records to supply from a paginated
	// request when no size param is supplied.
	PaginationDefault = 25

	// MediaType is the response header for JSONAPI documents.
	MediaType = "application/vnd.api+json"

	// KeyNextLink is the name of the key that contains the HREF for the next
	// document in a paginated response.
	KeyNextLink = "next"
	// KeyPreviousLink is the name of the key that contains the HREF for the
	// previous document in a paginated response.
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

func nextLink(url url.URL, size, offset int) jsonapi.Link {
	query := url.Query()
	query.Add("size", strconv.Itoa(size))
	query.Add("offset", strconv.Itoa(offset+size))
	url.RawQuery = query.Encode()
	return jsonapi.Link{Href: url.String()}
}

func prevLink(url url.URL, size, offset int) jsonapi.Link {
	query := url.Query()
	query.Add("size", strconv.Itoa(size))
	query.Add("offset", strconv.Itoa(offset-size))
	url.RawQuery = query.Encode()
	return jsonapi.Link{Href: url.String()}
}

// NewPaginatedResponse returns a jsonapi.Document with links to next and previous collection pages
func NewPaginatedResponse(url url.URL, size, offset, count int, resource interface{}) ([]byte, error) {
	document, err := jsonapi.MarshalToStruct(resource, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource to struct: %+v", err)
	}

	document.Links = make(jsonapi.Links)
	if count > size {
		if offset+size < count {
			document.Links[KeyNextLink] = nextLink(url, size, offset)
		}
		if offset > 0 {
			document.Links[KeyPreviousLink] = prevLink(url, size, offset)
		}
	}

	return json.Marshal(document)
}

// ParsePaginatedResponse parse a JSONAPI response
func ParsePaginatedResponse(input []byte, resource interface{}, links *jsonapi.Links) error {
	// First unmarshal using the jsonAPI into the Resource, whatever it may be,
	// as is api2go will discard the links
	err := jsonapi.Unmarshal(input, resource)
	if err != nil {
		return fmt.Errorf("unable to unmarshal Data record: %+v", err)
	}

	// Unmarshal using the stdlib Unmarshal to extract the Links part of the document
	document := jsonapi.Document{}
	err = json.Unmarshal(input, &document)
	if err != nil {
		return fmt.Errorf("unable to unmarshal Links: %+v", err)
	}
	*links = document.Links

	return nil
}
