package web_test

import (
	"chainlink/core/auth"
	"chainlink/core/store"
	"chainlink/core/web"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func authError(_ *store.Store, _ *gin.Context) error {
	return errors.New("random error")
}

func authFailure(_ *store.Store, _ *gin.Context) error {
	return auth.ErrorAuthFailed
}

func authSuccess(_ *store.Store, _ *gin.Context) error {
	return nil
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
