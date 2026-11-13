// Package jobsync reconciles the node's jobs against the Dynamic Config file
// the dynamic-config-sidecar renders. It replaces the Job Distributor
// websocket as the way a NOP delivers jobs: the sidecar verifies the signed
// config, writes it to disk and calls POST /reload, and this package makes the
// node's jobs match it.
//
// Only jobs whose name starts with the configured prefix are managed, so the
// node keeps every job an operator created by hand.
package jobsync

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"gopkg.in/guregu/null.v4"
)

const (
	// hashKey is the job KV key holding the hash of the TOML that created the
	// job, so a reload that changes nothing does not restart it.
	hashKey = "dynamic-config/toml-hash"
	// listLimit bounds the jobs read per reconcile; the node is not expected to
	// hold more.
	listLimit = 10000
)

var (
	// ErrInvalidConfig is returned when the config file cannot be applied. The
	// node keeps its current jobs.
	ErrInvalidConfig = errors.New("invalid dynamic config")
	// ErrDuplicateJob is returned when two entries share a name.
	ErrDuplicateJob = errors.New("duplicate job name")
)

// Config is the JSON document the sidecar renders.
type Config struct {
	Revision uint64 `json:"revision"`
	Jobs     []Job  `json:"jobs"`
}

// Job is one job the node must run.
type Job struct {
	// Name identifies the job and must start with the managed prefix.
	Name string `json:"name"`
	// TOML is the job spec, exactly as POST /v2/jobs takes it.
	TOML string `json:"toml"`
}

// Result reports what a reconcile changed.
type Result struct {
	Revision  uint64
	Created   []string
	Updated   []string
	Deleted   []string
	Unchanged []string
}

// Applied is the outcome of a reload: what changed and the exact bytes that
// were applied, so the caller can report the hash the sidecar verifies.
type Applied struct {
	Result Result
	Raw    []byte
}

// Validator turns a job spec in TOML into a job.Job.
type Validator interface {
	Validate(ctx context.Context, tomlString string) (job.Job, error)
}

// hashStore remembers, per job, the hash of the TOML that created it.
type hashStore interface {
	Get(ctx context.Context, jobID int32) (string, error)
	Store(ctx context.Context, jobID int32, hash string) error
}

// kvHashStore keeps the hashes in the job KV store, so they survive a restart
// with the node's database.
type kvHashStore struct {
	ds sqlutil.DataSource
}

func (s kvHashStore) Get(ctx context.Context, jobID int32) (string, error) {
	raw, err := job.NewKVStore(jobID, s.ds).Get(ctx, hashKey)

	return string(raw), err
}

func (s kvHashStore) Store(ctx context.Context, jobID int32, hash string) error {
	return job.NewKVStore(jobID, s.ds).Store(ctx, hashKey, []byte(hash))
}

// Reconciler applies a Config to the node's jobs.
type Reconciler struct {
	orm      job.ORM
	spawner  job.Spawner
	validate Validator
	hashes   hashStore
	prefix   string
	logger   commonlogger.Logger
}

// New returns a Reconciler that manages jobs whose name starts with prefix.
func New(orm job.ORM, spawner job.Spawner, validate Validator, prefix string, logger commonlogger.Logger) *Reconciler {
	if logger == nil {
		logger = commonlogger.Nop()
	}

	return &Reconciler{
		orm: orm, spawner: spawner, validate: validate,
		hashes: kvHashStore{ds: orm.DataSource()}, prefix: prefix, logger: logger,
	}
}

// Reload reads the config at path and applies it. It returns ErrInvalidConfig
// without touching any job when the file or any spec is invalid, so a bad
// revision never leaves the node half-configured.
func (r *Reconciler) Reload(ctx context.Context, path string) (Applied, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Applied{}, fmt.Errorf("%w: read %s: %w", ErrInvalidConfig, path, err)
	}
	cfg, err := Parse(raw)
	if err != nil {
		return Applied{}, err
	}
	res, err := r.Apply(ctx, cfg)
	if err != nil {
		return Applied{}, err
	}

	return Applied{Result: res, Raw: raw}, nil
}

// Parse decodes and validates the shape of a config document.
func Parse(raw []byte) (Config, error) {
	var cfg Config
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}
	if cfg.Revision == 0 {
		return Config{}, fmt.Errorf("%w: revision must be >= 1", ErrInvalidConfig)
	}
	seen := make(map[string]struct{}, len(cfg.Jobs))
	for _, j := range cfg.Jobs {
		if j.Name == "" {
			return Config{}, fmt.Errorf("%w: a job has no name", ErrInvalidConfig)
		}
		if _, ok := seen[j.Name]; ok {
			return Config{}, fmt.Errorf("%w: %s", ErrDuplicateJob, j.Name)
		}
		seen[j.Name] = struct{}{}
	}

	return cfg, nil
}

// Apply makes the node's managed jobs match cfg.
func (r *Reconciler) Apply(ctx context.Context, cfg Config) (Result, error) {
	desired, err := r.validateAll(ctx, cfg)
	if err != nil {
		return Result{}, err
	}

	existing, err := r.managedJobs(ctx)
	if err != nil {
		return Result{}, err
	}

	res := Result{Revision: cfg.Revision}
	for name, want := range desired {
		have, ok := existing[name]
		if !ok {
			if err := r.create(ctx, want); err != nil {
				return res, err
			}
			res.Created = append(res.Created, name)

			continue
		}
		if have.hash == want.hash {
			res.Unchanged = append(res.Unchanged, name)

			continue
		}
		if err := r.replace(ctx, have.id, want); err != nil {
			return res, err
		}
		res.Updated = append(res.Updated, name)
	}
	for name, have := range existing {
		if _, ok := desired[name]; ok {
			continue
		}
		if err := r.spawner.DeleteJob(ctx, nil, have.id); err != nil {
			return res, fmt.Errorf("delete job %s: %w", name, err)
		}
		res.Deleted = append(res.Deleted, name)
	}

	r.logger.Infow("dynamic config applied",
		"revision", cfg.Revision,
		"created", len(res.Created), "updated", len(res.Updated),
		"deleted", len(res.Deleted), "unchanged", len(res.Unchanged))

	return res, nil
}

// desiredJob is a validated job and the hash of the TOML that produced it.
type desiredJob struct {
	spec job.Job
	hash string
}

// validateAll validates every entry before anything is applied, so one bad
// spec cannot leave the node with a partial config.
func (r *Reconciler) validateAll(ctx context.Context, cfg Config) (map[string]desiredJob, error) {
	desired := make(map[string]desiredJob, len(cfg.Jobs))
	for _, j := range cfg.Jobs {
		if !strings.HasPrefix(j.Name, r.prefix) {
			return nil, fmt.Errorf("%w: job %q does not start with the managed prefix %q", ErrInvalidConfig, j.Name, r.prefix)
		}
		spec, err := r.validate.Validate(ctx, j.TOML)
		if err != nil {
			return nil, fmt.Errorf("%w: job %s: %w", ErrInvalidConfig, j.Name, err)
		}
		if spec.Name.Valid && spec.Name.String != j.Name {
			return nil, fmt.Errorf("%w: job %s: TOML name is %q", ErrInvalidConfig, j.Name, spec.Name.String)
		}
		spec.Name = null.StringFrom(j.Name)
		desired[j.Name] = desiredJob{spec: spec, hash: hashTOML(j.TOML)}
	}

	return desired, nil
}

// existingJob is a managed job already in the node.
type existingJob struct {
	id   int32
	hash string
}

// managedJobs lists the jobs whose name starts with the prefix, with the hash
// of the TOML that created each one (empty when unknown).
func (r *Reconciler) managedJobs(ctx context.Context) (map[string]existingJob, error) {
	jobs, _, err := r.orm.FindJobs(ctx, 0, listLimit)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	out := make(map[string]existingJob, len(jobs))
	for _, j := range jobs {
		if !j.Name.Valid || !strings.HasPrefix(j.Name.String, r.prefix) {
			continue
		}
		hash, err := r.hashes.Get(ctx, j.ID)
		if err != nil {
			// No record: the job predates this reconciler or its state was
			// lost. Treat the hash as unknown so the next apply converges.
			r.logger.Debugw("no recorded hash for managed job", "job", j.Name.String, "job_id", j.ID)
		}
		out[j.Name.String] = existingJob{id: j.ID, hash: hash}
	}

	return out, nil
}

func (r *Reconciler) create(ctx context.Context, want desiredJob) error {
	if err := r.spawner.CreateJob(ctx, nil, &want.spec); err != nil {
		return fmt.Errorf("create job %s: %w", want.spec.Name.String, err)
	}
	r.recordHash(ctx, want.spec.ID, want.hash)

	return nil
}

func (r *Reconciler) replace(ctx context.Context, id int32, want desiredJob) error {
	if err := r.spawner.DeleteJob(ctx, nil, id); err != nil {
		return fmt.Errorf("delete job %s before update: %w", want.spec.Name.String, err)
	}

	return r.create(ctx, want)
}

// recordHash stores the applied hash on the job. A failure only costs an extra
// restart on the next reload, so it is logged, not returned.
func (r *Reconciler) recordHash(ctx context.Context, jobID int32, hash string) {
	if err := r.hashes.Store(ctx, jobID, hash); err != nil {
		r.logger.Warnw("cannot record the applied job hash", "job_id", jobID, "error", err)
	}
}

func hashTOML(toml string) string {
	sum := sha256.Sum256([]byte(toml))

	return hex.EncodeToString(sum[:])
}
