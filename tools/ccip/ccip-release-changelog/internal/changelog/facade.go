package changelog

import (
	"context"
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

// PostToSlack posts the summary message into a thread and uploads the full
// markdown report as a file in the same thread.
func PostToSlack(ctx context.Context, token string, thread SlackThread, summary, filename, title, markdown string) error {
	c := newSlackClient(token)
	if err := c.PostMessage(ctx, thread, summary); err != nil {
		return err
	}
	return c.UploadFile(ctx, thread, filename, title, []byte(markdown), "")
}

var filenameUnsafe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// SanitizeForFilename makes a git ref safe for use in a filename.
func SanitizeForFilename(ref string) string {
	return strings.Trim(filenameUnsafe.ReplaceAllString(ref, "-"), "-")
}
