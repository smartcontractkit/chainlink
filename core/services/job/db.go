package job

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"time"
)

type (
	TaskDBRow struct {
		ID        uint64 `gorm:"primary_key;auto_increment;"`
		CreatedAt time.Time
		UpdatedAt time.Time

		InputTasks []*TaskDBRow `gorm:"many2many:task_dag"`

		HttpFetcher          *HttpFetcherDBRow          `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		BridgeFetcher        *BridgeFetcherDBRow        `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		MedianFetcher        *MedianFetcherDBRow        `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		MultiplyTransformer  *MultiplyTransformerDBRow  `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		JSONParseTransformer *JSONParseTransformerDBRow `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
	}

	HttpFetcherDBRow struct {
		TaskID       uint64
		CreatedAt    time.Time
		UpdatedAt    time.Time
		*HttpFetcher `gorm:"embedded;"`
	}

	BridgeFetcherDBRow struct {
		TaskID         uint64
		CreatedAt      time.Time
		UpdatedAt      time.Time
		*BridgeFetcher `gorm:"embedded;"`
	}

	MedianFetcherDBRow struct {
		TaskID         uint64
		CreatedAt      time.Time
		UpdatedAt      time.Time
		*MedianFetcher `gorm:"embedded;"`
	}

	MultiplyTransformerDBRow struct {
		TaskID               uint64
		CreatedAt            time.Time
		UpdatedAt            time.Time
		*MultiplyTransformer `gorm:"embedded;"`
	}

	JSONParseTransformerDBRow struct {
		TaskID                uint64
		CreatedAt             time.Time
		UpdatedAt             time.Time
		*JSONParseTransformer `gorm:"embedded;"`
	}
)

func (TaskDBRow) TableName() string                 { return "tasks" }
func (HttpFetcherDBRow) TableName() string          { return "http_fetchers" }
func (BridgeFetcherDBRow) TableName() string        { return "bridge_fetchers" }
func (MedianFetcherDBRow) TableName() string        { return "median_fetchers" }
func (MultiplyTransformerDBRow) TableName() string  { return "multiply_transformers" }
func (JSONParseTransformerDBRow) TableName() string { return "jsonparse_transformers" }

func (t TaskDBRow) Task() Task {
	if t.BridgeFetcher != nil {
		return t.BridgeFetcher.BridgeFetcher
	} else if t.HttpFetcher != nil {
		return t.HttpFetcher.HttpFetcher
	} else if t.MedianFetcher != nil {
		return t.MedianFetcher.MedianFetcher
	} else if t.MultiplyTransformer != nil {
		return t.MultiplyTransformer.MultiplyTransformer
	} else if t.JSONParseTransformer != nil {
		return t.JSONParseTransformer.JSONParseTransformer
	}
	return nil
}

func WrapTasksForDB(tasks ...Task) []*TaskDBRow {
	var dbRows []*TaskDBRow
	for _, task := range tasks {
		inputTaskRows := WrapTasksForDB(task.InputTasks()...)

		switch t := task.(type) {
		case *HttpFetcher:
			dbRows = append(dbRows, &TaskDBRow{HttpFetcher: &HttpFetcherDBRow{HttpFetcher: t}, InputTasks: inputTaskRows})
		case *BridgeFetcher:
			dbRows = append(dbRows, &TaskDBRow{BridgeFetcher: &BridgeFetcherDBRow{BridgeFetcher: t}, InputTasks: inputTaskRows})
		case *MedianFetcher:
			dbRows = append(dbRows, &TaskDBRow{MedianFetcher: &MedianFetcherDBRow{MedianFetcher: t}, InputTasks: inputTaskRows})
		case *MultiplyTransformer:
			dbRows = append(dbRows, &TaskDBRow{MultiplyTransformer: &MultiplyTransformerDBRow{MultiplyTransformer: t}, InputTasks: inputTaskRows})
		case *JSONParseTransformer:
			dbRows = append(dbRows, &TaskDBRow{JSONParseTransformer: &JSONParseTransformerDBRow{JSONParseTransformer: t}, InputTasks: inputTaskRows})
		}
	}
	return dbRows
}

func UnwrapTasksFromDB(rows ...*TaskDBRow) Tasks {
	var tasks Tasks
	for _, row := range rows {
		if row.Task() == nil {
			logger.Warnw("TaskDBRow has nil task",
				"id", row.ID,
			)
			continue
		}
		inputTasks := UnwrapTasksFromDB(row.InputTasks...)
		row.Task().SetInputTasks(inputTasks)
	}
	return tasks
}
