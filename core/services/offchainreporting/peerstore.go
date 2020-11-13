package offchainreporting

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	p2pPeer struct {
		ID        p2ppeer.ID
		Addr      ma.Multiaddr
		JobID     int32
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	Pstorewrapper struct {
		utils.StartStopOnce
		Peerstore     p2ppeerstore.Peerstore
		jobID         int32
		db            *gorm.DB
		writeInterval time.Duration
		ctx           context.Context
		ctxCancel     context.CancelFunc
		chDone        chan struct{}
	}
)

// NewPeerstoreWrapper creates a new database-backed peerstore wrapper scoped to the given jobID
// Multiple peerstore wrappers should not be instantiated with the same jobID
func NewPeerstoreWrapper(ctx context.Context, db *gorm.DB, writeInterval time.Duration, jobID int32) (*Pstorewrapper, error) {
	ctx, cancel := context.WithCancel(ctx)

	return &Pstorewrapper{
		utils.StartStopOnce{},
		pstoremem.NewPeerstore(),
		jobID,
		db,
		writeInterval,
		ctx,
		cancel,
		make(chan struct{}),
	}, nil
}

func (p *Pstorewrapper) Start() error {
	if !p.OkayToStart() {
		return errors.New("cannot start")
	}
	err := p.readFromDB()
	if err != nil {
		return errors.Wrap(err, "could not start peerstore wrapper")
	}
	go p.dbLoop()
	return nil
}

func (p *Pstorewrapper) dbLoop() {
	defer close(p.chDone)
	ticker := time.NewTicker(p.writeInterval)
	defer ticker.Stop()
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.writeToDB()
		}
	}
}

func (p *Pstorewrapper) Close() error {
	if !p.OkayToStop() {
		return errors.New("cannot stop")
	}
	p.ctxCancel()
	<-p.chDone
	return p.Peerstore.Close()
}

func (p *Pstorewrapper) readFromDB() error {
	peers, err := p.getPeers()
	if err != nil {
		return err
	}
	for _, peer := range peers {
		p.Peerstore.AddAddr(peer.ID, peer.Addr, p2ppeerstore.PermanentAddrTTL)
	}
	return nil
}

func (p *Pstorewrapper) getPeers() (peers []p2pPeer, err error) {
	rows, err := p.db.DB().QueryContext(p.ctx, `SELECT id, addr FROM p2p_peers WHERE job_id = $1`, p.jobID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying peers")
	}
	defer logger.ErrorIfCalling(rows.Close)

	peers = make([]p2pPeer, 0)

	for rows.Next() {
		peer := p2pPeer{}
		var maddr, peerID string
		rows.Scan(&peerID, &maddr)
		peer.ID, err = p2ppeer.Decode(peerID)
		if err != nil {
			return nil, errors.Wrap(err, "unexpectedly failed to decode peer ID")
		}
		peer.Addr, err = ma.NewMultiaddr(maddr)
		if err != nil {
			return nil, errors.Wrap(err, "unexpectedly failed to decode peer multiaddr")
		}
		peers = append(peers, peer)
	}

	return peers, nil
}

func (p *Pstorewrapper) writeToDB() error {
	err := postgres.GormTransaction(p.ctx, p.db, func(tx *gorm.DB) error {
		err := tx.Exec(`DELETE FROM p2p_peerstore WHERE job_id = ?`, p.jobID).Error
		if err != nil {
			return err
		}
		peers := make([]p2pPeer, 0)
		for _, pid := range p.Peerstore.PeersWithAddrs() {
			addrs := p.Peerstore.Addrs(pid)
			for _, addr := range addrs {
				p := p2pPeer{
					ID:    pid,
					Addr:  addr,
					JobID: p.jobID,
				}
				peers = append(peers, p)
			}
		}
		return tx.Create(&peers).Error
	})
	return errors.Wrap(err, "could not write peers to DB")
}
