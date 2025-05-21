package tracking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Mode string

const (
	Mode_Offline Mode = "offline"
	Mode_Online  Mode = "online"

	DISABLE_TRACKING_ENV_VAR = "DISABLE_DX_TRACKING"
	EnvVarLogLevel           = "DX_LOG_LEVEL"
)

type DxTracker struct {
	mode     Mode
	testMode bool
	noOp     bool

	logger zerolog.Logger

	apiToken       string
	githubUsername string
}

func NewDxTracker() (*DxTracker, error) {
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

	if os.Getenv(DISABLE_TRACKING_ENV_VAR) == "true" {
		t.noOp = true
		return t, nil
	}

	c, isConfigAvailable, configErr := openConfig()
	if configErr != nil {
		return nil, errors.Wrap(configErr, "failed to open local config")
	}

	// if local config is available read it and set mode to online
	if isConfigAvailable {
		t.mode = Mode_Online
	} else {
		// if local config is not available check if gh cli is available
		// try to configure tracker with gh cli
		if t.checkIfGhCLIAvailable() {
			var userNameErr error
			c = &config{}
			c.GithubUsername, userNameErr = t.readGHUsername()
			if userNameErr != nil {
				return nil, errors.Wrap(userNameErr, "failed to read github username")
			}

			var apiTokenErr error
			c.DxAPIToken, apiTokenErr = t.readDXAPIToken()
			if apiTokenErr != nil {
				return nil, errors.Wrap(apiTokenErr, "failed to read github api token")
			}

			saveErr := saveConfig(c)
			if saveErr != nil {
				return nil, errors.Wrap(saveErr, "failed to save config")
			}

			t.mode = Mode_Online
		} else {
			// if gh cli is not available, set mode to offline
			t.mode = Mode_Offline
		}
	}

	if t.mode == Mode_Online {
		t.apiToken = c.DxAPIToken
		t.githubUsername = c.GithubUsername

		go func() {
			sendErr := t.sendSavedEvents()
			if sendErr != nil {
				log.Debug().Msgf("failed to send saved events: %s\n", sendErr)
			}
		}()
	}

	t.logger.Debug().Msgf("DxTracker initialized with mode: %s", t.mode)

	return t, nil
}

func (t *DxTracker) Track(event string, metadata map[string]any) error {
	if t.noOp {
		return nil
	}

	if validateErr := validateEvent(event, time.Now().Unix(), metadata); validateErr != nil {
		return errors.Wrap(validateErr, "failed to validate event")
	}

	timestamp := time.Now().Unix()

	if t.mode == Mode_Online {
		sendErr := t.sendEvent(event, timestamp, metadata)
		if sendErr != nil {
			t.logger.Debug().Msgf("failed to send event: %s", sendErr)
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

func (t *DxTracker) sendEvent(name string, timestamp int64, metadata map[string]any) error {
	url := "https://api.getdx.com/events.track"

	body := map[string]any{
		"name":            name,
		"metadata":        metadata,
		"timestamp":       fmt.Sprintf("%d", timestamp),
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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiToken))

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

func (t *DxTracker) checkIfGhCLIAvailable() bool {
	cmd := exec.Command("gh", "auth", "status")
	_, err := cmd.Output()

	t.logger.Info().Msgf("gh CLI available: %t", err == nil)

	return err == nil
}

func (t *DxTracker) readGHUsername() (string, error) {
	cmd := exec.Command("gh", "api", "user", "--jq", ".login")
	output, err := cmd.Output()
	if err != nil {
		t.logger.Debug().Msgf("Failed to run GH CLI: %s", err)
		return "", errors.Wrap(err, "failed to run GH CLI")
	}

	username := strings.Trim(strings.TrimSpace(string(output)), "\n\r")
	if username == "" {
		t.logger.Debug().Msg("Github username not found")
		return "", errors.New("Github username not found")
	}

	t.logger.Debug().Msgf("Github username found: %s", username)

	return strings.Trim(strings.TrimSpace(string(output)), "\n\r"), nil
}

func (t *DxTracker) readDXAPIToken() (string, error) {
	cmd := exec.Command("gh", "variable", "get", "DX_API_TOKEN", "--repo", "smartcontractkit/local-cre-dx-tracking")
	output, err := cmd.Output()
	if err != nil {
		t.logger.Debug().Msgf("failed to run GH CLI: %s", err)
		return "", errors.Wrap(err, "failed to run GH CLI")
	}

	if len(output) == 0 {
		t.logger.Debug().Msg("DX API token not found")
		return "", errors.New("DX API token not found")
	}

	t.logger.Debug().Msg("DX API token found")

	return strings.Trim(strings.TrimSpace(string(output)), "\n\r"), nil
}

type config struct {
	DxAPIToken     string `json:"dx_api_token"`
	GithubUsername string `json:"github_username"`
}

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

type event struct {
	Name      string         `json:"name"`
	Timestamp int64          `json:"timestamp"`
	Metadata  map[string]any `json:"metadata"`
}

func (t *DxTracker) saveEvent(name string, timestamp int64, metadata map[string]any) error {
	storagePath, pathErr := storagePath()
	if pathErr != nil {
		return errors.Wrap(pathErr, "failed to get storage path")
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

	if err := os.WriteFile(storagePath, jsonData, 0644); err != nil {
		return errors.Wrap(err, "failed to write event to storage file")
	}

	return nil
}

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

	for _, event := range events {
		sendErr := t.sendEvent(event.Name, event.Timestamp, event.Metadata)
		if sendErr != nil {
			return errors.Wrap(sendErr, "failed to send event")
		}
	}

	clearErr := t.clearSavedEvents()
	if clearErr != nil {
		return errors.Wrap(clearErr, "failed to clear saved events")
	}

	return nil
}

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

func storagePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home directory")
	}

	return filepath.Join(homeDir, ".dx", "events.json"), nil
}

func configPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home directory")
	}

	return filepath.Join(homeDir, ".dx", "config.json"), nil
}
