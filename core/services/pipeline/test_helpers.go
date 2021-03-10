package pipeline

import (
	"reflect"

	"gorm.io/gorm"
)

const (
	DotStr = `
    // data source 1
    ds1          [type=bridge name=voter_turnout];
    ds1_parse    [type=jsonparse path="one,two"];
    ds1_multiply [type=multiply times=1.23];

    // data source 2
    ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\"hi\": \"hello\"}"];
    ds2_parse    [type=jsonparse path="three,four"];
    ds2_multiply [type=multiply times=4.56];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

    answer1 [type=median                      index=0];
    answer2 [type=bridge name=election_winner index=1];
`
)

func NewBaseTask(dotID string, t Task, index int32, nPreds int) BaseTask {
	return BaseTask{dotID: dotID, outputTask: t, Index: index, nPreds: nPreds}
}

func (t *BridgeTask) HelperSetConfigAndTxDB(config Config, txdb *gorm.DB) {
	t.config = config
	t.txdb = txdb
}

func (t *HTTPTask) HelperSetConfig(config Config) {
	t.config = config
}

func (t ResultTask) ExportedEquals(otherTask Task) bool {
	other, ok := otherTask.(*ResultTask)
	if !ok {
		return false
	} else if t.Index != other.Index {
		return false
	}
	return true
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
	} else if t.URL != other.URL {
		return false
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
