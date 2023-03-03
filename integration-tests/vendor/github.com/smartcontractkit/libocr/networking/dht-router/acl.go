package dhtrouter

import (
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
)

type PermitListACL interface {
	ACL

	// Activate ACL for a protocol with the provided allowlist.
	// No-op if ACL is already activated for that protocol
	Activate(protocol protocol.ID, permitted ...peer.ID)

	// Deactivate ACL for a protocol. No-op if not already activated
	Deactivate(protocol protocol.ID)

	// add ids to the allowlist for protocol. Error if not already activated
	Permit(protocol protocol.ID, ids ...peer.ID) error

	// remove id from the allowlist. Error if not already activated.
	Reject(protocol protocol.ID, id peer.ID) error
}

type permitList struct {
	allowed map[protocol.ID][]peer.ID // access control list. For protocols NOT in the table, the default decision is Permit.

	logger loghelper.LoggerWithContext
}

func NewPermitListACL(logger loghelper.LoggerWithContext) PermitListACL {
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

	acl.logger.Debug("New ACL activated", commontypes.LogFields{
		"id":         "DHT_ACL",
		"protocolID": protocol,
		"acl":        acl.allowed,
	})
}

func (acl permitList) Deactivate(protocol protocol.ID) {
	// delete is a no-op for non-existent
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
	// Found in the enforced map = ACL is enforced for this protocol
	allowed, enforced := acl.allowed[protocol]
	if !enforced {
		return true
	}

	// only Permit if id is in the white list
	for _, p := range allowed {
		if p == id {
			return true
		}
	}
	acl.logger.Debug("ACL: denied access", commontypes.LogFields{
		"id":         "DHT_ACL",
		"peerID":     id,
		"protocolID": protocol,
	})
	return false
}

func (acl permitList) IsACLEnforced(protocol protocol.ID) bool {
	_, found := acl.allowed[protocol]
	// Not found in the enforced map = ACL not enforced for this protocol
	return found
}

func (acl permitList) String() string {
	s := ""
	list := make(map[string][]string)
	for protocolId, aclMap := range acl.allowed {
		var permittedIds []string
		s += fmt.Sprintf("Protocol %s permits following nodes: ", protocolId)
		for _, peerId := range aclMap {
			permittedIds = append(permittedIds, peerId.Pretty())
		}
		list[string(protocolId)] = permittedIds
	}

	return fmt.Sprint(list)
}
