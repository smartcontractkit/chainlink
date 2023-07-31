package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/flakeytests"
)

const numReruns = 2

func main() {
	grafanaHost := flag.String("grafana_host", "", "grafana host URL")
	grafanaAuth := flag.String("grafana_auth", "", "grafana basic auth for Loki API")
	command := flag.String("command", "", "test command being rerun; used to tag metrics")
	ghSHA := flag.String("gh_sha", "", "commit sha for which we're rerunning tests")
	ghEventPath := flag.String("gh_event_path", "", "path to associated gh event")
	flag.Parse()

	if *grafanaHost == "" {
		log.Fatal("Error re-running flakey tests: `grafana_host` is required")
	}

	if *grafanaAuth == "" {
		log.Fatal("Error re-running flakey tests: `grafana_auth` is required")
	}

	if *command == "" {
		log.Fatal("Error re-running flakey tests: `command` is required")
	}

	args := flag.Args()

	log.Printf("Parsing output at: %v", strings.Join(args, ", "))
	readers := []io.Reader{}
	for _, f := range args {
		r, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}

		readers = append(readers, r)
	}

	ctx := flakeytests.GetGithubMetadata(*ghSHA, *ghEventPath)
	rep := flakeytests.NewLokiReporter(*grafanaHost, *grafanaAuth, *command, ctx)
	r := flakeytests.NewRunner(readers, rep, numReruns)
	err := r.Run()
	if err != nil {
		log.Fatalf("Error re-running flakey tests: %s", err)
	}
}
