package flakeytests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	panicRe = regexp.MustCompile(`^panic:`)
)

type Runner struct {
	readers   []io.Reader
	numReruns int
	runTestFn runTestCmd
	parse     parseFn
	reporter  reporter
}

type reporter interface {
	Report(map[string][]string) error
}

type runTestCmd func(pkg string, testNames []string, numReruns int, w io.Writer) error
type parseFn func(readers ...io.Reader) (map[string]map[string]int, error)

func NewRunner(readers []io.Reader, reporter reporter, numReruns int) *Runner {
	return &Runner{
		readers:   readers,
		numReruns: numReruns,
		runTestFn: runGoTest,
		parse:     parseOutput,
		reporter:  reporter,
	}
}

func runGoTest(pkg string, tests []string, numReruns int, w io.Writer) error {
	pkg = strings.Replace(pkg, "github.com/smartcontractkit/chainlink/v2", "", -1)
	testFilter := strings.Join(tests, "|")
	cmd := exec.Command("./tools/bin/go_core_tests", fmt.Sprintf(".%s", pkg)) //#nosec
	cmd.Env = append(os.Environ(), fmt.Sprintf("TEST_FLAGS=-count %d -run %s", numReruns, testFilter))
	cmd.Stdout = io.MultiWriter(os.Stdout, w)
	cmd.Stderr = io.MultiWriter(os.Stderr, w)
	return cmd.Run()
}

type TestEvent struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func newEvent(b []byte) (*TestEvent, error) {
	e := &TestEvent{}
	err := json.Unmarshal(b, e)
	return e, err
}

func parseOutput(readers ...io.Reader) (map[string]map[string]int, error) {
	tests := map[string]map[string]int{}
	for _, r := range readers {
		s := bufio.NewScanner(r)
		for s.Scan() {
			t := s.Bytes()
			if len(t) == 0 {
				continue
			}

			// Skip the line if doesn't start with a "{" --
			// this mean it isn't JSON output.
			if !strings.HasPrefix(string(t), "{") {
				continue
			}

			e, err := newEvent(t)
			if err != nil {
				return nil, err
			}

			// We're only interested in test failures, for which
			// both Package and Test would be present.
			if e.Package == "" || e.Test == "" {
				continue
			}

			switch e.Action {
			case "fail":
				if tests[e.Package] == nil {
					tests[e.Package] = map[string]int{}
				}
				tests[e.Package][e.Test]++
			case "output":
				if panicRe.MatchString(e.Output) {
					if tests[e.Package] == nil {
						tests[e.Package] = map[string]int{}
					}
					tests[e.Package][e.Test]++
				}
			}
		}

		if err := s.Err(); err != nil {
			return nil, err
		}
	}
	return tests, nil
}

type exitCoder interface {
	ExitCode() int
}

func (r *Runner) runTests(failedTests map[string]map[string]int) (io.Reader, error) {
	var out bytes.Buffer
	for pkg, tests := range failedTests {
		ts := []string{}
		for test := range tests {
			ts = append(ts, test)
		}

		log.Printf("Executing test command with parameters: pkg=%s, tests=%+v, numReruns=%d\n", pkg, ts, r.numReruns)
		err := r.runTestFn(pkg, ts, r.numReruns, &out)
		if err != nil {
			log.Printf("Test command errored: %s\n", err)
			// There was an error because the command failed with a non-zero
			// exit code. This could just mean that the test failed again, so let's
			// keep going.
			var exErr exitCoder
			if errors.As(err, &exErr) && exErr.ExitCode() > 0 {
				continue
			}
			return &out, err
		}
	}

	return &out, nil
}

func (r *Runner) Run() error {
	failedTests, err := r.parse(r.readers...)
	if err != nil {
		return err
	}

	output, err := r.runTests(failedTests)
	if err != nil {
		return err
	}

	failedReruns, err := r.parse(output)
	if err != nil {
		return err
	}

	suspectedFlakes := map[string][]string{}
	// A test is flakey if it appeared in the list of original flakey tests
	// and doesn't appear in the reruns, or if it hasn't failed each additional
	// run, i.e. if it hasn't twice after being re-run.
	for pkg, t := range failedTests {
		for test := range t {
			if failedReruns[pkg][test] != r.numReruns {
				suspectedFlakes[pkg] = append(suspectedFlakes[pkg], test)
			}
		}
	}

	if len(suspectedFlakes) > 0 {
		log.Printf("ERROR: Suspected flakes found: %+v\n", suspectedFlakes)
	} else {
		log.Print("SUCCESS: No suspected flakes detected")
	}

	return r.reporter.Report(suspectedFlakes)
}
