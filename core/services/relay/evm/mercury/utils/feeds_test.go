package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	v1FeedID       = (FeedID)([32]uint8{00, 01, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114})
	v2FeedID       = (FeedID)([32]uint8{00, 02, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114})
	v3FeedID       = (FeedID)([32]uint8{00, 03, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114})
	keystonev2Feed = (FeedID)([32]uint8{01, 12, 34, 56, 78, 00, 02, 04, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00})
	keystonev3Feed = (FeedID)([32]uint8{01, 12, 34, 56, 78, 00, 03, 04, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00})
	keystonev4Feed = (FeedID)([32]uint8{01, 12, 34, 56, 78, 00, 04, 04, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00, 00})
)

func Test_FeedID_Version(t *testing.T) {
	t.Run("versioned feed ID", func(t *testing.T) {
		assert.Equal(t, REPORT_V1, v1FeedID.Version())
		assert.True(t, v1FeedID.IsV1())
		assert.False(t, v1FeedID.IsV2())
		assert.False(t, v1FeedID.IsV3())

		assert.Equal(t, REPORT_V2, v2FeedID.Version())
		assert.False(t, v2FeedID.IsV1())
		assert.True(t, v2FeedID.IsV2())
		assert.False(t, v2FeedID.IsV3())

		assert.Equal(t, REPORT_V3, v3FeedID.Version())
		assert.False(t, v3FeedID.IsV1())
		assert.False(t, v3FeedID.IsV2())
		assert.True(t, v3FeedID.IsV3())

		assert.Equal(t, REPORT_V2, keystonev2Feed.Version())
		assert.False(t, keystonev2Feed.IsV1())
		assert.True(t, keystonev2Feed.IsV2())
		assert.False(t, keystonev2Feed.IsV3())
		assert.False(t, keystonev2Feed.IsV4())

		assert.Equal(t, REPORT_V3, keystonev3Feed.Version())
		assert.False(t, keystonev3Feed.IsV1())
		assert.False(t, keystonev3Feed.IsV2())
		assert.True(t, keystonev3Feed.IsV3())
		assert.False(t, keystonev3Feed.IsV4())

		assert.Equal(t, REPORT_V4, keystonev4Feed.Version())
		assert.False(t, keystonev4Feed.IsV1())
		assert.False(t, keystonev4Feed.IsV2())
		assert.False(t, keystonev4Feed.IsV3())
		assert.True(t, keystonev4Feed.IsV4())
	})
	t.Run("legacy special cases", func(t *testing.T) {
		for _, feedID := range legacyV1FeedIDs {
			assert.Equal(t, REPORT_V1, feedID.Version())
		}
	})
}
