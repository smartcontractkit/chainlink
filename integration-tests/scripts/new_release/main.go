package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Release represents a GitHub release
type Release struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
}

var repoURL = "https://api.github.com/repos/%s/releases/latest"

var client = &http.Client{}

func main() {
	if err := validateInputs(); err != nil {
		panic(err)
	}

	release, err := getLatestRelease(fmt.Sprintf(repoURL, os.Args[1]), client)
	if err != nil {
		panic(fmt.Errorf("error fetching release: %v\n", err))
	}

	days, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(fmt.Errorf("error parsing days: %v\n", err))
	}

	if isReleaseRecent(release.PublishedAt, days) {
		fmt.Println(release.TagName)
	} else {
		fmt.Println("none")
	}
}

// getLatestRelease fetches the latest release from the given repository URL
func getLatestRelease(url string, client *http.Client) (*Release, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release Release
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// isReleaseRecent checks if the release date is within the given number of days from today
func isReleaseRecent(publishedAt time.Time, days int) bool {
	duration := time.Since(publishedAt)
	return duration.Hours() <= float64(days*24)
}

func validateInputs() error {
	if len(os.Args) < 3 {
		return errors.New("usage: go run main.go <repository_name> <days>")
	}

	if os.Args[1] == "" {
		return errors.New("error: repository_name cannot be empty")
	}

	if len(strings.Split(os.Args[1], "/")) != 2 {
		return errors.New("error: repository_name must be in the format <org>/<repo>")
	}

	if _, err := strconv.Atoi(os.Args[2]); err != nil {
		return fmt.Errorf("error: days must be an integer, but '%s' is not an integer", os.Args[2])
	}

	return nil
}
