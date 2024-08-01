// bagdb: Simple datastorage
// Copyright 2023 billy authors
// SPDX-License-Identifier: BSD-3-Clause

package billy

// Infos contains a set of statistics about the underlying datastore.
type Infos struct {
	Shelves []*ShelfInfos
}

// ShelfInfos contains some statistics about the data stored in a single shelf.
type ShelfInfos struct {
	SlotSize    uint32
	FilledSlots uint64
	GappedSlots uint64
}

// Infos gathers and returns some stats about the database.
func (db *database) Infos() *Infos {
	infos := new(Infos)
	for _, shelf := range db.shelves {
		slots, gaps := shelf.stats()

		infos.Shelves = append(infos.Shelves, &ShelfInfos{
			SlotSize:    shelf.slotSize,
			FilledSlots: slots - gaps,
			GappedSlots: gaps,
		})
	}
	return infos
}
