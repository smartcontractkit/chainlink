package doer

import "net/http"

type Authed struct {
	cookie  string
	wrapped *http.Client
}

func NewAuthed(cookie string) *Authed {
	return &Authed{
		cookie:  cookie,
		wrapped: http.DefaultClient,
	}
}

func (a *Authed) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("cookie", a.cookie)
	return a.wrapped.Do(req)
}
