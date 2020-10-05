package dhtrouter

import (
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type PermitListACL interface {
	ACL

	Activate(protocol protocol.ID, permitted ...peer.ID)

	Deactivate(protocol protocol.ID)

	Permit(protocol protocol.ID, ids ...peer.ID) error

	Reject(protocol protocol.ID, id peer.ID) error
}

type permitList struct {
	allowed map[protocol.ID][]peer.ID
	logger  types.Logger
}

func NewPermitListACL(logger types.Logger) PermitListACL {
	return permitList{
		allowed: make(map[protocol.ID][]peer.ID),
		logger:  logger,
	}
}

func (acl permitList) Activate(protocol protocol.ID, permitted ...peer.ID) {
	_, found := acl.allowed[protocol]
	if found {
		return
	}

	acl.allowed[protocol] = make([]peer.ID, len(permitted))
	copy(acl.allowed[protocol], permitted)

	acl.logger.Debug("New ACL activated", types.LogFields{
		"id":         "DHT_ACL",
		"protocolID": protocol,
	})
}

func (acl permitList) Deactivate(protocol protocol.ID) {
	delete(acl.allowed, protocol)
}

func (acl permitList) Permit(protocol protocol.ID, ids ...peer.ID) error {
	list, found := acl.allowed[protocol]
	if !found {
		return errors.New("protocol not activated")
	}

	acl.allowed[protocol] = append(list, ids...)
	return nil
}

func (acl permitList) Reject(protocol protocol.ID, id peer.ID) error {
	panic("implement me")
}

func (acl permitList) IsAllowed(id peer.ID, protocol protocol.ID) bool {
	allowed, enforced := acl.allowed[protocol]
	if !enforced {
		return true
	}

	for _, p := range allowed {
		if p == id {
			return true
		}
	}
	acl.logger.Debug("ACL: denied access", types.LogFields{
		"id":         "DHT_ACL",
		"peerID":     id,
		"protocolID": protocol,
	})
	return false
}

func (acl permitList) IsACLEnforced(protocol protocol.ID) bool {
	_, found := acl.allowed[protocol]
	return found
}

func (acl permitList) String() string {
	s := ""
	for protocolId, aclMap := range acl.allowed {
		s += fmt.Sprintf("Protocol %s permits following nodes:\n", protocolId)
		for i, peerId := range aclMap {
			s += fmt.Sprintf("[%2d]\t%s\n", i, peerId.Pretty())
		}
	}

	return s
}
