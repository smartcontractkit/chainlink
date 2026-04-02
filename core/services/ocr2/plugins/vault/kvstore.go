package vault

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

const (
	keyPrefix              = "Key::"
	metadataPrefix         = "Metadata::"
	pendingQueueIndex      = "PendingQueue::Index"
	pendingQueueItemPrefix = "PendingQueue::Item::"
)

type KVStore struct {
	reader  ocr3_1types.KeyValueStateReader
	writer  ocr3_1types.KeyValueStateReadWriter
	metrics *pluginMetrics
}

func (s *KVStore) trackDuration(ctx context.Context, method string, start time.Time) {
	s.metrics.trackKVOperation(ctx, method, time.Since(start).Milliseconds())
}

type ReadKVStore interface {
	GetSecret(ctx context.Context, id *vault.SecretIdentifier) (*vault.StoredSecret, error)
	GetMetadata(ctx context.Context, owner string) (*vault.StoredMetadata, error)
	GetSecretIdentifiersCountForOwner(ctx context.Context, owner string) (int, error)
	GetPendingQueue(ctx context.Context) ([]*vault.StoredPendingQueueItem, error)
}

type WriteKVStore interface {
	ReadKVStore
	WriteSecret(ctx context.Context, id *vault.SecretIdentifier, secret *vault.StoredSecret) error
	WriteMetadata(ctx context.Context, owner string, metadata *vault.StoredMetadata) error
	DeleteSecret(ctx context.Context, id *vault.SecretIdentifier) error
	WritePendingQueue(ctx context.Context, pending []*vault.StoredPendingQueueItem) error
}

func NewReadStore(reader ocr3_1types.KeyValueStateReader, metrics *pluginMetrics) *KVStore {
	return &KVStore{reader: reader, metrics: metrics}
}

func NewWriteStore(writer ocr3_1types.KeyValueStateReadWriter, metrics *pluginMetrics) *KVStore {
	return &KVStore{reader: writer, writer: writer, metrics: metrics}
}

func (s *KVStore) GetSecret(ctx context.Context, id *vault.SecretIdentifier) (*vault.StoredSecret, error) {
	defer s.trackDuration(ctx, "GetSecret", time.Now())
	if id == nil {
		return nil, errors.New("id cannot be nil")
	}
	found, err := s.metadataContainsID(ctx, id)
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

func (s *KVStore) GetMetadata(ctx context.Context, owner string) (*vault.StoredMetadata, error) {
	defer s.trackDuration(ctx, "GetMetadata", time.Now())
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

func (s *KVStore) GetSecretIdentifiersCountForOwner(ctx context.Context, owner string) (int, error) {
	defer s.trackDuration(ctx, "GetSecretIdentifiersCountForOwner", time.Now())
	md, err := s.GetMetadata(ctx, owner)
	if err != nil {
		return 0, fmt.Errorf("failed to get metadata for owner %s: %w", owner, err)
	}

	count := 0
	if md != nil {
		count = len(md.SecretIdentifiers)
	}
	return count, nil
}

func (s *KVStore) WriteMetadata(ctx context.Context, owner string, metadata *vault.StoredMetadata) error {
	defer s.trackDuration(ctx, "WriteMetadata", time.Now())
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

func (s *KVStore) metadataContainsID(ctx context.Context, id *vault.SecretIdentifier) (bool, error) {
	if id == nil {
		return false, errors.New("id cannot be nil")
	}
	md, err := s.GetMetadata(ctx, id.Owner)
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

func (s *KVStore) addIDToMetadata(ctx context.Context, id *vault.SecretIdentifier) error {
	if id == nil {
		return errors.New("id cannot be nil")
	}
	md, err := s.GetMetadata(ctx, id.Owner)
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

	err = s.WriteMetadata(ctx, id.Owner, md)
	if err != nil {
		return fmt.Errorf("failed to write metadata for owner %s: %w", id.Owner, err)
	}

	return nil
}

func (s *KVStore) removeIDFromMetadata(ctx context.Context, id *vault.SecretIdentifier) error {
	if id == nil {
		return errors.New("id cannot be nil")
	}
	md, err := s.GetMetadata(ctx, id.Owner)
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
	err = s.WriteMetadata(ctx, id.Owner, newMd)
	if err != nil {
		return fmt.Errorf("failed to write metadata for owner %s: %w", id.Owner, err)
	}

	return nil
}

func (s *KVStore) WriteSecret(ctx context.Context, id *vault.SecretIdentifier, secret *vault.StoredSecret) error {
	defer s.trackDuration(ctx, "WriteSecret", time.Now())
	if id == nil {
		return errors.New("id cannot be nil")
	}
	b, err := proto.Marshal(secret)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}

	err = s.writer.Write([]byte(keyPrefix+vaulttypes.KeyFor(id)), b)
	if err != nil {
		return fmt.Errorf("failed to write secret: %w", err)
	}

	if err := s.addIDToMetadata(ctx, id); err != nil {
		return fmt.Errorf("failed to add id to metadata: %w", err)
	}

	return nil
}

func (s *KVStore) DeleteSecret(ctx context.Context, id *vault.SecretIdentifier) error {
	defer s.trackDuration(ctx, "DeleteSecret", time.Now())
	if id == nil {
		return errors.New("id cannot be nil")
	}
	err := s.removeIDFromMetadata(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to remove id from metadata: %w", err)
	}

	err = s.writer.Delete([]byte(keyPrefix + vaulttypes.KeyFor(id)))
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

func (s *KVStore) GetPendingQueue(ctx context.Context) ([]*vault.StoredPendingQueueItem, error) {
	defer s.trackDuration(ctx, "GetPendingQueue", time.Now())
	indexBytes, err := s.reader.Read([]byte(pendingQueueIndex))
	if err != nil {
		return nil, fmt.Errorf("failed to read pending queue index: %w", err)
	}

	if indexBytes == nil {
		return nil, nil
	}

	index := &vault.StoredPendingQueueIndex{}
	err = proto.Unmarshal(indexBytes, index)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal pending queue index: %w", err)
	}

	items := make([]*vault.StoredPendingQueueItem, 0, index.Length)
	for i := range index.Length {
		itemBytes, err := s.reader.Read([]byte(pendingQueueItemPrefix + strconv.Itoa(int(i))))
		if err != nil {
			return nil, fmt.Errorf("failed to read pending queue item at index %d: %w", i, err)
		}

		if itemBytes == nil {
			return nil, fmt.Errorf("pending queue item at index %d not found", i)
		}

		item := &vault.StoredPendingQueueItem{}
		err = proto.Unmarshal(itemBytes, item)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal pending queue item at index %d: %w", i, err)
		}

		if item.Item == nil {
			return nil, fmt.Errorf("pending queue item at index %d has nil Item", i)
		}

		items = append(items, item)
	}

	return items, nil
}

func (s *KVStore) deletePendingQueue() error {
	indexBytes, err := s.reader.Read([]byte(pendingQueueIndex))
	if err != nil {
		return fmt.Errorf("failed to read existing pending queue index: %w", err)
	}

	if indexBytes != nil {
		index := &vault.StoredPendingQueueIndex{}
		if err = proto.Unmarshal(indexBytes, index); err != nil {
			return fmt.Errorf("failed to unmarshal existing pending queue index: %w", err)
		}

		for i := 0; i < int(index.Length); i++ {
			if err := s.writer.Delete([]byte(pendingQueueItemPrefix + strconv.Itoa(i))); err != nil {
				return fmt.Errorf("failed to delete pending queue item at index %d: %w", i, err)
			}
		}
	}

	return nil
}

func (s *KVStore) WritePendingQueue(ctx context.Context, pending []*vault.StoredPendingQueueItem) error {
	defer s.trackDuration(ctx, "WritePendingQueue", time.Now())
	err := s.deletePendingQueue()
	if err != nil {
		return fmt.Errorf("failed to delete pending requests: %w", err)
	}

	for i, item := range pending {
		itemBytes, ierr := proto.Marshal(item)
		if ierr != nil {
			return fmt.Errorf("failed to marshal pending queue item at index %d: %w", i, ierr)
		}

		if ierr = s.writer.Write([]byte(pendingQueueItemPrefix+strconv.Itoa(i)), itemBytes); ierr != nil {
			return fmt.Errorf("failed to write pending queue item at index %d: %w", i, ierr)
		}
	}

	newIndex := &vault.StoredPendingQueueIndex{
		Length: int64(len(pending)),
	}
	newIndexBytes, err := proto.Marshal(newIndex)
	if err != nil {
		return fmt.Errorf("failed to marshal new pending queue index: %w", err)
	}

	if err := s.writer.Write([]byte(pendingQueueIndex), newIndexBytes); err != nil {
		return fmt.Errorf("failed to write new pending queue index: %w", err)
	}

	return nil
}
