package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
)

const bufferCapacity = 2048
const webRequestTimeout = 10

type Data = map[string]any

type AuditLogger interface {
	services.ServiceCtx

	Audit(eventID EventID, data Data)
}

type AuditLoggerConfig struct {
	ForwardToUrl   *string
	Environment    *string
	JsonWrapperKey *string
	Headers        []ServiceHeader
}

type HTTPAuditLoggerInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewAuditLoggerConfig(forwardToUrl string, isDev bool, jsonWrapperKey string, encodedHeaders string) (AuditLoggerConfig, error) {
	if forwardToUrl == "" {
		return AuditLoggerConfig{}, errors.Errorf("No forwardToURL provided")
	}

	if _, err := models.ParseURL(forwardToUrl); err != nil {
		return AuditLoggerConfig{}, errors.Errorf("forwardToURL value is not a valid URL")
	}

	environment := "production"
	if isDev {
		environment = "develop"
	}

	// Split and prepare optional service client headers from env variable
	headers := []ServiceHeader{}
	if encodedHeaders != "" {
		headerLines := strings.Split(encodedHeaders, "\\")
		for _, header := range headerLines {
			keyValue := strings.Split(header, "||")
			if len(keyValue) != 2 {
				return AuditLoggerConfig{}, errors.Errorf("Invalid headers provided for the audit logger. Value, single pair split on || required, got: %s", keyValue)
			}
			headers = append(headers, ServiceHeader{
				Header: keyValue[0],
				Value:  keyValue[1],
			})
		}
	}

	return AuditLoggerConfig{
		ForwardToUrl:   &forwardToUrl,
		Environment:    &environment,
		JsonWrapperKey: &jsonWrapperKey,
		Headers:        headers,
	}, nil
}

type AuditLoggerService struct {
	logger          logger.Logger            // The standard logger configured in the node
	enabled         bool                     // Whether the audit logger is enabled or not
	forwardToUrl    string                   // Location we are going to send logs to
	headers         []ServiceHeader          // Headers to be sent along with logs for identification/authentication
	jsonWrapperKey  string                   // Wrap audit data as a map under this key if present
	environmentName string                   // Decorate the environment this is coming from
	hostname        string                   // The self-reported hostname of the machine
	localIP         string                   // A non-loopback IP address as reported by the machine
	loggingClient   HTTPAuditLoggerInterface // Abstract type for sending logs onward

	loggingChannel chan wrappedAuditLog
	ctx            context.Context
	cancel         context.CancelFunc
	chDone         chan struct{}
}

// Configurable headers to include in POST to log service
type ServiceHeader struct {
	Header string
	Value  string
}

type wrappedAuditLog struct {
	eventID EventID
	data    Data
}

// NewAuditLogger returns a buffer push system that ingests audit log events and
// asynchronously pushes them up to an HTTP log service.
// Parses and validates the AUDIT_LOGS_* environment values and returns an enabled
// AuditLogger instance. If the environment variables are not set, the logger
// is disabled and short circuits execution via enabled flag.
func NewAuditLogger(logger logger.Logger, config *AuditLoggerConfig) (AuditLogger, error) {
	if config == nil {
		return &AuditLoggerService{}, errors.Errorf("Audit Log initialization error - no configuration")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return &AuditLoggerService{}, errors.Errorf("Audit Log initialization error - unable to get hostname: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	loggingChannel := make(chan wrappedAuditLog, bufferCapacity)

	// Create new AuditLoggerService
	auditLogger := AuditLoggerService{
		logger:          logger.Helper(1),
		enabled:         true,
		forwardToUrl:    *config.ForwardToUrl,
		headers:         config.Headers,
		jsonWrapperKey:  *config.JsonWrapperKey,
		environmentName: *config.Environment,
		hostname:        hostname,
		localIP:         getLocalIP(),
		loggingClient:   &http.Client{Timeout: time.Second * webRequestTimeout},

		loggingChannel: loggingChannel,
		ctx:            ctx,
		cancel:         cancel,
		chDone:         make(chan struct{}),
	}

	return &auditLogger, nil
}

func (l *AuditLoggerService) SetLoggingClient(newClient HTTPAuditLoggerInterface) {
	l.loggingClient = newClient
}

// Entrypoint for new audit logs. This buffers all logs that come in they will
// sent out by the goroutine that was started when the AuditLoggerService was
// created. If this service was not enabled, this immeidately returns.
//
// This function never blocks.
func (l *AuditLoggerService) Audit(eventID EventID, data Data) {
	fmt.Println("An audit log is being sent!")
	if !l.enabled {
		fmt.Println("Audit logger is not enabled?")
		return
	}

	wrappedLog := wrappedAuditLog{
		eventID: eventID,
		data:    data,
	}

	l.logger.Errorf("SEnding!!!")

	select {
	case l.loggingChannel <- wrappedLog:
	default:
		if l.loggingChannel == nil {
			l.logger.Errorw("Could not send log to audit subsystem because it has gone away!")
		} else {
			l.logger.Errorw("Audit log buffer is full. Dropping log with eventID: %s", eventID)
		}
	}
}

// Start the audit logger and begin processing logs on the channel
func (l *AuditLoggerService) Start(context.Context) error {
	if !l.enabled {
		return errors.Errorf("The audit logger is not enabled")
	}
	go l.runLoop()
	fmt.Println("Started the runloop")
	return nil
}

// Stops the logger and will close the channel.
func (l *AuditLoggerService) Close() error {
	if !l.enabled {
		return errors.Errorf("The audit logger is not enabled")
	}

	l.logger.Warnf("Disabled the audit logger service")
	l.cancel()
	<-l.chDone

	return nil
}

func (l *AuditLoggerService) Healthy() error {
	if !l.enabled {
		return errors.Errorf("The audit logger is not enabled")
	}

	if len(l.loggingChannel) == bufferCapacity {
		return errors.Errorf("The audit log buffer is full")
	}

	return nil
}

func (l *AuditLoggerService) Ready() error {
	if !l.enabled {
		return errors.Errorf("The audit logger is not enabled")
	}

	return nil
}

// Entrypoint for our log handling goroutine. This waits on the channel and sends out
// logs as they come in.
//
// This function calls postLogToLogService which blocks.
func (l *AuditLoggerService) runLoop() {
	defer close(l.chDone)

	for {
		select {
		case <-l.ctx.Done():
			// I've made this an error since we expect it should never happen.
			l.logger.Errorf("The audit logger has been requested to shut down!")
			return
		case event := <-l.loggingChannel:
			l.postLogToLogService(event.eventID, event.data)
		}
	}
}

// Takes an EventID and associated data and sends it to the configured logging
// endpoint. This function blocks on the send by timesout after a period of
// several seconds. This helps us prevent getting stuck on a single log
// due to transient network errors.
//
// This function blocks when called.
func (l *AuditLoggerService) postLogToLogService(eventID EventID, data Data) {
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

	// Send to remote service
	req, err := http.NewRequest("POST", l.forwardToUrl, bytes.NewReader(serializedLog))
	if err != nil {
		l.logger.Errorf("Failed to create request to remote logging service!")
	}
	for _, header := range l.headers {
		req.Header.Add(header.Header, header.Value)
	}
	resp, err := l.loggingClient.Do(req)
	if err != nil {
		l.logger.Errorw("Failed to send audit log to HTTP log service", "err", err, "logItem", logItem)
		return
	}
	if resp.StatusCode != 200 {
		if resp.Body == nil {
			l.logger.Errorw("There was no body to read. Possibly an error occured sending", "logItem", logItem)
			return
		}

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
