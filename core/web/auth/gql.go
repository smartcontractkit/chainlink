package auth

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthenticateGQL middleware checks the session cookie for a user and sets it
// on the context if it exists. It is the responsiblity of each resolver to
// validate whether it requires an authenticated user.
//
// We currently only support GQL authentication by session cookie.
func AuthenticateGQL(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID, ok := session.Get(SessionIDKey).(string)
		if !ok {
			return
		}

		user, err := authenticator.AuthorizedUserWithSession(sessionID)
		if err != nil {
			return
		}

		c.Set(SessionUserKey, &user)
	}
}
