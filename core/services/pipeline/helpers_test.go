package pipeline

import (
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
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

        ds1 -> ds1_parse -> ds1_multiply;
        ds2 -> ds2_parse -> ds2_multiply -> answer1;

        answer1 [type=median	index=0 input1="$(ds1_multiply)" input2="$(ds2_multiply)"];
        answer2 [type=bridge 	name=election_winner	index=1];
    `
)

func (t *BridgeTask) HelperSetDependencies(
	config Config,
	orm bridges.ORM,
	specId int32,
	id uuid.UUID,
	httpClient *http.Client) {
	t.config = config
	t.orm = orm
	t.uuid = id
	t.httpClient = httpClient
	t.specId = specId
}

func (t *HTTPTask) HelperSetDependencies(config Config, restrictedHTTPClient, unrestrictedHTTPClient *http.Client) {
	t.config = config
	t.httpClient = restrictedHTTPClient
	t.unrestrictedHTTPClient = unrestrictedHTTPClient
}

func (t *ETHCallTask) HelperSetDependencies(cc evm.ChainSet, config Config, specGasLimit *uint32, jobType string) {
	t.chainSet = cc
	t.config = config
	t.specGasLimit = specGasLimit
	t.jobType = jobType
}

func (t *ETHTxTask) HelperSetDependencies(cc evm.ChainSet, keyStore ETHKeyStore, specGasLimit *uint32, jobType string) {
	t.chainSet = cc
	t.keyStore = keyStore
	t.specGasLimit = specGasLimit
	t.jobType = jobType
}

func (t *ETHGetBlockTask) HelperSetDependencies(cc evm.ChainSet, config Config) {
	t.chainSet = cc
	t.config = config
}
