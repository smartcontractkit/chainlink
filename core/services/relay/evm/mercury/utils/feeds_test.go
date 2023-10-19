package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	v1FeedId = (FeedID)([32]uint8{00, 01, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114})
	v2FeedId = (FeedID)([32]uint8{00, 02, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114})
	v3FeedId = (FeedID)([32]uint8{00, 03, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114})
)

func Test_FeedID_Version(t *testing.T) {
	t.Run("versioned feed ID", func(t *testing.T) {
		assert.Equal(t, REPORT_V1, v1FeedId.Version())
		assert.True(t, v1FeedId.IsV1())
		assert.False(t, v1FeedId.IsV2())
		assert.False(t, v1FeedId.IsV3())

		assert.Equal(t, REPORT_V2, v2FeedId.Version())
		assert.False(t, v2FeedId.IsV1())
		assert.True(t, v2FeedId.IsV2())
		assert.False(t, v2FeedId.IsV3())

		assert.Equal(t, REPORT_V3, v3FeedId.Version())
		assert.False(t, v3FeedId.IsV1())
		assert.False(t, v3FeedId.IsV2())
		assert.True(t, v3FeedId.IsV3())
	})
	t.Run("legacy special cases", func(t *testing.T) {
		for _, feedID := range legacyV1FeedIDs {
			assert.Equal(t, REPORT_V1, feedID.Version())
		}
	})
}
