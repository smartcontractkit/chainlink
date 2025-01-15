package auth_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	webauth "github.com/smartcontractkit/chainlink/v2/core/web/auth"
)

func authError(*gin.Context, webauth.Authenticator) error {
	return errors.New("random error")
}

func authFailure(*gin.Context, webauth.Authenticator) error {
	return auth.ErrorAuthFailed
}

func authSuccess(*gin.Context, webauth.Authenticator) error {
	return nil
}

type userFindFailer struct {
	sessions.AuthenticationProvider
	err error
}

func (u userFindFailer) FindUser(ctx context.Context, email string) (sessions.User, error) {
	return sessions.User{}, u.err
}

func (u userFindFailer) FindUserByAPIToken(ctx context.Context, token string) (sessions.User, error) {
	return sessions.User{}, u.err
}

type userFindSuccesser struct {
	sessions.AuthenticationProvider
	user sessions.User
}

func (u userFindSuccesser) FindUser(ctx context.Context, email string) (sessions.User, error) {
	return u.user, nil
}

func (u userFindSuccesser) FindUserByAPIToken(ctx context.Context, token string) (sessions.User, error) {
	return u.user, nil
}

func TestAuthenticateByToken_Success(t *testing.T) {
	user := cltest.MustRandomUser(t)
	key, secret := uuid.New().String(), uuid.New().String()
	apiToken := auth.Token{AccessKey: key, Secret: secret}
	err := user.SetAuthToken(&apiToken)
	require.NoError(t, err)
	authr := userFindSuccesser{user: user}

	called := false
	router := gin.New()
	router.Use(webauth.Authenticate(authr, webauth.AuthenticateByToken))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	req.Header.Set(webauth.APIKey, key)
	req.Header.Set(webauth.APISecret, secret)
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusText(http.StatusOK), http.StatusText(w.Code))
}

func TestAuthenticateByToken_AuthFailed(t *testing.T) {
	authr := userFindFailer{err: auth.ErrorAuthFailed}

	called := false
	router := gin.New()
	router.Use(webauth.Authenticate(authr, webauth.AuthenticateByToken))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	req.Header.Set(webauth.APIKey, "bad-key")
	req.Header.Set(webauth.APISecret, "bad-secret")
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), http.StatusText(w.Code))
}

func TestAuthenticateByToken_RejectsBlankAccessKey(t *testing.T) {
	user := cltest.MustRandomUser(t)
	key, secret := "", uuid.New().String()
	apiToken := auth.Token{AccessKey: key, Secret: secret}
	err := user.SetAuthToken(&apiToken)
	require.NoError(t, err)
	authr := userFindSuccesser{user: user}

	called := false
	router := gin.New()
	router.Use(webauth.Authenticate(authr, webauth.AuthenticateByToken))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	req.Header.Set(webauth.APIKey, key)
	req.Header.Set(webauth.APISecret, secret)
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), http.StatusText(w.Code))
}

func TestRequireAuth_NoneRequired(t *testing.T) {
	called := false
	var authr webauth.Authenticator

	router := gin.New()
	router.Use(webauth.Authenticate(authr))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusText(http.StatusOK), http.StatusText(w.Code))
}

func TestRequireAuth_AuthFailed(t *testing.T) {
	called := false
	var authr webauth.Authenticator
	router := gin.New()
	router.Use(webauth.Authenticate(authr, authFailure))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), http.StatusText(w.Code))
}

func TestRequireAuth_LastAuthSuccess(t *testing.T) {
	called := false
	var authr webauth.Authenticator
	router := gin.New()
	router.Use(webauth.Authenticate(authr, authFailure, authSuccess))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusText(http.StatusOK), http.StatusText(w.Code))
}

func TestRequireAuth_Error(t *testing.T) {
	called := false
	var authr webauth.Authenticator
	router := gin.New()
	router.Use(webauth.Authenticate(authr, authError, authSuccess))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), http.StatusText(w.Code))
}

// Test RBAC (Role based access control) of each route and their required user roles
// Admin is omitted from the fields here since admin should be able to access all routes
type routeRules struct {
	verb               string
	path               string
	viewOnlyAllowed    bool
	editMinimalAllowed bool
	EditAllowed        bool
}

// The following are admin only routes
var routesRolesMap = [...]routeRules{
	{"GET", "/v2/users", false, false, false},
	{"POST", "/v2/users", false, false, false},
	{"PATCH", "/v2/users", false, false, false},
	{"DELETE", "/v2/users/MOCK", false, false, false},
	{"PATCH", "/v2/user/password", true, true, true},
	{"POST", "/v2/user/token", true, true, true},
	{"POST", "/v2/user/token/delete", true, true, true},
	{"GET", "/v2/enroll_webauthn", true, true, true},
	{"POST", "/v2/enroll_webauthn", true, true, true},
	{"GET", "/v2/external_initiators", true, true, true},
	{"POST", "/v2/external_initiators", false, false, true},
	{"DELETE", "/v2/external_initiators/MOCK", false, false, true},
	{"GET", "/v2/bridge_types", true, true, true},
	{"POST", "/v2/bridge_types", false, false, true},
	{"GET", "/v2/bridge_types/MOCK", true, true, true},
	{"PATCH", "/v2/bridge_types/MOCK", false, false, true},
	{"DELETE", "/v2/bridge_types/MOCK", false, false, true},
	{"POST", "/v2/transfers", false, false, false},
	{"POST", "/v2/transfers/evm", false, false, false},
	{"POST", "/v2/transfers/cosmos", false, false, false},
	{"POST", "/v2/transfers/solana", false, false, false},
	{"GET", "/v2/config", true, true, true},
	{"GET", "/v2/config/v2", true, true, true},
	{"GET", "/v2/tx_attempts", true, true, true},
	{"GET", "/v2/tx_attempts/evm", true, true, true},
	{"GET", "/v2/transactions/evm", true, true, true},
	{"GET", "/v2/transactions/evm/MOCK", true, true, true},
	{"GET", "/v2/transactions", true, true, true},
	{"GET", "/v2/transactions/MOCK", true, true, true},
	{"POST", "/v2/replay_from_block/MOCK", false, true, true},
	{"GET", "/v2/keys/csa", true, true, true},
	{"POST", "/v2/keys/csa", false, false, true},
	{"POST", "/v2/keys/csa/import", false, false, false},
	{"POST", "/v2/keys/csa/export/MOCK", false, false, false},
	{"GET", "/v2/keys/eth", true, true, true},
	{"POST", "/v2/keys/eth", false, false, true},
	{"DELETE", "/v2/keys/eth/MOCK", false, false, false},
	{"POST", "/v2/keys/eth/import", false, false, false},
	{"POST", "/v2/keys/eth/export/MOCK", false, false, false},
	{"GET", "/v2/keys/ocr", true, true, true},
	{"POST", "/v2/keys/ocr", false, false, true},
	{"DELETE", "/v2/keys/ocr/:MOCKkeyID", false, false, false},
	{"POST", "/v2/keys/ocr/import", false, false, false},
	{"POST", "/v2/keys/ocr/export/MOCK", false, false, false},
	{"GET", "/v2/keys/ocr2", true, true, true},
	{"POST", "/v2/keys/ocr2/MOCK", false, false, true},
	{"DELETE", "/v2/keys/ocr2/MOCK", false, false, false},
	{"POST", "/v2/keys/ocr2/import", false, false, false},
	{"POST", "/v2/keys/ocr2/export/MOCK", false, false, false},
	{"GET", "/v2/keys/p2p", true, true, true},
	{"POST", "/v2/keys/p2p", false, false, true},
	{"DELETE", "/v2/keys/p2p/MOCK", false, false, false},
	{"POST", "/v2/keys/p2p/import", false, false, false},
	{"POST", "/v2/keys/p2p/export/MOCK", false, false, false},
	{"GET", "/v2/keys/solana", true, true, true},
	{"GET", "/v2/keys/cosmos", true, true, true},
	{"GET", "/v2/keys/starknet", true, true, true},
	{"GET", "/v2/keys/aptos", true, true, true},
	{"GET", "/v2/keys/tron", true, true, true},
	{"POST", "/v2/keys/solana", false, false, true},
	{"POST", "/v2/keys/cosmos", false, false, true},
	{"POST", "/v2/keys/starknet", false, false, true},
	{"POST", "/v2/keys/aptos", false, false, true},
	{"POST", "/v2/keys/tron", false, false, true},
	{"DELETE", "/v2/keys/solana/MOCK", false, false, false},
	{"DELETE", "/v2/keys/cosmos/MOCK", false, false, false},
	{"DELETE", "/v2/keys/starknet/MOCK", false, false, false},
	{"DELETE", "/v2/keys/aptos/MOCK", false, false, false},
	{"DELETE", "/v2/keys/tron/MOCK", false, false, false},
	{"POST", "/v2/keys/solana/import", false, false, false},
	{"POST", "/v2/keys/cosmos/import", false, false, false},
	{"POST", "/v2/keys/starknet/import", false, false, false},
	{"POST", "/v2/keys/aptos/import", false, false, false},
	{"POST", "/v2/keys/tron/import", false, false, false},
	{"POST", "/v2/keys/solana/export/MOCK", false, false, false},
	{"POST", "/v2/keys/cosmos/export/MOCK", false, false, false},
	{"POST", "/v2/keys/starknet/export/MOCK", false, false, false},
	{"POST", "/v2/keys/aptos/export/MOCK", false, false, false},
	{"POST", "/v2/keys/tron/export/MOCK", false, false, false},
	{"GET", "/v2/keys/vrf", true, true, true},
	{"POST", "/v2/keys/vrf", false, false, true},
	{"DELETE", "/v2/keys/vrf/MOCK", false, false, false},
	{"POST", "/v2/keys/vrf/import", false, false, false},
	{"POST", "/v2/keys/vrf/export/MOCK", false, false, false},
	{"GET", "/v2/jobs", true, true, true},
	{"GET", "/v2/jobs/MOCK", true, true, true},
	{"POST", "/v2/jobs", false, false, true},
	{"DELETE", "/v2/jobs/MOCK", false, false, true},
	{"GET", "/v2/pipeline/runs", true, true, true},
	{"GET", "/v2/jobs/MOCK/runs", true, true, true},
	{"GET", "/v2/jobs/MOCK/runs/MOCK", true, true, true},
	{"GET", "/v2/features", true, true, true},
	{"DELETE", "/v2/pipeline/job_spec_errors/MOCK", false, false, true},
	{"GET", "/v2/log", true, true, true},
	{"PATCH", "/v2/log", false, false, false},
	{"GET", "/v2/chains/evm", true, true, true},
	{"GET", "/v2/chains/solana", true, true, true},
	{"GET", "/v2/chains/cosmos", true, true, true},
	{"GET", "/v2/chains/evm/MOCK", true, true, true},
	{"GET", "/v2/chains/cosmos/MOCK", true, true, true},
	{"GET", "/v2/nodes/", true, true, true},
	{"GET", "/v2/nodes/evm", true, true, true},
	{"GET", "/v2/nodes/solana", true, true, true},
	{"GET", "/v2/nodes/cosmos", true, true, true},
	{"GET", "/v2/chains/evm/MOCK/nodes", true, true, true},
	{"GET", "/v2/chains/solana/MOCK/nodes", true, true, true},
	{"GET", "/v2/chains/cosmos/MOCK/nodes", true, true, true},
	{"GET", "/v2/nodes/evm/forwarders", true, true, true},
	{"POST", "/v2/nodes/evm/forwarders/track", false, false, true},
	{"DELETE", "/v2/nodes/evm/forwarders/MOCK", false, false, true},
	{"GET", "/v2/build_info", true, true, true},
	{"GET", "/v2/ping", true, true, true},
	{"POST", "/v2/jobs/MOCK/runs", false, true, true},
}

// The following test implementations work by asserting only that "Unauthorized/Forbidden" errors are not returned (success case),
// because hitting the handler are not mocked and will crash as expected
// Iterate over the above routesRolesMap and assert each path is wrapped and
// the user role is enforced with the correct middleware
func TestRBAC_Routemap_Admin(t *testing.T) {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Assert all admin routes
	// no endpoint should return StatusUnauthorized
	client := app.NewHTTPClient(nil)
	for _, route := range routesRolesMap {
		func() {
			var resp *http.Response
			var cleanup func()

			switch route.verb {
			case "GET":
				resp, cleanup = client.Get(route.path)
			case "POST":
				resp, cleanup = client.Post(route.path, nil)
			case "DELETE":
				resp, cleanup = client.Delete(route.path)
			case "PATCH":
				resp, cleanup = client.Patch(route.path, nil)
			case "PUT":
				resp, cleanup = client.Put(route.path, nil)
			default:
				t.Fatalf("Unknown HTTP verb %s\n", route.verb)
			}
			defer cleanup()

			assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
			assert.NotEqual(t, http.StatusForbidden, resp.StatusCode)
		}()
	}
}

func TestRBAC_Routemap_Edit(t *testing.T) {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create a test edit user to work with
	u := &cltest.User{Role: sessions.UserRoleEdit}
	client := app.NewHTTPClient(u)

	// Assert all edit routes
	for _, route := range routesRolesMap {
		func() {
			var resp *http.Response
			var cleanup func()

			switch route.verb {
			case "GET":
				resp, cleanup = client.Get(route.path)
			case "POST":
				resp, cleanup = client.Post(route.path, nil)
			case "DELETE":
				resp, cleanup = client.Delete(route.path)
			case "PATCH":
				resp, cleanup = client.Patch(route.path, nil)
			case "PUT":
				resp, cleanup = client.Put(route.path, nil)
			default:
				t.Fatalf("Unknown HTTP verb %s\n", route.verb)
			}
			defer cleanup()

			// If this route allows up to an edit role, don't expect an unauthorized response
			if route.EditAllowed || route.editMinimalAllowed || route.viewOnlyAllowed {
				assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
				assert.NotEqual(t, http.StatusForbidden, resp.StatusCode)
			} else if !route.EditAllowed {
				assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			}
		}()
	}
}

func TestRBAC_Routemap_Run(t *testing.T) {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create a test run user to work with
	u := &cltest.User{Role: sessions.UserRoleRun}
	client := app.NewHTTPClient(u)

	// Assert all run routes
	for _, route := range routesRolesMap {
		func() {
			var resp *http.Response
			var cleanup func()

			switch route.verb {
			case "GET":
				resp, cleanup = client.Get(route.path)
			case "POST":
				resp, cleanup = client.Post(route.path, nil)
			case "DELETE":
				resp, cleanup = client.Delete(route.path)
			case "PATCH":
				resp, cleanup = client.Patch(route.path, nil)
			case "PUT":
				resp, cleanup = client.Put(route.path, nil)
			default:
				t.Fatalf("Unknown HTTP verb %s\n", route.verb)
			}
			defer cleanup()

			// If this route allows up to an edit minimal role, don't expect an unauthorized response
			if route.editMinimalAllowed || route.viewOnlyAllowed {
				assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
				assert.NotEqual(t, http.StatusForbidden, resp.StatusCode)
			} else if !route.EditAllowed {
				assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			}
		}()
	}
}

func TestRBAC_Routemap_ViewOnly(t *testing.T) {
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	router := web.Router(t, app, nil)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create a test run user to work with
	u := &cltest.User{Role: sessions.UserRoleView}
	client := app.NewHTTPClient(u)

	// Assert all view only routes
	for i, route := range routesRolesMap {
		route := route
		t.Run(fmt.Sprintf("%d-%s-%s", i, route.verb, route.path), func(t *testing.T) {
			var resp *http.Response
			var cleanup func()

			switch route.verb {
			case "GET":
				resp, cleanup = client.Get(route.path)
			case "POST":
				resp, cleanup = client.Post(route.path, nil)
			case "DELETE":
				resp, cleanup = client.Delete(route.path)
			case "PATCH":
				resp, cleanup = client.Patch(route.path, nil)
			case "PUT":
				resp, cleanup = client.Put(route.path, nil)
			default:
				t.Fatalf("Unknown HTTP verb %s\n", route.verb)
			}
			defer cleanup()

			// If this route only allows view only, don't expect an unauthorized response
			if route.viewOnlyAllowed {
				assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
				assert.NotEqual(t, http.StatusForbidden, resp.StatusCode)
			} else if !route.EditAllowed {
				assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			}
		})
	}
}

func mustRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	ctx := testutils.Context(t)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	require.NoError(t, err)
	return req
}
