package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/nanmu42/etherscan-api"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args[1]) != 2 {
		panic("need api key")
	}
	c := etherscan.New(etherscan.Mainnet, os.Args[1])
	txes, err := c.NormalTxByAddress("0xeb6ebaa55fa3d431f0046e2101ea1dd26d9658ff", nil, nil, 0, 0, false)
	panicErr(err)
	fmt.Printf("Scaping info for %v txes\n", len(txes))
	// For every txhash, scrape the gaslimit passed into the first call.
	type gv struct {
		hash          string
		val           int
		totalGasLimit int
		callString    string
	}
	var gs []gv
	var gasValues []int
	for _, tx := range txes {
		if tx.TxReceiptStatus == "0" {
			fmt.Println("skipping tx status", tx.TxReceiptStatus, tx.Hash)
			continue
		}
		if len(tx.Input) < 10 || tx.Input[:10] != "0x5e1c1059" {
			fmt.Println("skipping wrong method", tx.Input, tx.Hash)
			continue
		}
		cl := colly.NewCollector()
		cl.OnHTML("#ContentPlaceHolder1_divinternaltable > table > tbody > tr:nth-child(1) > td:nth-child(6)", func(e *colly.HTMLElement) {
			callbackLimit, err2 := strconv.Atoi(strings.Replace(e.Text, ",", "", 1))
			panicErr(err2)
			preCallbackGasUsed := tx.Gas - callbackLimit
			gasValues = append(gasValues, preCallbackGasUsed)
			gs = append(gs, gv{tx.Hash, preCallbackGasUsed, tx.Gas, e.Text})
		})
		err = cl.Visit(fmt.Sprintf("https://etherscan.io/tx/%s/advanced#internal", tx.Hash))
		panicErr(err)
	}
	s := 0.0
	for _, v := range gasValues {
		s += float64(v)
	}
	avg := s / float64(len(gasValues))
	sort.Slice(gasValues, func(i, j int) bool {
		return gasValues[i] < gasValues[j]
	})
	fmt.Println(gasValues)
	fmt.Println(gs)
	fmt.Println(avg)
}
