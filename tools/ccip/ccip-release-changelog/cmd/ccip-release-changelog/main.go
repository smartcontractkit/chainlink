// Command ccip-release-changelog generates a CCIP-focused changelog between
// two refs of the core chainlink repo (release branches, tags, or SHAs).
//
// It diffs the CCIP-relevant pins in go.mod and plugins/plugins.public.yaml,
// produces per-repo commit changelogs via the GitHub compare API (or local
// git for the core repo), flags release risks, and optionally posts the
// result into a Slack thread.
//
// Usage:
//
//	ccip-release-changelog --old v2.55.0 --new release/2.56.0 \
//	    [--repo .] [--out changelog.md] [--slack-thread https://...]
//
// Environment:
//
//	GITHUB_TOKEN / GH_TOKEN - for the GitHub compare API (falls back to
//	                          `gh auth token`; all tracked repos are public)
//	SLACK_BOT_TOKEN         - required when --slack-thread is given
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/smartcontractkit/chainlink/v2/tools/ccip/ccip-release-changelog/internal/changelog"
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	oldRef := flag.String("old", "", "old git ref (SHA, tag, or branch) that built the current release image")
	newRef := flag.String("new", "", "new git ref (SHA, tag, or branch) for the new release image")
	repoDir := flag.String("repo", ".", "path to the chainlink core repo checkout")
	outPath := flag.String("out", "", "write the full markdown report to this file (default: stdout)")
	slackThread := flag.String("slack-thread", "", "optional Slack thread URL to post the summary and report into")
	flag.Parse()

	if *oldRef == "" || *newRef == "" {
		fmt.Fprintln(os.Stderr, "error: --old and --new are required")
		flag.Usage()
		return 2
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx, *repoDir, *oldRef, *newRef, *outPath, *slackThread); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}

func run(ctx context.Context, repoDir, oldRef, newRef, outPath, slackThreadURL string) error {
	report, err := changelog.Generate(ctx, repoDir, oldRef, newRef)
	if err != nil {
		return err
	}

	markdown := changelog.RenderMarkdown(report)

	if outPath != "" {
		if err := os.WriteFile(outPath, []byte(markdown), 0o600); err != nil {
			return fmt.Errorf("writing %s: %w", outPath, err)
		}
		fmt.Fprintf(os.Stderr, "report written to %s\n", outPath)
	} else if slackThreadURL == "" {
		fmt.Print(markdown)
	}

	if slackThreadURL != "" {
		token := os.Getenv("SLACK_BOT_TOKEN")
		if token == "" {
			return errors.New("--slack-thread requires SLACK_BOT_TOKEN in the environment")
		}
		thread, err := changelog.ParseSlackThreadURL(slackThreadURL)
		if err != nil {
			return err
		}
		summary := changelog.RenderSlackSummary(report)
		filename := fmt.Sprintf("ccip-release-changelog-%s-%s.md",
			changelog.SanitizeForFilename(oldRef), changelog.SanitizeForFilename(newRef))
		title := fmt.Sprintf("CCIP Release Changelog %s → %s", oldRef, newRef)
		if err := changelog.PostToSlack(ctx, token, thread, summary, filename, title, markdown); err != nil {
			return fmt.Errorf("posting to Slack: %w", err)
		}
		fmt.Fprintf(os.Stderr, "posted summary and full report to Slack thread %s#%s\n", thread.Channel, thread.ThreadTS)
	}
	return nil
}
