package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Job struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
	URL   string `json:"html_url"`
}

type Step struct {
	Name       string `json:"name"`
	Conclusion string `json:"conclusion"`
}

type GitHubResponse struct {
	Jobs []Job `json:"jobs"`
}

type ParsedResult struct {
	Conclusion string `json:"conclusion"`
	Cap        string `json:"cap"`
	URL        string `json:"html_url"`
}

type ResultsMap map[string][]ParsedResult

func main() {
	// Define flags
	githubToken := flag.String("githubToken", "", "GitHub token for authentication")
	githubRepo := flag.String("githubRepo", "", "GitHub repository in the format owner/repo")
	workflowRunID := flag.String("workflowRunID", "", "ID of the GitHub Actions workflow run")
	jobNameRegex := flag.String("jobNameRegex", "", "Regex pattern to match job names")
	namedKey := flag.String("namedKey", "", "Optional named key under which results will be stored")
	outputFile := flag.String("outputFile", "", "Optional output file to save results")

	// Parse flags
	flag.Parse()

	// Validate flags
	if *githubToken == "" || *githubRepo == "" || *workflowRunID == "" || *jobNameRegex == "" {
		fmt.Println("Please provide all required flags: --githubToken, --githubRepo, --workflowRunID, --jobNameRegex")
		return
	}

	// Make GitHub API request
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs/%s/jobs?per_page=100", *githubRepo, *workflowRunID)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+*githubToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("GitHub API request failed with status:", resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Parse GitHub API response
	var githubResponse GitHubResponse
	err = json.Unmarshal(body, &githubResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON response:", err)
		return
	}

	// Process GitHub jobs
	var parsedResults []ParsedResult
	re := regexp.MustCompile(*jobNameRegex)
	for _, job := range githubResponse.Jobs {
		if re.MatchString(job.Name) {
			for _, step := range job.Steps {
				if step.Name == "Run Tests" {
					conclusion := ":x:"
					if step.Conclusion == "success" {
						conclusion = ":white_check_mark:"
					}
					cap := fmt.Sprintf("%s", re.FindStringSubmatch(job.Name)[1])
					parsedResults = append(parsedResults, ParsedResult{
						Conclusion: conclusion,
						Cap:        cap,
						URL:        job.URL,
					})
				}
			}
		}
	}

	// Create a map to store results
	results := ResultsMap{}

	// Check if output file exists and load existing data if it does
	if *outputFile != "" {
		if _, err := os.Stat(*outputFile); err == nil {
			existingData, err := ioutil.ReadFile(*outputFile)
			if err == nil {
				json.Unmarshal(existingData, &results)
			}
		}
	}

	// Append results under the named key if provided, otherwise use a default key
	key := "results"
	if *namedKey != "" {
		key = *namedKey
	}
	results[key] = append(results[key], parsedResults...)

	// Convert results to JSON format
	formattedResults, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling formatted results:", err)
		return
	}

	// Save results to file or print to stdout
	if *outputFile != "" {
		err = ioutil.WriteFile(*outputFile, formattedResults, 0644)
		if err != nil {
			fmt.Println("Error writing results to file:", err)
		} else {
			fmt.Println("Results saved to", *outputFile)
		}
	} else {
		fmt.Println(string(formattedResults))
	}
}
