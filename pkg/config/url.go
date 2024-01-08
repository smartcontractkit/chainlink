package config

import "net/url"

// URL extends url.URL to implement encoding.TextMarshaler.
type URL url.URL

func ParseURL(s string) (*URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return (*URL)(u), nil
}

func MustParseURL(s string) *URL {
	u, err := ParseURL(s)
	if err != nil {
		panic(err)
	}
	return u
}

func (u *URL) String() string {
	return (*url.URL)(u).String()
}

// URL returns a copy of u as a *url.URL
func (u *URL) URL() *url.URL {
	if u == nil {
		return nil
	}
	// defensive copy
	r := url.URL(*u)
	if u.User != nil {
		r.User = new(url.Userinfo)
		*r.User = *u.User
	}
	return &r
}

func (u *URL) IsZero() bool {
	return (url.URL)(*u) == url.URL{}
}

func (u *URL) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

func (u *URL) UnmarshalText(input []byte) error {
	v, err := url.Parse(string(input))
	if err != nil {
		return err
	}
	*u = URL(*v)
	return nil
}
