package presenters

import "fmt"

// HALResource is a wrapper for any collection of resources that is paginated.
type HALResource struct {
	Data  interface{} `json:"data"`
	Links Links       `json:"_links"`
}

// Links refers to the relations to the currently returned record in HAL spec.
type Links struct {
	Next *Ref `json:"next,omitempty"`
	Self *Ref `json:"self,omitempty"`
	Prev *Ref `json:"prev,omitempty"`
}

// Ref is a HAL JSON structure for a link
type Ref struct {
	Href string `json:"href"`
}

// NewPaginatedResponse returns a HALResource with links to next and previous collection pages
func NewHALResponse(path string, size, offset, count int, resource interface{}) *HALResource {
	var next *Ref
	var prev *Ref

	if count > size {
		if offset+size < count {
			nextURI := fmt.Sprintf("%s?size=%d&offset=%d", path, size, offset+size)
			next = &Ref{Href: nextURI}
		}
		if offset > 0 {
			prevURI := fmt.Sprintf("%s?size=%d&offset=%d", path, size, offset-size)
			prev = &Ref{Href: prevURI}
		}
	}

	return &HALResource{
		Data: resource,
		Links: Links{
			Next: next,
			Prev: prev,
		},
	}
}
