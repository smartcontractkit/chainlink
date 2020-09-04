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

		HTTPTask      *HTTPTaskDBRow      `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		BridgeTask    *BridgeTaskDBRow    `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		MedianTask    *MedianTaskDBRow    `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		MultiplyTask  *MultiplyTaskDBRow  `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
		JSONParseTask *JSONParseTaskDBRow `gorm:"foreignkey:task_id;save_association:true;association_autoupdate:true;association_autocreate:true"`
	}

	HTTPTaskDBRow struct {
		TaskID    uint64
		CreatedAt time.Time
		UpdatedAt time.Time
		*HTTPTask `gorm:"embedded;"`
	}

	BridgeTaskDBRow struct {
		TaskID      uint64
		CreatedAt   time.Time
		UpdatedAt   time.Time
		*BridgeTask `gorm:"embedded;"`
	}

	MedianTaskDBRow struct {
		TaskID      uint64
		CreatedAt   time.Time
		UpdatedAt   time.Time
		*MedianTask `gorm:"embedded;"`
	}

	MultiplyTaskDBRow struct {
		TaskID        uint64
		CreatedAt     time.Time
		UpdatedAt     time.Time
		*MultiplyTask `gorm:"embedded;"`
	}

	JSONParseTaskDBRow struct {
		TaskID         uint64
		CreatedAt      time.Time
		UpdatedAt      time.Time
		*JSONParseTask `gorm:"embedded;"`
	}
)

func (TaskDBRow) TableName() string          { return "tasks" }
func (HTTPTaskDBRow) TableName() string      { return "tasks_http" }
func (BridgeTaskDBRow) TableName() string    { return "tasks_bridge" }
func (MedianTaskDBRow) TableName() string    { return "tasks_median" }
func (MultiplyTaskDBRow) TableName() string  { return "tasks_multiply" }
func (JSONParseTaskDBRow) TableName() string { return "tasks_jsonparse" }

func (t TaskDBRow) Task() Task {
	if t.HTTPTask != nil {
		return t.HTTPTask.HTTPTask
	} else if t.BridgeTask != nil {
		return t.BridgeTask.BridgeTask
	} else if t.MedianTask != nil {
		return t.MedianTask.MedianTask
	} else if t.MultiplyTask != nil {
		return t.MultiplyTask.MultiplyTask
	} else if t.JSONParseTask != nil {
		return t.JSONParseTask.JSONParseTask
	}
	return nil
}

func (t *TaskDBRow) SetTask(task Task) {
	switch x := task.(type) {
	case *HTTPTask:
		t.HTTPTask = &HTTPTaskDBRow{HTTPTask: x}
	case *BridgeTask:
		t.BridgeTask = &BridgeTaskDBRow{BridgeTask: x}
	case *MedianTask:
		t.MedianTask = &MedianTaskDBRow{MedianTask: x}
	case *MultiplyTask:
		t.MultiplyTask = &MultiplyTaskDBRow{MultiplyTask: x}
	case *JSONParseTask:
		t.JSONParseTask = &JSONParseTaskDBRow{JSONParseTask: x}
	}
}

func WrapTasksForDB(tasks ...Task) []*TaskDBRow {
	var dbRows []*TaskDBRow
	for _, task := range tasks {
		inputTaskRows := WrapTasksForDB(task.InputTasks()...)

		switch t := task.(type) {
		case *HTTPTask:
			dbRows = append(dbRows, &TaskDBRow{HTTPTask: &HTTPTaskDBRow{HTTPTask: t}, InputTasks: inputTaskRows})
		case *BridgeTask:
			dbRows = append(dbRows, &TaskDBRow{BridgeTask: &BridgeTaskDBRow{BridgeTask: t}, InputTasks: inputTaskRows})
		case *MedianTask:
			dbRows = append(dbRows, &TaskDBRow{MedianTask: &MedianTaskDBRow{MedianTask: t}, InputTasks: inputTaskRows})
		case *MultiplyTask:
			dbRows = append(dbRows, &TaskDBRow{MultiplyTask: &MultiplyTaskDBRow{MultiplyTask: t}, InputTasks: inputTaskRows})
		case *JSONParseTask:
			dbRows = append(dbRows, &TaskDBRow{JSONParseTask: &JSONParseTaskDBRow{JSONParseTask: t}, InputTasks: inputTaskRows})
		}
	}
	return dbRows
}

func UnwrapTasksFromDB(rows ...*TaskDBRow) []Task {
	var tasks []Task
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
