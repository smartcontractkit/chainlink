package flakeytests

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	failedTestRe = regexp.MustCompile(`^--- FAIL: (Test\w+)`)
	logPanicRe   = regexp.MustCompile(`^panic: Log in goroutine after (Test\w+)`)

	failedPkgRe = regexp.MustCompile(`^FAIL\s+github\.com\/smartcontractkit\/chainlink\/v2\/(\S+)`)
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
	testFilter := strings.Join(tests, "|")
	cmd := exec.Command("./tools/bin/go_core_tests", fmt.Sprintf("./%s", pkg)) //#nosec
	cmd.Env = append(os.Environ(), fmt.Sprintf("TEST_FLAGS=-count %d -run %s", numReruns, testFilter))
	cmd.Stdout = io.MultiWriter(os.Stdout, w)
	cmd.Stderr = io.MultiWriter(os.Stderr, w)
	return cmd.Run()
}

func parseOutput(readers ...io.Reader) (map[string]map[string]int, error) {
	testsWithoutPackage := []string{}
	tests := map[string]map[string]int{}
	for _, r := range readers {
		s := bufio.NewScanner(r)
		for s.Scan() {
			t := s.Text()
			switch {
			case failedTestRe.MatchString(t):
				m := failedTestRe.FindStringSubmatch(t)
				testsWithoutPackage = append(testsWithoutPackage, m[1])
			case logPanicRe.MatchString(t):
				m := logPanicRe.FindStringSubmatch(t)
				testsWithoutPackage = append(testsWithoutPackage, m[1])
			case failedPkgRe.MatchString(t):
				p := failedPkgRe.FindStringSubmatch(t)
				for _, t := range testsWithoutPackage {
					if tests[p[1]] == nil {
						tests[p[1]] = map[string]int{}
					}
					tests[p[1]][t]++
				}
				testsWithoutPackage = []string{}
			}
		}

		if err := s.Err(); err != nil {
			return nil, err
		}
	}
	return tests, nil
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
