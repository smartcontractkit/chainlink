package ring

import (
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

// TransitionStateFromBool converts a proto bool (in_transition) to TransitionState
func TransitionStateFromBool(inTransition bool) shardorchestrator.TransitionState {
	if inTransition {
		return shardorchestrator.StateTransitioning
	}
	return shardorchestrator.StateSteady
}

// TransitionStateFromRoutingState returns the TransitionState based on RoutingState
func TransitionStateFromRoutingState(state *ringpb.RoutingState) shardorchestrator.TransitionState {
	if IsInSteadyState(state) {
		return shardorchestrator.StateSteady
	}
	return shardorchestrator.StateTransitioning
}

func IsInSteadyState(state *ringpb.RoutingState) bool {
	if state == nil {
		return false
	}
	_, ok := state.State.(*ringpb.RoutingState_RoutableShards)
	return ok
}

func NextStateFromSteady(currentID uint64, currentShards, wantShards uint32, now time.Time, timeToSync time.Duration) *ringpb.RoutingState {
	if currentShards == wantShards {
		return &ringpb.RoutingState{
			Id:    currentID,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: currentShards},
		}
	}

	return &ringpb.RoutingState{
		Id: currentID + 1,
		State: &ringpb.RoutingState_Transition{
			Transition: &ringpb.Transition{
				WantShards:       wantShards,
				LastStableCount:  currentShards,
				ChangesSafeAfter: timestamppb.New(now.Add(timeToSync)),
			},
		},
	}
}

func NextStateFromTransition(currentID uint64, transition *ringpb.Transition, now time.Time) *ringpb.RoutingState {
	safeAfter := transition.ChangesSafeAfter.AsTime()

	if now.Before(safeAfter) {
		return &ringpb.RoutingState{
			Id: currentID,
			State: &ringpb.RoutingState_Transition{
				Transition: transition,
			},
		}
	}

	return &ringpb.RoutingState{
		Id: currentID + 1,
		State: &ringpb.RoutingState_RoutableShards{
			RoutableShards: transition.WantShards,
		},
	}
}

func NextState(current *ringpb.RoutingState, wantShards uint32, now time.Time, timeToSync time.Duration) (*ringpb.RoutingState, error) {
	if current == nil {
		return nil, errors.New("current state is nil")
	}

	switch s := current.State.(type) {
	case *ringpb.RoutingState_RoutableShards:
		return NextStateFromSteady(current.Id, s.RoutableShards, wantShards, now, timeToSync), nil

	case *ringpb.RoutingState_Transition:
		return NextStateFromTransition(current.Id, s.Transition, now), nil

	// coverage:ignore
	default:
		return nil, errors.New("unknown state type")
	}
}
