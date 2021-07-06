package feeds

import "github.com/smartcontractkit/chainlink/core/services/feeds/proto"

// SetFMSClient allows us to manually set the FMS client. Only used for testing.
//
// This allows us to avoid having to connect to a real server.
func (s *service) SetFMSClient(c proto.FeedsManagerClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.fmsClient = c
}
