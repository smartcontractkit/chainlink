package doer

import "net/http"

type authed struct {
	cookie  string
	wrapped *http.Client
}

func NewAuthed(cookie string) *authed {
	return &authed{
		cookie:  cookie,
		wrapped: http.DefaultClient,
	}
}

func (a *authed) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("cookie", a.cookie)
	return a.wrapped.Do(req)
}
