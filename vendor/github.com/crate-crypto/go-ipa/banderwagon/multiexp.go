package banderwagon

import (
	"github.com/crate-crypto/go-ipa/bandersnatch"
	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
)

// MultiExpConfig enables to set optional configuration attribute to a call to MultiExp
type MultiExpConfig struct {
	NbTasks     int  // go routines to be used in the multiexp. can be larger than num cpus.
	ScalarsMont bool // indicates if the scalars are in montgomery form. Default to false.
}

// MultiExp calculates the multi exponentiation of points and scalars.
func (p *Element) MultiExp(points []Element, scalars []fr.Element, config MultiExpConfig) (*Element, error) {
	var projPoints = make([]bandersnatch.PointProj, len(points))
	for i := range points {
		projPoints[i] = points[i].inner
	}
	affinePoints := batchProjToAffine(projPoints)

	// NOTE: This is fine as long MultiExp does not use Equal functionality
	_, err := bandersnatch.MultiExp(&p.inner, affinePoints, scalars, bandersnatch.MultiExpConfig{
		NbTasks:     config.NbTasks,
		ScalarsMont: config.ScalarsMont,
	})

	return p, err
}
