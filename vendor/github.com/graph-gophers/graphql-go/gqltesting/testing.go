package gqltesting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/errors"
)

// Test is a GraphQL test case to be used with RunTest(s).
type Test struct {
	Context        context.Context
	Schema         *graphql.Schema
	Query          string
	OperationName  string
	Variables      map[string]interface{}
	ExpectedResult string
	ExpectedErrors []*errors.QueryError
	RawResponse    bool
}

// RunTests runs the given GraphQL test cases as subtests.
func RunTests(t *testing.T, tests []*Test) {
	if len(tests) == 1 {
		RunTest(t, tests[0])
		return
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			RunTest(t, test)
		})
	}
}

// RunTest runs a single GraphQL test case.
func RunTest(t *testing.T, test *Test) {
	if test.Context == nil {
		test.Context = context.Background()
	}
	result := test.Schema.Exec(test.Context, test.Query, test.OperationName, test.Variables)

	checkErrors(t, test.ExpectedErrors, result.Errors)

	if test.ExpectedResult == "" {
		if result.Data != nil {
			t.Fatalf("got: %s", result.Data)
			t.Fatalf("want: null")
		}
		return
	}

	// Verify JSON to avoid red herring errors.
	var got []byte

	if test.RawResponse {
		value, err := result.Data.MarshalJSON()
		if err != nil {
			t.Fatalf("got: unable to marshal JSON response: %s", err)
		}
		got = value
	} else {
		value, err := formatJSON(result.Data)
		if err != nil {
			t.Fatalf("got: invalid JSON: %s", err)
		}
		got = value
	}

	want, err := formatJSON([]byte(test.ExpectedResult))
	if err != nil {
		t.Fatalf("want: invalid JSON: %s", err)
	}

	if !bytes.Equal(got, want) {
		t.Logf("got:  %s", got)
		t.Logf("want: %s", want)
		t.Fail()
	}
}

func formatJSON(data []byte) ([]byte, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	formatted, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return formatted, nil
}

func checkErrors(t *testing.T, want, got []*errors.QueryError) {
	sortErrors(want)
	sortErrors(got)

	// Clear the underlying error before the DeepEqual check.  It's too
	// much to ask the tester to include the raw failing error.
	for _, err := range got {
		err.Err = nil
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected error: got %+v, want %+v", got, want)
	}
}

func sortErrors(errors []*errors.QueryError) {
	if len(errors) <= 1 {
		return
	}
	sort.Slice(errors, func(i, j int) bool {
		return fmt.Sprintf("%s", errors[i].Path) < fmt.Sprintf("%s", errors[j].Path)
	})
}
