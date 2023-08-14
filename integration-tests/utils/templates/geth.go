package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var InitGethScript = `
#!/bin/bash
if [ ! -d /root/.ethereum/keystore ]; then
	echo "/root/.ethereum/keystore not found, running 'geth init'..."
	geth init /root/genesis.json
	echo "...done!"
fi

geth "$@"
`

var GenesisJson = `
{
	"config": {
	  "chainId": {{ .ChainId }},
	  "homesteadBlock": 0,
	  "eip150Block": 0,
	  "eip155Block": 0,
	  "eip158Block": 0,
	  "eip160Block": 0,
	  "byzantiumBlock": 0,
	  "constantinopleBlock": 0,
	  "petersburgBlock": 0,
	  "istanbulBlock": 0,
	  "muirGlacierBlock": 0,
	  "berlinBlock": 0,
	  "londonBlock": 0
	},
	"nonce": "0x0000000000000042",
	"mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"difficulty": "1",
	"coinbase": "0x3333333333333333333333333333333333333333",
	"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"extraData": "0x",
	"gasLimit": "8000000000",
	"alloc": {
	  "{{ .AccountAddr }}": {
		"balance": "20000000000000000000000"
	  }
	}
  }`

func BuildGenesisJson(chainId, accountAddr string) (string, error) {
	data := struct {
		AccountAddr string
		ChainId     string
	}{
		AccountAddr: accountAddr,
		ChainId:     chainId,
	}

	t, err := template.New("genesis-json").Parse(GenesisJson)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err
}
