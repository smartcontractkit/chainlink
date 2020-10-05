package networking

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/peerstore"

	p2pnetwork "github.com/libp2p/go-libp2p-core/network"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2pprotocol "github.com/libp2p/go-libp2p-core/protocol"
	swarm "github.com/libp2p/go-libp2p-swarm"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/pkg/errors"
	dhtrouter "github.com/smartcontractkit/chainlink/libocr/networking/dht-router"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

var (
	_ types.BinaryNetworkEndpoint = &ocrEndpoint{}
)

type EndpointConfig struct {
	IncomingMessageBufferSize int

	OutgoingMessageBufferSize int

	NewStreamTimeout time.Duration

	DHTLookupInterval int
}

type ocrEndpointState int

type ocrEndpoint struct {
	config              EndpointConfig
	peerMapping         map[types.OracleID]p2ppeer.ID
	reversedPeerMapping map[p2ppeer.ID]types.OracleID
	peerAllowlist       map[p2ppeer.ID]struct{}
	peer                *concretePeer
	rhost               *rhost.RoutedHost
	routing             dhtrouter.PeerDiscoveryRouter
	configDigest        types.ConfigDigest
	protocolID          p2pprotocol.ID
	bootstrapperAddrs   []p2ppeer.AddrInfo
	failureThreshold    int
	ownOracleID         types.OracleID

	chRecvs      map[types.OracleID](chan []byte)
	chSends      map[types.OracleID](chan []byte)
	muSends      map[types.OracleID]*sync.Mutex
	chSendToSelf chan types.BinaryMessageWithSender
	chClose      chan struct{}
	state        ocrEndpointState
	stateMu      *sync.RWMutex
	wg           *sync.WaitGroup
	ctx          context.Context
	ctxCancel    context.CancelFunc

	recv chan types.BinaryMessageWithSender

	logger types.Logger
}

const (
	ocrEndpointUnstarted = iota
	ocrEndpointStarted
	ocrEndpointClosed

	sendToSelfBufferSize = 20

	protocolBaseName = "cl_offchainreporting"
	protocolVersion  = "1.0.0"
)

func newOCREndpoint(
	logger types.Logger,
	configDigest types.ConfigDigest,
	peer *concretePeer,
	peerIDs []p2ppeer.ID,
	bootstrappers []p2ppeer.AddrInfo,
	config EndpointConfig,
	failureThreshold int,
) (*ocrEndpoint, error) {
	peerMapping := make(map[types.OracleID]p2ppeer.ID)
	for i, peerID := range peerIDs {
		peerMapping[types.OracleID(i)] = peerID
	}
	reversedPeerMapping := reverseMapping(peerMapping)
	ownOracleID, ok := reversedPeerMapping[peer.ID()]
	if !ok {
		return nil, errors.Errorf("host peer ID 0x%x is not present in given peerMapping", peer.ID())
	}

	chRecvs := make(map[types.OracleID]chan []byte)
	chSends := make(map[types.OracleID]chan []byte)
	muSends := make(map[types.OracleID]*sync.Mutex)
	for oid := range peerMapping {
		if oid != ownOracleID {
			chRecvs[oid] = make(chan []byte, config.IncomingMessageBufferSize)
			chSends[oid] = make(chan []byte, config.OutgoingMessageBufferSize)
			muSends[oid] = new(sync.Mutex)
		}
	}

	chSendToSelf := make(chan types.BinaryMessageWithSender, sendToSelfBufferSize)

	protocolID := genProtocolID(configDigest)

	logger = loghelper.MakeLoggerWithContext(logger, types.LogFields{
		"protocolID":   protocolID,
		"configDigest": configDigest.Hex(),
		"oracleID":     ownOracleID,
		"id":           "OCREndpoint",
	})

	ctx, cancel := context.WithCancel(context.Background())

	allowlist := make(map[p2ppeer.ID]struct{})
	for pid := range reversedPeerMapping {
		allowlist[pid] = struct{}{}
	}
	for _, b := range bootstrappers {
		allowlist[b.ID] = struct{}{}
	}

	return &ocrEndpoint{
		config,
		peerMapping,
		reversedPeerMapping,
		allowlist,
		peer,
		nil,
		nil,
		configDigest,
		protocolID,
		bootstrappers,
		failureThreshold,
		ownOracleID,
		chRecvs,
		chSends,
		muSends,
		chSendToSelf,
		make(chan struct{}),
		ocrEndpointUnstarted,
		new(sync.RWMutex),
		new(sync.WaitGroup),
		ctx,
		cancel,
		make(chan types.BinaryMessageWithSender),
		logger,
	}, nil
}

func reverseMapping(m map[types.OracleID]p2ppeer.ID) map[p2ppeer.ID]types.OracleID {
	n := make(map[p2ppeer.ID]types.OracleID)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func genProtocolID(configDigest types.ConfigDigest) p2pprotocol.ID {
	return p2pprotocol.ID(fmt.Sprintf("/%s/%s/%x/1.0.0", protocolBaseName, protocolVersion, configDigest))
}

func (o *ocrEndpoint) Start() (err error) {
	o.stateMu.Lock()
	defer o.stateMu.Unlock()

	if o.state != ocrEndpointUnstarted {
		panic("ocrEndpoint has already been started")
	}
	o.state = ocrEndpointStarted

	if err := o.peer.register(o); err != nil {
		return err
	}

	if err := o.setupDHT(); err != nil {
		return errors.Wrap(err, "error setting up DHT")
	}

	o.rhost.SetStreamHandler(o.protocolID, o.streamReceiver)

	o.wg.Add(len(o.chRecvs))
	for oid := range o.chRecvs {
		go o.runRecv(oid)
	}
	o.wg.Add(len(o.chSends))
	for oid := range o.chSends {
		go o.runSend(oid)
	}
	o.wg.Add(1)
	go o.runSendToSelf()

	o.logger.Info("OCREndpoint: Started listening", nil)

	return nil
}

func (o *ocrEndpoint) setupDHT() (err error) {
	config := dhtrouter.BuildConfig(
		o.bootstrapperAddrs,
		dhtPrefix,
		o.configDigest,
		o.logger,
		o.failureThreshold,
		false,
	)

	acl := dhtrouter.NewPermitListACL(o.logger)

	acl.Activate(config.ProtocolID(), o.allowlist()...)
	aclHost := dhtrouter.WrapACL(o.peer, acl, o.logger)

	o.routing, err = dhtrouter.NewDHTRouter(
		o.ctx,
		config,
		aclHost,
	)
	if err != nil {
		return errors.Wrap(err, "could not initialize DHTRouter")
	}

	o.routing.Start()

	o.rhost = rhost.Wrap(o.peer, o.routing)

	return nil
}

func (o *ocrEndpoint) runRecv(oid types.OracleID) {
	defer o.wg.Done()
	var chRecv <-chan []byte = o.chRecvs[oid]
	for {
		select {
		case payload := <-chRecv:
			msg := types.BinaryMessageWithSender{
				Msg:    payload,
				Sender: oid,
			}

			select {
			case o.recv <- msg:
				continue
			case <-o.chClose:
				return
			}
		case <-o.chClose:
			return
		}
	}
}

func (o *ocrEndpoint) runSend(oid types.OracleID) {
	defer o.wg.Done()

	var chSend <-chan []byte = o.chSends[oid]
	destPeerID, err := o.oracleID2PeerID(oid)
	if err != nil {
		panic("error getting destination peer ID")
	}

	for {
		shouldRetry := o.sendOnStream(destPeerID, chSend)
		if !shouldRetry {
			return
		}
	}
}

func (o *ocrEndpoint) sendOnStream(destPeerID p2ppeer.ID, chSend <-chan []byte) (shouldRetry bool) {
	var stream p2pnetwork.Stream

	nRetry := 0

	for {
		var err error
		stream, err = func() (p2pnetwork.Stream, error) {
			var ctx context.Context
			if o.config.NewStreamTimeout == 0 {
				ctx = o.ctx
			} else {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(o.ctx, o.config.NewStreamTimeout)
				defer cancel()
			}
			return o.peer.NewStream(ctx, destPeerID, o.protocolID)
		}()

		if err == nil {
			break
		}

		if errors.Is(err, context.Canceled) {
			select {
			case <-o.chClose:
				return false
			default:
			}
		}

		if !errors.Is(err, swarm.ErrDialBackoff) {
			o.logger.Debug("Peer unreachable", types.LogFields{
				"err":            err,
				"remotePeerID":   destPeerID,
				"nRetry":         nRetry,
				"remoteOracleId": o.reversedPeerMapping[destPeerID],
			})
		}

		if nRetry > 0 && nRetry%o.config.DHTLookupInterval == 0 {
			pAddr, err := o.routing.FindPeer(o.ctx, destPeerID)
			switch {
			case err == nil:
				o.logger.Debug("DHT lookup finished", types.LogFields{
					"result": pAddr,
					"nRetry": nRetry,
				})
				o.peer.Peerstore().AddAddrs(destPeerID, pAddr.Addrs, peerstore.TempAddrTTL)
			case errors.Is(err, context.Canceled):
				return false
			default:
				o.logger.Error("DHT lookup failed", types.LogFields{
					"err":            err,
					"remoteOracleId": o.reversedPeerMapping[destPeerID],
					"nRetry":         nRetry,
					"remotePeerID":   destPeerID,
				})
			}
		}

		nRetry++

		waitms := time.Duration(int64((4+rand.Float64()*2)*1000)) * time.Millisecond
		waitCh := time.After(waitms)

		select {
		case <-waitCh:
		case <-o.chClose:
			return false
		}
	}

	defer stream.Reset()

	o.logger.Debug("Opened stream", types.LogFields{
		"remotePeerID": destPeerID,
	})

	for {
		select {
		case <-o.chClose:
			return false
		case payload := <-chSend:
			b := wireEncode(payload)
			_, err := stream.Write(b)
			if err != nil {
				o.logger.Debug("Could not write to stream", types.LogFields{
					"err":          err,
					"remotePeerID": destPeerID,
				})

				return true
			}
		}
	}
}

func (o *ocrEndpoint) runSendToSelf() {
	defer o.wg.Done()
	for {
		select {
		case <-o.chClose:
			return
		case m := <-o.chSendToSelf:
			select {
			case o.recv <- m:
			case <-o.chClose:
				return
			}
		}
	}
}

func (o *ocrEndpoint) Close() error {
	o.stateMu.Lock()
	if o.state != ocrEndpointStarted {
		o.stateMu.Unlock()
		panic("cannot close ocrEndpoint that is not started")
	}
	o.state = ocrEndpointClosed
	o.stateMu.Unlock()

	o.logger.Debug("OCREndpoint: Closing", nil)

	o.logger.Debug("OCREndpoint: Removing stream handler", nil)
	o.peer.RemoveStreamHandler(o.protocolID)

	o.logger.Debug("OCREndpoint: Closing streams", nil)
	close(o.chClose)
	o.ctxCancel()
	o.wg.Wait()

	o.logger.Debug("OCREndpoint: Closing dht", nil)
	err := o.routing.Close()
	if err != nil {
		return errors.Wrap(err, "error closing OCREndpoint: could not close dht")
	}

	o.logger.Debug("OracleGroupMember: Deregister", nil)
	if err := o.peer.deregister(o); err != nil {
		return errors.Wrap(err, "error closing OCREndpoint: could not deregister")
	}

	o.logger.Debug("OCREndpoint: Closing o.recv", nil)
	close(o.recv)

	o.logger.Info("OCREndpoint: Closed", nil)
	return nil
}

func (o *ocrEndpoint) streamReceiver(s p2pnetwork.Stream) {
	exit := make(chan struct{})
	defer close(exit)

	go func() {
		defer s.Reset()
		select {
		case <-o.chClose:
		case <-exit:
		}
	}()

	remotePeerID := s.Conn().RemotePeer()

	o.logger.Debug("Got incoming stream", types.LogFields{
		"remotePeerID":    remotePeerID,
		"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
	})

	sender, err := o.peerID2OracleID(remotePeerID)
	if err != nil {
		o.logger.Error("Error getting sender", types.LogFields{
			"err":             err,
			"remotePeerID":    remotePeerID,
			"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
		})
		return
	}
	r := bufio.NewReader(s)
	for {
		payload, err := readOneFromWire(r)
		if err != nil {
			o.logger.Debug("Lost connection to peer", types.LogFields{
				"err":             err,
				"remotePeerID":    remotePeerID,
				"remoteOracleID":  sender,
				"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
			})
			return
		}

		chRecv := o.chRecvs[sender]
		select {
		case chRecv <- payload:
			continue
		default:
			o.logger.Warn("Incoming buffer is full, dropping message", types.LogFields{
				"remotePeerID":    remotePeerID,
				"remoteOracleID":  sender,
				"remoteMultiaddr": s.Conn().RemoteMultiaddr(),
			})
		}
	}
}

func (o *ocrEndpoint) peerID2OracleID(peerID p2ppeer.ID) (types.OracleID, error) {
	oracleID, ok := o.reversedPeerMapping[peerID]
	if !ok {
		return 0, errors.New("peer ID not found")
	}
	return oracleID, nil
}

func (o *ocrEndpoint) oracleID2PeerID(oracleID types.OracleID) (p2ppeer.ID, error) {
	peerID, ok := o.peerMapping[oracleID]
	if !ok {
		return "", errors.New("oracle ID not found")
	}
	return peerID, nil
}

func (o *ocrEndpoint) isStarted() bool {
	o.stateMu.RLock()
	defer o.stateMu.RUnlock()
	return o.state == ocrEndpointStarted
}

func (o *ocrEndpoint) SendTo(payload []byte, to types.OracleID) {
	if !o.isStarted() {
		panic("send on non-running ocrEndpoint")
	}

	if to == o.ownOracleID {
		o.sendToSelf(payload)
		return
	}

	var chSend chan []byte = o.chSends[to]

	mu := o.muSends[to]
	mu.Lock()
	defer mu.Unlock()

	select {
	case chSend <- payload:
	default:
		select {
		case <-chSend:
			o.logger.Warn("Send buffer full, dropping oldest message", types.LogFields{
				"remoteOracleID": to,
			})
			chSend <- payload
		default:
			chSend <- payload
		}
	}
}

func (o *ocrEndpoint) sendToSelf(payload []byte) {
	m := types.BinaryMessageWithSender{
		Msg:    payload,
		Sender: o.ownOracleID,
	}

	select {
	case o.chSendToSelf <- m:
	default:
		o.logger.Error("Send-to-self buffer is full, dropping message", types.LogFields{
			"remoteOracleID": o.ownOracleID,
		})
	}
}

func (o *ocrEndpoint) Broadcast(payload []byte) {
	var wg sync.WaitGroup
	for oracleID := range o.peerMapping {
		wg.Add(1)
		go func(oid types.OracleID) {
			o.SendTo(payload, oid)
			wg.Done()
		}(oracleID)
	}
	wg.Wait()
}

func (o *ocrEndpoint) Receive() <-chan types.BinaryMessageWithSender {
	return o.recv
}

func (o *ocrEndpoint) isAllowed(id p2ppeer.ID) bool {
	_, ok := o.peerAllowlist[id]
	return ok
}

func (o *ocrEndpoint) allowlist() (allowlist []p2ppeer.ID) {
	for k := range o.peerAllowlist {
		allowlist = append(allowlist, k)
	}
	return
}

func (o *ocrEndpoint) getConfigDigest() types.ConfigDigest {
	return o.configDigest
}
