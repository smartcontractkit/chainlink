package services

import "github.com/smartcontractkit/chainlink/core/store/models"

// headTrackableCallback is a simple wrapper around an On Connect callback
type headTrackableCallback struct {
	onConnect func()
}

func (c *headTrackableCallback) Connect(*models.Head) error {
	c.onConnect()
	return nil
}

func (c *headTrackableCallback) Disconnect()            {}
func (c *headTrackableCallback) OnNewHead(*models.Head) {}
