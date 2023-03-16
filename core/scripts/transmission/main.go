package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/transmission/handler"
)

func main() {
	switch os.Args[1] {
	default:
		postBody, _ := json.Marshal(map[string]any{
			"sender_address":  "0x82d72a7f2f08b07D9B90134D47A673A75c1916A7",
			"source_chain_id": 5,
			"token_address":   "0x82d72a7f2f08b07D9B90134D47A673A75c1916A7",
		})
		responseBody := bytes.NewBuffer(postBody)
		//Leverage Go's HTTP Post function to make request
		resp, err := http.Post("http://localhost:2000/get_nonce", "application/json", responseBody)
		helpers.PanicErr(err)
		var gnr handler.GetNonceResponse
		json.NewDecoder(resp.Body).Decode(&gnr)
		fmt.Printf("deserialized response: \n %+v", gnr)
	}
}
