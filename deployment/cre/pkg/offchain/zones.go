package offchain

import jdtypesv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

const ZoneLabel = "zone"

// Zone represents a logical grouping of dons/nodes within an environment.
type Zone string

const (
	DefaultZone Zone = "" // backward compatibility, no zone label
	ZoneB       Zone = "zone-b"
)

func (z Zone) String() string {
	return string(z)
}

func (z Zone) Label() *jdtypesv1.Label {
	if z == DefaultZone {
		return nil // backward compatibility that omits zone label for default zone
	}
	s := z.String()

	return &jdtypesv1.Label{
		Key:   ZoneLabel,
		Value: &s,
	}
}

func ZoneSelector(z Zone) *jdtypesv1.Selector {
	l := z.Label()
	if l == nil {
		return nil
	}

	return &jdtypesv1.Selector{
		Key:   l.Key,
		Op:    jdtypesv1.SelectorOp_EQ,
		Value: l.Value,
	}
}
