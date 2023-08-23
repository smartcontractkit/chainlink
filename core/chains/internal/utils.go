package internal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type PageToken struct {
	Page int
	Size int
}

var (
	ErrInvalidToken = errors.New("invalid page token")
	ErrOutOfRange   = errors.New("out of range")
)

func (pr *PageToken) Encode() string {
	// TODO something more sophisicated
	s := fmt.Sprintf("page=%d&size=%d", pr.Page, pr.Size)
	return base64.RawStdEncoding.EncodeToString([]byte(s))
}

// enc is base64 encoded token string
func NewPageToken(enc string) (*PageToken, error) {
	// empty is valid
	if enc == "" {
		return &PageToken{}, nil
	}

	b, err := base64.RawStdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}
	// todo regex check

	vals, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}
	if !(vals.Has("page") && vals.Has("size")) {
		return nil, ErrInvalidToken
	}
	page, err := strconv.Atoi(vals.Get("page"))
	if err != nil {
		return nil, fmt.Errorf("%w: bad page", &ErrInvalidToken)
	}
	size, err := strconv.Atoi(vals.Get("size"))
	if err != nil {
		return nil, fmt.Errorf("%w: bad size", &ErrInvalidToken)
	}
	return &PageToken{
		Page: page,
		Size: size,
	}, err
}

func ValidatePageToken(token string) (page int, err error) {

	return -1, fmt.Errorf("validate page token unimplemented")
}

// if start is out of range, must return ErrOutOfRange
type ListNodeStatusFn = func(start, end int) (stats []types.NodeStatus, total int, err error)

func ListNodeStatuses(page_size int, page_token string, listFn ListNodeStatusFn) (stats []types.NodeStatus, next_page_token string, err error) {
	page, err := ValidatePageToken(page_token)
	if err != nil {
		return nil, "", err
	}
	start, end := page*page_size, (page+1)*page_size
	stats, total, err := listFn(start, end)
	if err != nil {
		return stats, "", err
	}
	if total > end {
		t := &PageToken{Page: page + 1, Size: page_size}
		next_page_token = t.Encode()
	}
	return stats, next_page_token, nil
}
