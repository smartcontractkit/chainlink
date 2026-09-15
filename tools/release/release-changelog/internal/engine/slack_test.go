package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseSlackThreadURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		url         string
		wantChannel string
		wantTS      string
		wantErr     bool
	}{
		{"https://myorg.slack.com/archives/C0123ABCD/p1690000000123456", "C0123ABCD", "1690000000.123456", false},
		{"https://my-org.slack.com/archives/G9876ZYXW/p1753967999000001?thread_ts=1753967999.000001&cid=G9876ZYXW", "G9876ZYXW", "1753967999.000001", false},
		{"https://myorg.slack.com/messages/C0123ABCD", "", "", true},
		{"not a url", "", "", true},
		{"", "", "", true},
	}
	for _, c := range cases {
		th, err := ParseSlackThreadURL(c.url)
		if (err != nil) != c.wantErr {
			t.Errorf("ParseSlackThreadURL(%q) err = %v, wantErr %v", c.url, err, c.wantErr)
			continue
		}
		if !c.wantErr && (th.Channel != c.wantChannel || th.ThreadTS != c.wantTS) {
			t.Errorf("ParseSlackThreadURL(%q) = (%s, %s), want (%s, %s)",
				c.url, th.Channel, th.ThreadTS, c.wantChannel, c.wantTS)
		}
	}
}

func TestSlackPostMessage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat.postMessage" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer xoxb-test" {
			t.Error("missing auth")
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decoding request: %v", err)
		}
		if payload["channel"] != "C123" || payload["thread_ts"] != "1.2" || payload["text"] != "hello" {
			t.Errorf("unexpected payload: %v", payload)
		}
		fmt.Fprint(w, `{"ok":true}`)
	}))
	defer srv.Close()

	c := newSlackClient("xoxb-test")
	c.apiBase = srv.URL
	if err := c.PostMessage(context.Background(), SlackThread{Channel: "C123", ThreadTS: "1.2"}, "hello"); err != nil {
		t.Fatal(err)
	}
}

func TestSlackUploadFile(t *testing.T) {
	t.Parallel()

	var uploadReceived []byte
	var completePayload map[string]any
	var serverURL string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/files.getUploadURLExternal":
			fmt.Fprintf(w, `{"ok":true,"upload_url":%q,"file_id":"F123"}`, serverURL+"/upload")
		case "/upload":
			uploadReceived, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusOK)
		case "/files.completeUploadExternal":
			if err := json.NewDecoder(r.Body).Decode(&completePayload); err != nil {
				t.Errorf("decoding request: %v", err)
			}
			fmt.Fprint(w, `{"ok":true}`)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	}))
	defer srv.Close()
	serverURL = srv.URL

	c := newSlackClient("xoxb-test")
	c.apiBase = srv.URL
	content := []byte("# markdown report")
	err := c.UploadFile(context.Background(), SlackThread{Channel: "C123", ThreadTS: "1.2"}, "report.md", "Report", content, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(uploadReceived), "markdown report") {
		t.Errorf("upload bytes not received: %q", uploadReceived)
	}
	if completePayload["channel_id"] != "C123" || completePayload["thread_ts"] != "1.2" {
		t.Errorf("unexpected complete payload: %v", completePayload)
	}
	files, ok := completePayload["files"].([]any)
	if !ok || len(files) != 1 {
		t.Fatalf("unexpected files payload: %v", completePayload)
	}
	f0 := files[0].(map[string]any)
	if f0["id"] != "F123" || f0["title"] != "Report" {
		t.Errorf("unexpected file entry: %v", f0)
	}
}

func TestSlackAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":false,"error":"channel_not_found"}`)
	}))
	defer srv.Close()

	c := newSlackClient("xoxb-test")
	c.apiBase = srv.URL
	err := c.PostMessage(context.Background(), SlackThread{Channel: "C999", ThreadTS: "1.2"}, "hi")
	if err == nil || !strings.Contains(err.Error(), "channel_not_found") {
		t.Errorf("expected channel_not_found error, got %v", err)
	}
}
