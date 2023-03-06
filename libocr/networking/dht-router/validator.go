package dhtrouter

// AnnouncementValidator verifies peer announcements
type AnnouncementValidator struct{}

const ValidatorNamespace = "peerinfo"

// Validate 1) extracts the pk from DHT key; 2) checks value is properly signed with pk
func (v AnnouncementValidator) Validate(key string, value []byte) error {
	// 1) extract public key from DHT key
	peerId, err := dhtKeyToPeerId(key)
	if err != nil {
		return err
	}

	ann, err := deserializeSignedAnnouncement(value)
	if err != nil {
		return err
	}

	// check key is consistent with the embedded PublicKey
	if !peerId.MatchesPublicKey(ann.PublicKey) {
		return InvalidDhtKey
	}

	// check the payload is properly signed
	return ann.verify()
}

// Select returns the latest (i.e., the one with the largest counter)
func (v AnnouncementValidator) Select(_ string, values [][]byte) (int, error) {
	counter := announcementCounter{}
	latestRecord := 0
	for i := 0; i < len(values); i++ {
		ann, err := deserializeSignedAnnouncement(values[i])
		if err != nil {
			return 0, err
		}

		if ann.Counter.Gt(counter) {
			latestRecord = i
			counter = ann.Counter
		}
	}
	return latestRecord, nil
}
