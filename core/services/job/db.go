package job

import (
	"time"
)

type (
	FetcherDBRow struct {
		ID        uint64 `gorm:"primary_key;auto_increment;"`
		CreatedAt time.Time
		UpdatedAt time.Time

		HttpFetcher   *HttpFetcherDBRow
		BridgeFetcher *BridgeFetcherDBRow
		MedianFetcher *MedianFetcherDBRow
	}

	HttpFetcherDBRow struct {
		ID           uint64 `gorm:"primary_key;auto_increment;"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		*HttpFetcher `gorm:"embedded;"`
		Transformers []*TransformerDBRow
	}

	BridgeFetcherDBRow struct {
		ID             uint64 `gorm:"primary_key;auto_increment;"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		*BridgeFetcher `gorm:"embedded;"`
		Transformers   []*TransformerDBRow
	}

	MedianFetcherDBRow struct {
		ID             uint64 `gorm:"primary_key;auto_increment;"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		*MedianFetcher `gorm:"embedded;"`
		Fetchers       []*FetcherDBRow `gorm:"foreignkey:parent_id"`
		Transformers   []*TransformerDBRow
	}
)

func (FetcherDBRow) TableName() string       { return "fetchers" }
func (HttpFetcherDBRow) TableName() string   { return "http_fetchers" }
func (BridgeFetcherDBRow) TableName() string { return "bridge_fetchers" }
func (MedianFetcherDBRow) TableName() string { return "median_fetchers" }

type (
	TransformerDBRow struct {
		ID                   uint64 `gorm:"primary_key;auto_increment"`
		CreatedAt            time.Time
		UpdatedAt            time.Time
		MultiplyTransformer  *MultiplyTransformerDBRow
		JSONParseTransformer *JSONParseTransformerDBRow
	}

	MultiplyTransformerDBRow struct {
		ID                   uint64 `gorm:"primary_key;auto_increment"`
		CreatedAt            time.Time
		UpdatedAt            time.Time
		*MultiplyTransformer `gorm:"embedded;"`
	}

	JSONParseTransformerDBRow struct {
		ID                    uint64 `gorm:"primary_key;auto_increment"`
		CreatedAt             time.Time
		UpdatedAt             time.Time
		*JSONParseTransformer `gorm:"embedded;"`
	}
)

func (TransformerDBRow) TableName() string          { return "transformers" }
func (MultiplyTransformerDBRow) TableName() string  { return "multiply_transformers" }
func (JSONParseTransformerDBRow) TableName() string { return "jsonparse_transformers" }

func WrapFetchersForDB(fetchers ...Fetcher) []*FetcherDBRow {
	var dbRows []*FetcherDBRow
	for _, fetcher := range fetchers {
		switch f := fetcher.(type) {
		case *HttpFetcher:
			ts := WrapTransformersForDB(f.Transformers...)
			dbRows = append(dbRows, &FetcherDBRow{HttpFetcher: &HttpFetcherDBRow{HttpFetcher: f, Transformers: ts}})
		case *BridgeFetcher:
			ts := WrapTransformersForDB(f.Transformers...)
			dbRows = append(dbRows, &FetcherDBRow{BridgeFetcher: &BridgeFetcherDBRow{BridgeFetcher: f, Transformers: ts}})
		case *MedianFetcher:
			ts := WrapTransformersForDB(f.Transformers...)
			dbRows = append(dbRows, &FetcherDBRow{MedianFetcher: &MedianFetcherDBRow{MedianFetcher: f, Fetchers: WrapFetchersForDB(f.Fetchers...), Transformers: ts}})
		}
	}
	return dbRows
}

func WrapTransformersForDB(transformers ...Transformer) []*TransformerDBRow {
	var dbRows []*TransformerDBRow
	for _, transformer := range transformers {
		switch t := transformer.(type) {
		case *MultiplyTransformer:
			dbRows = append(dbRows, &TransformerDBRow{MultiplyTransformer: &MultiplyTransformerDBRow{MultiplyTransformer: t}})
		case *JSONParseTransformer:
			dbRows = append(dbRows, &TransformerDBRow{JSONParseTransformer: &JSONParseTransformerDBRow{JSONParseTransformer: t}})
		}
	}
	return dbRows
}
