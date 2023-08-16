package templates

import (
	"github.com/google/uuid"
)

type GenesisJsonTemplate struct {
	AccountAddr string
	ChainId     string
}

// String representation of the job
func (c GenesisJsonTemplate) String() (string, error) {
	tpl := `
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
	return MarshalTemplate(c, uuid.NewString(), tpl)
}

var InitGethScript = `
#!/bin/bash
if [ ! -d /root/.ethereum/keystore ]; then
	echo "/root/.ethereum/keystore not found, running 'geth init'..."
	geth init /root/genesis.json
	echo "...done!"
fi

geth "$@"
`
