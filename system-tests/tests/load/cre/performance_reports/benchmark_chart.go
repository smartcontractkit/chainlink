package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type BenchspyReport struct {
	TestName           string    `json:"test_name"`
	CommitOrTag        string    `json:"commit_or_tag"`
	TestStartTimestamp time.Time `json:"test_start_timestamp"`
	TestEndTimestamp   time.Time `json:"test_end_timestamp"`
	GeneratorConfigs   struct {
		Generator struct {
			GeneratorName string `json:"generator_name"`
			LoadType      string `json:"load_type"`
			Schedule      []struct {
				From      int       `json:"from"`
				Duration  int64     `json:"duration"`
				Type      string    `json:"type"`
				TimeStart time.Time `json:"time_start"`
				TimeEnd   time.Time `json:"time_end"`
			} `json:"schedule"`
			RateLimitUnitDuration int64 `json:"rate_limit_unit_duration"`
			CallTimeout           int64 `json:"call_timeout"`
		} `json:"Generator"`
	} `json:"generator_configs"`
	Directory      string `json:"directory"`
	QueryExecutors []struct {
		Kind            string `json:"kind"`
		GeneratorConfig struct {
			GeneratorName string `json:"generator_name"`
			LoadType      string `json:"load_type"`
			Schedule      []struct {
				From      int       `json:"from"`
				Duration  int64     `json:"duration"`
				Type      string    `json:"type"`
				TimeStart time.Time `json:"time_start"`
				TimeEnd   time.Time `json:"time_end"`
			} `json:"schedule"`
			RateLimitUnitDuration int64 `json:"rate_limit_unit_duration"`
			CallTimeout           int64 `json:"call_timeout"`
		} `json:"generator_config"`
		Queries      []string `json:"queries"`
		QueryResults struct {
			Nine5ThPercentileLatency float64 `json:"95th_percentile_latency"`
			ErrorRate                int     `json:"error_rate"`
			MaxLatency               float64 `json:"max_latency"`
			MedianLatency            float64 `json:"median_latency"`
		} `json:"query_results"`
	} `json:"query_executors"`
}

func bar(value, max, width float64, symbol string) string {
	filled := value * width / max
	return strings.Repeat(symbol, int(filled)) + strings.Repeat(" ", int(width-filled))
}

func printP95Chart(runs []BenchspyReport) string {
	// Fetch max
	maxValue := 0.0
	for _, run := range runs {
		if run.QueryExecutors[0].QueryResults.Nine5ThPercentileLatency > maxValue {
			maxValue = run.QueryExecutors[0].QueryResults.Nine5ThPercentileLatency
		}
	}

	// Add padding to max for better visualization
	maxValue = maxValue * 1.2

	var b strings.Builder
	b.WriteString("\n## p95\n\n")
	b.WriteString("| Date       | Value	|\n")
	b.WriteString("| --- | --- |\n")
	for _, run := range runs {
		p95LatencyInSeconds := run.QueryExecutors[0].QueryResults.Nine5ThPercentileLatency / 1000.0
		b.WriteString(fmt.Sprintf("| %s | %-15s %.3fs |\n",
			run.TestStartTimestamp.Format("2006-01-02"),
			bar(run.QueryExecutors[0].QueryResults.Nine5ThPercentileLatency, maxValue, 15, "█"), p95LatencyInSeconds))
	}
	return b.String()
}

func printP50Chart(runs []BenchspyReport) string {
	// Fetch max
	maxValue := 0.0
	for _, run := range runs {
		if run.QueryExecutors[0].QueryResults.MedianLatency > maxValue {
			maxValue = run.QueryExecutors[0].QueryResults.MedianLatency
		}
	}

	// Add padding to max for better visualization
	maxValue = maxValue * 1.2

	var b strings.Builder
	b.WriteString("\n## p95\n\n")
	b.WriteString("| Date       | Value |\n")
	b.WriteString("| --- | --- |\n")
	for _, run := range runs {
		p50LatencyInSeconds := run.QueryExecutors[0].QueryResults.MedianLatency / 1000.0
		b.WriteString(fmt.Sprintf("| %s | %-15s %.3fs |\n",
			run.TestStartTimestamp.Format("2006-01-02"),
			bar(run.QueryExecutors[0].QueryResults.MedianLatency, maxValue, 15, "█"), p50LatencyInSeconds))
	}
	return b.String()
}
func readBenchmarkFiles(folder string) ([]BenchspyReport, error) {
	var runs []BenchspyReport

	// Read directory entries
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Process each JSON file
	for _, entry := range entries {
		// Skip directories and non-JSON files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// Read file content
		data, err := os.ReadFile(filepath.Join(folder, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", entry.Name(), err)
		}

		// Unmarshal JSON data
		var run BenchspyReport
		if err := json.Unmarshal(data, &run); err != nil {
			return nil, fmt.Errorf("failed to parse JSON from %s: %w", entry.Name(), err)
		}

		runs = append(runs, run)
	}

	// Sort by date for consistent output
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].TestStartTimestamp.Before(runs[j].TestStartTimestamp)
	})

	return runs, nil
}

func main() {
	fetchMetrics()
	folder := "/Users/ionita/chainlink/system-tests/tests/load/cre/performance_reports" // Folder containing JSON files
	runs, err := readBenchmarkFiles(folder)
	if err != nil {
		fmt.Println("Error reading benchmark files:", err)
		return
	}

	content := "# Benchmark Report\n\n"
	content += printP95Chart(runs)
	content += printP50Chart(runs)

	// Write to .md file
	if err := ioutil.WriteFile("benchmark_report.md", []byte(content), 0644); err != nil {
		fmt.Println("Error writing markdown file:", err)
	}
}

const prometheusURL = "http://localhost:9090/api/v1/query"

/*
names ----> http://localhost:9090/api/v1/label/name/values
CPU ------> sum by (name) (rate(container_cpu_usage_seconds_total{name!=""}[10m]) * 100)
RAM Spike------> sum by (name) (max_over_time(container_memory_working_set_bytes{name!=""}[10m])/1048576)
RAM ------> sum by (name) (rate(container_memory_working_set_bytes{name!=""}[10m])/1048576)
DISK -----> sum by (name) (container_fs_io_time_seconds_total{name!=""})
Network ----> sum by (name) (container_network_receive_bytes_total{name!=""} + container_network_transmit_bytes_total{name!=""})/1048576
*/
var queries = map[string]string{
	"cpu_percent": `sum by (name) (rate(container_cpu_usage_seconds_total{name!=""}[10m]) * 100)`,
	"mem_peak":    `sum by (name) (max_over_time(container_memory_working_set_bytes{name!=""}[10m])/1048576)`,
	"mem_avg":     `sum by (name) (avg_over_time(container_memory_working_set_bytes{name!=""}[10m])/1048576)`,
	"disk":        `sum by (name) (container_fs_io_time_seconds_total{name!=""})`,
	"network_tx":  `sum by (name) (container_network_transmit_bytes_total{name!=""})/1048576`,
	"network_rx":  `sum by (name) (container_network_receive_bytes_total{name!=""})/1048576`,
}

type promResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string        `json:"resultType"`
		Result     []promMetrics `json:"result"`
	} `json:"data"`
}

type promMetrics struct {
	Metric map[string]string `json:"metric"`
	Value  [2]interface{}    `json:"value"` // [ timestamp, value-as-string ]
}

func query(q string) ([]promMetrics, error) {
	vals := url.Values{}
	vals.Set("query", q)
	resp, err := http.Get(prometheusURL + "?" + vals.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr promResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, fmt.Errorf("unmarshal %q: %w", body, err)
	}
	if pr.Status != "success" {
		return nil, fmt.Errorf("prometheus status %q", pr.Status)
	}
	return pr.Data.Result, nil
}
func fetchMetrics() {
	metrics := make(map[string]map[string]float64)
	for field, q := range queries {
		results, err := query(q)
		if err != nil {
			log.Fatalf("query %s failed: %v", field, err)
		}
		for _, m := range results {
			name, ok := m.Metric["name"]
			if !ok {
				continue
			}
			// only save metrics that start with workflow or capabilities
			if !strings.Contains(name, "workflow") && !strings.Contains(name, "capabilities") {
				continue
			}
			sval := fmt.Sprintf("%v", m.Value[1])
			fval, err := strconv.ParseFloat(sval, 64)
			if err != nil {
				log.Printf("parse %q as float: %v", sval, err)
				continue
			}
			if _, ok := metrics[name]; !ok {
				metrics[name] = make(map[string]float64)
			}
			metrics[name][field] = fval
		}
		spew.Dump(metrics)
	}
}
