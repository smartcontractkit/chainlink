package vault

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"google.golang.org/protobuf/proto"
)

const (
	keyPrefix      = "Key::"
	metadataPrefix = "Metadata::"
)

type kvStore struct {
	reader ocr3_1types.KeyValueReader
	writer ocr3_1types.KeyValueReadWriter
}

type readKVStore interface {
	getSecret(id *vault.SecretIdentifier) (*vault.StoredSecret, error)
	getMetadata(owner string) (*vault.StoredMetadata, error)
}

type writeKVStore interface {
	readKVStore
	writeSecret(id *vault.SecretIdentifier, secret *vault.StoredSecret) error
	writeMetadata(owner string, metadata *vault.StoredMetadata) error
	addKeyToMetadata(id *vault.SecretIdentifier) error
}

func newReadStore(reader ocr3_1types.KeyValueReader) readKVStore {
	return &kvStore{reader: reader}
}

func newWriteStore(writer ocr3_1types.KeyValueReadWriter) writeKVStore {
	return &kvStore{reader: writer, writer: writer}
}

func (s *kvStore) getSecret(id *vault.SecretIdentifier) (*vault.StoredSecret, error) {
	b, err := s.reader.Read([]byte(keyPrefix + keyFor(id)))
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}

	if b == nil {
		return nil, nil
	}

	secret := &vault.StoredSecret{}
	err = proto.Unmarshal(b, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}
	return secret, nil
}

func (s *kvStore) getMetadata(owner string) (*vault.StoredMetadata, error) {
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

func (s *kvStore) writeMetadata(owner string, metadata *vault.StoredMetadata) error {
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

func (s *kvStore) addKeyToMetadata(id *vault.SecretIdentifier) error {
	md, err := s.getMetadata(id.Owner)
	if err != nil {
		return fmt.Errorf("failed to get metadata for owner %s: %w", id.Owner, err)
	}

	if md == nil {
		md = &vault.StoredMetadata{
			Keys: []string{keyPrefix + keyFor(id)},
		}
	} else {
		md.Keys = append(md.Keys, keyPrefix+keyFor(id))
	}

	err = s.writeMetadata(id.Owner, md)
	if err != nil {
		return fmt.Errorf("failed to write metadata for owner %s: %w", id.Owner, err)
	}

	return nil
}

func (s *kvStore) writeSecret(id *vault.SecretIdentifier, secret *vault.StoredSecret) error {
	b, err := proto.Marshal(secret)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}

	err = s.writer.Write([]byte(keyPrefix+keyFor(id)), b)
	if err != nil {
		return fmt.Errorf("failed to write secret: %w", err)
	}

	return nil
}
