// Package mocks re-exports the DON-to-DON dispatcher mocks now living in
// capabilities/libs/x/don2don/types/mocks, so callers can keep importing this path unchanged.
package mocks

import (
	don2donmocks "github.com/smartcontractkit/capabilities/libs/x/don2don/types/mocks"
)

type (
	Dispatcher                              = don2donmocks.Dispatcher
	Dispatcher_Expecter                     = don2donmocks.Dispatcher_Expecter
	Dispatcher_Close_Call                   = don2donmocks.Dispatcher_Close_Call
	Dispatcher_HealthReport_Call            = don2donmocks.Dispatcher_HealthReport_Call
	Dispatcher_Name_Call                    = don2donmocks.Dispatcher_Name_Call
	Dispatcher_Ready_Call                   = don2donmocks.Dispatcher_Ready_Call
	Dispatcher_RemoveReceiver_Call          = don2donmocks.Dispatcher_RemoveReceiver_Call
	Dispatcher_RemoveReceiverForMethod_Call = don2donmocks.Dispatcher_RemoveReceiverForMethod_Call
	Dispatcher_Send_Call                    = don2donmocks.Dispatcher_Send_Call
	Dispatcher_SetReceiver_Call             = don2donmocks.Dispatcher_SetReceiver_Call
	Dispatcher_SetReceiverForMethod_Call    = don2donmocks.Dispatcher_SetReceiverForMethod_Call
	Dispatcher_Start_Call                   = don2donmocks.Dispatcher_Start_Call

	Receiver              = don2donmocks.Receiver
	Receiver_Expecter     = don2donmocks.Receiver_Expecter
	Receiver_Receive_Call = don2donmocks.Receiver_Receive_Call
)

var (
	NewDispatcher = don2donmocks.NewDispatcher
	NewReceiver   = don2donmocks.NewReceiver
)
