package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/flakeytests"
)

const numReruns = 2

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		stop() // restore default exit behavior
		log.Println("Cancelling... interrupt again to exit")
	}()

	grafanaHost := flag.String("grafana_host", "", "grafana host URL")
	grafanaAuth := flag.String("grafana_auth", "", "grafana basic auth for Loki API")
	grafanaOrgID := flag.String("grafana_org_id", "", "grafana org ID")
	command := flag.String("command", "", "test command being rerun; used to tag metrics")
	ghSHA := flag.String("gh_sha", "", "commit sha for which we're rerunning tests")
	ghEventPath := flag.String("gh_event_path", "", "path to associated gh event")
	ghEventName := flag.String("gh_event_name", "", "type of associated gh event")
	ghRepo := flag.String("gh_repo", "", "name of gh repository")
	ghRunID := flag.String("gh_run_id", "", "run id of the gh workflow")
	flag.Parse()

	runAttempt := os.Getenv("GITHUB_RUN_ATTEMPT")
	if runAttempt == "" {
		log.Fatalf("GITHUB_RUN_ATTEMPT is required")
	}

	if *grafanaHost == "" {
		log.Fatal("Error re-running flakey tests: `grafana_host` is required")
	}

	if *grafanaAuth == "" {
		log.Fatal("Error re-running flakey tests: `grafana_auth` is required")
	}

	if *grafanaOrgID == "" {
		log.Fatal("Error re-running flakey tests: `grafana_org_id` is required")
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

	meta := flakeytests.GetGithubMetadata(*ghRepo, *ghEventName, *ghSHA, *ghEventPath, *ghRunID, runAttempt)
	rep := flakeytests.NewLokiReporter(*grafanaHost, *grafanaAuth, *grafanaOrgID, *command, meta)
	r := flakeytests.NewRunner(readers, rep, numReruns)
	err := r.Run(ctx)
	if err != nil {
		log.Fatalf("Error re-running flakey tests: %s", err)
	}
}
