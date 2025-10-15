package vault

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

const (
	keyPrefix      = "Key::"
	metadataPrefix = "Metadata::"
)

type KVStore struct {
	reader ocr3_1types.KeyValueReader
	writer ocr3_1types.KeyValueReadWriter
}

type ReadKVStore interface {
	GetSecret(id *vault.SecretIdentifier) (*vault.StoredSecret, error)
	GetMetadata(owner string) (*vault.StoredMetadata, error)
	GetSecretIdentifiersCountForOwner(owner string) (int, error)
}

type WriteKVStore interface {
	ReadKVStore
	WriteSecret(id *vault.SecretIdentifier, secret *vault.StoredSecret) error
	WriteMetadata(owner string, metadata *vault.StoredMetadata) error
	DeleteSecret(id *vault.SecretIdentifier) error
}

func NewReadStore(reader ocr3_1types.KeyValueReader) *KVStore {
	return &KVStore{reader: reader}
}

func NewWriteStore(writer ocr3_1types.KeyValueReadWriter) *KVStore {
	return &KVStore{reader: writer, writer: writer}
}

func (s *KVStore) GetSecret(id *vault.SecretIdentifier) (*vault.StoredSecret, error) {
	found, err := s.metadataContainsID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to check if metadata contains id: %w", err)
	}

	if !found {
		return nil, nil
	}

	b, err := s.reader.Read([]byte(keyPrefix + vaulttypes.KeyFor(id)))
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}

	if b == nil {
		return nil, errors.New("invariant violation: metadata contains id but secret not found")
	}

	secret := &vault.StoredSecret{}
	err = proto.Unmarshal(b, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}
	return secret, nil
}

func (s *KVStore) GetMetadata(owner string) (*vault.StoredMetadata, error) {
	b, err := s.reader.Read([]byte(metadataPrefix + owner))
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	if b == nil {
		return nil, nil
	}

	md := &vault.StoredMetadata{}
	err = proto.Unmarshal(b, md)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal md: %w", err)
	}
	return md, nil
}

func (s *KVStore) GetSecretIdentifiersCountForOwner(owner string) (int, error) {
	md, err := s.GetMetadata(owner)
	if err != nil {
		return 0, fmt.Errorf("failed to get metadata for owner %s: %w", owner, err)
	}

	count := 0
	if md != nil {
		count = len(md.SecretIdentifiers)
	}
	return count, nil
}

func (s *KVStore) WriteMetadata(owner string, metadata *vault.StoredMetadata) error {
	if metadata == nil {
		return errors.New("metadata cannot be nil")
	}
	b, err := proto.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = s.writer.Write([]byte(metadataPrefix+owner), b)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

func (s *KVStore) metadataContainsID(id *vault.SecretIdentifier) (bool, error) {
	md, err := s.GetMetadata(id.Owner)
	if err != nil {
		return false, fmt.Errorf("failed to get metadata for owner %s: %w", id.Owner, err)
	}

	if md == nil {
		return false, nil
	}

	for _, i := range md.SecretIdentifiers {
		if vaulttypes.KeyFor(id) == vaulttypes.KeyFor(i) {
			return true, nil
		}
	}

	return false, nil
}

func (s *KVStore) addIDToMetadata(id *vault.SecretIdentifier) error {
	md, err := s.GetMetadata(id.Owner)
	if err != nil {
		return fmt.Errorf("failed to get metadata for owner %s: %w", id.Owner, err)
	}

	if md == nil {
		md = &vault.StoredMetadata{
			SecretIdentifiers: []*vault.SecretIdentifier{id},
		}
	} else {
		for _, i := range md.SecretIdentifiers {
			if vaulttypes.KeyFor(id) == vaulttypes.KeyFor(i) {
				// Nothing to do, early exit.
				return nil
			}
		}

		md.SecretIdentifiers = append(md.SecretIdentifiers, id)
	}

	err = s.WriteMetadata(id.Owner, md)
	if err != nil {
		return fmt.Errorf("failed to write metadata for owner %s: %w", id.Owner, err)
	}

	return nil
}

func (s *KVStore) removeIDFromMetadata(id *vault.SecretIdentifier) error {
	md, err := s.GetMetadata(id.Owner)
	if err != nil {
		return fmt.Errorf("failed to get metadata for owner %s: %w", id.Owner, err)
	}

	if md == nil {
		return fmt.Errorf("no metadata found for owner %s", id.Owner)
	}

	si := []*vault.SecretIdentifier{}
	var found bool
	for _, i := range md.SecretIdentifiers {
		if vaulttypes.KeyFor(id) == vaulttypes.KeyFor(i) {
			found = true
		} else {
			si = append(si, i)
		}
	}

	if !found {
		return fmt.Errorf("id %s not found in metadata for owner %s", vaulttypes.KeyFor(id), id.Owner)
	}

	newMd := &vault.StoredMetadata{
		SecretIdentifiers: si,
	}
	err = s.WriteMetadata(id.Owner, newMd)
	if err != nil {
		return fmt.Errorf("failed to write metadata for owner %s: %w", id.Owner, err)
	}

	return nil
}

func (s *KVStore) WriteSecret(id *vault.SecretIdentifier, secret *vault.StoredSecret) error {
	b, err := proto.Marshal(secret)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}

	err = s.writer.Write([]byte(keyPrefix+vaulttypes.KeyFor(id)), b)
	if err != nil {
		return fmt.Errorf("failed to write secret: %w", err)
	}

	if err := s.addIDToMetadata(id); err != nil {
		return fmt.Errorf("failed to add id to metadata: %w", err)
	}

	return nil
}

func (s *KVStore) DeleteSecret(id *vault.SecretIdentifier) error {
	err := s.removeIDFromMetadata(id)
	if err != nil {
		return fmt.Errorf("failed to remove id from metadata: %w", err)
	}

	err = s.writer.Delete([]byte(keyPrefix + vaulttypes.KeyFor(id)))
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}
