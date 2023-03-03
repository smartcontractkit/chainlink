// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"fmt"
	"sync"
)

// AliaserReader allows one to lookup the aliases given to an ID.
type AliaserReader interface {
	Lookup(alias string) (ID, error)
	PrimaryAlias(id ID) (string, error)
	Aliases(id ID) ([]string, error)
}

// Aliaser allows one to give an ID aliases. An ID can have arbitrarily many
// aliases; two IDs may not have the same alias.
type AliaserWriter interface {
	Alias(id ID, alias string) error
	RemoveAliases(id ID)
}

// Aliaser allows one to give an ID aliases and lookup the aliases given to an
// ID.
type Aliaser interface {
	AliaserReader
	AliaserWriter
	PrimaryAliasOrDefault(id ID) string
}

type aliaser struct {
	lock    sync.RWMutex
	dealias map[string]ID
	aliases map[ID][]string
}

func NewAliaser() Aliaser {
	return &aliaser{
		dealias: make(map[string]ID),
		aliases: make(map[ID][]string),
	}
}

// Lookup returns the ID associated with alias
func (a *aliaser) Lookup(alias string) (ID, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if ID, ok := a.dealias[alias]; ok {
		return ID, nil
	}
	return ID{}, fmt.Errorf("there is no ID with alias %s", alias)
}

// PrimaryAlias returns the first alias of [id]
func (a *aliaser) PrimaryAlias(id ID) (string, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	aliases := a.aliases[id]
	if len(aliases) == 0 {
		return "", fmt.Errorf("there is no alias for ID %s", id)
	}
	return aliases[0], nil
}

// PrimaryAliasOrDefault returns the first alias of [id], or ID string as default
func (a *aliaser) PrimaryAliasOrDefault(id ID) string {
	alias, err := a.PrimaryAlias(id)
	if err != nil {
		return id.String()
	}
	return alias
}

// Aliases returns the aliases of an ID
func (a *aliaser) Aliases(id ID) ([]string, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.aliases[id], nil
}

// Alias gives [id] the alias [alias]
func (a *aliaser) Alias(id ID, alias string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, exists := a.dealias[alias]; exists {
		return fmt.Errorf("%s is already used as an alias for an ID", alias)
	}

	a.dealias[alias] = id
	a.aliases[id] = append(a.aliases[id], alias)
	return nil
}

// RemoveAliases of the provided ID
func (a *aliaser) RemoveAliases(id ID) {
	a.lock.Lock()
	defer a.lock.Unlock()

	aliases := a.aliases[id]
	delete(a.aliases, id)
	for _, alias := range aliases {
		delete(a.dealias, alias)
	}
}

// GetRelevantAliases returns the aliases with the redundant identity alias
// removed (each id is aliased to at least itself).
func GetRelevantAliases(aliaser Aliaser, ids []ID) (map[ID][]string, error) {
	result := make(map[ID][]string, len(ids))
	for _, id := range ids {
		aliases, err := aliaser.Aliases(id)
		if err != nil {
			return nil, err
		}

		// remove the redundant alias where alias = id.
		relevantAliases := make([]string, 0, len(aliases)-1)
		for _, alias := range aliases {
			if alias != id.String() {
				relevantAliases = append(relevantAliases, alias)
			}
		}
		result[id] = relevantAliases
	}
	return result, nil
}
