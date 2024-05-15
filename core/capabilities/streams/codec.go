package streams

import (
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Codec struct {
}

func (c Codec) Unwrap(raw values.Value) ([]mercury.FeedReport, error) {
	dest := []mercury.FeedReport{}
	err := raw.UnwrapTo(&dest)
	// TODO (KS-196): validate reports
	return dest, err
}

func (c Codec) Wrap(reports []mercury.FeedReport) (values.Value, error) {
	return values.Wrap(reports)
}

func NewCodec() Codec {
	return Codec{}
}
