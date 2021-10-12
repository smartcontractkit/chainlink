package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func authError(web.AuthStorer, *gin.Context) error {
	return errors.New("random error")
}

func authFailure(web.AuthStorer, *gin.Context) error {
	return auth.ErrorAuthFailed
}

func authSuccess(web.AuthStorer, *gin.Context) error {
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
	store := userFindSuccesser{user: user}

	called := false
	router := gin.New()
	router.Use(web.RequireAuth(store, web.AuthenticateByToken))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(web.APIKey, cltest.APIKey)
	req.Header.Set(web.APISecret, cltest.APISecret)
	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusText(http.StatusOK), http.StatusText(w.Code))
}

func TestAuthenticateByToken_AuthFailed(t *testing.T) {
	store := userFindFailer{err: auth.ErrorAuthFailed}

	called := false
	router := gin.New()
	router.Use(web.RequireAuth(store, web.AuthenticateByToken))
	router.GET("/", func(c *gin.Context) {
		called = true
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(web.APIKey, cltest.APIKey)
	req.Header.Set(web.APISecret, "bad-secret")
	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), http.StatusText(w.Code))
}

func TestRequireAuth_NoneRequired(t *testing.T) {
	called := false
	var store web.AuthStorer
	router := gin.New()
	router.Use(web.RequireAuth(store))
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
	var store web.AuthStorer
	router := gin.New()
	router.Use(web.RequireAuth(store, authFailure))
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
	var store web.AuthStorer
	router := gin.New()
	router.Use(web.RequireAuth(store, authFailure, authSuccess))
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
	var store web.AuthStorer
	router := gin.New()
	router.Use(web.RequireAuth(store, authError, authSuccess))
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
