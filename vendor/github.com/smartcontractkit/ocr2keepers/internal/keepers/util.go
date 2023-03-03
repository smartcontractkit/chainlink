package keepers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/cmplx"
	rnd "math/rand"
	"sort"
	"strings"
	"sync"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2keepers/internal/util"
	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	ktypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
	"golang.org/x/crypto/sha3"
)

var (
	ErrNotEnoughInputs = fmt.Errorf("not enough inputs")
)

func filterUpkeeps(upkeeps ktypes.UpkeepResults, filter ktypes.UpkeepState) ktypes.UpkeepResults {
	ret := make(ktypes.UpkeepResults, 0, len(upkeeps))

	for _, up := range upkeeps {
		if up.State == filter {
			ret = append(ret, up)
		}
	}

	return ret
}

func keyList(upkeeps ktypes.UpkeepResults) []ktypes.UpkeepKey {
	ret := make([]ktypes.UpkeepKey, len(upkeeps))

	for i, up := range upkeeps {
		ret[i] = up.Key
	}

	sort.Sort(sortUpkeepKeys(ret))

	return ret
}

type shuffler[T any] interface {
	Shuffle([]T) []T
}

type cryptoShuffler[T any] struct{}

func (_ *cryptoShuffler[T]) Shuffle(a []T) []T {
	r := rnd.New(util.NewCryptoRandSource())
	r.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

type sortUpkeepKeys []ktypes.UpkeepKey

func (s sortUpkeepKeys) Less(i, j int) bool {
	return s[i].String() < s[j].String()
}

func (s sortUpkeepKeys) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortUpkeepKeys) Len() int {
	return len(s)
}

func filterAndDedupe[T fmt.Stringer](inputs [][]T, filters ...func(T) bool) ([]T, error) {
	var max int
	for _, input := range inputs {
		max += len(input)
	}

	output := make([]T, 0, max)
	matched := make(map[string]struct{})
	for _, input := range inputs {
		for _, val := range input {
			add := true
			for _, filter := range filters {
				if !filter(val) {
					add = false
					break
				}
			}

			if !add {
				continue
			}

			key := val.String()
			_, ok := matched[key]
			if !ok {
				matched[key] = struct{}{}
				output = append(output, val)
			}
		}
	}

	return output, nil
}

func filterDedupeShuffleObservations(upkeepKeys [][]ktypes.UpkeepKey, keyRandSource [16]byte, filters ...func(ktypes.UpkeepKey) bool) ([]ktypes.UpkeepKey, error) {
	uniqueKeys, err := filterAndDedupe(upkeepKeys, filters...)
	if err != nil {
		return nil, err
	}

	rnd.New(util.NewKeyedCryptoRandSource(keyRandSource)).Shuffle(len(uniqueKeys), func(i, j int) {
		uniqueKeys[i], uniqueKeys[j] = uniqueKeys[j], uniqueKeys[i]
	})

	return uniqueKeys, nil
}

func shuffleObservations(upkeepIdentifiers []ktypes.UpkeepIdentifier, keyRandSource [16]byte) []ktypes.UpkeepIdentifier {
	rnd.New(util.NewKeyedCryptoRandSource(keyRandSource)).Shuffle(len(upkeepIdentifiers), func(i, j int) {
		upkeepIdentifiers[i], upkeepIdentifiers[j] = upkeepIdentifiers[j], upkeepIdentifiers[i]
	})

	return upkeepIdentifiers
}

func observationsToUpkeepKeys(logger *log.Logger, observations []types.AttributedObservation, reportBlockLag int) ([][]ktypes.UpkeepKey, error) {
	var parseErrors int

	upkeepIDs := make([][]ktypes.UpkeepIdentifier, len(observations))

	var allBlockKeys []*big.Int
	for i, observation := range observations {
		// a single observation returning an error here can void all other
		// good observations. ensure this loop continues on error, but collect
		// them and throw an error if ALL observations fail at this point.
		// TODO we can't rely on this concrete type for decoding/encoding
		var upkeepObservation *chain.UpkeepObservation
		if err := decode(observation.Observation, &upkeepObservation); err != nil {
			logger.Printf("unable to decode observation: %s", err.Error())
			parseErrors++
			continue
		}

		blockKeyInt, ok := big.NewInt(0).SetString(upkeepObservation.BlockKey.String(), 10)
		if !ok {
			parseErrors++
			continue
		}
		allBlockKeys = append(allBlockKeys, blockKeyInt)

		// if we have a non-empty list of upkeep identifiers, limit the upkeeps we take to observationUpkeepsLimit
		if len(upkeepObservation.UpkeepIdentifiers) > 0 {
			upkeepIDs[i] = upkeepObservation.UpkeepIdentifiers[:observationUpkeepsLimit]
		}
	}

	if parseErrors == len(observations) {
		return nil, fmt.Errorf("%w: cannot prepare sorted key list; observations not properly encoded", ErrTooManyErrors)
	}

	// Here we calculate the median block that will be applied for all upkeep keys.
	// reportBlockLag is subtracted from the median block to ensure enough nodes have that block in their blockchain
	medianBlock := calculateMedianBlock(allBlockKeys, reportBlockLag)
	logger.Printf("calculated median block %s, accounting for reportBlockLag of %d", medianBlock.String(), reportBlockLag)

	upkeepKeys, err := createKeysWithMedianBlock(medianBlock, upkeepIDs)
	if err != nil {
		return nil, err
	}

	return upkeepKeys, nil
}

func createKeysWithMedianBlock(medianBlock ktypes.BlockKey, upkeepIDLists [][]ktypes.UpkeepIdentifier) ([][]ktypes.UpkeepKey, error) {
	var res = make([][]ktypes.UpkeepKey, len(upkeepIDLists))

	for i, upkeepIDs := range upkeepIDLists {
		var keys []ktypes.UpkeepKey
		for _, upkeepID := range upkeepIDs {
			keys = append(keys, chain.NewUpkeepKeyFromBlockAndID(medianBlock, upkeepID))
		}
		res[i] = keys
	}

	return res, nil
}

func calculateMedianBlock(blockNumbers []*big.Int, reportBlockLag int) ktypes.BlockKey {
	sort.Slice(blockNumbers, func(i, j int) bool {
		return blockNumbers[i].Cmp(blockNumbers[j]) < 0
	})

	// this is a crude median calculation; for a list of an odd number of elements, e.g. [10, 20, 30], the center value
	// is chosen as the median. for a list of an even number of elements, a true median calculation would average the
	// two center elements, e.g. [10, 20, 30, 40] = (20 + 30) / 2 = 25, but we want to constrain our median block to
	// one of the block numbers reported, e.g. either 20 or 30. right now we want to choose the higher block number, e.g.
	// 30. for this reason, the logic for selecting the median value from an odd number of elements is the same as the
	// logic for selecting the median value from an even number of elements
	var median *big.Int
	if l := len(blockNumbers); l == 0 {
		median = big.NewInt(0)
	} else {
		median = blockNumbers[l/2]
	}

	if reportBlockLag > 0 {
		median = median.Sub(median, big.NewInt(int64(reportBlockLag)))
	}

	return chain.BlockKey(median.String())
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

func lowest(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	return values[0]
}

type syncedArray[T any] struct {
	data []T
	mu   sync.RWMutex
}

func newSyncedArray[T any]() *syncedArray[T] {
	return &syncedArray[T]{
		data: []T{},
	}
}

func (a *syncedArray[T]) Append(vals ...T) *syncedArray[T] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.data = append(a.data, vals...)
	return a
}

func (a *syncedArray[T]) Values() []T {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.data
}

func limitedLengthEncode(obs *chain.UpkeepObservation, limit int) ([]byte, error) {
	if len(obs.UpkeepIdentifiers) == 0 {
		return encode(obs)
	}

	var res []byte
	for i := range obs.UpkeepIdentifiers {
		b, err := encode(&chain.UpkeepObservation{
			BlockKey:          obs.BlockKey,
			UpkeepIdentifiers: obs.UpkeepIdentifiers[:i+1],
		})
		if err != nil {
			return nil, err
		}
		if len(b) > limit {
			break
		}
		res = b
	}

	return res, nil
}

// Generates a randomness source derived from the report timestamp (config, epoch, round) so
// that it's the same across the network for the same round
func getRandomKeySource(rt types.ReportTimestamp) [16]byte {
	// similar key building as libocr transmit selector
	hash := sha3.NewLegacyKeccak256()
	hash.Write(rt.ConfigDigest[:])
	temp := make([]byte, 8)
	binary.LittleEndian.PutUint64(temp, uint64(rt.Epoch))
	hash.Write(temp)
	binary.LittleEndian.PutUint64(temp, uint64(rt.Round))
	hash.Write(temp)

	var keyRandSource [16]byte
	copy(keyRandSource[:], hash.Sum(nil))
	return keyRandSource
}

func upkeepKeysToString(keys []ktypes.UpkeepKey) string {
	keysStr := make([]string, len(keys))
	for i, key := range keys {
		keysStr[i] = key.String()
	}

	return strings.Join(keysStr, ", ")
}

func upkeepIdentifiersToString(ids []ktypes.UpkeepIdentifier) string {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = string(id)
	}

	return strings.Join(idsStr, ", ")
}

func createBatches[T any](b []T, size int) (batches [][]T) {
	for i := 0; i < len(b); i += size {
		j := i + size
		if j > len(b) {
			j = len(b)
		}
		batches = append(batches, b[i:j])
	}
	return
}

// buffer is a goroutine safe bytes.Buffer
type buffer struct {
	buffer bytes.Buffer
	mutex  sync.Mutex
}

// Write appends the contents of p to the buffer, growing the buffer as needed. It returns
// the number of bytes written.
func (s *buffer) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.Write(p)
}

// String returns the contents of the unread portion of the buffer
// as a string.  If the Buffer is a nil pointer, it returns "<nil>".
func (s *buffer) String() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.String()
}
