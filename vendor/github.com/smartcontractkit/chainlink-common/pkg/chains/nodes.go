package chains

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// pageToken is simple internal representation for coordination requests and responses in a paginated API
// It is inspired by the Google API Design patterns
// https://cloud.google.com/apis/design/design_patterns#list_pagination
// https://google.aip.dev/158
type pageToken struct {
	Page int
	Size int
}

var (
	ErrInvalidToken = errors.New("invalid page token")
	ErrOutOfRange   = errors.New("out of range")
	defaultSize     = 100
)

// Encode the token in base64 for transmission for the wire
func (pr *pageToken) Encode() string {
	if pr.Size == 0 {
		pr.Size = defaultSize
	}
	// this is a simple minded implementation and may benefit from something fancier
	// note that this is a valid url.Query string, which we leverage in decoding
	s := fmt.Sprintf("page=%d&size=%d", pr.Page, pr.Size)
	return base64.RawStdEncoding.EncodeToString([]byte(s))
}

// b64enc must be the base64 encoded token string, corresponding to [pageToken.Encode()]
func NewPageToken(b64enc string) (*pageToken, error) {
	// empty is valid
	if b64enc == "" {
		return &pageToken{Page: 0, Size: defaultSize}, nil
	}

	b, err := base64.RawStdEncoding.DecodeString(b64enc)
	if err != nil {
		return nil, err
	}
	// here too, this is simple minded and could be fancier

	vals, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}
	if !(vals.Has("page") && vals.Has("size")) {
		return nil, ErrInvalidToken
	}
	page, err := strconv.Atoi(vals.Get("page"))
	if err != nil {
		return nil, fmt.Errorf("%w: bad page", ErrInvalidToken)
	}
	size, err := strconv.Atoi(vals.Get("size"))
	if err != nil {
		return nil, fmt.Errorf("%w: bad size", ErrInvalidToken)
	}
	return &pageToken{
		Page: page,
		Size: size,
	}, err
}

// if start is out of range, must return ErrOutOfRange
type ListNodeStatusFn = func(start, end int) (stats []types.NodeStatus, total int, err error)

func ListNodeStatuses(pageSize int, pageTokenStr string, listFn ListNodeStatusFn) (stats []types.NodeStatus, nextPageToken string, total int, err error) {
	if pageSize == 0 {
		pageSize = defaultSize
	}
	t := &pageToken{Page: 0, Size: pageSize}
	if pageTokenStr != "" {
		t, err = NewPageToken(pageTokenStr)
		if err != nil {
			return nil, "", -1, err
		}
	}
	start, end := t.Page*t.Size, (t.Page+1)*t.Size
	stats, total, err = listFn(start, end)
	if err != nil {
		return stats, "", -1, err
	}
	if total > end {
		nextToken := &pageToken{Page: t.Page + 1, Size: t.Size}
		nextPageToken = nextToken.Encode()
	}
	return stats, nextPageToken, total, nil
}
