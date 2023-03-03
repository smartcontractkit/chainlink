package dht

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"

	"github.com/libp2p/go-libp2p-kad-dht/metrics"
	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"github.com/libp2p/go-libp2p-kad-dht/providers"
	"github.com/libp2p/go-libp2p-kad-dht/rtrefresh"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"github.com/libp2p/go-libp2p-kbucket/peerdiversity"
	record "github.com/libp2p/go-libp2p-record"
	recpb "github.com/libp2p/go-libp2p-record/pb"

	"github.com/gogo/protobuf/proto"
	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"github.com/multiformats/go-base32"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

var (
	logger     = logging.Logger("dht")
	baseLogger = logger.Desugar()

	rtFreezeTimeout = 1 * time.Minute
)

const (
	// BaseConnMgrScore is the base of the score set on the connection
	// manager "kbucket" tag. It is added with the common prefix length
	// between two peer IDs.
	baseConnMgrScore = 5
)

type mode int

const (
	modeServer mode = iota + 1
	modeClient
)

const (
	kad1 protocol.ID = "/kad/1.0.0"
)

const (
	kbucketTag       = "kbucket"
	protectedBuckets = 2
)

type addPeerRTReq struct {
	p         peer.ID
	queryPeer bool
}

// IpfsDHT is an implementation of Kademlia with S/Kademlia modifications.
// It is used to implement the base Routing module.
type IpfsDHT struct {
	host      host.Host // the network services we need
	self      peer.ID   // Local peer (yourself)
	selfKey   kb.ID
	peerstore peerstore.Peerstore // Peer Registry

	datastore ds.Datastore // Local data

	routingTable *kb.RoutingTable // Array of routing tables for differently distanced nodes
	// ProviderManager stores & manages the provider records for this Dht peer.
	ProviderManager *providers.ProviderManager

	// manages Routing Table refresh
	rtRefreshManager *rtrefresh.RtRefreshManager

	birth time.Time // When this peer started up

	Validator record.Validator

	ctx  context.Context
	proc goprocess.Process

	strmap map[peer.ID]*messageSender
	smlk   sync.Mutex

	plk sync.Mutex

	stripedPutLocks [256]sync.Mutex

	// DHT protocols we query with. We'll only add peers to our routing
	// table if they speak these protocols.
	protocols     []protocol.ID
	protocolsStrs []string

	// DHT protocols we can respond to.
	serverProtocols []protocol.ID

	auto   ModeOpt
	mode   mode
	modeLk sync.Mutex

	bucketSize int
	alpha      int // The concurrency parameter per path
	beta       int // The number of peers closest to a target that must have responded for a query path to terminate

	queryPeerFilter        QueryFilterFunc
	routingTablePeerFilter RouteTableFilterFunc
	rtPeerDiversityFilter  peerdiversity.PeerIPGroupFilter

	autoRefresh bool

	// A set of bootstrap peers to fallback on if all other attempts to fix
	// the routing table fail (or, e.g., this is the first time this node is
	// connecting to the network).
	bootstrapPeers []peer.AddrInfo

	maxRecordAge time.Duration

	// Allows disabling dht subsystems. These should _only_ be set on
	// "forked" DHTs (e.g., DHTs with custom protocols and/or private
	// networks).
	enableProviders, enableValues bool

	disableFixLowPeers bool
	fixLowPeersChan    chan struct{}

	addPeerToRTChan   chan addPeerRTReq
	refreshFinishedCh chan struct{}

	rtFreezeTimeout time.Duration

	// configuration variables for tests
	testAddressUpdateProcessing bool
}

// Assert that IPFS assumptions about interfaces aren't broken. These aren't a
// guarantee, but we can use them to aid refactoring.
var (
	_ routing.ContentRouting = (*IpfsDHT)(nil)
	_ routing.Routing        = (*IpfsDHT)(nil)
	_ routing.PeerRouting    = (*IpfsDHT)(nil)
	_ routing.PubKeyFetcher  = (*IpfsDHT)(nil)
	_ routing.ValueStore     = (*IpfsDHT)(nil)
)

// New creates a new DHT with the specified host and options.
// Please note that being connected to a DHT peer does not necessarily imply that it's also in the DHT Routing Table.
// If the Routing Table has more than "minRTRefreshThreshold" peers, we consider a peer as a Routing Table candidate ONLY when
// we successfully get a query response from it OR if it send us a query.
func New(ctx context.Context, h host.Host, options ...Option) (*IpfsDHT, error) {
	var cfg config
	if err := cfg.apply(append([]Option{defaults}, options...)...); err != nil {
		return nil, err
	}
	if err := cfg.applyFallbacks(h); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	dht, err := makeDHT(ctx, h, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create DHT, err=%s", err)
	}

	dht.autoRefresh = cfg.routingTable.autoRefresh

	dht.maxRecordAge = cfg.maxRecordAge
	dht.enableProviders = cfg.enableProviders
	dht.enableValues = cfg.enableValues
	dht.disableFixLowPeers = cfg.disableFixLowPeers

	dht.Validator = cfg.validator

	dht.testAddressUpdateProcessing = cfg.testAddressUpdateProcessing

	dht.auto = cfg.mode
	switch cfg.mode {
	case ModeAuto, ModeClient:
		dht.mode = modeClient
	case ModeAutoServer, ModeServer:
		dht.mode = modeServer
	default:
		return nil, fmt.Errorf("invalid dht mode %d", cfg.mode)
	}

	if dht.mode == modeServer {
		if err := dht.moveToServerMode(); err != nil {
			return nil, err
		}
	}

	// register for event bus and network notifications
	sn, err := newSubscriberNotifiee(dht)
	if err != nil {
		return nil, err
	}
	dht.proc.Go(sn.subscribe)
	// handle providers
	dht.proc.AddChild(dht.ProviderManager.Process())

	// go-routine to make sure we ALWAYS have RT peer addresses in the peerstore
	// since RT membership is decoupled from connectivity
	go dht.persistRTPeersInPeerStore()

	dht.proc.Go(dht.rtPeerLoop)

	// Fill routing table with currently connected peers that are DHT servers
	dht.plk.Lock()
	for _, p := range dht.host.Network().Peers() {
		dht.peerFound(dht.ctx, p, false)
	}
	dht.plk.Unlock()

	dht.proc.Go(dht.populatePeers)

	return dht, nil
}

// NewDHT creates a new DHT object with the given peer as the 'local' host.
// IpfsDHT's initialized with this function will respond to DHT requests,
// whereas IpfsDHT's initialized with NewDHTClient will not.
func NewDHT(ctx context.Context, h host.Host, dstore ds.Batching) *IpfsDHT {
	dht, err := New(ctx, h, Datastore(dstore))
	if err != nil {
		panic(err)
	}
	return dht
}

// NewDHTClient creates a new DHT object with the given peer as the 'local'
// host. IpfsDHT clients initialized with this function will not respond to DHT
// requests. If you need a peer to respond to DHT requests, use NewDHT instead.
func NewDHTClient(ctx context.Context, h host.Host, dstore ds.Batching) *IpfsDHT {
	dht, err := New(ctx, h, Datastore(dstore), Mode(ModeClient))
	if err != nil {
		panic(err)
	}
	return dht
}

func makeDHT(ctx context.Context, h host.Host, cfg config) (*IpfsDHT, error) {
	var protocols, serverProtocols []protocol.ID

	v1proto := cfg.protocolPrefix + kad1

	if cfg.v1ProtocolOverride != "" {
		v1proto = cfg.v1ProtocolOverride
	}

	protocols = []protocol.ID{v1proto}
	serverProtocols = []protocol.ID{v1proto}

	dht := &IpfsDHT{
		datastore:              cfg.datastore,
		self:                   h.ID(),
		selfKey:                kb.ConvertPeerID(h.ID()),
		peerstore:              h.Peerstore(),
		host:                   h,
		strmap:                 make(map[peer.ID]*messageSender),
		birth:                  time.Now(),
		protocols:              protocols,
		protocolsStrs:          protocol.ConvertToStrings(protocols),
		serverProtocols:        serverProtocols,
		bucketSize:             cfg.bucketSize,
		alpha:                  cfg.concurrency,
		beta:                   cfg.resiliency,
		queryPeerFilter:        cfg.queryPeerFilter,
		routingTablePeerFilter: cfg.routingTable.peerFilter,
		rtPeerDiversityFilter:  cfg.routingTable.diversityFilter,

		fixLowPeersChan: make(chan struct{}, 1),

		addPeerToRTChan:   make(chan addPeerRTReq),
		refreshFinishedCh: make(chan struct{}),
	}

	var maxLastSuccessfulOutboundThreshold time.Duration

	// The threshold is calculated based on the expected amount of time that should pass before we
	// query a peer as part of our refresh cycle.
	// To grok the Math Wizardy that produced these exact equations, please be patient as a document explaining it will
	// be published soon.
	if cfg.concurrency < cfg.bucketSize { // (alpha < K)
		l1 := math.Log(float64(1) / float64(cfg.bucketSize))                              //(Log(1/K))
		l2 := math.Log(float64(1) - (float64(cfg.concurrency) / float64(cfg.bucketSize))) // Log(1 - (alpha / K))
		maxLastSuccessfulOutboundThreshold = time.Duration(l1 / l2 * float64(cfg.routingTable.refreshInterval))
	} else {
		maxLastSuccessfulOutboundThreshold = cfg.routingTable.refreshInterval
	}

	// construct routing table
	// use twice the theoritical usefulness threhold to keep older peers around longer
	rt, err := makeRoutingTable(dht, cfg, 2*maxLastSuccessfulOutboundThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to construct routing table,err=%s", err)
	}
	dht.routingTable = rt
	dht.bootstrapPeers = cfg.bootstrapPeers

	// rt refresh manager
	rtRefresh, err := makeRtRefreshManager(dht, cfg, maxLastSuccessfulOutboundThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to construct RT Refresh Manager,err=%s", err)
	}
	dht.rtRefreshManager = rtRefresh

	// create a DHT proc with the given context
	dht.proc = goprocessctx.WithContextAndTeardown(ctx, func() error {
		return rtRefresh.Close()
	})

	// create a tagged context derived from the original context
	ctxTags := dht.newContextWithLocalTags(ctx)
	// the DHT context should be done when the process is closed
	dht.ctx = goprocessctx.WithProcessClosing(ctxTags, dht.proc)

	pm, err := providers.NewProviderManager(dht.ctx, h.ID(), cfg.datastore, cfg.providersOptions...)
	if err != nil {
		return nil, err
	}
	dht.ProviderManager = pm

	dht.rtFreezeTimeout = rtFreezeTimeout

	return dht, nil
}

func makeRtRefreshManager(dht *IpfsDHT, cfg config, maxLastSuccessfulOutboundThreshold time.Duration) (*rtrefresh.RtRefreshManager, error) {
	keyGenFnc := func(cpl uint) (string, error) {
		p, err := dht.routingTable.GenRandPeerID(cpl)
		return string(p), err
	}

	queryFnc := func(ctx context.Context, key string) error {
		_, err := dht.GetClosestPeers(ctx, key)
		return err
	}

	r, err := rtrefresh.NewRtRefreshManager(
		dht.host, dht.routingTable, cfg.routingTable.autoRefresh,
		keyGenFnc,
		queryFnc,
		cfg.routingTable.refreshQueryTimeout,
		cfg.routingTable.refreshInterval,
		maxLastSuccessfulOutboundThreshold,
		dht.refreshFinishedCh)

	return r, err
}

func makeRoutingTable(dht *IpfsDHT, cfg config, maxLastSuccessfulOutboundThreshold time.Duration) (*kb.RoutingTable, error) {
	// make a Routing Table Diversity Filter
	var filter *peerdiversity.Filter
	if dht.rtPeerDiversityFilter != nil {
		df, err := peerdiversity.NewFilter(dht.rtPeerDiversityFilter, "rt/diversity", func(p peer.ID) int {
			return kb.CommonPrefixLen(dht.selfKey, kb.ConvertPeerID(p))
		})

		if err != nil {
			return nil, fmt.Errorf("failed to construct peer diversity filter: %w", err)
		}

		filter = df
	}

	rt, err := kb.NewRoutingTable(cfg.bucketSize, dht.selfKey, time.Minute, dht.host.Peerstore(), maxLastSuccessfulOutboundThreshold, filter)
	if err != nil {
		return nil, err
	}

	cmgr := dht.host.ConnManager()

	rt.PeerAdded = func(p peer.ID) {
		commonPrefixLen := kb.CommonPrefixLen(dht.selfKey, kb.ConvertPeerID(p))
		if commonPrefixLen < protectedBuckets {
			cmgr.Protect(p, kbucketTag)
		} else {
			cmgr.TagPeer(p, kbucketTag, baseConnMgrScore)
		}
	}
	rt.PeerRemoved = func(p peer.ID) {
		cmgr.Unprotect(p, kbucketTag)
		cmgr.UntagPeer(p, kbucketTag)

		// try to fix the RT
		dht.fixRTIfNeeded()
	}

	return rt, err
}

// GetRoutingTableDiversityStats returns the diversity stats for the Routing Table.
func (dht *IpfsDHT) GetRoutingTableDiversityStats() []peerdiversity.CplDiversityStats {
	return dht.routingTable.GetDiversityStats()
}

// Mode allows introspection of the operation mode of the DHT
func (dht *IpfsDHT) Mode() ModeOpt {
	return dht.auto
}

func (dht *IpfsDHT) populatePeers(_ goprocess.Process) {
	if !dht.disableFixLowPeers {
		dht.fixLowPeers(dht.ctx)
	}

	if err := dht.rtRefreshManager.Start(); err != nil {
		logger.Error(err)
	}

	// listens to the fix low peers chan and tries to fix the Routing Table
	if !dht.disableFixLowPeers {
		dht.proc.Go(dht.fixLowPeersRoutine)
	}

}

// fixLowPeersRouting manages simultaneous requests to fixLowPeers
func (dht *IpfsDHT) fixLowPeersRoutine(proc goprocess.Process) {
	ticker := time.NewTicker(periodicBootstrapInterval)
	defer ticker.Stop()

	for {
		select {
		case <-dht.fixLowPeersChan:
		case <-ticker.C:
		case <-proc.Closing():
			return
		}

		dht.fixLowPeers(dht.Context())
	}

}

// fixLowPeers tries to get more peers into the routing table if we're below the threshold
func (dht *IpfsDHT) fixLowPeers(ctx context.Context) {
	if dht.routingTable.Size() > minRTRefreshThreshold {
		return
	}

	// we try to add all peers we are connected to to the Routing Table
	// in case they aren't already there.
	for _, p := range dht.host.Network().Peers() {
		dht.peerFound(ctx, p, false)
	}

	// TODO Active Bootstrapping
	// We should first use non-bootstrap peers we knew of from previous
	// snapshots of the Routing Table before we connect to the bootstrappers.
	// See https://github.com/libp2p/go-libp2p-kad-dht/issues/387.
	if dht.routingTable.Size() == 0 {
		if len(dht.bootstrapPeers) == 0 {
			// No point in continuing, we have no peers!
			return
		}

		found := 0
		for _, i := range rand.Perm(len(dht.bootstrapPeers)) {
			ai := dht.bootstrapPeers[i]
			err := dht.Host().Connect(ctx, ai)
			if err == nil {
				found++
			} else {
				logger.Warnw("failed to bootstrap", "peer", ai.ID, "error", err)
			}

			// Wait for two bootstrap peers, or try them all.
			//
			// Why two? In theory, one should be enough
			// normally. However, if the network were to
			// restart and everyone connected to just one
			// bootstrapper, we'll end up with a mostly
			// partitioned network.
			//
			// So we always bootstrap with two random peers.
			if found == maxNBoostrappers {
				break
			}
		}
	}

	// if we still don't have peers in our routing table(probably because Identify hasn't completed),
	// there is no point in triggering a Refresh.
	if dht.routingTable.Size() == 0 {
		return
	}

	if dht.autoRefresh {
		dht.rtRefreshManager.RefreshNoWait()
	}
}

// TODO This is hacky, horrible and the programmer needs to have his mother called a hamster.
// SHOULD be removed once https://github.com/libp2p/go-libp2p/issues/800 goes in.
func (dht *IpfsDHT) persistRTPeersInPeerStore() {
	tickr := time.NewTicker(peerstore.RecentlyConnectedAddrTTL / 3)
	defer tickr.Stop()

	for {
		select {
		case <-tickr.C:
			ps := dht.routingTable.ListPeers()
			for _, p := range ps {
				dht.peerstore.UpdateAddrs(p, peerstore.RecentlyConnectedAddrTTL, peerstore.RecentlyConnectedAddrTTL)
			}
		case <-dht.ctx.Done():
			return
		}
	}
}

// putValueToPeer stores the given key/value pair at the peer 'p'
func (dht *IpfsDHT) putValueToPeer(ctx context.Context, p peer.ID, rec *recpb.Record) error {
	pmes := pb.NewMessage(pb.Message_PUT_VALUE, rec.Key, 0)
	pmes.Record = rec
	rpmes, err := dht.sendRequest(ctx, p, pmes)
	if err != nil {
		logger.Debugw("failed to put value to peer", "to", p, "key", loggableRecordKeyBytes(rec.Key), "error", err)
		return err
	}

	if !bytes.Equal(rpmes.GetRecord().Value, pmes.GetRecord().Value) {
		logger.Infow("value not put correctly", "put-message", pmes, "get-message", rpmes)
		return errors.New("value not put correctly")
	}

	return nil
}

var errInvalidRecord = errors.New("received invalid record")

// getValueOrPeers queries a particular peer p for the value for
// key. It returns either the value or a list of closer peers.
// NOTE: It will update the dht's peerstore with any new addresses
// it finds for the given peer.
func (dht *IpfsDHT) getValueOrPeers(ctx context.Context, p peer.ID, key string) (*recpb.Record, []*peer.AddrInfo, error) {
	pmes, err := dht.getValueSingle(ctx, p, key)
	if err != nil {
		return nil, nil, err
	}

	// Perhaps we were given closer peers
	peers := pb.PBPeersToPeerInfos(pmes.GetCloserPeers())

	if rec := pmes.GetRecord(); rec != nil {
		// Success! We were given the value
		logger.Debug("got value")

		// make sure record is valid.
		err = dht.Validator.Validate(string(rec.GetKey()), rec.GetValue())
		if err != nil {
			logger.Debug("received invalid record (discarded)")
			// return a sentinal to signify an invalid record was received
			err = errInvalidRecord
			rec = new(recpb.Record)
		}
		return rec, peers, err
	}

	if len(peers) > 0 {
		return nil, peers, nil
	}

	return nil, nil, routing.ErrNotFound
}

// getValueSingle simply performs the get value RPC with the given parameters
func (dht *IpfsDHT) getValueSingle(ctx context.Context, p peer.ID, key string) (*pb.Message, error) {
	pmes := pb.NewMessage(pb.Message_GET_VALUE, []byte(key), 0)
	return dht.sendRequest(ctx, p, pmes)
}

// getLocal attempts to retrieve the value from the datastore
func (dht *IpfsDHT) getLocal(key string) (*recpb.Record, error) {
	logger.Debugw("finding value in datastore", "key", loggableRecordKeyString(key))

	rec, err := dht.getRecordFromDatastore(mkDsKey(key))
	if err != nil {
		logger.Warnw("get local failed", "key", loggableRecordKeyString(key), "error", err)
		return nil, err
	}

	// Double check the key. Can't hurt.
	if rec != nil && string(rec.GetKey()) != key {
		logger.Errorw("BUG: found a DHT record that didn't match it's key", "expected", loggableRecordKeyString(key), "got", rec.GetKey())
		return nil, nil

	}
	return rec, nil
}

// putLocal stores the key value pair in the datastore
func (dht *IpfsDHT) putLocal(key string, rec *recpb.Record) error {
	data, err := proto.Marshal(rec)
	if err != nil {
		logger.Warnw("failed to put marshal record for local put", "error", err, "key", loggableRecordKeyString(key))
		return err
	}

	return dht.datastore.Put(mkDsKey(key), data)
}

func (dht *IpfsDHT) rtPeerLoop(proc goprocess.Process) {
	bootstrapCount := 0
	isBootsrapping := false
	var timerCh <-chan time.Time

	for {
		select {
		case <-timerCh:
			dht.routingTable.MarkAllPeersIrreplaceable()
		case addReq := <-dht.addPeerToRTChan:
			prevSize := dht.routingTable.Size()
			if prevSize == 0 {
				isBootsrapping = true
				bootstrapCount = 0
				timerCh = nil
			}
			newlyAdded, err := dht.routingTable.TryAddPeer(addReq.p, addReq.queryPeer, isBootsrapping)
			if err != nil {
				// peer not added.
				continue
			}
			if !newlyAdded && addReq.queryPeer {
				// the peer is already in our RT, but we just successfully queried it and so let's give it a
				// bump on the query time so we don't ping it too soon for a liveliness check.
				dht.routingTable.UpdateLastSuccessfulOutboundQueryAt(addReq.p, time.Now())
			}
		case <-dht.refreshFinishedCh:
			bootstrapCount = bootstrapCount + 1
			if bootstrapCount == 2 {
				timerCh = time.NewTimer(dht.rtFreezeTimeout).C
			}

			old := isBootsrapping
			isBootsrapping = false
			if old {
				dht.rtRefreshManager.RefreshNoWait()
			}

		case <-proc.Closing():
			return
		}
	}
}

// peerFound signals the routingTable that we've found a peer that
// might support the DHT protocol.
// If we have a connection a peer but no exchange of a query RPC ->
//    LastQueriedAt=time.Now (so we don't ping it for some time for a liveliness check)
//    LastUsefulAt=0
// If we connect to a peer and then exchange a query RPC ->
//    LastQueriedAt=time.Now (same reason as above)
//    LastUsefulAt=time.Now (so we give it some life in the RT without immediately evicting it)
// If we query a peer we already have in our Routing Table ->
//    LastQueriedAt=time.Now()
//    LastUsefulAt remains unchanged
// If we connect to a peer we already have in the RT but do not exchange a query (rare)
//    Do Nothing.
func (dht *IpfsDHT) peerFound(ctx context.Context, p peer.ID, queryPeer bool) {
	if c := baseLogger.Check(zap.DebugLevel, "peer found"); c != nil {
		c.Write(zap.String("peer", p.String()))
	}
	b, err := dht.validRTPeer(p)
	if err != nil {
		logger.Errorw("failed to validate if peer is a DHT peer", "peer", p, "error", err)
	} else if b {
		select {
		case dht.addPeerToRTChan <- addPeerRTReq{p, queryPeer}:
		case <-dht.ctx.Done():
			return
		}
	}
}

// peerStoppedDHT signals the routing table that a peer is unable to responsd to DHT queries anymore.
func (dht *IpfsDHT) peerStoppedDHT(ctx context.Context, p peer.ID) {
	logger.Debugw("peer stopped dht", "peer", p)
	// A peer that does not support the DHT protocol is dead for us.
	// There's no point in talking to anymore till it starts supporting the DHT protocol again.
	dht.routingTable.RemovePeer(p)
}

func (dht *IpfsDHT) fixRTIfNeeded() {
	select {
	case dht.fixLowPeersChan <- struct{}{}:
	default:
	}
}

// FindLocal looks for a peer with a given ID connected to this dht and returns the peer and the table it was found in.
func (dht *IpfsDHT) FindLocal(id peer.ID) peer.AddrInfo {
	switch dht.host.Network().Connectedness(id) {
	case network.Connected, network.CanConnect:
		return dht.peerstore.PeerInfo(id)
	default:
		return peer.AddrInfo{}
	}
}

// findPeerSingle asks peer 'p' if they know where the peer with id 'id' is
func (dht *IpfsDHT) findPeerSingle(ctx context.Context, p peer.ID, id peer.ID) (*pb.Message, error) {
	pmes := pb.NewMessage(pb.Message_FIND_NODE, []byte(id), 0)
	return dht.sendRequest(ctx, p, pmes)
}

func (dht *IpfsDHT) findProvidersSingle(ctx context.Context, p peer.ID, key multihash.Multihash) (*pb.Message, error) {
	pmes := pb.NewMessage(pb.Message_GET_PROVIDERS, key, 0)
	return dht.sendRequest(ctx, p, pmes)
}

// nearestPeersToQuery returns the routing tables closest peers.
func (dht *IpfsDHT) nearestPeersToQuery(pmes *pb.Message, count int) []peer.ID {
	closer := dht.routingTable.NearestPeers(kb.ConvertKey(string(pmes.GetKey())), count)
	return closer
}

// betterPeersToQuery returns nearestPeersToQuery with some additional filtering
func (dht *IpfsDHT) betterPeersToQuery(pmes *pb.Message, from peer.ID, count int) []peer.ID {
	closer := dht.nearestPeersToQuery(pmes, count)

	// no node? nil
	if closer == nil {
		logger.Infow("no closer peers to send", from)
		return nil
	}

	filtered := make([]peer.ID, 0, len(closer))
	for _, clp := range closer {

		// == to self? thats bad
		if clp == dht.self {
			logger.Error("BUG betterPeersToQuery: attempted to return self! this shouldn't happen...")
			return nil
		}
		// Dont send a peer back themselves
		if clp == from {
			continue
		}

		filtered = append(filtered, clp)
	}

	// ok seems like closer nodes
	return filtered
}

func (dht *IpfsDHT) setMode(m mode) error {
	dht.modeLk.Lock()
	defer dht.modeLk.Unlock()

	if m == dht.mode {
		return nil
	}

	switch m {
	case modeServer:
		return dht.moveToServerMode()
	case modeClient:
		return dht.moveToClientMode()
	default:
		return fmt.Errorf("unrecognized dht mode: %d", m)
	}
}

// moveToServerMode advertises (via libp2p identify updates) that we are able to respond to DHT queries and sets the appropriate stream handlers.
// Note: We may support responding to queries with protocols aside from our primary ones in order to support
// interoperability with older versions of the DHT protocol.
func (dht *IpfsDHT) moveToServerMode() error {
	dht.mode = modeServer
	for _, p := range dht.serverProtocols {
		dht.host.SetStreamHandler(p, dht.handleNewStream)
	}
	return nil
}

// moveToClientMode stops advertising (and rescinds advertisements via libp2p identify updates) that we are able to
// respond to DHT queries and removes the appropriate stream handlers. We also kill all inbound streams that were
// utilizing the handled protocols.
// Note: We may support responding to queries with protocols aside from our primary ones in order to support
// interoperability with older versions of the DHT protocol.
func (dht *IpfsDHT) moveToClientMode() error {
	dht.mode = modeClient
	for _, p := range dht.serverProtocols {
		dht.host.RemoveStreamHandler(p)
	}

	pset := make(map[protocol.ID]bool)
	for _, p := range dht.serverProtocols {
		pset[p] = true
	}

	for _, c := range dht.host.Network().Conns() {
		for _, s := range c.GetStreams() {
			if pset[s.Protocol()] {
				if s.Stat().Direction == network.DirInbound {
					_ = s.Reset()
				}
			}
		}
	}
	return nil
}

func (dht *IpfsDHT) getMode() mode {
	dht.modeLk.Lock()
	defer dht.modeLk.Unlock()
	return dht.mode
}

// Context returns the DHT's context.
func (dht *IpfsDHT) Context() context.Context {
	return dht.ctx
}

// Process returns the DHT's process.
func (dht *IpfsDHT) Process() goprocess.Process {
	return dht.proc
}

// RoutingTable returns the DHT's routingTable.
func (dht *IpfsDHT) RoutingTable() *kb.RoutingTable {
	return dht.routingTable
}

// Close calls Process Close.
func (dht *IpfsDHT) Close() error {
	return dht.proc.Close()
}

func mkDsKey(s string) ds.Key {
	return ds.NewKey(base32.RawStdEncoding.EncodeToString([]byte(s)))
}

// PeerID returns the DHT node's Peer ID.
func (dht *IpfsDHT) PeerID() peer.ID {
	return dht.self
}

// PeerKey returns a DHT key, converted from the DHT node's Peer ID.
func (dht *IpfsDHT) PeerKey() []byte {
	return kb.ConvertPeerID(dht.self)
}

// Host returns the libp2p host this DHT is operating with.
func (dht *IpfsDHT) Host() host.Host {
	return dht.host
}

// Ping sends a ping message to the passed peer and waits for a response.
func (dht *IpfsDHT) Ping(ctx context.Context, p peer.ID) error {
	req := pb.NewMessage(pb.Message_PING, nil, 0)
	resp, err := dht.sendRequest(ctx, p, req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	if resp.Type != pb.Message_PING {
		return fmt.Errorf("got unexpected response type: %v", resp.Type)
	}
	return nil
}

// newContextWithLocalTags returns a new context.Context with the InstanceID and
// PeerID keys populated. It will also take any extra tags that need adding to
// the context as tag.Mutators.
func (dht *IpfsDHT) newContextWithLocalTags(ctx context.Context, extraTags ...tag.Mutator) context.Context {
	extraTags = append(
		extraTags,
		tag.Upsert(metrics.KeyPeerID, dht.self.Pretty()),
		tag.Upsert(metrics.KeyInstanceID, fmt.Sprintf("%p", dht)),
	)
	ctx, _ = tag.New(
		ctx,
		extraTags...,
	) // ignoring error as it is unrelated to the actual function of this code.
	return ctx
}

func (dht *IpfsDHT) maybeAddAddrs(p peer.ID, addrs []ma.Multiaddr, ttl time.Duration) {
	// Don't add addresses for self or our connected peers. We have better ones.
	if p == dht.self || dht.host.Network().Connectedness(p) == network.Connected {
		return
	}
	dht.peerstore.AddAddrs(p, addrs, ttl)
}
