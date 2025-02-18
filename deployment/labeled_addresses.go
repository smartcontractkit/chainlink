package deployment

type LabeledAddresses map[string]TypeAndVersion

func (la LabeledAddresses) And(labels ...string) LabeledAddresses {
	var (
		filtered        = make(LabeledAddresses, 0)
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
