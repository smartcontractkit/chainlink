// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ips

import (
	"encoding/json"
	"net"
	"sync"
)

var _ DynamicIPPort = &dynamicIPPort{}

// An IPPort that can change.
// Safe for use by multiple goroutines.
type DynamicIPPort interface {
	// Returns the IP + port pair.
	IPPort() IPPort
	// Changes the IP.
	SetIP(ip net.IP)
}

type dynamicIPPort struct {
	lock   sync.RWMutex
	ipPort IPPort
}

func NewDynamicIPPort(ip net.IP, port uint16) DynamicIPPort {
	return &dynamicIPPort{
		ipPort: IPPort{
			IP:   ip,
			Port: port,
		},
	}
}

func (i *dynamicIPPort) IPPort() IPPort {
	i.lock.RLock()
	defer i.lock.RUnlock()

	return i.ipPort
}

func (i *dynamicIPPort) SetIP(ip net.IP) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.ipPort.IP = ip
}

func (i *dynamicIPPort) MarshalJSON() ([]byte, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	return json.Marshal(i.ipPort)
}
