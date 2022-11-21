//go:build go1.18

package pipeline_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func FuzzParse(f *testing.F) {
	f.Add(`ds1 [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}> timeout="10s"];`)
	f.Add(`ds1 [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}>];`)
	f.Add(`ds1 [type=http allowunrestrictednetworkaccess=true method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}> timeout="10s"];`)
	f.Add(`ds1 [type=any failEarly=true];`)
	f.Add(`ds1 [type=any];`)
	f.Add(`ds1 [type=any retries=5];`)
	f.Add(`ds1 [type=http retries=10 minBackoff="1s" maxBackoff="30m"];`)
	f.Add(pipeline.DotStr)
	f.Add(CBORDietEmpty)
	f.Add(CBORStdString)
	f.Add(`
        a [type=bridge];
        b [type=multiply times=1.23];
        a -> b -> a;
    `)
	f.Add(`
a [type=multiply input="$(val)" times=2]
b1 [type=multiply input="$(a)" times=2]
b2 [type=multiply input="$(a)" times=3]
c [type=median values=<[ $(b1), $(b2) ]> index=0]
a->b1->c;
a->b2->c;`)
	f.Add(`
// data source 1
ds1 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds1_parse [type=jsonparse path="latest"];

// data source 2
ds2 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds2_parse [type=jsonparse path="latest"];

ds1 -> ds1_parse -> answer1;
ds2 -> ds2_parse -> answer1;

answer1 [type=median index=0];
`)
	f.Add(taskRunWithVars{
		bridgeName:        "testBridge",
		ds2URL:            "https://example.com/path/to/service?with=args&foo=bar",
		ds4URL:            "http://chain.link",
		submitBridgeName:  "testSubmitBridge",
		includeInputAtKey: "path.to.key",
	}.String())
	f.Add(`s->s`)
	f.Add(`0->s->s`)
	f.Fuzz(func(t *testing.T, spec string) {
		if len(spec) > 1_000_000 {
			t.Skip()
		}
		_, err := pipeline.Parse(spec)
		if err != nil {
			t.Skip()
		}
	})
}
