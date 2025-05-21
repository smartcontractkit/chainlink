package tracking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type Mode string

const (
	Mode_Offline Mode = "offline"
	Mode_Online  Mode = "online"

	DISABLE_TRACKING_ENV_VAR = "DISABLE_DX_TRACKING"
)

type DxTracker struct {
	mode     Mode
	testMode bool
	noOp     bool

	apiToken       string
	githubUsername string
}

func (t *DxTracker) New() (*DxTracker, error) {
	if os.Getenv(DISABLE_TRACKING_ENV_VAR) == "true" {
		return &DxTracker{
			noOp: true,
		}, nil
	}

	c, isConfigAvailable, err := openConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open local config")
	}

	// if local config is available read it and set mode to online
	if isConfigAvailable {
		t.mode = Mode_Online
	} else {
		// if local config is not available check if gh cli is available
		// try to configure tracker with gh cli
		if checkIfGhCLIAvailable() {
			var userNameErr error
			c.GithubUsername, userNameErr = readGHUsername()
			if userNameErr != nil {
				return nil, errors.Wrap(userNameErr, "failed to read github username")
			}

			var apiTokenErr error
			c.DxAPIToken, apiTokenErr = readDXAPIToken()
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
			storageErr := createLocalStorage()
			if storageErr != nil {
				return nil, errors.Wrap(storageErr, "failed to create local storage")
			}
		}
	}

	if t.mode == Mode_Online {
		t.apiToken = c.DxAPIToken
		t.githubUsername = c.GithubUsername

		go func() {
			sendErr := t.sendSavedEvents()
			if sendErr != nil {
				fmt.Fprintf(os.Stderr, "failed to send saved events: %s\n", sendErr)
			}
		}()
	}

	return t, nil
}

func (t *DxTracker) Track(event string, metadata map[string]any) error {
	if t.noOp {
		return nil
	}

	if validateErr := validateEvent(event, time.Now().Unix(), metadata); validateErr != nil {
		return validateErr
	}

	if t.mode == Mode_Online {
		return t.sendEvent(event, time.Now().Unix(), metadata)
	}

	return t.saveEvent(event, metadata)
}

func (t *DxTracker) sendEvent(eventName string, timestamp int64, metadata map[string]any) error {
	url := "https://api.getdx.com/events.track"

	body := map[string]any{
		"event":     eventName,
		"metadata":  metadata,
		"timestamp": timestamp,
	}

	if t.testMode {
		body["test_data"] = true
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "failed to marshal event")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

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

func checkIfGhCLIAvailable() bool {
	cmd := exec.Command("gh", "auth", "status")
	_, err := cmd.Output()

	return err == nil
}

func readGHUsername() (string, error) {
	cmd := exec.Command("gh", "api", "user", "--jq", ".login")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to run gh cli")
	}
	return string(output), nil
}

func readDXAPIToken() (string, error) {
	cmd := exec.Command("gh", "variable", "get", "DX_API_TOKEN", "--repo", "smartcontractkit/local-cre-dx-tracking")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to run gh cli")
	}
	return string(output), nil
}

type config struct {
	DxAPIToken     string `json:"dx_api_token"`
	GithubUsername string `json:"github_username"`
}

func openConfig() (*config, bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get user home directory")
	}

	configPath := filepath.Join(homeDir, ".dx", "config.json")

	if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
		return nil, false, nil
	}

	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to read config file")
	}

	var localConfig config
	err = json.Unmarshal(configContent, &localConfig)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to unmarshal config file")
	}

	return &localConfig, true, nil
}

func saveConfig(c *config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get user home directory")
	}

	configPath := filepath.Join(homeDir, ".dx", "config.json")

	err = os.MkdirAll(filepath.Dir(configPath), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to create config directory")
	}

	configFile, err := os.Create(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to create config file")
	}
	defer configFile.Close()

	jsonData, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	_, err = configFile.Write(jsonData)
	if err != nil {
		return errors.Wrap(err, "failed to write config file")
	}

	return nil
}

func createLocalStorage() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get user home directory")
	}

	storagePath := filepath.Join(homeDir, ".dx", "storage", "events.json")

	err = os.MkdirAll(storagePath, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create storage directory")
	}

	storageFile, err := os.Create(storagePath)
	if err != nil {
		return errors.Wrap(err, "failed to create storage file")
	}
	defer storageFile.Close()

	return nil
}

func (t *DxTracker) saveEvent(event string, metadata map[string]any) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get user home directory")
	}

	storagePath := filepath.Join(homeDir, ".dx", "storage", "events.json")

	storageFile, err := os.OpenFile(storagePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open storage file")
	}
	defer storageFile.Close()

	var events []map[string]any
	if _, err := os.Stat(storageFile.Name()); err == nil {
		existingData, err := os.ReadFile(storageFile.Name())
		if err != nil {
			return errors.Wrap(err, "failed to read existing storage file")
		}
		if len(existingData) > 0 {
			if err := json.Unmarshal(existingData, &events); err != nil {
				return errors.Wrap(err, "failed to parse existing JSON data")
			}
		}
	}

	events = append(events, map[string]any{
		"event":     event,
		"metadata":  metadata,
		"timestamp": time.Now().Unix(),
	})

	jsonData, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal events to JSON")
	}

	_, err = storageFile.Write(jsonData)
	if err != nil {
		return errors.Wrap(err, "failed to write event to storage file")
	}

	return nil
}

func (t *DxTracker) sendSavedEvents() error {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return errors.Wrap(homeErr, "failed to get user home directory")
	}

	storagePath := filepath.Join(homeDir, ".dx", "storage", "events.json")

	if _, statErr := os.Stat(storagePath); os.IsNotExist(statErr) {
		return nil
	}

	storageFile, storageErr := os.OpenFile(storagePath, os.O_RDONLY, 0644)
	if storageErr != nil {
		return errors.Wrap(storageErr, "failed to open storage file")
	}
	defer storageFile.Close()

	var events []any

	decoderErr := json.NewDecoder(storageFile).Decode(&events)
	if decoderErr != nil {
		return errors.Wrap(decoderErr, "failed to decode events from storage file")
	}

	var toEventFn = func(maybeEvent any) (eventName string, timestamp int64, metadata map[string]any, err error) {
		if maybeEvent == nil {
			return "", 0, nil, errors.New("event is nil")
		}

		if asMap, ok := maybeEvent.(map[string]any); ok {
			eventKey, ok := asMap["event"]
			if !ok {
				return "", 0, nil, fmt.Errorf("potential event doesn't have 'event' key: %v", maybeEvent)
			}

			metadataKey, ok := asMap["metadata"]
			if !ok {
				return "", 0, nil, fmt.Errorf("potential event doesn't have 'metadata' key: %v", maybeEvent)
			}

			timestampKey, ok := asMap["timestamp"]
			if !ok {
				return "", 0, nil, fmt.Errorf("potential event doesn't have 'timestamp' key: %v", maybeEvent)
			}

			if eventContent, ok := eventKey.(string); ok {
				eventName = eventContent
			}

			if metadataContent, ok := metadataKey.(map[string]any); ok {
				metadata = metadataContent
			}

			if timestampContent, ok := timestampKey.(int64); ok {
				timestamp = timestampContent
			}

			if validateErr := validateEvent(eventName, timestamp, metadata); validateErr != nil {
				return "", 0, nil, validateErr
			}

			return eventName, timestamp, metadata, nil
		}

		return "", 0, nil, fmt.Errorf("event is not a map[string]any, but %T", maybeEvent)

	}

	for _, event := range events {
		eventName, timestamp, metadata, toEventErr := toEventFn(event)
		if toEventErr != nil {
			fmt.Printf("failed to parse event: %s. Continuing...\n", toEventErr)
			continue
		}

		sendErr := t.sendEvent(eventName, timestamp, metadata)
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get user home directory")
	}

	storagePath := filepath.Join(homeDir, ".dx", "storage", "events.json")

	file, err := os.OpenFile(storagePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to truncate storage file")
	}
	defer file.Close()

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
