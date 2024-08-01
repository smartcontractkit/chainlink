package polling

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"strconv"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"
	"github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
)

// PollingObserverFactory ...
type PollingObserverFactory struct {
	Logger  *log.Logger
	Source  UpkeepProvider
	Heads   HeadProvider
	Runner  Runner
	Encoder Encoder
}

// NewConditionalObserver ...
func (f *PollingObserverFactory) NewConditionalObserver(oc config.OffchainConfig, c types.ReportingPluginConfig, coord ocr2keepers.Coordinator) (ocr2keepers.ConditionalObserver, error) {
	var (
		p      float64
		err    error
		sample sampleRatio
	)

	p, err = strconv.ParseFloat(oc.TargetProbability, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse configured probability", err)
	}

	sample, err = sampleFromProbability(oc.TargetInRounds, c.N-c.F, float32(p))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create plugin", err)
	}

	ob := NewPollingObserver(
		f.Logger,
		f.Source,
		f.Heads,
		f.Runner,
		f.Encoder,
		sample,
		time.Duration(oc.SamplingJobDuration)*time.Millisecond,
		coord,
		oc.MercuryLookup,
	)

	return ob, nil
}

func sampleFromProbability(rounds, nodes int, probability float32) (sampleRatio, error) {
	var ratio sampleRatio

	if rounds <= 0 {
		return ratio, fmt.Errorf("number of rounds must be greater than 0")
	}

	if nodes <= 0 {
		return ratio, fmt.Errorf("number of nodes must be greater than 0")
	}

	if probability > 1 || probability <= 0 {
		return ratio, fmt.Errorf("probability must be less than 1 and greater than 0")
	}

	r := complex(float64(rounds), 0)
	n := complex(float64(nodes), 0)
	p := complex(float64(probability), 0)

	g := -1.0 * (p - 1.0)
	x := cmplx.Pow(cmplx.Pow(g, 1.0/r), 1.0/n)
	rat := cmplx.Abs(-1.0 * (x - 1.0))
	rat = math.Round(rat/0.01) * 0.01
	ratio = sampleRatio(float32(rat))

	return ratio, nil
}

type sampleRatio float32

func (r sampleRatio) OfInt(count int) int {
	// rounds the result using basic rounding op
	return int(math.Round(float64(r) * float64(count)))
}

func (r sampleRatio) String() string {
	return fmt.Sprintf("%.8f", float32(r))
}
