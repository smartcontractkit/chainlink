package web_test

import (
	"chainlink/core/auth"
	"chainlink/core/internal/cltest"
	"chainlink/core/store"
	"chainlink/core/web"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func authError(*store.Store, *gin.Context) error {
	return errors.New("random error")
}

func authFailure(*store.Store, *gin.Context) error {
	return auth.ErrorAuthFailed
}

func authSuccess(*store.Store, *gin.Context) error {
	return nil
}

func TestAuthenticateByToken_Success(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	require.NoError(t, app.Start())
	app.Start()
	app.MustSeedUserAPIKey()

	called := false
	router := gin.New()
	router.Use(web.RequireAuth(app.GetStore(), web.AuthenticateByToken))
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
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	require.NoError(t, app.Start())
	app.Start()
	app.MustSeedUserAPIKey()

	called := false
	router := gin.New()
	router.Use(web.RequireAuth(app.GetStore(), web.AuthenticateByToken))
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
	var store *store.Store
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
	var store *store.Store
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
	var store *store.Store
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
	var store *store.Store
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
