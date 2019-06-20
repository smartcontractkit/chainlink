package web

import (
	"net/http"
)

// FindSession returns the cookie for the chainlink session
func FindSession(cookies []*http.Cookie) *http.Cookie {
	for _, c := range cookies {
		if c.Name == "clsession" {
			return c
		}
	}

	return nil
}
