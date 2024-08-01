package util

type Results struct {
	Successes int
	Failures  int
	Err       error
}

func (r *Results) Total() int {
	return r.Successes + r.Failures
}

func (r *Results) SuccessRate() float64 {
	if r.Total() == 0 {
		return 0
	}

	return float64(r.Successes) / float64(r.Total())
}

func (r *Results) FailureRate() float64 {
	if r.Total() == 0 {
		return 0
	}

	return float64(r.Failures) / float64(r.Total())
}
