package flakeytests

import (
	"encoding/json"
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

func GetGithubMetadata(sha string, path string) Context {
	event := map[string]interface{}{}
	if path != "" {
		r, err := os.Open(path)
		if err != nil {
			log.Fatalf("Error reading gh event at path: %s", path)
		}

		d, err := io.ReadAll(r)
		if err != nil {
			log.Fatal("Error reading gh event into string")
		}

		err = json.Unmarshal(d, &event)
		if err != nil {
			log.Fatalf("Error unmarshaling gh event at path: %s", path)
		}
	}

	prURL := ""
	url, err := DigString(event, []string{"pull_request", "_links", "html", "href"})
	if err == nil {
		prURL = url
	}
	ctx := Context{
		CommitSHA:      sha,
		PullRequestURL: prURL,
	}
	return ctx
}
