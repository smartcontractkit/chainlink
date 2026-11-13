package jobsync

import (
	"context"
	"errors"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	jobmocks "github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
)

// fakeValidator accepts any spec whose TOML carries a name, so the tests do not
// need a full node configuration.
type fakeValidator struct {
	err error
}

func (f fakeValidator) Validate(_ context.Context, tomlString string) (job.Job, error) {
	if f.err != nil {
		return job.Job{}, f.err
	}
	var jb job.Job
	if err := toml.Unmarshal([]byte(tomlString), &jb); err != nil {
		return job.Job{}, err
	}

	return jb, nil
}

// memHashes is an in-memory hashStore.
type memHashes map[int32]string

func (m memHashes) Get(_ context.Context, jobID int32) (string, error) {
	h, ok := m[jobID]
	if !ok {
		return "", errors.New("not found")
	}

	return h, nil
}

func (m memHashes) Store(_ context.Context, jobID int32, hash string) error {
	m[jobID] = hash

	return nil
}

func newTestReconciler(t *testing.T, existing []job.Job, hashes memHashes) (*Reconciler, *jobmocks.ORM, *jobmocks.Spawner) {
	t.Helper()
	orm := jobmocks.NewORM(t)
	spawner := jobmocks.NewSpawner(t)
	orm.On("FindJobs", mock.Anything, 0, listLimit).Return(existing, len(existing), nil).Maybe()

	return &Reconciler{
		orm: orm, spawner: spawner, validate: fakeValidator{},
		hashes: hashes, prefix: "dyncfg-", logger: logger.Test(t),
	}, orm, spawner
}

func TestParse(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		raw  string
		want error
	}{
		{"valid", `{"revision":2,"jobs":[{"name":"dyncfg-a","toml":"x"}]}`, nil},
		{"no revision", `{"jobs":[]}`, ErrInvalidConfig},
		{"unknown field", `{"revision":1,"jobs":[],"extra":true}`, ErrInvalidConfig},
		{"no name", `{"revision":1,"jobs":[{"toml":"x"}]}`, ErrInvalidConfig},
		{"duplicate", `{"revision":1,"jobs":[{"name":"a","toml":"x"},{"name":"a","toml":"y"}]}`, ErrDuplicateJob},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := Parse([]byte(tt.raw))
			if tt.want == nil {
				require.NoError(t, err)

				return
			}
			require.ErrorIs(t, err, tt.want)
		})
	}
}

func TestApplyCreatesMissingJobs(t *testing.T) {
	t.Parallel()
	r, _, spawner := newTestReconciler(t, nil, memHashes{})
	spawner.On("CreateJob", mock.Anything, mock.Anything, mock.MatchedBy(func(jb *job.Job) bool {
		return jb.Name.String == "dyncfg-a"
	})).Return(nil).Once()

	res, err := r.Apply(t.Context(), Config{Revision: 2, Jobs: []Job{{Name: "dyncfg-a", TOML: cronSpec("dyncfg-a")}}})
	require.NoError(t, err)
	assert.Equal(t, []string{"dyncfg-a"}, res.Created)
	assert.Empty(t, res.Updated)
	assert.Empty(t, res.Deleted)
}

func TestApplyLeavesUnchangedJobsAlone(t *testing.T) {
	t.Parallel()
	spec := cronSpec("dyncfg-a")
	existing := []job.Job{{ID: 7, Name: null.StringFrom("dyncfg-a")}}
	r, _, spawner := newTestReconciler(t, existing, memHashes{7: hashTOML(spec)})

	res, err := r.Apply(t.Context(), Config{Revision: 2, Jobs: []Job{{Name: "dyncfg-a", TOML: spec}}})
	require.NoError(t, err)
	assert.Equal(t, []string{"dyncfg-a"}, res.Unchanged)
	spawner.AssertNotCalled(t, "CreateJob", mock.Anything, mock.Anything, mock.Anything)
	spawner.AssertNotCalled(t, "DeleteJob", mock.Anything, mock.Anything, mock.Anything)
}

func TestApplyReplacesChangedJobs(t *testing.T) {
	t.Parallel()
	existing := []job.Job{{ID: 7, Name: null.StringFrom("dyncfg-a")}}
	r, _, spawner := newTestReconciler(t, existing, memHashes{7: "stale"})
	spawner.On("DeleteJob", mock.Anything, mock.Anything, int32(7)).Return(nil).Once()
	spawner.On("CreateJob", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	res, err := r.Apply(t.Context(), Config{Revision: 3, Jobs: []Job{{Name: "dyncfg-a", TOML: cronSpec("dyncfg-a")}}})
	require.NoError(t, err)
	assert.Equal(t, []string{"dyncfg-a"}, res.Updated)
}

func TestApplyDeletesJobsNoLongerListed(t *testing.T) {
	t.Parallel()
	existing := []job.Job{
		{ID: 7, Name: null.StringFrom("dyncfg-gone")},
		{ID: 8, Name: null.StringFrom("operator-job")},
	}
	r, _, spawner := newTestReconciler(t, existing, memHashes{})
	spawner.On("DeleteJob", mock.Anything, mock.Anything, int32(7)).Return(nil).Once()

	res, err := r.Apply(t.Context(), Config{Revision: 2})
	require.NoError(t, err)
	assert.Equal(t, []string{"dyncfg-gone"}, res.Deleted)
	spawner.AssertNotCalled(t, "DeleteJob", mock.Anything, mock.Anything, int32(8))
}

func TestApplyRejectsInvalidSpecWithoutTouchingJobs(t *testing.T) {
	t.Parallel()
	existing := []job.Job{{ID: 7, Name: null.StringFrom("dyncfg-gone")}}
	r, _, spawner := newTestReconciler(t, existing, memHashes{})
	r.validate = fakeValidator{err: errors.New("bad spec")}

	_, err := r.Apply(t.Context(), Config{Revision: 2, Jobs: []Job{{Name: "dyncfg-a", TOML: "x"}}})
	require.ErrorIs(t, err, ErrInvalidConfig)
	spawner.AssertNotCalled(t, "CreateJob", mock.Anything, mock.Anything, mock.Anything)
	spawner.AssertNotCalled(t, "DeleteJob", mock.Anything, mock.Anything, mock.Anything)
}

func TestApplyRejectsUnmanagedNames(t *testing.T) {
	t.Parallel()
	r, _, spawner := newTestReconciler(t, nil, memHashes{})

	_, err := r.Apply(t.Context(), Config{Revision: 2, Jobs: []Job{{Name: "operator-job", TOML: cronSpec("operator-job")}}})
	require.ErrorIs(t, err, ErrInvalidConfig)
	spawner.AssertNotCalled(t, "CreateJob", mock.Anything, mock.Anything, mock.Anything)
}

func TestApplyRejectsNameMismatch(t *testing.T) {
	t.Parallel()
	r, _, _ := newTestReconciler(t, nil, memHashes{})

	_, err := r.Apply(t.Context(), Config{Revision: 2, Jobs: []Job{{Name: "dyncfg-a", TOML: cronSpec("dyncfg-b")}}})
	require.ErrorIs(t, err, ErrInvalidConfig)
}

func cronSpec(name string) string {
	return "type = \"cron\"\nschemaVersion = 1\nname = \"" + name + "\"\nschedule = \"*/5 * * * *\"\n"
}
