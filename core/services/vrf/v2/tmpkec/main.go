package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	sigs := []string{
		"RandomWordsFulfilled(uint256,uint256,uint256,uint96,bool,bool)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint96,bool,bool,bool)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint256,bool,bool)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint256,bool,bool,bool)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint96,bool)",
		"RandomWordsFulfilled(uint256,uint256,uint96,bool)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint96,bool,bool,uint256)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint96,bool,bool,bytes)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint256,bool,bool,bytes)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint96,bool,bool,uint96)",
		"RandomWordsFulfilled(uint256,uint256,uint256,uint96,address,bool,bool)",
	}
	for _, s := range sigs {
		_ = 0
		
		fmt.Printf("0x%x  %s\n", crypto.Keccak256([]byte(s)), s)
	}
}
