package job_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	// "github.com/BurntSushi/toml"
	"gonum.org/v1/gonum/graph/encoding/dot"
)

func TestGraph_Decode(t *testing.T) {
	dotStr := `
        // data source 1
        ds1       [type=bridge name=voter_turnout];
        ds1_parse [type=jsonparse path="data,result"];

        // data source 2
        ds2       [type=http url="https://chain.link/voter_turnout/USA-2020" requestData="{\"hi\": \"hello\"}"];
        ds2_parse [type=jsonparse path="data,result"];

        answer1 [type=median];

        ds1 -> ds1_parse -> answer1;
        ds2 -> ds2_parse -> answer1;

        answer2 [type=bridge name=election_winner];
    `

	var graph job.TaskDAG
	err := dot.Unmarshal([]byte(dotStr), &graph)
	require.NoError(t, err)

	iter := graph.Nodes()
	for iter.Next() {
		n := iter.Node()
		n2 := n.(*node)
		// j, err := json.MarshalIndent(n2.attrs, "", "    ")
		// if err != nil {
		//  panic(err)
		// }
		// fmt.Println(n2.dotID, string(j))
		for k, v := range n2.attrs {
			fmt.Println("  +", k, ":", v)
		}

		iter2 := graph.To(n2.ID())
		for iter2.Next() {
			n := iter2.Node()
			n2 := n.(*node)
			fmt.Println("  -", n2.dotID)
		}
		fmt.Println()
	}

}
