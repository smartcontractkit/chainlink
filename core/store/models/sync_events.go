package models

import "time"

// SyncEvent represents an event sourcing style event, which is used to sync
// data upstream with another service
type SyncEvent struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
}
