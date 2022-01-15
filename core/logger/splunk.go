package logger

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type event struct {
	Time       int64           `json:"time"`                 // epoch time in nano-seconds
	Host       string          `json:"host"`                 // hostname
	Source     string          `json:"source,omitempty"`     // optional description of the source of the event; typically the app's name
	SourceType string          `json:"sourcetype,omitempty"` // optional name of a Splunk parsing configuration; this is usually inferred by Splunk
	Index      string          `json:"index,omitempty"`      // optional name of the Splunk index to store the event in; not required if the token has a default index set in Splunk
	Event      json.RawMessage `json:"event"`                // zap json logger
}

type splunkSink struct {
	HTTPClient    httpClient
	nowFn         func() int64
	URL           string // URL to Splunk HTTP Event Collector
	Hostname      string // Hostname is the human-readable machine name sending logs
	Token         string // Splunk HTTP Event Collector token
	Source        string // Splunk source field value, description of the source of the event
	SourceType    string // Splunk source type, optional name of a sourcetype field value
	Index         string // Splunk Index, optional name of the Splunk index to store the event in
	SkipTLSVerify bool   // Skip verifying the certificate of the HTTP Event Collector
	events        [][]byte
	newEvent      chan bool
	sync.Mutex
}

func envVarOrDefault(key string, defaultValue string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return defaultValue
	}
	return envVar
}

func initialize(s *splunkSink) *splunkSink {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: s.SkipTLSVerify}}
	s.HTTPClient = &http.Client{Timeout: time.Second * 20, Transport: tr}

	if s.Hostname == "" {
		s.Hostname, _ = os.Hostname()
	}
	if s.nowFn == nil {
		s.nowFn = func() int64 {
			return time.Now().UnixNano()
		}
	}
	s.newEvent = make(chan bool, 1)

	runSync := func() {
		err := s.Sync()
		if err != nil {
			log.Printf("Error sending to Splunk: %v\n", err)
		}
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		atLeastOneEvent := false
		for {
			select {
			case <-s.newEvent:
				atLeastOneEvent = true
				if len(s.events) > 50 {
					atLeastOneEvent = false
					go runSync()
				}
			case <-ticker.C:
				if atLeastOneEvent {
					go runSync()
				}
			}
		}
	}()

	return s
}

func newSplunkSink() (*splunkSink, error) {
	s := &splunkSink{}
	s.URL = os.Getenv("SPLUNK_URL")
	s.Source = envVarOrDefault("SPLUNK_SOURCE", "chainlink")
	s.SourceType = envVarOrDefault("SPLUNK_SOURCETYPE", "chainlink")
	s.Index = os.Getenv("SPLUNK_INDEX")
	s.Hostname = os.Getenv("SPLUNK_HOSTNAME")
	s.Token = os.Getenv("SPLUNK_TOKEN")
	s.SkipTLSVerify, _ = strconv.ParseBool(os.Getenv("SPLUNK_SKIP_TLS_VERIFY"))

	return initialize(s), nil
}

func (s *splunkSink) Close() error {
	return nil
}

// Write implement zap.Sink func Write
func (s *splunkSink) Write(b []byte) (n int, err error) {
	e := &event{
		Time:       s.nowFn(),
		Host:       s.Hostname,
		Source:     s.Source,
		SourceType: s.SourceType,
		Index:      s.Index,
		Event:      json.RawMessage(b),
	}
	encoded, err := json.Marshal(e)
	if err != nil {
		return 0, err
	}
	s.Lock()
	defer s.Unlock()
	s.events = append(s.events, encoded)
	go func() {
		s.newEvent <- true
	}()

	return len(b), nil
}

// Sync implement zap.Sink func Sync
func (s *splunkSink) Sync() error {
	eventsLen := len(s.events)
	if eventsLen > 0 {
		s.Lock()
		eventsToLog := s.events[0:eventsLen]
		s.events = s.events[eventsLen:]
		s.Unlock()
		return s.logEvents(eventsToLog)

	}
	return nil
}

func (s *splunkSink) logEvents(events [][]byte) error {
	buf := new(bytes.Buffer)
	for _, e := range events {
		buf.Write(e)
	}

	return s.doRequest(buf)
}

func (s *splunkSink) doRequest(b *bytes.Buffer) error {
	url := s.URL
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Splunk "+s.Token)

	res, err := s.HTTPClient.Do(req)
	if err != nil && err != io.EOF {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		responseBody := buf.String()
		err = errors.New(responseBody)
	} else {
		_, _ = io.Copy(ioutil.Discard, res.Body)
	}
	return err
}
