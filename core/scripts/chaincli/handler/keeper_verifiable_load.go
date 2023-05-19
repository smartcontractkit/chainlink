package handler

import (
	"context"
	"log"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/montanaflynn/stats"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/verifiable_load_upkeep_wrapper"
)

type UpkeepInfo struct {
	mu                    sync.Mutex
	ID                    *big.Int
	Bucket                uint16
	TimestampBucket       uint16
	DelayBuckets          map[uint16][]float64
	DelayTimestampBuckets map[uint16][]float64
	SortedAllDelays       []float64
	TotalDelayBlock       float64
	TotalPerforms         uint64
}

func (ui *UpkeepInfo) AddBucket(bucketNum uint16, bucketDelays []float64) {
	ui.mu.Lock()
	ui.DelayBuckets[bucketNum] = bucketDelays
	ui.mu.Unlock()
}

func (ui *UpkeepInfo) AddTimestampBucket(bucketNum uint16, bucketDelays []float64) {
	ui.mu.Lock()
	ui.DelayTimestampBuckets[bucketNum] = bucketDelays
	ui.mu.Unlock()
}

type UpkeepStats struct {
	BlockNumber     uint64
	AllInfos        []*UpkeepInfo
	TotalDelayBlock float64
	TotalPerforms   uint64
	SortedAllDelays []float64
}

func (k *Keeper) GetVerifiableLoadStats(ctx context.Context) {
	addr := common.HexToAddress(k.cfg.VerifiableLoadContractAddress)
	v, err := verifiable_load_upkeep_wrapper.NewVerifiableLoadUpkeep(addr, k.client)
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
	upkeepIds, err := v.GetActiveUpkeepIDs(opts, big.NewInt(0), big.NewInt(0))
	if err != nil {
		log.Fatalf("failed to get active upkeep IDs from %s: %v", k.cfg.VerifiableLoadContractAddress, err)
	}

	upkeepStats := &UpkeepStats{BlockNumber: blockNum}

	for _, id := range upkeepIds {
		// it's possible to do the following to get delays, but it may run into out of gas issues
		// performDelays, err := v.GetDelays(opts, uid)
		log.Println()
		log.Printf("================================== UPKEEP %s SUMMARY =======================================================", id.String())

		c, err := v.Counters(opts, id)
		if err != nil {
			log.Fatalf("failed to get counter for %s: %v", id.String(), err)
		}

		// get all the buckets of an upkeep. 100 performs is a bucket.
		b, err := v.Buckets(opts, id)
		if err != nil {
			log.Fatalf("failed to get current bucket count for %s: %v", id.String(), err)
		}
		log.Printf("upkeep ID %s has been performed %d times in %d buckets\n", id.String(), c, b+1)
		info := &UpkeepInfo{
			ID:                    id,
			Bucket:                b,
			TotalPerforms:         c.Uint64(),
			DelayBuckets:          map[uint16][]float64{},
			DelayTimestampBuckets: map[uint16][]float64{},
		}

		var delays []float64
		var wg sync.WaitGroup
		for i := uint16(0); i <= b; i++ {
			go k.getBucketData(v, opts, false, id, i, &wg, info)
		}
		wg.Wait()

		// get all the timestamp buckets of an upkeep. performs which happen every 1 hour after the first perform fall into the same bucket.
		t, err := v.TimestampBuckets(opts, id)
		if err != nil {
			log.Fatalf("failed to get timestamp bucket for %s: %v", id.String(), err)
		}
		info.TimestampBucket = t
		for i := uint16(0); i <= t; i++ {
			go k.getBucketData(v, opts, true, id, i, &wg, info)
		}
		wg.Wait()

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
		maxDelay := info.SortedAllDelays[len(info.SortedAllDelays)-1]

		log.Printf("%d performs in total. p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %d, average perform delay: %f\n", info.TotalPerforms, p50, p90, p95, p99, maxDelay, uint64(info.TotalDelayBlock), info.TotalDelayBlock/float64(info.TotalPerforms))
		log.Printf("All delays: %v", info.SortedAllDelays)
		upkeepStats.AllInfos = append(upkeepStats.AllInfos, info)
		upkeepStats.TotalPerforms += info.TotalPerforms
		upkeepStats.TotalDelayBlock += info.TotalDelayBlock
		upkeepStats.SortedAllDelays = append(upkeepStats.SortedAllDelays, info.SortedAllDelays...)
	}

	sort.Float64s(upkeepStats.SortedAllDelays)

	log.Println("\n\n================================== ALL UPKEEPS SUMMARY =======================================================")
	p50, _ := stats.Percentile(upkeepStats.SortedAllDelays, 50)
	p90, _ := stats.Percentile(upkeepStats.SortedAllDelays, 90)
	p95, _ := stats.Percentile(upkeepStats.SortedAllDelays, 95)
	p99, _ := stats.Percentile(upkeepStats.SortedAllDelays, 99)
	maxDelay := upkeepStats.SortedAllDelays[len(upkeepStats.SortedAllDelays)-1]
	log.Printf("For total %d upkeeps: total performs: %d, p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %f, average perform delay: %f\n", len(upkeepIds), upkeepStats.TotalPerforms, p50, p90, p95, p99, maxDelay, upkeepStats.TotalDelayBlock, upkeepStats.TotalDelayBlock/float64(upkeepStats.TotalPerforms))
	log.Printf("All STATS ABOVE ARE CALCULATED AT BLOCK %d", blockNum)
}

func (k *Keeper) getBucketData(v *verifiable_load_upkeep_wrapper.VerifiableLoadUpkeep, opts *bind.CallOpts, getTimestampBucket bool, id *big.Int, bucketNum uint16, wg *sync.WaitGroup, info *UpkeepInfo) {
	wg.Add(1)
	defer wg.Done()

	var bucketDelays []*big.Int
	var err error
	if getTimestampBucket {
		bucketDelays, err = v.GetTimestampDelays(opts, id, bucketNum)
		if err != nil {
			log.Fatalf("failed to get timestamp bucketed delays for upkeep id %s timestamp bucket %d: %v", id.String(), bucketNum, err)
		}
	} else {
		bucketDelays, err = v.GetBucketedDelays(opts, id, bucketNum)
		if err != nil {
			log.Fatalf("failed to get bucketed delays for upkeep id %s bucket %d: %v", id.String(), bucketNum, err)
		}
	}

	var floatBucketDelays []float64
	for _, d := range bucketDelays {
		floatBucketDelays = append(floatBucketDelays, float64(d.Uint64()))
	}
	sort.Float64s(floatBucketDelays)

	if getTimestampBucket {
		info.AddTimestampBucket(bucketNum, floatBucketDelays)
	} else {
		info.AddBucket(bucketNum, floatBucketDelays)
	}
}
