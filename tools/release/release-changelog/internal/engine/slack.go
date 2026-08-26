package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SlackThread identifies a Slack thread to post into.
type SlackThread struct {
	Channel  string
	ThreadTS string
}

// threadURLPattern parses Slack thread/message links like:
// https://workspace.slack.com/archives/C0123ABCD/p16900000001234567
var threadURLPattern = regexp.MustCompile(`^https://[a-z0-9-]+\.slack\.com/archives/([A-Z0-9]+)/p(\d{16})\b`)

// ParseSlackThreadURL extracts channel ID and thread timestamp from a Slack
// thread link.
func ParseSlackThreadURL(raw string) (SlackThread, error) {
	m := threadURLPattern.FindStringSubmatch(strings.TrimSpace(raw))
	if m == nil {
		return SlackThread{}, fmt.Errorf("not a Slack thread URL: %q (expected https://<workspace>.slack.com/archives/<channel>/p<ts>)", raw)
	}
	digits := m[2]
	ts := digits[:len(digits)-6] + "." + digits[len(digits)-6:]
	return SlackThread{Channel: m[1], ThreadTS: ts}, nil
}

// slackClient posts messages and files to Slack.
type slackClient struct {
	httpClient *http.Client
	token      string
	apiBase    string // overridable for tests
}

func newSlackClient(token string) *slackClient {
	return &slackClient{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		token:      token,
		apiBase:    "https://slack.com/api",
	}
}

type slackAPIResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// doJSON performs a Slack Web API call with a JSON body.
func (s *slackClient) doJSON(ctx context.Context, method string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiBase+"/"+method, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("slack %s: %w", method, err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var sr slackAPIResponse
	if err := json.Unmarshal(respBody, &sr); err != nil {
		return fmt.Errorf("slack %s: decoding response: %w (body: %s)", method, err, truncate(string(respBody), 200))
	}
	if !sr.OK {
		return fmt.Errorf("slack %s: API error: %s", method, sr.Error)
	}
	return nil
}

// PostMessage posts text into a thread.
func (s *slackClient) PostMessage(ctx context.Context, t SlackThread, text string) error {
	return s.doJSON(ctx, "chat.postMessage", map[string]any{
		"channel":      t.Channel,
		"thread_ts":    t.ThreadTS,
		"text":         text,
		"unfurl_links": false,
	})
}

// UploadFile uploads content as a file into a thread using the external
// upload flow (files.getUploadURLExternal -> raw upload ->
// files.completeUploadExternal).
func (s *slackClient) UploadFile(ctx context.Context, t SlackThread, filename, title string, content []byte, comment string) error {
	// Step 1: get an upload URL.
	form := url.Values{
		"filename": {filename},
		"length":   {strconv.Itoa(len(content))},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiBase+"/files.getUploadURLExternal", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("files.getUploadURLExternal: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var urlResp struct {
		slackAPIResponse
		UploadURL string `json:"upload_url"`
		FileID    string `json:"file_id"`
	}
	if unmarshalErr := json.Unmarshal(body, &urlResp); unmarshalErr != nil {
		return fmt.Errorf("files.getUploadURLExternal: decoding response: %w", unmarshalErr)
	}
	if !urlResp.OK {
		return fmt.Errorf("files.getUploadURLExternal: API error: %s", urlResp.Error)
	}

	// Step 2: upload the raw bytes to the provided URL (no auth header).
	upReq, err := http.NewRequestWithContext(ctx, http.MethodPost, urlResp.UploadURL, bytes.NewReader(content))
	if err != nil {
		return err
	}
	upResp, err := s.httpClient.Do(upReq)
	if err != nil {
		return fmt.Errorf("uploading file bytes: %w", err)
	}
	defer upResp.Body.Close()
	if upResp.StatusCode != http.StatusOK {
		upBody, _ := io.ReadAll(upResp.Body)
		return fmt.Errorf("uploading file bytes: HTTP %d: %s", upResp.StatusCode, truncate(string(upBody), 200))
	}

	// Step 3: complete the upload into the thread.
	payload := map[string]any{
		"files": []map[string]string{
			{"id": urlResp.FileID, "title": title},
		},
		"channel_id": t.Channel,
		"thread_ts":  t.ThreadTS,
	}
	if comment != "" {
		payload["initial_comment"] = comment
	}
	return s.doJSON(ctx, "files.completeUploadExternal", payload)
}
