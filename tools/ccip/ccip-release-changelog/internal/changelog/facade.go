package changelog

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Generate runs the full analysis for two refs of the core repo checkout at
// repoDir.
func Generate(ctx context.Context, repoDir, oldRef, newRef string) (*Report, error) {
	g := gitRunner{dir: repoDir}
	gh := newGHClient(ctx)
	return Analyze(ctx, g, gh, oldRef, newRef)
}

// PostSummary posts a message into a Slack thread.
func PostSummary(ctx context.Context, token string, thread SlackThread, text string) error {
	return newSlackClient(token).PostMessage(ctx, thread, text)
}

// UploadReport uploads the full markdown report as a file in a Slack thread.
// Requires the files:write scope on the bot token; callers should treat a
// failure here as non-fatal if the summary was already delivered.
func UploadReport(ctx context.Context, token string, thread SlackThread, filename, title, markdown string) error {
	return newSlackClient(token).UploadFile(ctx, thread, filename, title, []byte(markdown), "")
}

// ActionsRunURL returns the URL of the current GitHub Actions run, or "" when
// not running in CI.
func ActionsRunURL() string {
	server, repo, runID := os.Getenv("GITHUB_SERVER_URL"), os.Getenv("GITHUB_REPOSITORY"), os.Getenv("GITHUB_RUN_ID")
	if server == "" || repo == "" || runID == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/actions/runs/%s", server, repo, runID)
}

var filenameUnsafe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// SanitizeForFilename makes a git ref safe for use in a filename.
func SanitizeForFilename(ref string) string {
	return strings.Trim(filenameUnsafe.ReplaceAllString(ref, "-"), "-")
}
