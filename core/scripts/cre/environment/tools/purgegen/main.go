// purgegen writes N serialized CloudEvent payloads (plus a few corrupt ones) as
// psql COPY input for cre.chip_durable_events, with a weighted domain/subject
// mix, and prints the exact distribution so a purge can be checked against it.
package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
)

type kind struct {
	domain, subject string
	weight          int
}

func main() {
	n := flag.Int("n", 100000, "number of decodable events")
	corrupt := flag.Int("corrupt", 5, "number of undecodable payloads")
	age := flag.String("age", "3 minutes", "created_at backdate interval")
	out := flag.String("out", "events.tsv", "COPY input file")
	flag.Parse()

	kinds := []kind{
		{"platform", "workflow.execution.started", 30},
		{"platform", "workflow.execution.finished", 30},
		{"platform", "workflow.trigger", 12},
		{"billing", "meter.record", 15},
		{"billing", "workflow.receipt", 5},
		{"cre", "node.heartbeat", 4},
		{"cre", "capability.execution", 3},
		{"vault", "secrets.access", 1},
	}
	var total int
	for _, k := range kinds {
		total += k.weight
	}
	rng := rand.New(rand.NewSource(42))
	counts := map[string]int{}

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriterSize(f, 1<<20)
	body := []byte(`{"executionID":"0000000000000000000000000000000000000000000000000000000000000000","status":"completed","durationMs":1234,"payload":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}`)
	for i := 0; i < *n; i++ {
		r := rng.Intn(total)
		var k kind
		for _, c := range kinds {
			if r < c.weight {
				k = c
				break
			}
			r -= c.weight
		}
		ev, err := chipingress.NewEvent(k.domain, k.subject, body, map[string]any{"i": fmt.Sprint(i)})
		if err != nil {
			panic(err)
		}
		pb, err := chipingress.EventToProto(ev)
		if err != nil {
			panic(err)
		}
		raw, err := proto.Marshal(pb)
		if err != nil {
			panic(err)
		}
		counts[k.domain+"/"+k.subject]++
		fmt.Fprintf(w, "\\\\x%s\tnow() - interval '%s'\n", hex.EncodeToString(raw), *age)
	}
	for i := 0; i < *corrupt; i++ {
		fmt.Fprintf(w, "\\\\x%s\tnow() - interval '%s'\n", hex.EncodeToString([]byte{0xff, 0xfe, byte(i), 0x00, 0x99}), *age)
	}
	if err := w.Flush(); err != nil {
		panic(err)
	}
	f.Close()

	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Printf("generated %d decodable + %d corrupt payloads -> %s\n", *n, *corrupt, *out)
	for _, k := range keys {
		fmt.Printf("  %-40s %6d\n", k, counts[k])
	}
}
