//nolint
//nolint:nolintlint
package dkg

type ParticipantSecretKey []byte
type ParticipantPublicKey []byte

type MasterPublicKey []byte
type MasterPublicKeyShare []byte
type MasterSecretKeyShare []byte

type InstanceID string

type ParticipantsConfig struct {
	F          int                    // Maximum number of faulty participants.
	T          int                    // Reconstruction threshold for secret sharing, the minimal number of shares needed to reconstruct the master secret.
	PublicKeys []ParticipantPublicKey // Public keys of the participants.
}

type Result interface {
	InstanceID() InstanceID // Unique identifier for the DKG instance
	// DealersConfig() ParticipantsConfig // [TODO] Is there any scenario to expose dealers config in the result?
	RecipientsConfig() ParticipantsConfig // Public keys and fault tolerance parameters of the recipients.
	MasterPublicKey() (MasterPublicKey, error)
	MasterPublicKeyShares() ([]MasterPublicKeyShare, error)
	MasterSecretKeyShare(ParticipantSecretKey) (MasterSecretKeyShare, error)
}
