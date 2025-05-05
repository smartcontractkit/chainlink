package main

import (
	"errors"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/por"
)

type corpusEntry = struct {
	Parent     string
	Path       string
	Data       []byte
	Values     []any
	Generation int
	IsSeed     bool
}

var errMain = errors.New("testing: unexpected use of func Main")

type matchStringOnly func(pat, str string) (bool, error)

func (f matchStringOnly) MatchString(pat, str string) (bool, error)   { return f(pat, str) }
func (f matchStringOnly) StartCPUProfile(w io.Writer) error           { return errMain }
func (f matchStringOnly) StopCPUProfile()                             {}
func (f matchStringOnly) WriteProfileTo(string, io.Writer, int) error { return errMain }
func (f matchStringOnly) ImportPath() string                          { return "" }
func (f matchStringOnly) StartTestLog(io.Writer)                      {}
func (f matchStringOnly) StopTestLog() error                          { return errMain }
func (f matchStringOnly) SetPanicOnExit0(bool)                        {}
func (f matchStringOnly) CoordinateFuzzing(time.Duration, int64, time.Duration, int64, int, []corpusEntry, []reflect.Type, string, string) error {
	return errMain
}
func (f matchStringOnly) RunFuzzWorker(func(corpusEntry) error) error { return errMain }
func (f matchStringOnly) ReadCorpus(string, []reflect.Type) ([]corpusEntry, error) {
	return nil, errMain
}
func (f matchStringOnly) CheckCorpus([]any, []reflect.Type) error { return nil }
func (f matchStringOnly) ResetCoverage()                          {}
func (f matchStringOnly) SnapshotCoverage()                       {}

func (f matchStringOnly) InitRuntimeCoverage() (mode string, tearDown func(string, string) (string, error), snapcov func() float64) {
	return
}

func main() {

	if err := os.Chdir("/Users/matthewpendrey/Projects/chainlink/core/capabilities/integration_tests/por"); err != nil {
		panic(err)
	}

	err := os.Setenv("WASM_FILE_NAME", "fetchtrueusd.wasm")
	if err != nil {
		panic(err)
	}

	tests := []testing.InternalTest{
		{"RunWorkflow", por.RunWorkflow},
	}

	// Create a matcher (can be nil if not filtering tests)
	matchers := []testing.InternalBenchmark{}
	fuzzTargets := []testing.InternalFuzzTarget{}
	examples := []testing.InternalExample{}

	matchString := func(pat, str string) (bool, error) {
		return true, nil
	}

	os.Exit(testing.MainStart(matchStringOnly(matchString), tests, matchers, fuzzTargets, examples).Run())

}
