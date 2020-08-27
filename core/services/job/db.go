package job

import (
	"time"
)

type (
	FetcherDBRow struct {
		ID         uint64 `gorm:"primary_key;auto_increment;"`
		CreatedAt  time.Time
		UpdatedAt  time.Time
		ParentID   uint64
		ParentType string

		HttpFetcher   *HttpFetcherDBRow   `gorm:"preload:true;foreignkey:fetcher_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		BridgeFetcher *BridgeFetcherDBRow `gorm:"preload:true;foreignkey:fetcher_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		MedianFetcher *MedianFetcherDBRow `gorm:"preload:true;foreignkey:fetcher_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
	}

	HttpFetcherDBRow struct {
		ID        uint64 `gorm:"primary_key;auto_increment;"`
		FetcherID uint64
		CreatedAt time.Time
		UpdatedAt time.Time

		*HttpFetcher `gorm:"embedded;"`
		Transformers []*TransformerDBRow `gorm:"foreignkey:fetcher_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
	}

	BridgeFetcherDBRow struct {
		ID        uint64 `gorm:"primary_key;auto_increment;"`
		FetcherID uint64
		CreatedAt time.Time
		UpdatedAt time.Time

		*BridgeFetcher `gorm:"embedded;"`
		Transformers   []*TransformerDBRow `gorm:"foreignkey:fetcher_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
	}

	MedianFetcherDBRow struct {
		ID        uint64 `gorm:"primary_key;auto_increment;"`
		FetcherID uint64
		CreatedAt time.Time
		UpdatedAt time.Time

		*MedianFetcher `gorm:"embedded;"`
		Fetchers       []*FetcherDBRow     `gorm:"foreignkey:parent_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		Transformers   []*TransformerDBRow `gorm:"foreignkey:fetcher_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
	}
)

func (FetcherDBRow) TableName() string       { return "fetchers" }
func (HttpFetcherDBRow) TableName() string   { return "http_fetchers" }
func (BridgeFetcherDBRow) TableName() string { return "bridge_fetchers" }
func (MedianFetcherDBRow) TableName() string { return "median_fetchers" }

type (
	TransformerDBRow struct {
		ID                   uint64 `gorm:"primary_key;auto_increment"`
		FetcherID            uint64
		CreatedAt            time.Time
		UpdatedAt            time.Time
		MultiplyTransformer  *MultiplyTransformerDBRow  `gorm:"foreignkey:transformer_id;association_autoupdate:true;association_autocreate:true"`
		JSONParseTransformer *JSONParseTransformerDBRow `gorm:"foreignkey:transformer_id;association_autoupdate:true;association_autocreate:true"`
	}

	MultiplyTransformerDBRow struct {
		ID                   uint64 `gorm:"primary_key;auto_increment"`
		TransformerID        uint64
		CreatedAt            time.Time
		UpdatedAt            time.Time
		*MultiplyTransformer `gorm:"embedded;"`
	}

	JSONParseTransformerDBRow struct {
		ID                    uint64 `gorm:"primary_key;auto_increment"`
		TransformerID         uint64
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

func UnwrapFetchersFromDB(rows ...*FetcherDBRow) Fetchers {
	var fetchers Fetchers
	for _, row := range rows {
		if row.BridgeFetcher != nil {
			row.BridgeFetcher.BridgeFetcher.Transformers = UnwrapTransformersFromDB(row.BridgeFetcher.Transformers)
			fetchers = append(fetchers, row.BridgeFetcher.BridgeFetcher)
		} else if row.HttpFetcher != nil {
			row.HttpFetcher.HttpFetcher.Transformers = UnwrapTransformersFromDB(row.HttpFetcher.Transformers)
			fetchers = append(fetchers, row.HttpFetcher.HttpFetcher)
		} else if row.MedianFetcher != nil {
			row.MedianFetcher.MedianFetcher.Fetchers = UnwrapFetchersFromDB(row.MedianFetcher.Fetchers...)
			row.MedianFetcher.MedianFetcher.Transformers = UnwrapTransformersFromDB(row.MedianFetcher.Transformers)
			fetchers = append(fetchers, row.MedianFetcher.MedianFetcher)
		}
	}
	return fetchers
}

func UnwrapTransformersFromDB(rows []*TransformerDBRow) Transformers {
	var transformers Transformers
	for _, row := range rows {
		if row.MultiplyTransformer != nil {
			transformers = append(transformers, row.MultiplyTransformer)
		} else if row.JSONParseTransformer != nil {
			transformers = append(transformers, row.JSONParseTransformer)
		}
	}
	return transformers
}
