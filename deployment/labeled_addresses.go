package deployment

import cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

// LabeledAddresses is an alias to a map whose keys are contract addresses
type LabeledAddresses map[string]cldf.TypeAndVersion

// And filters the LabeledAddresses to only include those entries that contain all
// labels.  If labels is empty, then only unlabeled entries are returned.
func (la LabeledAddresses) And(labels ...string) LabeledAddresses {
	var (
		filtered        = make(LabeledAddresses)
		selectUnlabeled = len(labels) == 0
		filterByLabels  = len(labels) > 0
	)

	for addr, tv := range la {
		// ignore labeled contracts by default
		if selectUnlabeled && !tv.Labels.IsEmpty() {
			continue
		}

		// ingore unlabeled contracts if labels are received
		if filterByLabels && tv.Labels.IsEmpty() {
			continue
		}

		if selectUnlabeled && tv.Labels.IsEmpty() {
			filtered[addr] = tv
			continue
		}

		if filterByLabels {
			keep := true
			for _, label := range labels {
				keep = keep && tv.Labels.Contains(label)
			}
			if keep {
				filtered[addr] = tv
			}
		}
	}

	return filtered
}
