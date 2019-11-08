package store

import (
	"chainlink/core/store/models"
)

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being attached to HeadTracker.
type HeadTrackable interface {
	Connect(*models.Head) error
	Disconnect()
	OnNewHead(*models.Head)
}
