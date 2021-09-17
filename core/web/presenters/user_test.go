package presenters

import (
	"testing"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserResource(t *testing.T) {
	var (
		ts = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	user := sessions.User{
		Email:     "notreal@fakeemail.ch",
		CreatedAt: ts,
	}

	r := NewUserResource(user)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := `
	{
		"data": {
		   "type": "users",
		   "id": "notreal@fakeemail.ch",
		   "attributes": {
			  "email": "notreal@fakeemail.ch",
			  "createdAt": "2000-01-01T00:00:00Z"
		   }
		}
	 }
	`

	assert.JSONEq(t, expected, string(b))
}
