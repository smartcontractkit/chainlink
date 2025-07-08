package tracking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Mode determines whether events are sent immediately or stored for later transmission.
type Mode string

const (
	ModeOffline Mode = "offline"
	ModeOnline  Mode = "online"

	EnvVarLogLevel         = "DX_LOG_LEVEL"
	EnvVarTestMode         = "DX_TEST_MODE"
	EnvVarForceOfflineMode = "DX_FORCE_OFFLINE_MODE"
	EnvVarDisableTracking  = "DISABLE_DX_TRACKING"

	MinGHCLIVersion = "v2.50.0"
)

type Tracker interface {
	Track(event string, metadata map[string]any) error
}

type NoOpTracker struct{}

func (t *NoOpTracker) Track(event string, metadata map[string]any) error {
	return nil
}

// DxTracker manages event tracking with automatic retry and offline support.
type DxTracker struct {
	mode     Mode
	testMode bool

	logger zerolog.Logger

	apiToken       string
	githubUsername string
}

// NewDxTracker initializes a tracker with automatic GitHub CLI integration for authentication.
func NewDxTracker() (Tracker, error) {
	t := &DxTracker{}

	lvlStr := os.Getenv(EnvVarLogLevel)
	if lvlStr == "" {
		lvlStr = "info"
	}
	lvl, lvlErr := zerolog.ParseLevel(lvlStr)
	if lvlErr != nil {
		return nil, errors.Wrap(lvlErr, "failed to parse log level")
	}
	t.logger = log.With().Str("logger_name", "DxTracker").Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(lvl).With().Logger()
	t.logger.Debug().Msg("Initializing DxTracker")

	if os.Getenv(EnvVarDisableTracking) == "true" {
		t.logger.Debug().Msg("Tracking disabled by environment variable")

		return &NoOpTracker{}, nil
	}

	if os.Getenv(EnvVarTestMode) == "true" {
		t.testMode = true
		t.logger.Debug().Msg("Tracking in test mode")
	}

	c, isConfigAvailable, configErr := openConfig()
	if configErr != nil {
		return nil, errors.Wrap(configErr, "failed to open local config")
	}

	// if local config is available read it and set mode to online
	if isConfigAvailable && isConfigValid(c) {
		t.logger.Debug().Msg("Valid local config found")
		t.mode = ModeOnline
	} else {
		// if local config is not available check if GH CLI is available
		// and if so, try to configure tracker with it
		if t.checkIfGhCLIAvailable() {
			var configErr error
			c, configErr = t.buildConfigWithGhCLI()
			if configErr != nil {
				t.mode = ModeOffline
				t.logger.Warn().Msgf("Failed to build config with GH CLI: %s", configErr.Error())
			} else {
				t.mode = ModeOnline
				t.logger.Debug().Msg("Config created, setting mode to online")
			}
		} else {
			// if gh cli is not available, set mode to offline
			t.mode = ModeOffline
			t.logger.Debug().Msg("GH CLI not available, setting mode to offline")
		}
	}

	if os.Getenv(EnvVarForceOfflineMode) == "true" {
		t.mode = ModeOffline
		t.logger.Debug().Msg("Tracking forced to offline by environment variable")
	}

	if t.mode == ModeOnline {
		t.apiToken = c.DxAPIToken
		t.githubUsername = c.GithubUsername

		go func() {
			sendErr := t.sendSavedEvents()
			if sendErr != nil {
				log.Debug().Msgf("Failed to send saved events: %s\n", sendErr)
			}
		}()
	}

	t.logger.Debug().Msgf("DxTracker initialized with mode: %s", t.mode)

	return t, nil
}

func (t *DxTracker) buildConfigWithGhCLI() (*config, error) {
	var userNameErr error
	c := &config{}
	c.GithubUsername, userNameErr = t.readGHUsername()
	if userNameErr != nil {
		return nil, errors.Wrap(userNameErr, "failed to read github username")
	}

	var apiTokenErr error
	c.DxAPIToken, apiTokenErr = t.readDXAPIToken()
	if apiTokenErr != nil {
		return nil, errors.Wrap(apiTokenErr, "failed to read DX API token")
	}

	saveErr := saveConfig(c)
	if saveErr != nil {
		return nil, errors.Wrap(saveErr, "failed to save config")
	}

	return c, nil
}

// Track queues or sends an event, automatically handling offline scenarios.
func (t *DxTracker) Track(event string, metadata map[string]any) error {
	if validateErr := validateEvent(event, time.Now().Unix(), metadata); validateErr != nil {
		return errors.Wrap(validateErr, "failed to validate event")
	}

	timestamp := time.Now().Unix()

	if t.mode == ModeOnline {
		sendErr := t.sendEvent(event, timestamp, metadata)
		if sendErr != nil {
			saveErr := t.saveEvent(event, timestamp, metadata)
			if saveErr != nil {
				t.logger.Debug().Msgf("failed to save event: %s", saveErr)

				return sendErr
			}
		}

		return nil
	}

	return t.saveEvent(event, timestamp, metadata)
}

// sendEvent attempts to deliver an event to the DX API with a 15-second timeout.
func (t *DxTracker) sendEvent(name string, timestamp int64, metadata map[string]any) error {
	url := "https://api.getdx.com/events.track"

	truncatedMetadata, truncateErr := truncateMetadata(metadata, []string{"to_truncate", "debug"})
	if truncateErr != nil {
		return errors.Wrap(truncateErr, "failed to truncate metadata")
	}

	body := map[string]any{
		"name":            name,
		"metadata":        truncatedMetadata,
		"timestamp":       strconv.FormatInt(timestamp, 10),
		"github_username": t.githubUsername,
	}

	if t.testMode {
		body["test_data"] = true
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "failed to marshal event")
	}

	t.logger.Debug().Msgf("Sending event: %s", string(jsonData))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send event")
	}
	defer resp.Body.Close()

	type dxResponse struct {
		Ok    bool    `json:"ok"`
		Error *string `json:"error"`
	}

	var dxResp dxResponse
	err = json.NewDecoder(resp.Body).Decode(&dxResp)
	if err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if !dxResp.Ok {
		return fmt.Errorf("failed to send event, error: %s", *dxResp.Error)
	}

	return nil
}

// checkIfGhCLIAvailable determines if GitHub CLI is available for authentication.
func (t *DxTracker) checkIfGhCLIAvailable() bool {
	cmd := exec.Command("gh", "auth", "status")
	_, outputErr := cmd.Output()

	isAvailable := outputErr == nil
	if !isAvailable {
		return false
	}

	cmd = exec.Command("gh", "--version")
	output, outputErr := cmd.Output()
	if outputErr != nil {
		t.logger.Warn().Msgf("failed to get GH CLI version: %s", outputErr.Error())
		return false
	}

	re := regexp.MustCompile(`gh version (\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		t.logger.Warn().Msgf("failed to parse GH CLI version: %s", string(output))
		return false
	}

	version, versionErr := semver.NewVersion(matches[1])
	if versionErr != nil {
		t.logger.Warn().Msgf("failed to parse GH CLI version: %s", versionErr.Error())
		return false
	}

	isEnoughVersion := version.Compare(semver.MustParse(MinGHCLIVersion)) >= 0
	if !isEnoughVersion {
		t.logger.Warn().Msgf("GH CLI version is too old, please update to at least %s", MinGHCLIVersion)
	}

	t.logger.Debug().Msgf("GH CLI version found: %s", version)

	return true
}

// readGHUsername fetches the authenticated GitHub username via CLI.
func (t *DxTracker) readGHUsername() (string, error) {
	cmd := exec.Command("gh", "api", "user", "--jq", ".login")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to run GH CLI")
	}

	username := strings.Trim(strings.TrimSpace(string(output)), "\n\r")
	if username == "" {
		return "", errors.New("Github username not found")
	}

	t.logger.Debug().Msgf("Github username found: %s", username)

	return strings.Trim(strings.TrimSpace(string(output)), "\n\r"), nil
}

// readDXAPIToken retrieves the API token from GitHub repository secrets.
func (t *DxTracker) readDXAPIToken() (string, error) {
	cmd := exec.Command("gh", "variable", "get", "DX_API_TOKEN", "--repo", "smartcontractkit/local-cre-dx-tracking")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to run GH CLI")
	}

	if len(output) == 0 {
		return "", errors.New("DX API token not found")
	}

	t.logger.Debug().Msg("DX API token found")

	return strings.Trim(strings.TrimSpace(string(output)), "\n\r"), nil
}

// config stores authentication credentials for the DX API.
type config struct {
	DxAPIToken     string `json:"dx_api_token"`
	GithubUsername string `json:"github_username"`
}

// openConfig attempts to load existing configuration from the user's home directory.
func openConfig() (*config, bool, error) {
	configPath, pathErr := configPath()
	if pathErr != nil {
		return nil, false, errors.Wrap(pathErr, "failed to get config path")
	}

	if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
		return nil, false, nil
	}

	configContent, readErr := os.ReadFile(configPath)
	if readErr != nil {
		return nil, false, errors.Wrap(readErr, "failed to read config file")
	}

	var localConfig config
	unmarshalErr := json.Unmarshal(configContent, &localConfig)
	if unmarshalErr != nil {
		return nil, false, errors.Wrap(unmarshalErr, "failed to unmarshal config file")
	}

	return &localConfig, true, nil
}

// isConfigValid ensures both API token and GitHub username are present.
func isConfigValid(c *config) bool {
	return c.DxAPIToken != "" && c.GithubUsername != ""
}

// saveConfig persists configuration to the user's home directory with proper permissions.
func saveConfig(c *config) error {
	configPath, pathErr := configPath()
	if pathErr != nil {
		return errors.Wrap(pathErr, "failed to get config path")
	}

	mkdirErr := os.MkdirAll(filepath.Dir(configPath), 0755)
	if mkdirErr != nil {
		return errors.Wrap(mkdirErr, "failed to create config directory")
	}

	configFile, createErr := os.Create(configPath)
	if createErr != nil {
		return errors.Wrap(createErr, "failed to create config file")
	}
	defer configFile.Close()

	jsonData, marshalErr := json.Marshal(c)
	if marshalErr != nil {
		return errors.Wrap(marshalErr, "failed to marshal config")
	}

	_, writeErr := configFile.Write(jsonData)
	if writeErr != nil {
		return errors.Wrap(writeErr, "failed to write config file")
	}

	return nil
}

// event represents a tracking event with its associated metadata.
type event struct {
	Name      string         `json:"name"`
	Timestamp int64          `json:"timestamp"`
	Metadata  map[string]any `json:"metadata"`
}

// saveEvent stores an event locally for later transmission when offline.
func (t *DxTracker) saveEvent(name string, timestamp int64, metadata map[string]any) error {
	t.logger.Debug().Msgf("Saving event. Name: %s, Timestamp: %d, Metadata: %v", name, timestamp, metadata)

	storagePath, pathErr := storagePath()
	if pathErr != nil {
		return errors.Wrap(pathErr, "failed to get storage path")
	}

	mkdirErr := os.MkdirAll(filepath.Dir(storagePath), 0755)
	if mkdirErr != nil {
		return errors.Wrap(mkdirErr, "failed to create storage directory")
	}

	var events []event

	if _, statErr := os.Stat(storagePath); statErr == nil {
		content, err := os.ReadFile(storagePath)
		if err == nil && len(content) > 0 {
			if err := json.Unmarshal(content, &events); err != nil {
				t.logger.Debug().Msgf("Failed to parse JSON: %s", err)
				events = []event{}
			}
		}
	}

	newEvent := event{
		Name:      name,
		Timestamp: timestamp,
		Metadata:  metadata,
	}
	events = append(events, newEvent)

	jsonData, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal events to JSON")
	}

	if err := os.WriteFile(storagePath, jsonData, 0600); err != nil {
		return errors.Wrap(err, "failed to write event to storage file")
	}

	return nil
}

// sendSavedEvents attempts to send all queued events and clears the queue on success.
func (t *DxTracker) sendSavedEvents() error {
	storagePath, pathErr := storagePath()
	if pathErr != nil {
		return errors.Wrap(pathErr, "failed to get storage path")
	}

	stats, statErr := os.Stat(storagePath)
	if os.IsNotExist(statErr) {
		return nil
	}

	if stats.Size() == 0 {
		return nil
	}

	storageFile, storageErr := os.OpenFile(storagePath, os.O_RDONLY, 0644)
	if storageErr != nil {
		return errors.Wrap(storageErr, "failed to open storage file")
	}
	defer storageFile.Close()

	var events []event

	decoderErr := json.NewDecoder(storageFile).Decode(&events)
	if decoderErr != nil {
		return errors.Wrap(decoderErr, "failed to decode events from storage file")
	}

	t.logger.Debug().Msgf("Sending %d saved events", len(events))

	failedEvents := []event{}

	for _, event := range events {
		sendErr := t.sendEvent(event.Name, event.Timestamp, event.Metadata)
		if sendErr != nil {
			failedEvents = append(failedEvents, event)
			t.logger.Debug().Msgf("Failed to send event: %s", sendErr.Error())
		}
	}

	clearErr := t.clearSavedEvents()
	if clearErr != nil {
		return errors.Wrap(clearErr, "failed to clear saved events")
	}

	// if there are failed events, save them to the storage file to try again later
	if len(failedEvents) > 0 {
		t.logger.Warn().Msgf("Failed to send %d events", len(failedEvents))
		for _, event := range failedEvents {
			saveErr := t.saveEvent(event.Name, event.Timestamp, event.Metadata)
			if saveErr != nil {
				t.logger.Warn().Msgf("Failed to save failed event: %s", saveErr.Error())
			}
		}
	} else {
		t.logger.Debug().Msg("All saved events sent successfully")
	}

	return nil
}

// clearSavedEvents removes all queued events after successful transmission.
func (t *DxTracker) clearSavedEvents() error {
	storagePath, pathErr := storagePath()
	if pathErr != nil {
		return errors.Wrap(pathErr, "failed to get storage path")
	}

	storageFile, openErr := os.OpenFile(storagePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return errors.Wrap(openErr, "failed to truncate storage file")
	}
	defer storageFile.Close()

	return nil
}

// validateEvent ensures all required event fields are present and non-empty.
func validateEvent(event string, timestamp int64, metadata map[string]any) error {
	if event == "" {
		return errors.New("event is required")
	}

	if timestamp == 0 {
		return errors.New("timestamp is required")
	}

	if len(metadata) == 0 {
		return errors.New("metadata is required")
	}

	return nil
}

// storagePath returns the path to the events queue file in the user's home directory.
func storagePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home directory")
	}

	return filepath.Join(homeDir, ".local", "share", "dx", "events.json"), nil
}

// configPath returns the path to the configuration file in the user's home directory.
func configPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home directory")
	}

	return filepath.Join(homeDir, ".local", "share", "dx", "config.json"), nil
}

const maxMetadataSize = 1024
const truncatedNote = " (... truncated)"

// truncateMetadata ensures the JSON-serialized metadata fits within the 1024-byte limit
// by intelligently truncating string values while preserving as much data as possible.
//
// The function creates a deep copy of the input metadata to avoid mutations, then applies
// a priority-based truncation strategy:
//
//  1. If the metadata is already under the size limit, it's returned unchanged
//  2. A priority list is built starting with keys specified in truncateFirst,
//     followed by remaining keys sorted by string value length (descending)
//  3. For each key in priority order:
//     - First attempts to keep the full value if it fits
//     - If not, uses binary search to find the optimal truncation point
//     - Appends " (... truncated)" suffix to indicate truncation
//     - Non-string values are skipped during truncation
//
// Parameters:
//   - metadata: The original metadata map to potentially truncate
//   - truncateFirst: Keys to prioritize for truncation (useful for debug/verbose fields)
//
// Returns:
//   - A metadata map that fits within maxMetadataSize (1024 bytes) when JSON-serialized
//   - An error if the metadata cannot be reduced to fit even after aggressive truncation
//
// The function guarantees that the returned metadata, when marshaled to JSON,
// will not exceed 1024 bytes in size.
func truncateMetadata(metadata map[string]any, truncateFirst []string) (map[string]any, error) {
	// Deep copy to avoid mutating input
	original := maps.Clone(metadata)

	getSize := func(m map[string]any) int {
		b, _ := json.Marshal(m)
		return len(b)
	}

	if getSize(original) <= maxMetadataSize {
		return original, nil
	}

	// Build priority list
	priority := make([]string, 0, len(original))
	seen := map[string]bool{}
	for _, k := range truncateFirst {
		if _, ok := original[k]; ok {
			priority = append(priority, k)
			seen[k] = true
		}
	}
	for k := range original {
		if !seen[k] {
			priority = append(priority, k)
		}
	}

	// Sort by string length descending
	sort.SliceStable(priority, func(i, j int) bool {
		vi, okI := original[priority[i]].(string)
		vj, okJ := original[priority[j]].(string)
		if !okI {
			return false
		}
		if !okJ {
			return true
		}
		return len(vi) > len(vj)
	})

	// Map of best truncations so far
	bestSoFar := maps.Clone(original)

	for _, key := range priority {
		val, ok := original[key].(string)
		if !ok {
			continue
		}

		// First: try full value
		bestSoFar[key] = val
		if getSize(bestSoFar) <= maxMetadataSize {
			continue
		}

		// Binary search on prefix length
		low, high := 0, len(val)
		best := ""
		for low <= high {
			mid := (low + high) / 2
			trialVal := val[:mid] + truncatedNote

			trialMap := maps.Clone(bestSoFar)
			trialMap[key] = trialVal

			if getSize(trialMap) <= maxMetadataSize {
				best = trialVal
				low = mid + 1
			} else {
				high = mid - 1
			}
		}

		if best != "" {
			bestSoFar[key] = best
		} else {
			bestSoFar[key] = truncatedNote
		}
	}

	if getSize(bestSoFar) <= maxMetadataSize {
		return bestSoFar, nil
	}
	return map[string]any{}, errors.New("unable to fit metadata within 1024 bytes even after truncation")
}
