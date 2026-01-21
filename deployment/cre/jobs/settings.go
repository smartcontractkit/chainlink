package jobs

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/pelletier/go-toml"
	toml2 "github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
)

func verifyCRESettingsSpecInputs(inputs job_types.JobSpecInput) error {
	ji := &operations.ProposeCRESettingsJobsInput{}
	if err := inputs.UnmarshalTo(ji); err != nil {
		return fmt.Errorf("failed to unmarshal job spec input to StandardCapabilityJob: %w", err)
	}

	return VerifyCRESettings(ji.Settings)
}

// VerifyCRESettings ensures that each of the scoped overrides match the cresettings.Schema and every field value is a string,
// or returns an error.
func VerifyCRESettings(s string) error {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var data struct {
		Globals  *toml.Tree           `toml:"global"`
		Org      map[string]toml.Tree `toml:"org"`
		Owner    map[string]toml.Tree `toml:"owner"`
		Workflow map[string]toml.Tree `toml:"workflow"`
	}

	if err := toml.NewDecoder(strings.NewReader(s)).Decode(&data); err != nil {
		return fmt.Errorf("invalid toml settings: %w", err)
	}
	var errs error
	if data.Globals != nil {
		if err := ensureSchema(*data.Globals); err != nil {
			errs = errors.Join(errs, fmt.Errorf("invalid global: %w", err))
		}
	}
	for id, org := range data.Org {
		// TODO to be enforced after upgrading deployed nodes and settings
		// if strings.HasPrefix(id, "org_") {
		// 	 errs = errors.Join(errs, fmt.Errorf("invalid org id %s: must not be prefixed org_", id))
		// } else if strings.ToLower(id) != id {
		//   errs = errors.Join(errs, fmt.Errorf("invalid org id %s: must be lower case", id))
		if strings.HasPrefix(id, "0x") {
			errs = errors.Join(errs, fmt.Errorf("invalid org id %s: must not be prefixed 0x", id))
		} else if err := ensureSchema(org); err != nil {
			errs = errors.Join(errs, fmt.Errorf("invalid org %s: %w", id, err))
		}
	}
	for id, owner := range data.Owner {
		if strings.HasPrefix(id, "owner_") {
			errs = errors.Join(errs, fmt.Errorf("invalid owner id %s: must not be prefixed owner_", id))
		} else if strings.HasPrefix(id, "0x") {
			errs = errors.Join(errs, fmt.Errorf("invalid owner id %s: must not be prefixed 0x", id))
		} else if strings.ToLower(id) != id {
			errs = errors.Join(errs, fmt.Errorf("invalid owner id %s: must be lower case", id))
		} else if err := ensureSchema(owner); err != nil {
			errs = errors.Join(errs, fmt.Errorf("invalid owner %s: %w", id, err))
		}
	}
	for id, wf := range data.Workflow {
		if strings.HasPrefix(id, "0x") {
			errs = errors.Join(errs, fmt.Errorf("invalid wf id %s: must not be prefixed 0x", id))
		} else if strings.ToLower(id) != id {
			errs = errors.Join(errs, fmt.Errorf("invalid wf id %s: must be lower case", id))
		} else if err := ensureSchema(wf); err != nil {
			errs = errors.Join(errs, fmt.Errorf("invalid wf %s: %w", id, err))
		}
	}
	return errs
}

func ensureSchema(tree toml.Tree) (errs error) {
	if err := ensureStrings(nil, tree); err != nil {
		return err
	}
	b, err := tree.Marshal()
	if err != nil {
		return fmt.Errorf("failed to re-marshal: %w", err)
	}
	schema := cresettings.Default // copy
	if err := toml2.NewDecoder(strings.NewReader(string(b))).DisallowUnknownFields().Decode(&schema); err != nil {
		var missingErr *toml2.StrictMissingError
		if errors.As(err, &missingErr) {
			return fmt.Errorf("unknown fields - if these are new fields, then chainlink-common must be updated:\n%s", missingErr.String())
		}
		return fmt.Errorf("invalid toml settings: %w", err)
	}
	return nil
}

func ensureStrings(parent []string, tree toml.Tree) (errs error) {
	for k, v := range tree.Values() {
		key := slices.Clone(parent)
		key = append(key, k)
		switch t := v.(type) {
		case *toml.PubTOMLValue:
			val := t.Value()
			if _, ok := val.(string); !ok {
				errs = errors.Join(errs, fmt.Errorf("%s is not a string: %T", strings.Join(key, "."), val))
			}
		case *toml.PubTree:
			if err := ensureStrings(key, *t); err != nil {
				errs = errors.Join(errs, err)
			}
		case []*toml.PubTree:
			errs = errors.Join(errs, fmt.Errorf("unsupported type for %s: list", key))
		default:
			errs = errors.Join(errs, fmt.Errorf("unsupported type for %s: %T", key, v))
		}
	}
	return
}
