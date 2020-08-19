package job

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

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
	ID           uint64                 `json:"-" gorm:"primary_key;auto_increment"`
	BridgeName   string                 `json:"name"`
	RequestData  map[string]interface{} `json:"requestData"`
	Transformers Transformers           `json:"transformPipeline"`
}

func (f BridgeFetcher) Fetch() (interface{}, error) {
	// ...

	return f.Transformers.Run(nil)
}

type HttpFetcher struct {
	ID           uint64                 `json:"-" gorm:"primary_key;auto_increment"`
	URL          string                 `json:"url"`
	Method       string                 `json:"method"`
	RequestData  map[string]interface{} `json:"requestData"`
	Transformers Transformers           `json:"transformPipeline"`
}

func (f HttpFetcher) Fetch() (interface{}, error) {
	// ...

	return f.Transformers.Run(nil)
}

type MedianFetcher struct {
	ID       uint64   `json:"-" gorm:"primary_key;auto_increment"`
	Fetchers Fetchers `json:"fetchers"`
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

type Fetchers []Fetcher

func (f *Fetchers) UnmarshalJSON(bs []byte) (err error) {
	defer withStack(&err)

	var spec []json.RawMessage
	err = json.Unmarshal(bs, &spec)
	if err != nil {
		return err
	}

	for _, fetcherBytes := range spec {
		fetcher, err := UnmarshalFetcherJSON([]byte(fetcherBytes))
		if err != nil {
			return err
		}
		*f = append(*f, fetcher)
	}
	return nil
}

func UnmarshalFetcherJSON(bs []byte) (_ Fetcher, err error) {
	defer withStack(&err)

	var header struct {
		Type FetcherType `json:"type"`
	}
	err = json.Unmarshal(bs, &header)
	if err != nil {
		return nil, err
	}

	var fetcher Fetcher
	switch header.Type {
	case FetcherTypeBridge:
		bridgeFetcher := BridgeFetcher{}
		err = json.Unmarshal(bs, &bridgeFetcher)
		if err != nil {
			return nil, err
		}
		fetcher = bridgeFetcher

	case FetcherTypeHttp:
		httpFetcher := HttpFetcher{}
		err = json.Unmarshal(bs, &httpFetcher)
		if err != nil {
			return nil, err
		}
		fetcher = httpFetcher

	case FetcherTypeMedian:
		medianFetcher := MedianFetcher{}
		err = json.Unmarshal(bs, &medianFetcher)
		if err != nil {
			return nil, err
		}
		fetcher = medianFetcher

	default:
		return nil, errors.New("unknown fetcher type")
	}

	return fetcher, nil
}
