package flakeytests

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

func DigString(mp map[string]interface{}, path []string) (string, error) {
	var val interface{}
	val = mp
	for _, p := range path {
		v, ok := val.(map[string]interface{})[p]
		if !ok {
			return "", errors.New("could not find string")
		}

		val = v
	}

	vs, ok := val.(string)
	if !ok {
		return "", errors.Errorf("could not coerce value to string: %v", val)
	}

	return vs, nil
}

func getGithubMetadata(repo string, eventName string, sha string, e io.Reader, runID string, runAttempt string) Context {
	d, err := io.ReadAll(e)
	if err != nil {
		log.Fatal("Error reading gh event into string")
	}

	event := map[string]interface{}{}
	err = json.Unmarshal(d, &event)
	if err != nil {
		log.Fatalf("Error unmarshaling gh event at path")
	}

	runURL := fmt.Sprintf("github.com/%s/actions/runs/%s/attempts/%s", repo, runID, runAttempt)
	basicCtx := &Context{Repository: repo, CommitSHA: sha, Type: eventName, RunURL: runURL}
	switch eventName {
	case "pull_request":
		prURL := ""
		url, err := DigString(event, []string{"pull_request", "_links", "html", "href"})
		if err == nil {
			prURL = url
		}

		basicCtx.PullRequestURL = prURL

		// For pull request events, the $GITHUB_SHA variable doesn't actually
		// contain the sha for the latest commit, as documented here:
		// https://stackoverflow.com/a/68068674
		var newSha string
		s, err := DigString(event, []string{"pull_request", "head", "sha"})
		if err == nil {
			newSha = s
		}

		basicCtx.CommitSHA = newSha
		return *basicCtx
	default:
		return *basicCtx
	}
}

func GetGithubMetadata(repo string, eventName string, sha string, path string, runID string, runAttempt string) Context {
	event, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error reading gh event at path: %s", path)
	}
	return getGithubMetadata(repo, eventName, sha, event, runID, runAttempt)
}
