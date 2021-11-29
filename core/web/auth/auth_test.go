package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/sessions"
	webauth "github.com/smartcontractkit/chainlink/core/web/auth"
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
	sessions.ORM
	err error
}

func (u userFindFailer) FindUser() (sessions.User, error) {
	return sessions.User{}, u.err
}

type userFindSuccesser struct {
	sessions.ORM
	user sessions.User
}

func (u userFindSuccesser) FindUser() (sessions.User, error) {
	return u.user, nil
}

func TestAuthenticateByToken_Success(t *testing.T) {
	user := cltest.MustRandomUser(t)
	apiToken := auth.Token{AccessKey: cltest.APIKey, Secret: cltest.APISecret}
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
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(webauth.APIKey, cltest.APIKey)
	req.Header.Set(webauth.APISecret, cltest.APISecret)
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
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(webauth.APIKey, cltest.APIKey)
	req.Header.Set(webauth.APISecret, "bad-secret")
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
	req, _ := http.NewRequest("GET", "/", nil)
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
	req, _ := http.NewRequest("GET", "/", nil)
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
	req, _ := http.NewRequest("GET", "/", nil)
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
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), http.StatusText(w.Code))
}
