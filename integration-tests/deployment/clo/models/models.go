// TODO: KS-455: Refactor this package to use chainlink-common
// Fork from: https://github.com/smartcontractkit/feeds-manager/tree/develop/api/models
// until it can be refactored in cahinlink-common.

// Package models provides generated go types that reflect the GraphQL types
// defined in the API schemas.
//
// To maintain compatibility with the default JSON interfaces, any necessary model
// overrides can be defined this package.
package models

import "go/types"

// Generic Error model to override existing GQL Error types and allow unmarshaling of GQL Error unions
type Error struct {
	Typename string    `json:"__typename,omitempty"`
	Message  string    `json:"message,omitempty"`
	Path     *[]string `json:"path,omitempty"`
}

func (e *Error) Underlying() types.Type { return e }
func (e *Error) String() string         { return "Error" }

// Unmarshal GQL Time fields into go strings because by default, time.Time results in zero-value
// timestamps being present in the CLI output, e.g. "createdAt": "0001-01-01T00:00:00Z".
// This is because the default JSON interfaces don't recognize it as an empty value for Time.time
// and fail to omit it when using `json:"omitempty"` tags.
type Time string
