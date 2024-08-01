# Reflectx [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/scylladb/go-reflectx) [![Go Report Card](https://goreportcard.com/badge/github.com/scylladb/go-reflectx)](https://goreportcard.com/report/github.com/scylladb/go-reflectx) [![Build Status](https://travis-ci.org/scylladb/go-reflectx.svg?branch=master)](https://travis-ci.org/scylladb/go-reflectx)

Package reflectx implements extensions to the standard reflect lib suitable for implementing marshalling and unmarshalling packages.
The main Mapper type allows for Go-compatible named attribute access, including accessing embedded struct attributes and the ability to use functions and struct tags to customize field names.

This is a standalone version of `reflectx` package that originates from an SQL row to struct mapper [sqlx](https://github.com/jmoiron/sqlx).
We are using it at [Scylla gocqlx](https://github.com/scylladb/gocqlx) for scanning of CQL results to structs and slices.

## Example

This example demonstrates usage of the reflectx package to automatically bind URL parameters to a request model.

```go
type RequestContext struct {
	SessionID string `http:"sid"`
}

type SearchRequest struct {
	RequestContext
	Labels     []string `http:"l"`
	MaxResults int      `http:"max"`
	Exact      bool     `http:"x"`
}

func Search(w http.ResponseWriter, r *http.Request) {
	// URL /search?sid=id&l=foo&l=bar&max=100&x=true
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var data SearchRequest
	if err := bindParams(r, &data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("%+v", data) // "RequestContext:{SessionID:id} Labels:[foo bar] MaxResults:100 Exact:true}"
}
```

See the full example in [example_test.go](example_test.go).

## License

Copyright (C) 2019 ScyllaDB

This project is distributed under the Apache 2.0 license. See the [LICENSE](https://github.com/scylladb/go-reflectx/blob/master/LICENSE) file for details.
It contains software from:

* [sqlx project](https://github.com/jmoiron/sqlx), licensed under the MIT license
