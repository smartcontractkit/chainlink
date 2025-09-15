package common

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// ExtraDataCodecMap is a map of chain family to SourceChainExtraDataCodec
// Deprecated: use ccipocr3.ExtraDataCodecMap or ExtraDataCodecRegistry to manage codecs instead
type ExtraDataCodecMap = ccipocr3.ExtraDataCodecMap
