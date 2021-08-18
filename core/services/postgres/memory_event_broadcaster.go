package postgres

import (
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
	"gorm.io/gorm"
)

var _ EventBroadcaster = &memoryEventBroadcaster{}

type memoryEventBroadcaster struct {
	mu            sync.RWMutex
	subscriptions map[string][]*memoryEventBroadcasterSub
}

func NewMemoryEventBroadcaster() EventBroadcaster {
	return &memoryEventBroadcaster{sync.RWMutex{}, make(map[string][]*memoryEventBroadcasterSub)}
}
func (m *memoryEventBroadcaster) Subscribe(channel, payloadFilter string) (Subscription, error) {
	fmt.Println("BALLS Subscribe!", channel, payloadFilter)
	sub := &memoryEventBroadcasterSub{channel, payloadFilter, make(chan Event, 1000)}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscriptions[channel] = append(m.subscriptions[channel], sub)
	return sub, nil
}
func (m *memoryEventBroadcaster) Notify(channel string, payload string) error {
	fmt.Println("BALLS Notify!", channel, payload)
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, sub := range m.subscriptions[channel] {
		fmt.Println("BALLS", sub)
		ev := Event{Channel: channel, Payload: payload}
		if sub.InterestedIn(ev) {
			sub.Send(ev)
		}
	}
	return nil
}
func (m *memoryEventBroadcaster) NotifyInsideGormTx(_ *gorm.DB, channel string, payload string) error {
	m.Notify(channel, payload)
	return nil
}
func (m *memoryEventBroadcaster) Start() error   { return nil }
func (m *memoryEventBroadcaster) Close() error   { return nil }
func (m *memoryEventBroadcaster) Ready() error   { return nil }
func (m *memoryEventBroadcaster) Healthy() error { return nil }

type memoryEventBroadcasterSub struct {
	channel       string
	payloadFilter string
	ch            chan (Event)
}

var _ Subscription = &memoryEventBroadcasterSub{}

func (m *memoryEventBroadcasterSub) Events() <-chan Event {
	return m.ch
}
func (m *memoryEventBroadcasterSub) Close() {
	close(m.ch)
}
func (m *memoryEventBroadcasterSub) ChannelName() string {
	return m.channel
}
func (m *memoryEventBroadcasterSub) InterestedIn(event Event) bool {
	return m.payloadFilter == event.Payload || m.payloadFilter == ""
}
func (m *memoryEventBroadcasterSub) Send(event Event) {
	fmt.Println("BALLS Send!", event)
	select {
	case m.ch <- event:
	default:
		logger.Warn("memoryEventBroadcasterSub timed out sending event")
	}
}
