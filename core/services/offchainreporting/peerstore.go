package offchainreporting

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type (
	P2PPeer struct {
		ID        string
		Addr      string
		PeerID    string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	Pstorewrapper struct {
		utils.StartStopOnce
		Peerstore     p2ppeerstore.Peerstore
		peerID        string
		db            *gorm.DB
		writeInterval time.Duration
		ctx           context.Context
		ctxCancel     context.CancelFunc
		chDone        chan struct{}
	}
)

func (P2PPeer) TableName() string {
	return "p2p_peers"
}

// NewPeerstoreWrapper creates a new database-backed peerstore wrapper scoped to the given jobID
// Multiple peerstore wrappers should not be instantiated with the same jobID
func NewPeerstoreWrapper(db *gorm.DB, writeInterval time.Duration, peerID p2pkey.PeerID) (*Pstorewrapper, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &Pstorewrapper{
		utils.StartStopOnce{},
		pstoremem.NewPeerstore(),
		peerID.Raw(),
		db,
		writeInterval,
		ctx,
		cancel,
		make(chan struct{}),
	}, nil
}

func (p *Pstorewrapper) Start() error {
	return p.StartOnce("PeerStore", func() error {
		err := p.readFromDB()
		if err != nil {
			return errors.Wrap(err, "could not start peerstore wrapper")
		}
		go gracefulpanic.WrapRecover(func() {
			p.dbLoop()
		})
		return nil
	})
}

func (p *Pstorewrapper) dbLoop() {
	defer close(p.chDone)
	ticker := time.NewTicker(utils.WithJitter(p.writeInterval))
	defer ticker.Stop()
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			if err := p.WriteToDB(); err != nil {
				logger.Errorw("Error writing peerstore to DB", "err", err)
			}
		}
	}
}

func (p *Pstorewrapper) Close() error {
	return p.StopOnce("PeerStore", func() error {
		p.ctxCancel()
		<-p.chDone
		return p.Peerstore.Close()
	})
}

func (p *Pstorewrapper) readFromDB() error {
	peers, err := p.getPeers()
	if err != nil {
		return err
	}
	for _, peer := range peers {
		peerID, err := p2ppeer.Decode(peer.ID)
		if err != nil {
			return errors.Wrapf(err, "unexpectedly failed to decode peer ID '%s'", peer.ID)
		}
		peerAddr, err := ma.NewMultiaddr(peer.Addr)
		if err != nil {
			return errors.Wrapf(err, "unexpectedly failed to decode peer multiaddr '%s'", peer.Addr)
		}
		p.Peerstore.AddAddr(peerID, peerAddr, p2ppeerstore.PermanentAddrTTL)
	}
	return nil
}

func (p *Pstorewrapper) getPeers() (peers []P2PPeer, err error) {
	rows, err := p.db.WithContext(p.ctx).Raw(`SELECT id, addr FROM p2p_peers WHERE peer_id = ?`, p.peerID).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "error querying peers")
	}
	defer logger.ErrorIfCalling(rows.Close)

	peers = make([]P2PPeer, 0)

	for rows.Next() {
		peer := P2PPeer{}
		if err = rows.Scan(&peer.ID, &peer.Addr); err != nil {
			return nil, errors.Wrap(err, "unexpected error scanning row")
		}
		peers = append(peers, peer)
	}

	return peers, nil
}

func (p *Pstorewrapper) WriteToDB() error {
	ctx, cancel := context.WithTimeout(p.ctx, postgres.DefaultQueryTimeout)
	defer cancel()
	err := postgres.GormTransaction(ctx, p.db, func(tx *gorm.DB) error {
		err := tx.Exec(`DELETE FROM p2p_peers WHERE peer_id = ?`, p.peerID).Error
		if err != nil {
			return err
		}
		peers := make([]P2PPeer, 0)
		for _, pid := range p.Peerstore.PeersWithAddrs() {
			addrs := p.Peerstore.Addrs(pid)
			for _, addr := range addrs {
				p := P2PPeer{
					ID:     pid.String(),
					Addr:   addr.String(),
					PeerID: p.peerID,
				}
				peers = append(peers, p)
			}
		}
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, p := range peers {
			valueStrings = append(valueStrings, "(?, ?, ?, NOW(), NOW())")
			valueArgs = append(valueArgs, p.ID)
			valueArgs = append(valueArgs, p.Addr)
			valueArgs = append(valueArgs, p.PeerID)
		}

		/* #nosec G201 */
		stmt := fmt.Sprintf("INSERT INTO p2p_peers (id, addr, peer_id, created_at, updated_at) VALUES %s", strings.Join(valueStrings, ","))
		return tx.Exec(stmt, valueArgs...).Error
	})
	return errors.Wrap(err, "could not write peers to DB")
}
