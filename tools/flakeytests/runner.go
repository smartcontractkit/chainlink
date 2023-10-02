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
	readers     []io.Reader
	testCommand tester
	numReruns   int
	parse       parseFn
	reporter    reporter
}

type tester interface {
	test(pkg string, tests []string, w io.Writer) error
}

type reporter interface {
	Report(map[string]map[string]struct{}) error
}

type parseFn func(readers ...io.Reader) (map[string]map[string]int, error)

func NewRunner(readers []io.Reader, reporter reporter, numReruns int) *Runner {
	tc := &testCommand{
		repo:      "github.com/smartcontractkit/chainlink/v2",
		command:   "./tools/bin/go_core_tests",
		overrides: func(*exec.Cmd) {},
	}
	return &Runner{
		readers:     readers,
		numReruns:   numReruns,
		testCommand: tc,
		parse:       parseOutput,
		reporter:    reporter,
	}
}

type testCommand struct {
	command   string
	repo      string
	overrides func(*exec.Cmd)
}

func (t *testCommand) test(pkg string, tests []string, w io.Writer) error {
	replacedPkg := strings.Replace(pkg, t.repo, "", -1)
	testFilter := strings.Join(tests, "|")
	cmd := exec.Command(t.command, fmt.Sprintf(".%s", replacedPkg)) //#nosec
	cmd.Env = append(os.Environ(), fmt.Sprintf("TEST_FLAGS=-run %s", testFilter))
	cmd.Stdout = io.MultiWriter(os.Stdout, w)
	cmd.Stderr = io.MultiWriter(os.Stderr, w)
	t.overrides(cmd)
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

func (r *Runner) runTests(failedTests map[string]map[string]int) (map[string]map[string]struct{}, error) {
	suspectedFlakes := map[string]map[string]struct{}{}

	for pkg, tests := range failedTests {
		ts := []string{}
		for test := range tests {
			ts = append(ts, test)
		}

		log.Printf("Executing test command with parameters: pkg=%s, tests=%+v, numReruns=%d\n", pkg, ts, r.numReruns)
		for i := 0; i < r.numReruns; i++ {
			var out bytes.Buffer

			err := r.testCommand.test(pkg, ts, &out)
			if err != nil {
				log.Printf("Test command errored: %s\n", err)
				// There was an error because the command failed with a non-zero
				// exit code. This could just mean that the test failed again, so let's
				// keep going.
				var exErr exitCoder
				if errors.As(err, &exErr) && exErr.ExitCode() > 0 {
					continue
				}
				return suspectedFlakes, err
			}

			fr, err := r.parse(&out)
			if err != nil {
				return nil, err
			}

			for t := range tests {
				failures := fr[pkg][t]
				if failures == 0 {
					if suspectedFlakes[pkg] == nil {
						suspectedFlakes[pkg] = map[string]struct{}{}
					}
					suspectedFlakes[pkg][t] = struct{}{}
				}
			}
		}
	}

	return suspectedFlakes, nil
}

func (r *Runner) Run() error {
	failedTests, err := r.parse(r.readers...)
	if err != nil {
		return err
	}

	suspectedFlakes, err := r.runTests(failedTests)
	if err != nil {
		return err
	}

	if len(suspectedFlakes) > 0 {
		log.Printf("ERROR: Suspected flakes found: %+v\n", suspectedFlakes)
	} else {
		log.Print("SUCCESS: No suspected flakes detected")
	}

	return r.reporter.Report(suspectedFlakes)
}
