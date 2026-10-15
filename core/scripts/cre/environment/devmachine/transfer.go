package devmachine

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
)

// Transfer implements phase P4: refresh the managed SSH host entry and copy the
// agent binary and image tarballs to the remote machine.
func (o *Orchestrator) Transfer(ctx context.Context) error {
	o.rep.Phase(phaseTitle(phaseTransfer))

	if o.cfg.Machine.InstanceID == "" {
		return errors.New("no instance id available; run phase P1 first")
	}
	if err := o.ensureSSHConfig(); err != nil {
		return err
	}

	cache := o.cacheDir
	if cache == "" {
		var err error
		cache, err = CacheDir()
		if err != nil {
			return err
		}
	}

	tars, err := filepath.Glob(filepath.Join(cache, "*.tar"))
	if err != nil {
		return errors.Wrap(err, "failed to list cached image tarballs")
	}

	files := make([]string, 0, 1+len(tars))
	files = append(files, filepath.Join(cache, "remote-agent"))
	files = append(files, tars...)

	args := make([]string, 0, 4+len(files)+1)
	args = append(args, "-avz", "--progress", "-e", "ssh -o StrictHostKeyChecking=accept-new")
	args = append(args, files...)
	args = append(args, o.cfg.SSH.Alias+":~/")

	o.rep.Info("copying %d file(s) to %s", len(files), o.cfg.SSH.Alias)
	return o.run.Stream(ctx, "rsync", args...)
}
