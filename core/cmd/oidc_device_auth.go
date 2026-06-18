package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

// OIDC device authorization grant client for the CLI.
//
// When a node is configured with AuthenticationMethod = 'oidc', the local
// users table only holds break-glass admins; interactive operators authenticate
// through the identity provider. The CLI cannot drive a browser redirect, so it
// uses the RFC 8628 device flow the node brokers at /oidc-device/start and
// /oidc-device/poll. The CLI never talks to the identity provider directly and
// never holds any token; it receives only an opaque handle and, on success, the
// same session cookie the browser flow produces.

// oidcDeviceStartResponse mirrors oidcauth.DeviceStartResponse.
type oidcDeviceStartResponse struct {
	DeviceHandle            string `json:"device_handle"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int64  `json:"expires_in"`
	Interval                int64  `json:"interval"`
}

// oidcDevicePollResponse mirrors oidcauth.DevicePollResponse.
type oidcDevicePollResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// OIDCDeviceCookieAuthenticator obtains a session cookie via the node-brokered
// device authorization flow and persists it to the same cookie store the
// password authenticator uses.
type OIDCDeviceCookieAuthenticator struct {
	config ClientOpts
	store  CookieStore
	lggr   logger.SugaredLogger
	// out is where user-facing prompts are written. Defaults to os.Stdout in
	// production; overridable in tests.
	out io.Writer
}

func NewOIDCDeviceCookieAuthenticator(config ClientOpts, store CookieStore, out io.Writer, lggr logger.Logger) *OIDCDeviceCookieAuthenticator {
	return &OIDCDeviceCookieAuthenticator{
		config: config,
		store:  store,
		lggr:   logger.Sugared(lggr),
		out:    out,
	}
}

// NodeHasOIDCEnabled probes the node's unauthenticated /oidc-enabled endpoint to
// decide whether the device flow applies. A non-OIDC node 404s the route.
func NodeHasOIDCEnabled(ctx context.Context, config ClientOpts, lggr logger.Logger) bool {
	client := newHttpClient(lggr, config.InsecureSkipVerify)
	u := config.RemoteNodeURL.String() + "/oidc-enabled"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// Login runs the full device flow: start, prompt the operator, poll until the
// node reports completion, then save the returned session cookie.
func (o *OIDCDeviceCookieAuthenticator) Login(ctx context.Context) error {
	start, err := o.start(ctx)
	if err != nil {
		return err
	}

	o.promptUser(start)

	interval := time.Duration(start.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}
	deadline := time.Now().Add(time.Duration(start.ExpiresIn) * time.Second)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return errors.New("device authorization timed out before approval")
			}
			cookie, status, perr := o.poll(ctx, start.DeviceHandle)
			if perr != nil {
				return perr
			}
			switch status {
			case "complete":
				if cookie == nil {
					return errors.New("node reported login complete but returned no session cookie")
				}
				return o.store.Save(cookie)
			case "denied":
				return errors.New("device authorization was denied or expired")
			default: // "pending"
				continue
			}
		}
	}
}

func (o *OIDCDeviceCookieAuthenticator) start(ctx context.Context) (*oidcDeviceStartResponse, error) {
	u := o.config.RemoteNodeURL.String() + "/oidc-device/start"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}
	client := newHttpClient(o.lggr, o.config.InsecureSkipVerify)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to start device authorization (status %d): %s", resp.StatusCode, string(body))
	}
	var sr oidcDeviceStartResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, errors.Wrap(err, "failed to parse device authorization response")
	}
	return &sr, nil
}

// poll asks the node for the flow status. On "complete" the node sets the
// session cookie on the response, which is extracted and returned.
func (o *OIDCDeviceCookieAuthenticator) poll(ctx context.Context, handle string) (*http.Cookie, string, error) {
	b, err := json.Marshal(map[string]string{"device_handle": handle})
	if err != nil {
		return nil, "", err
	}
	u := o.config.RemoteNodeURL.String() + "/oidc-device/poll"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := newHttpClient(o.lggr, o.config.InsecureSkipVerify)
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, "denied", nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.Errorf("device poll failed (status %d): %s", resp.StatusCode, string(body))
	}
	var pr oidcDevicePollResponse
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, "", errors.Wrap(err, "failed to parse device poll response")
	}
	if pr.Status == "complete" {
		return web.FindSessionCookie(resp.Cookies()), "complete", nil
	}
	return nil, pr.Status, nil
}

func (o *OIDCDeviceCookieAuthenticator) promptUser(start *oidcDeviceStartResponse) {
	uri := start.VerificationURIComplete
	if uri == "" {
		uri = start.VerificationURI
	}
	fmt.Fprintln(o.out, "To log in, open the following URL in a browser and sign in with your identity provider:")
	fmt.Fprintf(o.out, "\n    %s\n\n", uri)
	if start.VerificationURIComplete == "" {
		fmt.Fprintf(o.out, "When prompted, enter the code: %s\n\n", start.UserCode)
	} else {
		fmt.Fprintf(o.out, "Verify the code shown matches: %s\n\n", start.UserCode)
	}
	fmt.Fprintln(o.out, "Waiting for approval...")
}
