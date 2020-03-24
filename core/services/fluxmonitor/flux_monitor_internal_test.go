package fluxmonitor

import (
	"sort"
	"testing"

	"chainlink/core/store/models"

	"github.com/stretchr/testify/require"
)

type jobAdder []string

func (j *jobAdder) AddJob(arg models.JobSpec) error {
	*j = append(*j, arg.ID.String())
	return nil
}

func TestBackgroundJobLoader(t *testing.T) {
	id1 := models.NewID()
	id2 := models.NewID()

	tests := []struct {
		Name string
		Jobs []*models.JobSpec
		Exp  []string
	}{
		{
			Name: "Nil ok",
			Jobs: nil,
			Exp:  []string{},
		},
		{
			Name: "All successful",
			Jobs: []*models.JobSpec{
				&models.JobSpec{ID: id1},
				&models.JobSpec{ID: id2},
			},
			Exp: []string{
				id1.String(),
				id2.String(),
			},
		},
		{
			Name: "Expect nil skipped",
			Jobs: []*models.JobSpec{
				nil,
				&models.JobSpec{ID: id2},
			},
			Exp: []string{
				id2.String(),
			},
		},
	}
	for _, test := range tests {
		t.Log(test.Name)

		ret := &jobAdder{}
		obj := backgroundJobLoader{fm: ret}
		for _, j := range test.Jobs {
			ok := obj.AddJob(j)
			require.True(t, ok)
		}
		obj.Wait()

		sort.Sort(sort.StringSlice(([]string)(test.Exp)))
		sort.Sort(sort.StringSlice(([]string)(*ret)))
		require.Equal(t, test.Exp, ([]string)(*ret))
	}
}
