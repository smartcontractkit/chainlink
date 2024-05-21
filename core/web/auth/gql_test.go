package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/web/auth"
)

func Test_AuthenticateGQL_Unauthenticated(t *testing.T) {
	t.Parallel()

	sessionORM := mocks.NewAuthenticationProvider(t)
	sessionStore := cookie.NewStore([]byte("secret"))

	r := gin.Default()
	r.Use(sessions.Sessions(auth.SessionName, sessionStore))
	r.Use(auth.AuthenticateGQL(sessionORM, logger.TestLogger(t)))

	r.GET("/", func(c *gin.Context) {
		session, ok := auth.GetGQLAuthenticatedSession(c)
		assert.False(t, ok)
		assert.Nil(t, session)

		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	r.ServeHTTP(w, req)
}

func Test_AuthenticateGQL_Authenticated(t *testing.T) {
	t.Parallel()

	sessionORM := mocks.NewAuthenticationProvider(t)
	sessionStore := cookie.NewStore([]byte(cltest.SessionSecret))
	sessionID := "sessionID"

	r := gin.Default()
	r.Use(sessions.Sessions(auth.SessionName, sessionStore))
	r.Use(auth.AuthenticateGQL(sessionORM, logger.TestLogger(t)))

	r.GET("/", func(c *gin.Context) {
		session, ok := auth.GetGQLAuthenticatedSession(c.Request.Context())
		assert.True(t, ok)
		assert.NotNil(t, session)

		c.String(http.StatusOK, "")
	})

	sessionORM.On("AuthorizedUserWithSession", mock.Anything, sessionID).Return(clsessions.User{Email: cltest.APIEmailAdmin, Role: clsessions.UserRoleAdmin}, nil)

	w := httptest.NewRecorder()
	req := mustRequest(t, "GET", "/", nil)
	cookie := cltest.MustGenerateSessionCookie(t, sessionID)
	req.AddCookie(cookie)

	r.ServeHTTP(w, req)
}

func Test_GetAndSetGQLAuthenticatedSession(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	user := clsessions.User{Email: cltest.APIEmailAdmin, Role: clsessions.UserRoleAdmin}

	ctx = auth.WithGQLAuthenticatedSession(ctx, user, "sessionID")

	actual, ok := auth.GetGQLAuthenticatedSession(ctx)
	assert.True(t, ok)
	assert.Equal(t, &user, actual.User)
	assert.Equal(t, "sessionID", actual.SessionID)
}
