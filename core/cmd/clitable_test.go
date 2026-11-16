package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/mattn/go-runewidth"
)

//go:embed testdata/tablewriter-golden.jsonl
var cliTableGoldenDefault []byte

//go:embed testdata/tablewriter-golden-eastasian.jsonl
var cliTableGoldenEastAsian []byte

type cliTableGoldenCase struct {
	Name    string
	Headers []string
	Rows    [][]string
}

type cliTableGoldenRecord struct {
	Name   string
	Output []byte
}

var cliTableGoldenCases = []cliTableGoldenCase{
	{Name: "service", Headers: []string{"ID", "Service", "LogLevel"}, Rows: [][]string{{"svc-1", "evm", "debug"}, {"svc-2", "a much longer service name", "warn"}}},
	{Name: "external-initiator", Headers: []string{"Name", "URL", "AccessKey", "Secret", "OutgoingToken", "OutgoingSecret"}, Rows: [][]string{{"bitcoin", "http://localhost:8888", "accesskey", "secret", "outgoingToken", "outgoingSecret"}}},
	{Name: "pipeline", Headers: []string{"ID", "Created At", "Finished At"}, Rows: [][]string{{"42", "2026-11-12 10:00:00 +0000 UTC", "2026-11-12 10:01:02 +0000 UTC"}}},
	{Name: "headers-only", Headers: []string{"First", "Second"}},
	{Name: "empty-row", Headers: []string{"First", "Second"}, Rows: [][]string{{"", ""}}},
	{Name: "multiple", Headers: []string{"ID", "Name", "Created At"}, Rows: [][]string{{"1", "alpha", "10"}, {"22", "β", "100000"}}},
	{Name: "camel-and-dots", Headers: []string{"LogLevel", "AccessKey", "version_2", "schema.1", "HTTPServer"}},
	{Name: "unicode", Headers: []string{"ID", "ServiceName", "Combined"}, Rows: [][]string{{"é", "東京", "👩‍💻"}, {"é", "Αθήνα", "\x1b[31mred\x1b[0m"}}},
	{Name: "multiline", Headers: []string{"MultiLine", "Other"}, Rows: [][]string{{"first\n  second  \n\n fourth ", "short"}, {"x\ny", "longer value"}}},
	{Name: "inner-whitespace", Headers: []string{"Key", "Value"}, Rows: [][]string{{"key", "  padded\nnext  "}}},
	{Name: "single-column", Headers: []string{"ID"}, Rows: [][]string{{"a"}, {"bc"}}},
	{Name: "irregular", Headers: []string{"A", "B", "C"}, Rows: [][]string{{"one"}, {"one", "two", "three", "ignored"}}},
	{Name: "invalid-utf8", Headers: []string{"Valid", "Value"}, Rows: [][]string{{"x", string([]byte{0xff, 'a'})}}},
	{Name: "crlf", Headers: []string{"A", "B"}, Rows: [][]string{{"first\r\nsecond", "x"}}},
}

func decodeCLITableGolden(t *testing.T, data []byte) map[string][]byte {
	t.Helper()
	records := make(map[string][]byte)
	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		var record cliTableGoldenRecord
		if err := decoder.Decode(&record); err != nil {
			t.Fatalf("decode golden record: %v", err)
		}
		if _, exists := records[record.Name]; exists {
			t.Fatalf("duplicate golden record %q", record.Name)
		}
		records[record.Name] = record.Output
	}
	if len(records) != len(cliTableGoldenCases) {
		t.Fatalf("decoded %d golden records, want %d", len(records), len(cliTableGoldenCases))
	}
	return records
}

func TestCLITableGolden(t *testing.T) {
	suites := []struct {
		name      string
		data      []byte
		eastAsian bool
	}{
		{name: "default", data: cliTableGoldenDefault},
		{name: "east-asian", data: cliTableGoldenEastAsian, eastAsian: true},
	}
	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			records := decodeCLITableGolden(t, suite.data)
			for _, testCase := range cliTableGoldenCases {
				t.Run(testCase.Name, func(t *testing.T) {
					expected, ok := records[testCase.Name]
					if !ok {
						t.Fatalf("missing golden record")
					}
					var output bytes.Buffer
					table := newCLITable(&output, testCase.Headers)
					table.widthCondition = &runewidth.Condition{
						EastAsianWidth:     suite.eastAsian,
						StrictEmojiNeutral: true,
					}
					for _, row := range testCase.Rows {
						if err := table.Append(row); err != nil {
							t.Fatalf("Append returned an error: %v", err)
						}
					}
					if err := table.render(); err != nil {
						t.Fatalf("render returned an error: %v", err)
					}
					if !bytes.Equal(output.Bytes(), expected) {
						t.Errorf("output mismatch\n got: %s\nwant: %s", fmt.Sprintf("%q", output.Bytes()), fmt.Sprintf("%q", expected))
					}
				})
			}
		})
	}
}

type cliTableErrorWriter struct{}

func (cliTableErrorWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestCLITableRenderIgnoresWriterErrors(t *testing.T) {
	table := newCLITable(cliTableErrorWriter{}, []string{"ID", "Value"})
	if err := table.Append([]string{"1", "value"}); err != nil {
		t.Fatalf("Append returned an error: %v", err)
	}
	if err := table.render(); err != nil {
		t.Fatalf("render returned an error: %v", err)
	}
}
