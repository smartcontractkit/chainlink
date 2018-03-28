package web

import (
	"github.com/smartcontractkit/chainlink/services"
)

// StatusCodeForError returns an http status code for an error type.
func StatusCodeForError(err interface{}) int {
	switch err.(type) {
	case *services.ValidationError:
		return 400
	default:
		return 500
	}
}
