package benchspy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// MetricConfig represents configuration for a single metric
type MetricConfig struct {
	Group       string  // Which group this metric belongs to
	UnitLabel   string  // Unit label (e.g. "ms", "MB")
	ScaleFactor float64 // Scale factor for display (e.g. 1000.0 to convert ms to s)
}

type report struct {
	TestName           string    `json:"test_name"`
	CommitOrTag        string    `json:"commit_or_tag"`
	TestStartTimestamp time.Time `json:"test_start_timestamp"`
	TestEndTimestamp   time.Time `json:"test_end_timestamp"`
	QueryExecutors     []struct {
		Kind         string         `json:"kind"`
		QueryResults map[string]any `json:"query_results"`
	} `json:"query_executors"`
}

func bar(value, maxValue, width float64, symbol string) string {
	// Check if max is zero or close to zero
	if maxValue <= 0.000001 { // Use small epsilon for floating point comparison
		return strings.Repeat(" ", int(width)) // Return empty bar if max is zero
	}
	filled := value * width / maxValue
	return strings.Repeat(symbol, int(filled)) + strings.Repeat(" ", int(width-filled))
}

func createBarChart(
	title string,
	dates []time.Time,
	values []float64,
	maxValue float64, // If <= 0, calculated from values with 20% padding
	unitLabel string,
	scaleFactor float64, // e.g., 1000.0 to convert ms to s
) string {
	// Calculate max if not provided
	if maxValue <= 0 {
		for _, value := range values {
			if value > maxValue {
				maxValue = value
			}
		}
		maxValue *= 1.2 // Add 20% padding
	}

	// Ensure scaleFactor is never zero to prevent divide-by-zero
	if scaleFactor <= 0 {
		scaleFactor = 1.0 // Default to 1.0 if invalid
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n## %s\n\n", title))
	b.WriteString("| Date       | Value |\n")
	b.WriteString("| --- | --- |\n")

	for i := 0; i < len(dates) && i < len(values); i++ {
		scaledValue := values[i] / scaleFactor
		b.WriteString(fmt.Sprintf("| %s | %-15s %.3f%s |\n",
			dates[i].Format("2006-01-02"),
			bar(values[i], maxValue, 15, "█"),
			scaledValue,
			unitLabel))
	}

	return b.String()
}

func readBenchmarkFiles(folder string) ([]report, error) {
	// Read directory entries
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	runs := make([]report, 0, len(entries))

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
		var run report
		if err := json.Unmarshal(data, &run); err != nil {
			return nil, fmt.Errorf("failed to parse JSON from %s: %w", entry.Name(), err)
		}

		runs = append(runs, run)
	}

	// Sort by date (newest first)
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].TestStartTimestamp.After(runs[j].TestStartTimestamp)
	})

	return runs, nil
}

func createMultiSeriesBarChart(
	title string,
	dates []time.Time,
	metricNames []string,
	metricValues [][]float64,
	maxValues []float64,
	unitLabels []string,
	scaleFactors []float64,
) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n## %s\n\n", title))
	b.WriteString("| Date |")

	// Add header row with metric names
	for _, name := range metricNames {
		b.WriteString(fmt.Sprintf(" %s |", name))
	}
	b.WriteString("\n")

	// Add separator row
	b.WriteString("| --- |")
	for range metricNames {
		b.WriteString(" --- |")
	}
	b.WriteString("\n")

	// Add data rows
	for i := range dates {
		b.WriteString(fmt.Sprintf("| %s |", dates[i].Format("2006-01-02 15:04:05")))

		for j := range metricNames {
			if i < len(metricValues[j]) {
				value := metricValues[j][i]
				if j < len(scaleFactors) && scaleFactors[j] > 0 {
					value /= scaleFactors[j]
				}
				unitLabel := ""
				if j < len(unitLabels) {
					unitLabel = unitLabels[j]
				}
				b.WriteString(fmt.Sprintf(" %s %.3f%s |",
					bar(metricValues[j][i], maxValues[j], 10, "█"),
					value,
					unitLabel))
			} else {
				b.WriteString(" N/A |")
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

// GenerateMarkdownReport generates a markdown report from benchmark data with configurable metric grouping
func GenerateMarkdownReport(folder string, outputPath string, metricConfigs map[string]MetricConfig) error {
	runs, err := readBenchmarkFiles(folder)
	if err != nil {
		return fmt.Errorf("error reading benchmark files: %w", err)
	}

	runMetrics, metricMaxValues := extractMetrics(runs)
	content := generateReportContent(runs, runMetrics, metricMaxValues, metricConfigs)

	// Write to .md file
	if err := os.WriteFile(outputPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("error writing markdown file: %w", err)
	}

	return nil
}

// extractMetrics processes benchmark runs to extract metrics
func extractMetrics(runs []report) ([]map[string]float64, map[string]float64) {
	runMetrics := make([]map[string]float64, len(runs))
	// Extract metrics from each run
	for i, run := range runs {
		runMetrics[i] = make(map[string]float64)
		// For every query executor
		for _, query := range run.QueryExecutors {
			// For every metric in the query results
			for metricName, metricValue := range query.QueryResults {
				switch query.Kind {
				case "direct":
					processDirectMetric(runMetrics[i], metricName, metricValue)
				case "prometheus":
					processPrometheusMetric(runMetrics[i], metricName, metricValue)
				}
			}
		}
	}

	// Get max of every runMetric
	metricMaxValues := findMaxMetricValues(runMetrics)

	return runMetrics, metricMaxValues
}

// processDirectMetric processes a direct query metric
func processDirectMetric(metrics map[string]float64, metricName string, metricValue any) {
	if strValue, ok := metricValue.(string); ok {
		floatValue, err := strconv.ParseFloat(strValue, 64)
		if err == nil {
			metrics[metricName] = floatValue
		} else {
			fmt.Printf("Metrics %s: failed to convert string value: %v\n", metricName, err)
		}
	} else if floatValue, ok := metricValue.(float64); ok {
		metrics[metricName] = floatValue
	} else {
		fmt.Printf("Metrics %s: unexpected value type: %T\n", metricName, metricValue)
	}
}

// processPrometheusMetric processes a prometheus query metric
func processPrometheusMetric(metrics map[string]float64, metricName string, metricValue any) {
	if metricMap, ok := metricValue.(map[string]any); ok {
		if valueSlice, ok := metricMap["value"].([]any); ok && len(valueSlice) > 0 {
			if valueMap, ok := valueSlice[0].(map[string]any); ok {
				if innerValue, ok := valueMap["value"].([]any); ok && len(innerValue) > 1 {
					if strValue, ok := innerValue[1].(string); ok {
						floatValue, err := strconv.ParseFloat(strValue, 64)
						if err == nil {
							metrics[metricName] = floatValue
						} else {
							fmt.Printf("Metrics %s: failed to convert string value: %v\n", metricName, err)
						}
					} else if floatValue, ok := innerValue[1].(float64); ok {
						metrics[metricName] = floatValue
					} else {
						fmt.Printf("Metrics %s: unexpected value type: %T\n", metricName, innerValue[1])
					}
				}
			}
		}
	}
}

// findMaxMetricValues finds the maximum value for each metric across all runs
func findMaxMetricValues(runMetrics []map[string]float64) map[string]float64 {
	metricMaxValues := make(map[string]float64)
	if len(runMetrics) == 0 {
		return metricMaxValues
	}

	for metricName := range runMetrics[0] {
		maxValue := 0.0
		for i := range runMetrics {
			if value, exists := runMetrics[i][metricName]; exists && value > maxValue {
				maxValue = value
			}
		}
		metricMaxValues[metricName] = maxValue
	}

	return metricMaxValues
}

func generateReportContent(runs []report, runMetrics []map[string]float64,
	metricMaxValues map[string]float64, metricConfigs map[string]MetricConfig) string {
	content := "# Benchmark Report\n\n"

	// Generate actual chart content after collecting metrics
	dates := make([]time.Time, len(runs))
	for i, run := range runs {
		dates[i] = run.TestStartTimestamp
	}

	content += generateGroupedCharts(dates, runs, runMetrics, metricMaxValues, metricConfigs)
	content += generateUngroupedCharts(dates, runs, runMetrics, metricMaxValues, metricConfigs)

	return content
}

// generateGroupedCharts generates charts for grouped metrics
func generateGroupedCharts(dates []time.Time, runs []report, runMetrics []map[string]float64, metricMaxValues map[string]float64, metricConfigs map[string]MetricConfig) string {
	if len(runMetrics) == 0 {
		return ""
	}

	var content string

	// Group metrics by their configured group
	groupedMetrics := make(map[string][]string)
	for metricName, config := range metricConfigs {
		if config.Group != "" {
			groupedMetrics[config.Group] = append(groupedMetrics[config.Group], metricName)
		}
	}

	// Prepare for chart generation
	ungroupedMetrics := make(map[string]bool)
	for metricName := range runMetrics[0] {
		ungroupedMetrics[metricName] = true
	}

	// Process each defined group
	for groupName, groupMetrics := range groupedMetrics {
		content += processMetricGroup(groupName, groupMetrics, dates, runs,
			runMetrics, metricMaxValues, ungroupedMetrics, metricConfigs)
	}

	return content
}

// processMetricGroup processes a single metric group to generate a chart
func processMetricGroup(groupName string, groupMetrics []string, dates []time.Time, runs []report,
	runMetrics []map[string]float64, metricMaxValues map[string]float64,
	ungroupedMetrics map[string]bool, metricConfigs map[string]MetricConfig) string {
	// Find which metrics in this group actually exist in the data
	var existingMetrics []string
	for _, metricName := range groupMetrics {
		if _, exists := runMetrics[0][metricName]; exists {
			existingMetrics = append(existingMetrics, metricName)
			delete(ungroupedMetrics, metricName) // Remove from ungrouped
		}
	}

	if len(existingMetrics) == 0 {
		return ""
	}

	// Prepare data for multi-series chart
	metricValues := make([][]float64, len(existingMetrics))
	maxValues := make([]float64, len(existingMetrics))
	unitLabels := make([]string, len(existingMetrics))
	scaleFactors := make([]float64, len(existingMetrics))

	for j, metricName := range existingMetrics {
		metricValues[j] = make([]float64, len(runs))
		for i := range runs {
			if val, exists := runMetrics[i][metricName]; exists {
				metricValues[j][i] = val
			}
		}
		maxValues[j] = metricMaxValues[metricName]

		// Get unit label and scale factor from config
		config, exists := metricConfigs[metricName]
		if exists {
			unitLabels[j] = config.UnitLabel
			scaleFactors[j] = config.ScaleFactor
		} else {
			unitLabels[j] = ""
			scaleFactors[j] = 1.0
		}
	}

	return createMultiSeriesBarChart(
		groupName,
		dates,
		existingMetrics,
		metricValues,
		maxValues,
		unitLabels,
		scaleFactors,
	)
}

// generateUngroupedCharts generates individual charts for ungrouped metrics
func generateUngroupedCharts(dates []time.Time, runs []report, runMetrics []map[string]float64,
	metricMaxValues map[string]float64, metricConfigs map[string]MetricConfig) string {
	if len(runMetrics) == 0 {
		return ""
	}

	var content string

	// Identify ungrouped metrics
	ungroupedMetrics := make(map[string]bool)
	for metricName := range runMetrics[0] {
		if config, exists := metricConfigs[metricName]; !exists || config.Group == "" {
			ungroupedMetrics[metricName] = true
		}
	}

	// Create individual charts for ungrouped metrics
	for metricName := range ungroupedMetrics {
		values := make([]float64, len(runs))
		for i := range runs {
			values[i] = runMetrics[i][metricName]
		}

		// Get unit label and scale factor
		unitLabel := ""
		scaleFactor := 1.0
		if config, exists := metricConfigs[metricName]; exists {
			unitLabel = config.UnitLabel
			scaleFactor = config.ScaleFactor
		}

		content += createBarChart(
			"Metric: "+metricName,
			dates,
			values,
			metricMaxValues[metricName],
			unitLabel,
			scaleFactor,
		)
	}

	return content
}
