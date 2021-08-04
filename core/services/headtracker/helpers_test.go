package headtracker

import "sync"

func GetHeadListenerConnectedMutex(hl *HeadListener) *sync.RWMutex {
	return &hl.connectedMutex
}
