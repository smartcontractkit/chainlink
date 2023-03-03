package sync

import (
	"sync"

	ds "github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
)

// MutexDatastore contains a child datastore and a mutex.
// used for coarse sync
type MutexDatastore struct {
	sync.RWMutex

	child ds.Datastore
}

// MutexWrap constructs a datastore with a coarse lock around the entire
// datastore, for every single operation.
func MutexWrap(d ds.Datastore) *MutexDatastore {
	return &MutexDatastore{child: d}
}

// Children implements Shim
func (d *MutexDatastore) Children() []ds.Datastore {
	return []ds.Datastore{d.child}
}

// Put implements Datastore.Put
func (d *MutexDatastore) Put(key ds.Key, value []byte) (err error) {
	d.Lock()
	defer d.Unlock()
	return d.child.Put(key, value)
}

// Sync implements Datastore.Sync
func (d *MutexDatastore) Sync(prefix ds.Key) error {
	d.Lock()
	defer d.Unlock()
	return d.child.Sync(prefix)
}

// Get implements Datastore.Get
func (d *MutexDatastore) Get(key ds.Key) (value []byte, err error) {
	d.RLock()
	defer d.RUnlock()
	return d.child.Get(key)
}

// Has implements Datastore.Has
func (d *MutexDatastore) Has(key ds.Key) (exists bool, err error) {
	d.RLock()
	defer d.RUnlock()
	return d.child.Has(key)
}

// GetSize implements Datastore.GetSize
func (d *MutexDatastore) GetSize(key ds.Key) (size int, err error) {
	d.RLock()
	defer d.RUnlock()
	return d.child.GetSize(key)
}

// Delete implements Datastore.Delete
func (d *MutexDatastore) Delete(key ds.Key) (err error) {
	d.Lock()
	defer d.Unlock()
	return d.child.Delete(key)
}

// Query implements Datastore.Query
func (d *MutexDatastore) Query(q dsq.Query) (dsq.Results, error) {
	d.RLock()
	defer d.RUnlock()

	// Apply the entire query while locked. Non-sync datastores may not
	// allow concurrent queries.

	results, err := d.child.Query(q)
	if err != nil {
		return nil, err
	}

	entries, err1 := results.Rest()
	err2 := results.Close()
	switch {
	case err1 != nil:
		return nil, err1
	case err2 != nil:
		return nil, err2
	}
	return dsq.ResultsWithEntries(q, entries), nil
}

func (d *MutexDatastore) Batch() (ds.Batch, error) {
	d.RLock()
	defer d.RUnlock()
	bds, ok := d.child.(ds.Batching)
	if !ok {
		return nil, ds.ErrBatchUnsupported
	}

	b, err := bds.Batch()
	if err != nil {
		return nil, err
	}
	return &syncBatch{
		batch: b,
		mds:   d,
	}, nil
}

func (d *MutexDatastore) Close() error {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()
	return d.child.Close()
}

// DiskUsage implements the PersistentDatastore interface.
func (d *MutexDatastore) DiskUsage() (uint64, error) {
	d.RLock()
	defer d.RUnlock()
	return ds.DiskUsage(d.child)
}

type syncBatch struct {
	batch ds.Batch
	mds   *MutexDatastore
}

func (b *syncBatch) Put(key ds.Key, val []byte) error {
	b.mds.Lock()
	defer b.mds.Unlock()
	return b.batch.Put(key, val)
}

func (b *syncBatch) Delete(key ds.Key) error {
	b.mds.Lock()
	defer b.mds.Unlock()
	return b.batch.Delete(key)
}

func (b *syncBatch) Commit() error {
	b.mds.Lock()
	defer b.mds.Unlock()
	return b.batch.Commit()
}

func (d *MutexDatastore) Check() error {
	if c, ok := d.child.(ds.CheckedDatastore); ok {
		d.RWMutex.Lock()
		defer d.RWMutex.Unlock()
		return c.Check()
	}
	return nil
}

func (d *MutexDatastore) Scrub() error {
	if c, ok := d.child.(ds.ScrubbedDatastore); ok {
		d.RWMutex.Lock()
		defer d.RWMutex.Unlock()
		return c.Scrub()
	}
	return nil
}

func (d *MutexDatastore) CollectGarbage() error {
	if c, ok := d.child.(ds.GCDatastore); ok {
		d.RWMutex.Lock()
		defer d.RWMutex.Unlock()
		return c.CollectGarbage()
	}
	return nil
}
