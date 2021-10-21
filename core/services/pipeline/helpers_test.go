package pipeline

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/services/eth"
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

func (t *BridgeTask) HelperSetDependencies(config Config, db *gorm.DB, id uuid.UUID) {
	t.config = config
	t.db = db
	t.uuid = id
}

func (t *HTTPTask) HelperSetDependencies(config Config) {
	t.config = config
}

func (t *ETHCallTask) HelperSetDependencies(client eth.Client) {
	t.ethClient = client
}

func (t *ETHTxTask) HelperSetDependencies(db *gorm.DB, config Config, keyStore ETHKeyStore, txManager TxManager) {
	t.db = db
	t.config = config
	t.keyStore = keyStore
	t.txManager = txManager
}
