package devmachine

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

// Load implements phase P5: load the image tarballs into Docker on the remote
// machine and retag the job-distributor image to the name the environment
// expects.
func (o *Orchestrator) Load(ctx context.Context) error {
	o.rep.Phase(phaseTitle(phaseLoad))

	if o.cfg.Machine.InstanceID == "" {
		return errors.New("no instance id available; run phase P1 first")
	}

	jdTar := filepath.Base(tarName(o.cfg.Images.JobDistributor))
	clTar := filepath.Base(tarName(o.cfg.Images.Chainlink))

	script := fmt.Sprintf(`set -e
docker load -i %s
docker tag %s %s
docker load -i %s
docker images`,
		jdTar, o.cfg.Images.JobDistributor, o.cfg.Images.JDLoadTag, clTar)

	o.rep.Info("loading images on %s", o.cfg.SSH.Alias)
	return o.run.Stream(ctx, "ssh", o.cfg.SSH.Alias, script)
}
