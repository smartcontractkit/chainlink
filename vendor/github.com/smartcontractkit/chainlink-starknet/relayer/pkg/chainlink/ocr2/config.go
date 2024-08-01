package ocr2

import "time"

type Config interface {
	OCR2CachePollPeriod() time.Duration
	OCR2CacheTTL() time.Duration
}
