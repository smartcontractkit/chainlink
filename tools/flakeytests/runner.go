package flakeytests

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
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
	Report(ctx context.Context, r *Report) error
}

type parseFn func(readers ...io.Reader) (*Report, error)

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
	cmd := exec.Command(t.command, fmt.Sprintf(".%s", replacedPkg)) //#nosec
	cmd.Env = os.Environ()

	if len(tests) > 0 {
		testFilter := strings.Join(tests, "|")
		cmd.Env = append(cmd.Env, fmt.Sprintf("TEST_FLAGS=-run %s", testFilter))
	}

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

func parseOutput(readers ...io.Reader) (*Report, error) {
	report := NewReport()
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

			switch e.Action {
			case "fail":
				// Fail logs come in two forms:
				// - with e.Package && e.Test, in which case it indicates a test failure.
				// - with e.Package only, which indicates that the package test has failed,
				// or possible that there has been a panic in an out-of-process goroutine running
				// as part of the tests.
				//
				// We can ignore the last case because a package failure will be accounted elsewhere, either
				// in the form of a failing test entry, or in the form of a panic output log, covered below.
				if e.Test == "" {
					continue
				}

				report.IncTest(e.Package, e.Test)
			case "output":
				if panicRe.MatchString(e.Output) {
					// Similar to the above, a panic can come in two forms:
					// - attached to a test (i.e. with e.Test != ""), in which case
					// we'll treat it like a failing test.
					// - package-scoped, in which case we'll treat it as a package panic.
					if e.Test != "" {
						report.IncTest(e.Package, e.Test)
					} else {
						report.IncPackagePanic(e.Package)
					}
				}
			}
		}

		if err := s.Err(); err != nil {
			return nil, err
		}
	}

	return report, nil
}

type exitCoder interface {
	ExitCode() int
}

type Report struct {
	tests         map[string]map[string]int
	packagePanics map[string]int
}

func NewReport() *Report {
	return &Report{
		tests:         map[string]map[string]int{},
		packagePanics: map[string]int{},
	}
}

func (r *Report) HasFlakes() bool {
	return len(r.tests) > 0 || len(r.packagePanics) > 0
}

func (r *Report) SetTest(pkg, test string, val int) {
	if r.tests[pkg] == nil {
		r.tests[pkg] = map[string]int{}
	}
	r.tests[pkg][test] = val
}

func (r *Report) IncTest(pkg string, test string) {
	if r.tests[pkg] == nil {
		r.tests[pkg] = map[string]int{}
	}
	r.tests[pkg][test]++
}

func (r *Report) IncPackagePanic(pkg string) {
	r.packagePanics[pkg]++
}

func (r *Runner) runTest(pkg string, tests []string) (*Report, error) {
	var out bytes.Buffer
	err := r.testCommand.test(pkg, tests, &out)
	if err != nil {
		log.Printf("Test command errored: %s\n", err)
		// There was an error because the command failed with a non-zero
		// exit code. This could just mean that the test failed again, so let's
		// keep going.
		var exErr exitCoder
		if errors.As(err, &exErr) && exErr.ExitCode() > 0 {
			return r.parse(&out)
		}
		return nil, err
	}

	return r.parse(&out)
}

func (r *Runner) runTests(rep *Report) (*Report, error) {
	report := NewReport()

	// We need to deal with two types of flakes here:
	// - flakes where we know the test that failed; in this case, we just rerun the failing test in question
	// - flakes where we don't know what test failed. These are flakes where a panic occurred in an out-of-process goroutine,
	// thus failing the package as a whole. For these, we'll rerun the whole package again.
	for pkg, tests := range rep.tests {
		ts := []string{}
		for test := range tests {
			ts = append(ts, test)
		}

		for i := 0; i < r.numReruns; i++ {
			log.Printf("[FLAKEY_TEST] Executing test command with parameters: pkg=%s, tests=%+v, numReruns=%d currentRun=%d\n", pkg, ts, r.numReruns, i)
			pr, err := r.runTest(pkg, ts)
			if err != nil {
				return report, err
			}

			for t := range tests {
				failures := pr.tests[pkg][t]
				if failures == 0 {
					report.SetTest(pkg, t, 1)
				}
			}
		}
	}

	for pkg := range rep.packagePanics {
		for i := 0; i < r.numReruns; i++ {
			log.Printf("[PACKAGE_PANIC]: Executing test command with parameters: pkg=%s, numReruns=%d currentRun=%d\n", pkg, r.numReruns, i)
			pr, err := r.runTest(pkg, []string{})
			if err != nil {
				return report, err
			}

			if pr.packagePanics[pkg] == 0 {
				report.IncPackagePanic(pkg)
			}
		}
	}

	return report, nil
}

func isSubtest(tn string) bool {
	return strings.Contains(tn, "/")
}

func isSubtestOf(st, mt string) bool {
	return isSubtest(st) && strings.Contains(st, mt)
}

func dedupeEntries(report *Report) (*Report, error) {
	out := NewReport()
	out.packagePanics = report.packagePanics
	for pkg, tests := range report.tests {
		// Sort the test names
		testNames := make([]string, 0, len(tests))
		for t := range tests {
			testNames = append(testNames, t)
		}

		sort.Strings(testNames)

		for i, tn := range testNames {
			// Is this the last element? If it is, then add it to the deduped set.
			// This is because a) it's a main test, in which case we add it because
			// it has no subtests following it, or b) it's a subtest, which we always add.
			if i == len(testNames)-1 {
				out.SetTest(pkg, tn, report.tests[pkg][tn])
				continue
			}

			// Next, let's compare the current item to the next one in the alphabetical order.
			// In all cases we want to add the current item, UNLESS the current item is a main test,
			// and the following one is a subtest of the current item.
			nextItem := testNames[i+1]
			if !isSubtest(tn) && isSubtestOf(nextItem, tn) {
				continue
			}

			out.SetTest(pkg, tn, report.tests[pkg][tn])
		}
	}

	return out, nil
}

func (r *Runner) Run(ctx context.Context) error {
	parseReport, err := r.parse(r.readers...)
	if err != nil {
		return err
	}

	report, err := r.runTests(parseReport)
	if err != nil {
		return err
	}

	if report.HasFlakes() {
		log.Printf("ERROR: Suspected flakes found: %+v\n", report)
	} else {
		log.Print("SUCCESS: No suspected flakes detected")
	}

	// Before reporting the errors, let's dedupe some entries:
	// In actuality, a failing subtest will produce two failing test entries,
	// namely one for the test as a whole, and one for the subtest.
	// This leads to inaccurate metrics since a failing subtest is double-counted.
	report, err = dedupeEntries(report)
	if err != nil {
		return err
	}

	return r.reporter.Report(ctx, report)
}
