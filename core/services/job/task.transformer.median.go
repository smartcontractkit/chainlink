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

type MedianFetcher struct {
	BaseTask
}

func (f *MedianFetcher) Run(inputs []Result) (out interface{}, err error) {
	answers := []decimal.Decimal{}
	fetchErrors := []error{}

	for _, input := range inputs {
		if input.Error != nil {
			fetchErrors = append(fetchErrors, input.Error)
			continue
		}

		answer, err := utils.ToDecimal(input.Value)
		if err != nil {
			logger.Error(err)
			fetchErrors = append(fetchErrors, err)
			continue
		}

		answers = append(answers, answer)
	}

	errorRate := float64(len(fetchErrors)) / float64(len(answers)+len(fetchErrors))
	if errorRate >= 0.5 {
		return nil, errors.Wrap(multierr.Combine(fetchErrors...), "majority of fetchers in median failed")
	}

	sort.Slice(answers, func(i, j int) bool {
		return answers[i].LessThan(answers[j])
	})
	k := len(answers) / 2
	if len(answers)%2 == 1 {
		return answers[k], nil
	}
	median := answers[k].Add(answers[k-1]).Div(decimal.NewFromInt(2))
	return median, nil
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
