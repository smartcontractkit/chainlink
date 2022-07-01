package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/pkg/errors"
)

type AuditLogger interface {
	Audit(ctx context.Context, eventID EventID, data map[string]interface{})
}

type AuditLoggerService struct {
	logger          logger.Logger
	enabled         bool
	serviceURL      string
	serviceHeaders  []serviceHeader
	jsonWrapperKey  string
	environmentName string
	hostname        string
	localIP         string
}

// Configurable headers to include in POST to log service
type serviceHeader struct {
	header string
	value  string
}

// NewAuditLogger returns a buffer push system that ingests audit log events and
// asynchronously pushes them up to an HTTP log service.
// Parses and validates the AUDIT_LOGS_* environment values and returns an enabled
// AuditLogger instance. If the environment variables are not set, the logger
// is disabled and short circuits execution via enabled flag.
func NewAuditLogger(logger logger.Logger) (AuditLogger, error) {

	// Start parsing environment variables for audit logger
	auditLogsURL := os.Getenv("AUDIT_LOGS_FORWARDER_URL")
	if auditLogsURL == "" {
		// Unset, return a disabled audit logger
		logger.Info("No AUDIT_LOGS_FORWARDER_URL environment set, audit log events will not be captured")

		return &AuditLoggerService{}, nil
	}

	env := "production"
	if os.Getenv("CHAINLINK_DEV") == "true" {
		env = "develop"
	}
	hostname, err := os.Hostname()
	if err != nil {
		return &AuditLoggerService{}, errors.Errorf("Audit Log initialization error - unable to get hostname", "err", err)
	}

	// Split and prepare optional service client headers from env variable
	headers := []serviceHeader{}
	headersEncoded := os.Getenv("AUDIT_LOGS_FORWARDER_HEADERS")
	if headersEncoded != "" {
		headerLines := strings.Split(headersEncoded, "\\")
		for _, header := range headerLines {
			keyValue := strings.Split(header, "||")
			if len(keyValue) != 2 {
				return &AuditLoggerService{}, errors.Errorf("Invalid AUDIT_LOGS_FORWARDER_HEADERS value, single pair split on || required, got: %s", keyValue)
			}
			headers = append(headers, serviceHeader{
				header: keyValue[0],
				value:  keyValue[1],
			})
		}
	}

	// Finally, create new auditLogger with parameters
	auditLogger := AuditLoggerService{
		logger:          logger.Helper(1),
		enabled:         true,
		serviceURL:      auditLogsURL,
		serviceHeaders:  headers,
		jsonWrapperKey:  os.Getenv("AUDIT_LOGS_FORWARDER_JSON_WRAPPER_KEY"),
		environmentName: env,
		hostname:        hostname,
		localIP:         getLocalIP(),
	}
	return &auditLogger, nil
}

func (l *AuditLoggerService) Audit(ctx context.Context, eventID EventID, data map[string]interface{}) {
	if !l.enabled {
		return
	}
	l.postLogToLogService(eventID, data)
}

func (l *AuditLoggerService) postLogToLogService(eventID EventID, data map[string]interface{}) {
	// Audit log JSON data
	logItem := map[string]interface{}{
		"eventID":  eventID,
		"hostname": l.hostname,
		"localIP":  l.localIP,
		"env":      l.environmentName,
		"data":     data,
	}

	// Optionally wrap audit log data into JSON object to help dynamically structure for an HTTP log service call
	if l.jsonWrapperKey != "" {
		logItem = map[string]interface{}{l.jsonWrapperKey: logItem}
	}

	serializedLog, err := json.Marshal(logItem)
	if err != nil {
		l.logger.Errorw("Unable to serialize wrapped audit log item to JSON", "err", err, "logItem", logItem)
		return
	}

	// Send up to HEC log collector
	httpClient := &http.Client{Timeout: time.Second * 10}
	req, _ := http.NewRequest("POST", l.serviceURL, bytes.NewReader(serializedLog))
	for _, header := range l.serviceHeaders {
		req.Header.Add(header.header, header.value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		l.logger.Errorw("Failed to send audit log to HTTP log service", "err", err, "logItem", logItem)
		return
	}
	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			l.logger.Errorw("Error reading errored HTTP log service webhook response body", "err", err, "logItem", logItem)
			return
		}
		l.logger.Errorw("Error sending log to HTTP log service", "statusCode", resp.StatusCode, "bodyString", string(bodyBytes))
		return
	}
}

// getLocalIP returns the first non-loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// filter and return address types for first non loopback address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
