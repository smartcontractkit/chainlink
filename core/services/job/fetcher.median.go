package job

import (
	"encoding/json"
	"sort"

	"github.com/shopspring/decimal"
)

type MedianFetcher struct {
	BaseFetcher

	ID        models.ID `json:"-" gorm:"primary_key;auto_increment"`
	JobSpecID models.ID `json:"jobSpecID"`
	Fetchers  Fetchers  `json:"fetchers"`
}

func (f *MedianFetcher) SetNotifiee(n Notifiee) {
	f.BaseFetcher.SetNotifiee(n)
	for _, fetcher := range f.Fetchers {
		fetcher.SetNotifiee(n)
	}
}

func (f MedianFetcher) Fetch() (interface{}, error) {
	defer func() { f.notifiee.OnEndStage(f, out, err) }()
	f.notifiee.OnBeginStage(f, nil)

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
	median := answers[k].Add(answers[k-1]).Div(decimal.NewFromInt(2))
	return f.Transformers.Transform(median)
}

func (f MedianFetcher) MarshalJSON() ([]byte, error) {
	type preventInfiniteRecursion MedianFetcher
	type fetcherWithType struct {
		Type FetcherType `json:"type"`
		preventInfiniteRecursion
	}
	f2 := fetcherWithType{FetcherTypeMedian, preventInfiniteRecursion(f)}
	return json.Marshal(f2)
}
