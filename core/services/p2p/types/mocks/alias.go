// Package mocks re-exports the DON-to-DON peer mocks now living in
// capabilities/libs/x/rage/mocks, so callers can keep importing this path unchanged.
package mocks

import (
	ragemocks "github.com/smartcontractkit/capabilities/libs/x/rage/mocks"
)

type (
	Peer                        = ragemocks.Peer
	Peer_Expecter               = ragemocks.Peer_Expecter
	Peer_Close_Call             = ragemocks.Peer_Close_Call
	Peer_HealthReport_Call      = ragemocks.Peer_HealthReport_Call
	Peer_ID_Call                = ragemocks.Peer_ID_Call
	Peer_IsBootstrap_Call       = ragemocks.Peer_IsBootstrap_Call
	Peer_Name_Call              = ragemocks.Peer_Name_Call
	Peer_Ready_Call             = ragemocks.Peer_Ready_Call
	Peer_Receive_Call           = ragemocks.Peer_Receive_Call
	Peer_Send_Call              = ragemocks.Peer_Send_Call
	Peer_Start_Call             = ragemocks.Peer_Start_Call
	Peer_UpdateConnections_Call = ragemocks.Peer_UpdateConnections_Call

	PeerWrapper                   = ragemocks.PeerWrapper
	PeerWrapper_Expecter          = ragemocks.PeerWrapper_Expecter
	PeerWrapper_Close_Call        = ragemocks.PeerWrapper_Close_Call
	PeerWrapper_GetPeer_Call      = ragemocks.PeerWrapper_GetPeer_Call
	PeerWrapper_HealthReport_Call = ragemocks.PeerWrapper_HealthReport_Call
	PeerWrapper_Name_Call         = ragemocks.PeerWrapper_Name_Call
	PeerWrapper_Ready_Call        = ragemocks.PeerWrapper_Ready_Call
	PeerWrapper_Start_Call        = ragemocks.PeerWrapper_Start_Call

	SharedPeer                              = ragemocks.SharedPeer
	SharedPeer_Expecter                     = ragemocks.SharedPeer_Expecter
	SharedPeer_Close_Call                   = ragemocks.SharedPeer_Close_Call
	SharedPeer_HealthReport_Call            = ragemocks.SharedPeer_HealthReport_Call
	SharedPeer_ID_Call                      = ragemocks.SharedPeer_ID_Call
	SharedPeer_IsBootstrap_Call             = ragemocks.SharedPeer_IsBootstrap_Call
	SharedPeer_Name_Call                    = ragemocks.SharedPeer_Name_Call
	SharedPeer_Ready_Call                   = ragemocks.SharedPeer_Ready_Call
	SharedPeer_Receive_Call                 = ragemocks.SharedPeer_Receive_Call
	SharedPeer_Send_Call                    = ragemocks.SharedPeer_Send_Call
	SharedPeer_Start_Call                   = ragemocks.SharedPeer_Start_Call
	SharedPeer_UpdateConnections_Call       = ragemocks.SharedPeer_UpdateConnections_Call
	SharedPeer_UpdateConnectionsByDONs_Call = ragemocks.SharedPeer_UpdateConnectionsByDONs_Call

	Signer                 = ragemocks.Signer
	Signer_Expecter        = ragemocks.Signer_Expecter
	Signer_Initialize_Call = ragemocks.Signer_Initialize_Call
	Signer_Sign_Call       = ragemocks.Signer_Sign_Call
)

var (
	NewPeer        = ragemocks.NewPeer
	NewPeerWrapper = ragemocks.NewPeerWrapper
	NewSharedPeer  = ragemocks.NewSharedPeer
	NewSigner      = ragemocks.NewSigner
)
