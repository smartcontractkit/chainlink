package dhtrouter

type AnnouncementValidator struct{}

const ValidatorNamespace = "peerinfo"

func (v AnnouncementValidator) Validate(key string, value []byte) error {
	pid, err := dhtKeyToPeerId(key)
	if err != nil {
		return err
	}

	var ann Announcement
	err = ann.UnmarshalJSON(value)
	if err != nil {
		return err
	}

	if !pid.MatchesPublicKey(ann.Pk) {
		return InvalidDhtKey
	}

	ok, err := ann.SelfVerify()
	if err != nil {
		return err
	} else if !ok {
		return InvalidSignature
	} else {
		return nil
	}
}

func (v AnnouncementValidator) Select(key string, values [][]byte) (int, error) {
	strs := make([]string, len(values))
	latestTime := int64(0)
	latestRecord := 0
	for i := 0; i < len(values); i++ {
		ann := Announcement{}
		if err := ann.UnmarshalJSON(values[i]); err != nil {
			return 0, err
		}

		if ann.timestamp > latestTime {
			latestRecord = i
			latestTime = ann.timestamp
		}
		strs[i] = string(values[i])
	}
	return latestRecord, nil
}
