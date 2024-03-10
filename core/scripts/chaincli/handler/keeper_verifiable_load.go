package handler

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/montanaflynn/stats"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/verifiable_load_upkeep_wrapper"
)

const (
	// workerNum is the total number of workers calculating upkeeps' delay summary
	workerNum = 5
	// retryDelay is the time the go routine will wait before calling the same contract function
	retryDelay = 1 * time.Second
	// retryNum defines how many times the go routine will attempt the same contract call
	retryNum = 3
	// maxUpkeepNum defines the size of channels. Increase if there are lots of upkeeps.
	maxUpkeepNum = 100
)

type upkeepInfo struct {
	mu              sync.Mutex
	ID              *big.Int
	Bucket          uint16
	DelayBuckets    map[uint16][]float64
	SortedAllDelays []float64
	TotalDelayBlock float64
	TotalPerforms   uint64
}

type verifiableLoad interface {
	GetAllActiveUpkeepIDsOnRegistry(opts *bind.CallOpts, startIndex *big.Int, maxCount *big.Int) ([]*big.Int, error)
	Counters(opts *bind.CallOpts, upkeepId *big.Int) (*big.Int, error)
	GetBucketedDelays(opts *bind.CallOpts, upkeepId *big.Int, bucket uint16) ([]*big.Int, error)
	Buckets(opts *bind.CallOpts, arg0 *big.Int) (uint16, error)
}

func (ui *upkeepInfo) AddBucket(bucketNum uint16, bucketDelays []float64) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	ui.DelayBuckets[bucketNum] = bucketDelays
}

type upkeepStats struct {
	BlockNumber     uint64
	AllInfos        []*upkeepInfo
	TotalDelayBlock float64
	TotalPerforms   uint64
	SortedAllDelays []float64
}

func (k *Keeper) PrintVerifiableLoadStats(ctx context.Context, csv bool) {
	var v verifiableLoad
	var err error
	addr := common.HexToAddress(k.cfg.VerifiableLoadContractAddress)
	v, err = verifiable_load_upkeep_wrapper.NewVerifiableLoadUpkeep(addr, k.client)
	if err != nil {
		log.Fatalf("failed to create a new verifiable load upkeep from address %s: %v", k.cfg.VerifiableLoadContractAddress, err)
	}

	// get all the stats from this block
	blockNum, err := k.client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("failed to get block number: %v", err)
	}

	opts := &bind.CallOpts{
		From:        k.fromAddr,
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNum)),
	}

	// get all active upkeep IDs on this verifiable load contract
	upkeepIds, err := v.GetAllActiveUpkeepIDsOnRegistry(opts, big.NewInt(0), big.NewInt(0))
	if err != nil {
		log.Fatalf("failed to get active upkeep IDs from %s: %v", k.cfg.VerifiableLoadContractAddress, err)
	}

	if csv {
		fmt.Println("upkeep ID,total performs,p50,p90,p95,p99,max delay,total delay blocks,average perform delay")
	}

	us := &upkeepStats{BlockNumber: blockNum}

	resultsChan := make(chan *upkeepInfo, maxUpkeepNum)
	idChan := make(chan *big.Int, maxUpkeepNum)

	var wg sync.WaitGroup

	// create a number of workers to process the upkeep ids in batch
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go k.fetchUpkeepInfo(idChan, resultsChan, v, opts, &wg, csv)
	}

	for _, id := range upkeepIds {
		idChan <- id
	}

	close(idChan)
	wg.Wait()

	close(resultsChan)

	for info := range resultsChan {
		us.AllInfos = append(us.AllInfos, info)
		us.TotalPerforms += info.TotalPerforms
		us.TotalDelayBlock += info.TotalDelayBlock
		us.SortedAllDelays = append(us.SortedAllDelays, info.SortedAllDelays...)
	}

	sort.Float64s(us.SortedAllDelays)

	log.Println("\n\n================================== ALL UPKEEPS SUMMARY =======================================================")
	p50, _ := stats.Percentile(us.SortedAllDelays, 50)
	p90, _ := stats.Percentile(us.SortedAllDelays, 90)
	p95, _ := stats.Percentile(us.SortedAllDelays, 95)
	p99, _ := stats.Percentile(us.SortedAllDelays, 99)

	maxDelay := float64(0)
	if len(us.SortedAllDelays) > 0 {
		maxDelay = us.SortedAllDelays[len(us.SortedAllDelays)-1]
	}
	log.Printf("For total %d upkeeps: total performs: %d, p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %f, average perform delay: %f\n", len(upkeepIds), us.TotalPerforms, p50, p90, p95, p99, maxDelay, us.TotalDelayBlock, us.TotalDelayBlock/float64(us.TotalPerforms))
	log.Printf("All STATS ABOVE ARE CALCULATED AT BLOCK %d", blockNum)
}

func (k *Keeper) fetchUpkeepInfo(idChan chan *big.Int, resultsChan chan *upkeepInfo, v verifiableLoad, opts *bind.CallOpts, wg *sync.WaitGroup, csv bool) {
	defer wg.Done()

	for id := range idChan {
		// fetch how many times this upkeep has been executed
		c, err := v.Counters(opts, id)
		if err != nil {
			log.Fatalf("failed to get counter for %s: %v", id.String(), err)
		}

		// get all the buckets of an upkeep. 100 performs is a bucket.
		b, err := v.Buckets(opts, id)
		if err != nil {
			log.Fatalf("failed to get current bucket count for %s: %v", id.String(), err)
		}

		info := &upkeepInfo{
			ID:            id,
			Bucket:        b,
			TotalPerforms: c.Uint64(),
			DelayBuckets:  map[uint16][]float64{},
		}

		var delays []float64
		var wg1 sync.WaitGroup
		for i := uint16(0); i <= b; i++ {
			wg1.Add(1)
			go k.fetchBucketData(v, opts, id, i, &wg1, info)
		}
		wg1.Wait()

		for i := uint16(0); i <= b; i++ {
			bucketDelays := info.DelayBuckets[i]
			delays = append(delays, bucketDelays...)
			for _, d := range bucketDelays {
				info.TotalDelayBlock += d
			}
		}
		sort.Float64s(delays)
		info.SortedAllDelays = delays
		info.TotalPerforms = uint64(len(info.SortedAllDelays))

		p50, _ := stats.Percentile(info.SortedAllDelays, 50)
		p90, _ := stats.Percentile(info.SortedAllDelays, 90)
		p95, _ := stats.Percentile(info.SortedAllDelays, 95)
		p99, _ := stats.Percentile(info.SortedAllDelays, 99)

		maxDelay := float64(0)

		if len(info.SortedAllDelays) > 0 {
			maxDelay = info.SortedAllDelays[len(info.SortedAllDelays)-1]
		}

		if csv {
			fmt.Printf("%s,%d,%f,%f,%f,%f,%f,%d,%f\n", id, info.TotalPerforms, p50, p90, p95, p99, maxDelay, uint64(info.TotalDelayBlock), info.TotalDelayBlock/float64(info.TotalPerforms))
		} else {
			log.Printf("upkeep ID %s has %d performs in total. p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %d, average perform delay: %f\n", id, info.TotalPerforms, p50, p90, p95, p99, maxDelay, uint64(info.TotalDelayBlock), info.TotalDelayBlock/float64(info.TotalPerforms))
		}
		resultsChan <- info
	}
}

func (k *Keeper) fetchBucketData(v verifiableLoad, opts *bind.CallOpts, id *big.Int, bucketNum uint16, wg *sync.WaitGroup, info *upkeepInfo) {
	defer wg.Done()

	var bucketDelays []*big.Int
	var err error
	for i := 0; i < retryNum; i++ {
		bucketDelays, err = v.GetBucketedDelays(opts, id, bucketNum)
		if err == nil {
			break
		}
		log.Printf("failed to get bucketed delays for upkeep id %s bucket %d: %v, retrying...", id.String(), bucketNum, err)
		time.Sleep(retryDelay)
	}

	var floatBucketDelays []float64
	for _, d := range bucketDelays {
		floatBucketDelays = append(floatBucketDelays, float64(d.Uint64()))
	}
	sort.Float64s(floatBucketDelays)
	info.AddBucket(bucketNum, floatBucketDelays)
}
