// Package schema is used to read schema files
// go:generate go-bindata -ignore=\.go -pkg=schema -o=bindata.go ./...
package schema

import (
	"bytes"
	"embed"
	"fmt"
)

//go:embed *.graphql type/*.graphql
var fs embed.FS

// GetRootSchema reads the schema files and combines them into a single schema.
func GetRootSchema() (string, error) {
	b, err := fs.ReadFile("schema.graphql")
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	buf.Write(b)

	types, err := fs.ReadDir("type")
	if err != nil {
		return "", err
	}

	for _, t := range types {
		b, err = fs.ReadFile(fmt.Sprintf("type/%s", t.Name()))
		if err != nil {
			return "", err
		}

		buf.Write(b)

		// Add a newline if the file does not end in a newline.
		if len(b) > 0 && b[len(b)-1] != '\n' {
			buf.WriteByte('\n')
		}
	}

	return buf.String(), nil
}

// MustGetRootSchema reads the schema files and combines them into a single
// schema. It panics if there are any errors.
func MustGetRootSchema() string {
	s, err := GetRootSchema()
	if err != nil {
		panic(err)
	}

	return s
}
