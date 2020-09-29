package pipeline

import (
	"reflect"
)

func NewBaseTask(dotID string, t Task, index int32) BaseTask {
	return BaseTask{dotID: dotID, outputTask: t, Index: index}
}

func (t *HTTPTask) HelperSetConfig(config Config) {
	t.config = config
}

func (t *BridgeTask) HelperSetConfigAndORM(config Config, orm ORM) {
	t.config = config
	t.orm = orm
}

func (t MultiplyTask) ExportedEquals(otherTask Task) bool {
	other, ok := otherTask.(*MultiplyTask)
	if !ok {
		return false
	} else if t.Index != other.Index {
		return false
	} else if !t.Times.Equal(other.Times) {
		return false
	}
	return true
}

func (t MedianTask) ExportedEquals(otherTask Task) bool {
	other, ok := otherTask.(*MedianTask)
	if !ok {
		return false
	} else if t.Index != other.Index {
		return false
	}
	return true
}

func (t JSONParseTask) ExportedEquals(otherTask Task) bool {
	other, ok := otherTask.(*JSONParseTask)
	if !ok {
		return false
	} else if t.Index != other.Index {
		return false
	} else if !reflect.DeepEqual(t.Path, other.Path) {
		return false
	}
	return true
}

func (t HTTPTask) ExportedEquals(otherTask Task) bool {
	other, ok := otherTask.(*HTTPTask)
	if !ok {
		return false
	} else if t.Index != other.Index {
		return false
	} else if t.Method != other.Method {
		return false
		// } else if !reflect.DeepEqual(t.ExtendedPath, other.ExtendedPath) {
		// 	return false
		// } else if !reflect.DeepEqual(t.Headers, other.Headers) {
		// 	return false
		// } else if !reflect.DeepEqual(t.QueryParams, other.QueryParams) {
		// 	return false
	} else if !reflect.DeepEqual(t.RequestData, other.RequestData) {
		return false
	}
	return true
}

func (t BridgeTask) ExportedEquals(otherTask Task) bool {
	other, ok := otherTask.(*BridgeTask)
	if !ok {
		return false
	} else if t.Index != other.Index {
		return false
	} else if t.Name != other.Name {
		return false
	} else if !reflect.DeepEqual(t.RequestData, other.RequestData) {
		return false
	}
	return true
}
