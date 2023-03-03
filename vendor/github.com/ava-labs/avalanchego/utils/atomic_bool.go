// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import "sync/atomic"

type AtomicBool struct {
	value uint32
}

func (a *AtomicBool) GetValue() bool {
	return atomic.LoadUint32(&a.value) != 0
}

func (a *AtomicBool) SetValue(b bool) {
	var value uint32
	if b {
		value = 1
	}
	atomic.StoreUint32(&a.value, value)
}
