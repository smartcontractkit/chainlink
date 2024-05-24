package flakeytests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigString(t *testing.T) {
	in := map[string]interface{}{
		"pull_request": map[string]interface{}{
			"url": "some-url",
		},
	}
	out, err := DigString(in, []string{"pull_request", "url"})
	require.NoError(t, err)
	assert.Equal(t, "some-url", out)
}

var prEventTemplate = `
{
  "pull_request": {
    "head": {
      "sha": "%s"
    },
    "_links": {
      "html": {
        "href": "%s"
      }
    }
  }
}
`

func TestGetGithubMetadata(t *testing.T) {
	repo, eventName, sha, event, runID, runAttempt := "chainlink", "merge_group", "a-sha", `{}`, "1234", "1"
	expectedRunURL := fmt.Sprintf("github.com/%s/actions/runs/%s/attempts/%s", repo, runID, runAttempt)
	ctx := getGithubMetadata(repo, eventName, sha, strings.NewReader(event), runID, runAttempt)
	assert.Equal(t, Context{Repository: repo, CommitSHA: sha, Type: eventName, RunURL: expectedRunURL}, ctx)

	anotherSha, eventName, url := "another-sha", "pull_request", "a-url"
	event = fmt.Sprintf(prEventTemplate, anotherSha, url)
	sha = "302eb05d592132309b264e316f443f1ceb81b6c3"
	ctx = getGithubMetadata(repo, eventName, sha, strings.NewReader(event), runID, runAttempt)
	assert.Equal(t, Context{Repository: repo, CommitSHA: anotherSha, Type: eventName, PullRequestURL: url, RunURL: expectedRunURL}, ctx)
}
