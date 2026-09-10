package matrix

import "errors"

// runsOnParams holds resolved parameters used to construct runs-on labels.
type runsOnParams struct {
	RunID      string
	RunAttempt string
	SpotFlag   string
}

// resolveRunsOnParams validates common runs-on label parameters and applies
// defaults for run attempt and spot flag. Run ID has no default because empty
// run IDs produce invalid runner labels.
func resolveRunsOnParams(runID, runAttempt, spotFlag string) (runsOnParams, error) {
	if runID == "" {
		return runsOnParams{}, errors.New("run ID is required to generate runs-on labels (provide --run-id, action input run_id, or GITHUB_RUN_ID)")
	}
	if runAttempt == "" {
		runAttempt = "1"
	}
	if spotFlag == "" {
		spotFlag = "spot=co"
	}
	return runsOnParams{RunID: runID, RunAttempt: runAttempt, SpotFlag: spotFlag}, nil
}
