package models

import (
	"encoding"
	"fmt"
	"net/url"
)

/*
	Use NewNonRedactedSecret or NewNonRedactedSecretURL when you need to marshal your secrets as it
	Used in docker e2e tests
*/

type MaybeRedactedSecret interface {
	fmt.Stringer
	encoding.TextMarshaler
}

var (
	_ fmt.Stringer           = (*NonRedactedSecret)(nil)
	_ encoding.TextMarshaler = (*NonRedactedSecret)(nil)
)

type NonRedactedSecret string

func NewNonRedactedSecret(s string) *NonRedactedSecret { return (*NonRedactedSecret)(&s) }

func (s NonRedactedSecret) String() string { return string(s) }

func (s NonRedactedSecret) GoString() string { return string(s) }

func (s NonRedactedSecret) MarshalText() ([]byte, error) { return []byte(s), nil }

type MaybeRedactedSecretURL interface {
	fmt.Stringer
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	URL() *url.URL
}

var (
	_ fmt.Stringer             = (*NonRedactedSecretURL)(nil)
	_ encoding.TextMarshaler   = (*NonRedactedSecretURL)(nil)
	_ encoding.TextUnmarshaler = (*NonRedactedSecretURL)(nil)
)

type NonRedactedSecretURL URL

func NewNonRedactedSecretURL(u *URL) *NonRedactedSecretURL { return (*NonRedactedSecretURL)(u) }

func MustNonRedactedSecretURL(u string) *NonRedactedSecretURL {
	return NewNonRedactedSecretURL(MustParseURL(u))
}

func (s *NonRedactedSecretURL) String() string { return (*URL)(s).String() }

func (s *NonRedactedSecretURL) GoString() string { return (*URL)(s).String() }

func (s *NonRedactedSecretURL) URL() *url.URL { return (*URL)(s).URL() }

func (s *NonRedactedSecretURL) MarshalText() ([]byte, error) { return []byte((*URL)(s).String()), nil }

func (s *NonRedactedSecretURL) UnmarshalText(text []byte) error {
	if err := (*URL)(s).UnmarshalText(text); err != nil {
		return fmt.Errorf("failed to parse url: %s", (*URL)(s).String())
	}
	return nil
}
