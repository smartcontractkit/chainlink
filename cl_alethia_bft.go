package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	Index     int
	Timestamp int64
	Data      string
	PrevHash  string
	Hash      string
}

func CalculateHash(b Block) string {
	record := fmt.Sprintf("%d%d%s%s", b.Index, b.Timestamp, b.Data, b.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
	fmt.Println("--- Alethia-BFT Consensus: Genesis ---")
	genesis := Block{0, time.Now().Unix(), "AURA_GENESIS", "0", ""}
	genesis.Hash = CalculateHash(genesis)
	
	fmt.Printf("Block: %d | Hash: %s\n", genesis.Index, genesis.Hash)
	fmt.Println("Consensus_Status: REACHED | Validator: The_Keeper")
}
