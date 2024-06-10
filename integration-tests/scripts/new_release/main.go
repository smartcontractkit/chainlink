package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Release represents a GitHub release
type Release struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
}

const repoURL = "https://api.github.com/repos/%s/releases/latest"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <Github org/repository> <days>")
		os.Exit(1)
	}

	release, err := getLatestRelease(fmt.Sprintf(repoURL, os.Args[1]))
	if err != nil {
		fmt.Printf("Error fetching release: %v\n", err)
		os.Exit(1)
	}

	days, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error parsing days: %v\n", err)
		os.Exit(1)
	}

	if isReleaseRecent(release.PublishedAt, days) {
		fmt.Println(release.TagName)
	} else {
		fmt.Println("none")
	}
}

// getLatestRelease fetches the latest release from the given repository URL
func getLatestRelease(url string) (*Release, error) {
	resp, err := http.Get(url)
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
