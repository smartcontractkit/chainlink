package models

import "time"

// SyncEvent represents an event sourcing style event, which is used to sync
// data upstream with another service
type SyncEvent struct {
	ID        int       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt" gorm:"index"`
	Body      string
}
