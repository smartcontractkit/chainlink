package job

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type Fetchers []Fetcher

func (f *Fetchers) UnmarshalJSON(bs []byte) error {
	var spec struct {
		Fetchers []json.RawMessage `json:"fetchers"`
	}
	err := json.Unmarshal(bs, &spec)
	if err != nil {
		return err
	}

	for _, fetcherBytes := range spec.Fetchers {
		var fetcherSpec struct {
			Type FetcherType `json:"type"`
		}
		err := json.Unmarshal([]byte(fetcherBytes), &fetcherSpec)
		if err != nil {
			return err
		}

		var fetcher Fetcher
		switch fetcherSpec.Type {
		case FetcherTypeBridge:
			fetcher = BridgeFetcher{}
		case FetcherTypeHttp:
			fetcher = HttpFetcher{}
		case FetcherTypeMedian:
			fetcher = MedianFetcher{}
		default:
			return errors.New("unknown fetcher type")
		}

		err = json.Unmarshal([]byte(fetcherBytes), &fetcher)
		if err != nil {
			return err
		}

		j.Fetchers = append(j.Fetchers, fetcher)
	}
}

type Fetcher interface {
	Fetch() (interface{}, error)
}

type FetcherType string

var (
	FetcherTypeBridge FetcherType = "bridge"
	FetcherTypeHttp   FetcherType = "http"
	FetcherTypeMedian FetcherType = "median"
)

type BridgeFetcher struct {
	Transformers Transformers           `json:"transformPipeline"`
	BridgeName   string                 `json:"name"`
	RequestData  map[string]interface{} `json:"requestData"`
}

func (f BridgeFetcher) Fetch() (interface{}, error) {
	// ...

	return f.Transformers.Run(value)
}

type HttpFetcher struct {
	Transformers Transformers           `json:"transformPipeline"`
	URL          string                 `json:"url"`
	Method       string                 `json:"method"`
	RequestData  map[string]interface{} `json:"requestData"`
}

func (f HttpFetcher) Fetch() (interface{}, error) {
	// ...

	return f.Transformers.Run(value)
}

type MedianFetcher struct {
	Fetchers []Fetchers `json:"fetchers"`
}

func (f MedianFetcher) Fetch() (interface{}, error) {
	answers := []decimal.Decimal{}
	fetchErrors := []error{}

	type result struct {
		answer decimal.Decimal
		err    error
	}

	chResults := make(chan result)
	for _, fetcher := range f.Fetchers {
		fetcher := fetcher
		go func() {
			fetchedVal, err := fetcher.Fetch()
			if err != nil {
				logger.Error(err)
				chResults <- result{err: err}
				return
			}

			answer, err := utils.ToDecimal(fetchedVal)
			if err != nil {
				logger.Error(err)
				chResults <- result{err: err}
				return
			}

			chResults <- result{answer: answer}
		}()
	}

	for i := 0; i < len(f.Fetchers); i++ {
		r := <-chResults
		if r.err != nil {
			fetchErrors = append(fetchErrors, r.err)
		} else {
			answers = append(answers, r.answer)
		}
	}

	errorRate := float64(len(fetchErrors)) / float64(len(f.Fetchers))
	if errorRate >= 0.5 {
		return decimal.Decimal{}, errors.Wrap(multierr.Combine(fetchErrors...), "majority of fetchers in median failed")
	}

	sort.Slice(answers, func(i, j int) bool {
		return answers[i].LessThan(answers[j])
	})
	k := len(answers) / 2
	if len(answers)%2 == 1 {
		return answers[k], nil
	}
	return answers[k].Add(answers[k-1]).Div(decimal.NewFromInt(2)), nil
}
