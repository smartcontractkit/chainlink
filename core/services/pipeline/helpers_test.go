package pipeline

import (
	"gorm.io/gorm"
)

var (
	NewKeypathFromString = newKeypathFromString
)

const (
	DotStr = `
        // data source 1
        ds1          [type=bridge name=voter_turnout];
        ds1_parse    [type=jsonparse path="one,two"];
        ds1_multiply [type=multiply times=1.23];

        // data source 2
        ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}>];
        ds2_parse    [type=jsonparse path="three,four"];
        ds2_multiply [type=multiply times=4.56];

        ds1 -> ds1_parse -> ds1_multiply -> answer1;
        ds2 -> ds2_parse -> ds2_multiply -> answer1;

        answer1 [type=median                      index=0];
        answer2 [type=bridge name=election_winner index=1];
    `
)

func (t *BridgeTask) HelperSetConfigAndTxDB(config Config, txdb *gorm.DB) {
	t.config = config
	t.tx = txdb
}

func (t *HTTPTask) HelperSetConfig(config Config) {
	t.config = config
}
