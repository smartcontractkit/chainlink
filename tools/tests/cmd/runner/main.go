package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/smartcontractkit/chainlink/v2/tools/flakeytests"
)

func test(w io.Writer) error {
	cmd := exec.Command("./tools/bin/go_core_tests", "./...") //#nosec
	cmd.Env = append(os.Environ(), fmt.Sprintf("TEST_FLAGS=-p 1 -parallel 1"))
	cmd.Stdout = io.MultiWriter(os.Stdout, w)
	cmd.Stderr = io.MultiWriter(os.Stderr, w)
	return cmd.Run()
}

type exitCoder interface {
	ExitCode() int
}

func main() {
	count := flag.Int("count", 0, "number of times to run the tests")
	grafanaHost := flag.String("grafana_host", "", "grafana host URL")
	grafanaAuth := flag.String("grafana_auth", "", "grafana basic auth for Loki API")
	flag.Parse()

	if *grafanaHost == "" {
		log.Fatal("Error running tests: `grafana_host` is required")
	}

	if *grafanaAuth == "" {
		log.Fatal("Error running tests: `grafana_auth` is required")
	}

	if *count == 0 {
		*count = 10
	}

	outputs := []io.ReadWriter{}
	testErrors := []error{}
	for i := 0; i < *count; i++ {
		var out bytes.Buffer
		err := test(&out)
		if err != nil {
			log.Printf("Test command errored: %s\n", err)
			var exErr exitCoder
			if errors.As(err, &exErr) && exErr.ExitCode() > 0 {
				testErrors = append(testErrors, err)
			}
		}

		outputs = append(outputs, &out)
	}

	reportErrors := []error{}
	for i, o := range outputs {
		results, err := flakeytests.ParseOutput(o)
		if err != nil {
			log.Fatalf("Error parsing output: %s", err)
		}

		failures := map[string]map[string]struct{}{}
		for pkg, tests := range results {
			if failures[pkg] == nil {
				failures[pkg] = map[string]struct{}{}
			}

			for test := range tests {
				failures[pkg][test] = struct{}{}
			}
		}

		reporter := flakeytests.NewLokiReporter(*grafanaHost, *grafanaAuth, "./tools/bin/go_core_tests", flakeytests.Context{PullRequestURL: fmt.Sprintf("test run %d", i)})
		err = reporter.Report(failures)
		if err != nil {
			reportErrors = append(reportErrors, err)
		}
	}

	if len(reportErrors) > 0 {
		log.Printf("Report errors: %+v\n", reportErrors)
	}
	if len(testErrors) > 0 {
		log.Printf("Test errors: %+v\n", testErrors)
		log.Printf("Test errors count: %d\n", len(testErrors))
	}
}
