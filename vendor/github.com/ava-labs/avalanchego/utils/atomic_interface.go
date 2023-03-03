// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"sync"
)

type AtomicInterface struct {
	value interface{}
	lock  sync.RWMutex
}

func NewAtomicInterface(v interface{}) *AtomicInterface {
	mutexInterface := AtomicInterface{}
	mutexInterface.SetValue(v)
	return &mutexInterface
}

func (a *AtomicInterface) GetValue() interface{} {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.value
}

func (a *AtomicInterface) SetValue(v interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.value = v
}
