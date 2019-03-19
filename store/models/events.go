package models

import "github.com/jinzhu/gorm"

// SyncEvent represents an event sourcing style event, which is used to sync
// data upstream with another service
type SyncEvent struct {
	gorm.Model
	Body string
}
