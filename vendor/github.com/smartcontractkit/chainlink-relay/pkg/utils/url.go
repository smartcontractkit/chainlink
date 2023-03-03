package utils

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

func (u *URL) MarshalText() ([]byte, error) {
	return []byte((*url.URL)(u).String()), nil
}

func (u *URL) UnmarshalText(input []byte) error {
	v, err := url.Parse(string(input))
	if err != nil {
		return err
	}
	*u = URL(*v)
	return nil
}
