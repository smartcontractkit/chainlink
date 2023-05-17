package handler

import (
	"context"
	"log"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/montanaflynn/stats"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/verifiable_load_upkeep_wrapper"
)

func (k *Keeper) GetVerifiableLoadStats(ctx context.Context) {
	//chainId, err := k.client.ChainID(ctx)
	//log.Println(chainId)
	//addr, tx, v, err := verifiable_load_upkeep_wrapper.DeployVerifiableLoadUpkeep(k.buildTxOpts(ctx), k.client, common.HexToAddress(k.cfg.Registrar), false)
	//if err != nil {
	//	log.Fatalf("failed to deploy verifiable load upkeep: %v", err)
	//}
	//
	//k.waitDeployment(ctx, tx)
	//log.Println("verifiable load upkeep deployed:", addr.Hex(), "-", helpers.ExplorerLink(k.cfg.ChainID, tx.Hash()))
	//
	//tx, err = k.linkToken.Approve(k.buildTxOpts(ctx), k.fromAddr, big.NewInt(3000000000000000000))
	//if err != nil {
	//	log.Fatalf("failed to approve: %v", err)
	//}
	//k.waitTx(ctx, tx)
	//
	//tx, err = k.linkToken.Transfer(k.buildTxOpts(ctx), addr, big.NewInt(3000000000000000000))
	//if err != nil {
	//	log.Fatalf("failed to transfer: %v", err)
	//}
	//k.waitTx(ctx, tx)
	//
	//tx, err = v.BatchRegisterUpkeeps(k.buildTxOpts(ctx), 1, 3000000, big.NewInt(1000000000000000000), big.NewInt(1000000), big.NewInt(1000000))
	//if err != nil {
	//	log.Fatalf("failed to batch register: %v", err)
	//}
	//k.waitTx(ctx, tx)
	//
	////v.BatchRegisterUpkeeps(txOpts, 1, 3000000)

	addr := common.HexToAddress(k.cfg.VerifiableLoadContractAddress)
	v, err := verifiable_load_upkeep_wrapper.NewVerifiableLoadUpkeep(addr, k.client)
	if err != nil {
		log.Fatalf("failed to create a new verifiable load upkeep from address %s: %v", k.cfg.VerifiableLoadContractAddress, err)
	}

	// get all active upkeep IDs on this verifiable load contract
	opts := &bind.CallOpts{
		From:    k.fromAddr,
		Context: ctx,
	}
	uid := big.NewInt(0)
	uid, _ = uid.SetString("101926416817634971218389061868275311034641546399220249757058712333177045691390", 10)

	upkeepIds, err := v.GetActiveUpkeepIDs(opts, big.NewInt(0), big.NewInt(0))
	if err != nil {
		log.Fatalf("failed to get active upkeep IDs from %s: %v", k.cfg.VerifiableLoadContractAddress, err)
	}

	var allUpkeepsDelays []float64
	var allUpkeepsTotalDelay, allUpkeepsTotalPerforms uint64

	for _, id := range upkeepIds {
		// it's possible to do the following to get delays, but it may run into out of gas issues
		// performDelays, err := v.GetDelays(opts, uid)
		log.Println()
		log.Printf("================================== UPKEEP %s SUMMARY =======================================================", id.String())

		c, err := v.Counters(opts, id)
		if err != nil {
			log.Fatalf("failed to get counter for %s: %v", id.String(), err)
		}

		b, err := v.Buckets(opts, id)
		if err != nil {
			log.Fatalf("failed to get current bucket count for %s: %v", id.String(), err)
		}
		log.Printf("upkeep ID %s has been performed %d times in %d buckets\n", id.String(), c, b+1)

		var delays []float64
		var totalDelay float64
		var totalPerforms uint64
		for i := uint16(0); i <= b; i++ {
			bucketDelays, err := v.GetBucketedDelays(opts, id, i)
			if err != nil {
				log.Fatalf("failed to get bucketed delays for upkeep id %s bucket %d: %v", id.String(), i, err)
			}

			var floatBucketDelays []float64
			var totalBucketDelay float64
			for _, d := range bucketDelays {
				delays = append(delays, float64(d.Uint64()))
				floatBucketDelays = append(floatBucketDelays, float64(d.Uint64()))
				totalDelay += float64(d.Uint64())
				totalBucketDelay += float64(d.Uint64())
				totalPerforms++
			}
			allUpkeepsTotalDelay += uint64(totalBucketDelay)

			p50, _ := stats.Percentile(floatBucketDelays, 50)
			p90, _ := stats.Percentile(floatBucketDelays, 90)
			p95, _ := stats.Percentile(floatBucketDelays, 95)
			p99, _ := stats.Percentile(floatBucketDelays, 99)
			sort.Float64s(floatBucketDelays)
			log.Printf("bucket %d th 100 performs p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %d, average perform delay: %f\n", i, p50, p90, p95, p99, floatBucketDelays[len(floatBucketDelays)-1], uint64(totalBucketDelay), totalBucketDelay/float64(len(bucketDelays)))
		}
		allUpkeepsTotalPerforms += totalPerforms

		p50, _ := stats.Percentile(delays, 50)
		p90, _ := stats.Percentile(delays, 90)
		p95, _ := stats.Percentile(delays, 95)
		p99, _ := stats.Percentile(delays, 99)
		sort.Float64s(delays)
		maxDelay := delays[len(delays)-1]

		allUpkeepsDelays = append(allUpkeepsDelays, delays...)

		log.Printf("%d performs in total. p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %d, average perform delay: %f\n", totalPerforms, p50, p90, p95, p99, maxDelay, uint64(totalDelay), totalDelay/float64(totalPerforms))

		t, err := v.TimestampBuckets(opts, id)
		if err != nil {
			log.Fatalf("failed to get timestamp bucket for %s: %v", id.String(), err)
		}

		delays = nil
		totalDelay = 0
		totalPerforms = 0
		for i := uint16(0); i <= t; i++ {
			timestampDelays, err := v.GetTimestampDelays(opts, id, i)
			if err != nil {
				log.Fatalf("failed to get timestamped delays for upkeep id %s timestamp bucket %d: %v", id.String(), i, err)
			}

			var floatTimestampDelays []float64
			var totalTimestampDelay float64
			for _, d := range timestampDelays {
				delays = append(delays, float64(d.Uint64()))
				floatTimestampDelays = append(floatTimestampDelays, float64(d.Uint64()))
				totalDelay += float64(d.Uint64())
				totalTimestampDelay += float64(d.Uint64())
				totalPerforms++
			}
			p50, _ = stats.Percentile(floatTimestampDelays, 50)
			p90, _ = stats.Percentile(floatTimestampDelays, 90)
			p95, _ = stats.Percentile(floatTimestampDelays, 95)
			p99, _ = stats.Percentile(floatTimestampDelays, 99)
			sort.Float64s(floatTimestampDelays)
			log.Printf("timestamp bucket %d th hour performs p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %d, average perform delay: %f\n", i, p50, p90, p95, p99, floatTimestampDelays[len(floatTimestampDelays)-1], uint64(totalTimestampDelay), totalTimestampDelay/float64(len(timestampDelays)))
		}

		p50, _ = stats.Percentile(delays, 50)
		p90, _ = stats.Percentile(delays, 90)
		p95, _ = stats.Percentile(delays, 95)
		p99, _ = stats.Percentile(delays, 99)
		sort.Float64s(delays)
		maxDelay = delays[len(delays)-1]

		log.Printf("%d performs in total. p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %d, average perform delay: %f\n", totalPerforms, p50, p90, p95, p99, maxDelay, uint64(totalDelay), totalDelay/float64(totalPerforms))
	}

	log.Println("================================== ALL UPKEEPS SUMMARY =======================================================")
	sort.Float64s(allUpkeepsDelays)
	p50, _ := stats.Percentile(allUpkeepsDelays, 50)
	p90, _ := stats.Percentile(allUpkeepsDelays, 90)
	p95, _ := stats.Percentile(allUpkeepsDelays, 95)
	p99, _ := stats.Percentile(allUpkeepsDelays, 99)
	maxDelay := allUpkeepsDelays[len(allUpkeepsDelays)-1]
	log.Printf("For total %d upkeeps: total performs: %d, p50: %f, p90: %f, p95: %f, p99: %f, max delay: %f, total delay blocks: %d, average perform delay: %f\n", len(upkeepIds), allUpkeepsTotalPerforms, p50, p90, p95, p99, maxDelay, allUpkeepsTotalDelay, float64(allUpkeepsTotalDelay)/float64(allUpkeepsTotalPerforms))

}
