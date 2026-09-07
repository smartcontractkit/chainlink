// Command release-changelog generates a CCIP-focused changelog between
// two refs of the core chainlink repo (release branches, tags, or SHAs).
//
// It diffs the CCIP-relevant pins in go.mod and plugins/plugins.public.yaml,
// produces per-repo commit changelogs via the GitHub compare API (or local
// git for the core repo), flags release risks, and optionally posts the
// result into a Slack thread.
//
// The tool is a thin wrapper around internal/engine with the CCIP product
// definition (internal/products/ccip) wired in — the engine itself is
// product-agnostic.
//
// Usage:
//
//	release-changelog --old v2.55.0 --new release/2.56.0 \
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
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/release/release-changelog/internal/engine"
	"github.com/smartcontractkit/chainlink/v2/tools/release/release-changelog/internal/products/ccip"
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	oldRef := flag.String("old", "", "old git ref (SHA, tag, branch, or image tag) that built the current release image")
	newRef := flag.String("new", "", "new git ref (SHA, tag, branch, or image tag) for the new release image")
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
	report, err := engine.Generate(ctx, repoDir, ccip.Product, oldRef, newRef)
	if err != nil {
		return err
	}

	markdown := engine.RenderMarkdown(report)

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
		thread, err := engine.ParseSlackThreadURL(slackThreadURL)
		if err != nil {
			return err
		}
		summary := engine.RenderSlackSummary(report)
		productSlug := strings.ToLower(engine.SanitizeForFilename(ccip.Product.DisplayName))
		filename := fmt.Sprintf("%s-release-changelog-%s-%s.md",
			productSlug, engine.SanitizeForFilename(oldRef), engine.SanitizeForFilename(newRef))
		title := fmt.Sprintf("%s Release Changelog %s → %s", ccip.Product.DisplayName, oldRef, newRef)

		// The summary is the audit payload: if it can't be delivered, fail.
		if err := engine.PostSummary(ctx, token, thread, summary); err != nil {
			return fmt.Errorf("posting summary to Slack: %w", err)
		}

		// The file upload needs the files:write scope, which the bot token
		// may not have. Degrade gracefully: point the thread at the CI
		// artifact instead, and don't fail the run — the report content is
		// already on stdout / in --out.
		if err := engine.UploadReport(ctx, token, thread, filename, title, markdown); err != nil {
			fallback := fmt.Sprintf("⚠️ Full report upload failed (%v).", err)
			if url := engine.ActionsRunURL(); url != "" {
				fallback += " The markdown report is attached to this CI run as an artifact: " + url
			}
			if ferr := engine.PostSummary(ctx, token, thread, fallback); ferr != nil {
				return fmt.Errorf("uploading report: %w; posting fallback message: %w", err, ferr)
			}
			fmt.Fprintf(os.Stderr, "warning: report upload failed (%v); summary posted to thread %s#%s\n", err, thread.Channel, thread.ThreadTS)
			return nil
		}
		fmt.Fprintf(os.Stderr, "posted summary and full report to Slack thread %s#%s\n", thread.Channel, thread.ThreadTS)
	}
	return nil
}
