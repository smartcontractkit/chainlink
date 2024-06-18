package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	TotalCount int   `json:"total_count"`
	Jobs       []Job `json:"jobs"`
}

type ParsedResult struct {
	Conclusion string `json:"conclusion"`
	Cap        string `json:"cap"`
	URL        string `json:"html_url"`
}

type ResultsMap map[string][]ParsedResult

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func fetchGitHubJobs(apiURL, token string, client HTTPClient) ([]Job, error) {
	var allJobs []Job

	for {
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating HTTP request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making HTTP request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GitHub API request failed with status: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		_ = resp.Body.Close()

		var githubResponse GitHubResponse
		err = json.Unmarshal(body, &githubResponse)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
		}

		allJobs = append(allJobs, githubResponse.Jobs...)

		re := regexp.MustCompile(`&page=(\d+)`)
		matches := re.FindStringSubmatch(apiURL)
		pageNum := 0
		if len(matches) > 1 {
			_, err = fmt.Sscanf(matches[1], "%d", &pageNum)
			if err != nil {
				return nil, fmt.Errorf("error parsing page number: %w", err)
			}
		}

		if (len(githubResponse.Jobs) < 100) || (len(githubResponse.Jobs) == 100 && githubResponse.TotalCount/100 == pageNum) {
			break
		}

		if pageNum == 0 {
			apiURL = apiURL + "&page=2"
		} else {
			apiURL = re.ReplaceAllString(apiURL, fmt.Sprintf("&page=%d", pageNum+1))
		}
	}

	return allJobs, nil
}

func parseResults(jobNameRegex, workflowRunID *string, jobs []Job) ([]ParsedResult, error) {
	var parsedResults []ParsedResult
	re := regexp.MustCompile(*jobNameRegex)
	for _, job := range jobs {
		if re.MatchString(job.Name) {
			for _, step := range job.Steps {
				if step.Name == "Run Tests" {
					conclusion := ":x:"
					if step.Conclusion == "success" {
						conclusion = ":white_check_mark:"
					}
					captureGroup := re.FindStringSubmatch(job.Name)[1]
					parsedResults = append(parsedResults, ParsedResult{
						Conclusion: conclusion,
						Cap:        captureGroup,
						URL:        job.URL,
					})
				}
			}
		}
	}

	if len(parsedResults) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "no results found for '%s' regex in workflow id %s", *jobNameRegex, *workflowRunID)
	}

	return parsedResults, nil
}

func processResults(parsedResults []ParsedResult, namedKey, jobNameRegex, workflowRunID, outputFile *string) error {
	results := ResultsMap{}

	if *outputFile != "" {
		if _, statErr := os.Stat(*outputFile); statErr == nil {
			existingData, readErr := os.ReadFile(*outputFile)
			if readErr == nil {
				jsonErr := json.Unmarshal(existingData, &results)
				if jsonErr != nil {
					return fmt.Errorf("error unmarshalling existing data: %w", jsonErr)
				}
			}
		}
	}

	key := "results"
	if *namedKey != "" {
		key = *namedKey
	}
	results[key] = append(results[key], parsedResults...)

	formattedResults, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling formatted results: %w", err)
	}

	if *outputFile != "" {
		err = os.WriteFile(*outputFile, formattedResults, 0600)
		if err != nil {
			return fmt.Errorf("error writing results to file: %w", err)
		}
		fmt.Printf("Results for '%s' regex and workflow id %s saved to %s\n", *jobNameRegex, *workflowRunID, *outputFile)
	}

	fmt.Println(string(formattedResults))
	return nil
}

func execute(githubToken, githubRepo, workflowRunID, jobNameRegex, namedKey, outputFile *string, client HTTPClient) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs/%s/jobs?per_page=100", *githubRepo, *workflowRunID)

	jobs, err := fetchGitHubJobs(apiURL, *githubToken, client)
	if err != nil {
		return err
	}

	parsedResults, err := parseResults(jobNameRegex, workflowRunID, jobs)
	if err != nil {
		return err
	}

	return processResults(parsedResults, namedKey, jobNameRegex, workflowRunID, outputFile)
}

func main() {
	githubToken := flag.String("githubToken", "", "GitHub token for authentication")
	githubRepo := flag.String("githubRepo", "", "GitHub repository in the format owner/repo")
	workflowRunID := flag.String("workflowRunID", "", "ID of the GitHub Actions workflow run")
	jobNameRegex := flag.String("jobNameRegex", "", "Regex pattern to match job names")
	namedKey := flag.String("namedKey", "", "Optional named key under which results will be stored")
	outputFile := flag.String("outputFile", "", "Optional output file to save results")

	flag.Parse()

	if *githubToken == "" || *githubRepo == "" || *workflowRunID == "" || *jobNameRegex == "" {
		panic(fmt.Errorf("Please provide all required flags: --githubToken, --githubRepo, --workflowRunID, --jobNameRegex"))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	if err := execute(githubToken, githubRepo, workflowRunID, jobNameRegex, namedKey, outputFile, client); err != nil {
		panic(err)
	}
}
