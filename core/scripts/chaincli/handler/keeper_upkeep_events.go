package handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/upkeep_counter_wrapper"
)

// UpkeepCounterEvents print out emitted events and write to csv file
func (k *Keeper) UpkeepCounterEvents(ctx context.Context, hexAddr string, fromBlock, toBlock uint64) {
	contractAddress := common.HexToAddress(hexAddr)
	upkeepCounter, err := upkeep_counter_wrapper.NewUpkeepCounter(contractAddress, k.client)
	if err != nil {
		log.Fatalln("Failed to create a new upkeep counter", err)
	}
	filterOpts := bind.FilterOpts{
		Start:   fromBlock,
		End:     &toBlock,
		Context: ctx,
	}
	upkeepIterator, err := upkeepCounter.FilterPerformingUpkeep(&filterOpts, nil)
	if err != nil {
		log.Fatalln("Failed to get upkeep iterator", err)
	}
	filename := fmt.Sprintf("%s.csv", hexAddr)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	fmt.Println("From, InitialBlock, LastBlock, PreviousBlock, Counter")
	row := []string{"From", "InitialBlock", "LastBlock", "PreviousBlock", "Counter"}
	if err = w.Write(row); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for upkeepIterator.Next() {
		fmt.Printf("%s,%s,%s,%s,%s\n",
			upkeepIterator.Event.From,
			upkeepIterator.Event.InitialBlock,
			upkeepIterator.Event.LastBlock,
			upkeepIterator.Event.PreviousBlock,
			upkeepIterator.Event.Counter,
		)
		row = []string{upkeepIterator.Event.From.String(),
			upkeepIterator.Event.InitialBlock.String(),
			upkeepIterator.Event.LastBlock.String(),
			upkeepIterator.Event.PreviousBlock.String(),
			upkeepIterator.Event.Counter.String()}
		if err = w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
